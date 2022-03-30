package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/acarl005/stripansi"
	"github.com/bramvdbogaerde/go-scp"
	"github.com/bramvdbogaerde/go-scp/auth"
	"github.com/direktiv/direktiv-apps/pkg/reusable"
	"github.com/google/uuid"

	"golang.org/x/crypto/ssh"
)

const (
	sshType     = "ssh"
	scpType     = "scp"
	scpFromType = "scpfrom"
)

type sshscp struct {
	Type     string          `json:"type"` // scp/ssh
	Host     string          `json:"host"` // user@server:/tmp
	Files    []reusable.File `json:"files"`
	Auth     string          `json:"auth"`
	Port     int             `json:"port"`
	Continue bool            `json:"continue"`
	Output   string          `json:"output"`

	Args []string `json:"args"`
}

type requestInput struct {
	SSHSCP   []sshscp `json:"actions"`
	Continue bool     `json:"continue"`
	Silent   bool     `json:"silent"`
}

type connector struct {
	user, host, path string
}

func parseURL(url, scpssh string) (*connector, error) {

	// parse url
	u := strings.SplitN(url, "@", 2)
	if len(u) != 2 {
		return nil, fmt.Errorf("host %s is not the right format: user@hostname", url)
	}
	p := strings.SplitN(u[1], ":", 2)

	c := &connector{
		user: u[0],
		host: p[0],
	}

	if len(p) == 2 {
		c.path = p[1]
	}

	return c, nil

}

func generateAuth(attached []reusable.File, pwd, user string, ri *reusable.RequestInfo) (ssh.ClientConfig, bool, error) {

	var isCert bool
	var cc ssh.ClientConfig
	var err error

	for a := range attached {
		f := attached[a]
		if f.Name == pwd {

			ri.Logger().Infof("using %s as certficate", f.Name)
			cert, err := f.AsFile(ri, 0400)
			if err != nil {
				return cc, isCert, err
			}

			// with cert
			cc, err = auth.PrivateKey(user, cert.Name(), ssh.InsecureIgnoreHostKey())
			if err != nil {
				return cc, isCert, err
			}
			isCert = true
		}
	}

	if len(cc.Auth) == 0 {
		ri.Logger().Infof("using password authentication")
		// TODO: we probably have to add something to use known_hosts as input
		cc, err = auth.PasswordKey(user, pwd, ssh.InsecureIgnoreHostKey())
		if err != nil {
			return cc, isCert, err
		}
	}

	return cc, isCert, nil

}

func connect(svr string, cc ssh.ClientConfig) (scp.Client, error) {
	cc.Timeout = time.Duration(10 * time.Second)
	client := scp.NewClient(svr, &cc)
	err := client.Connect()
	return client, err
}

func scpFromExec(s sshscp, c *connector, ri *reusable.RequestInfo) ([]string, error) {

	var (
		err error
		cc  ssh.ClientConfig
	)

	files := []string{}

	cc, _, err = generateAuth(s.Files, s.Auth, c.user, ri)
	if err != nil {
		return files, err
	}

	if len(c.path) == 0 {
		return files, fmt.Errorf("no path specified for scp")
	}

	svr := fmt.Sprintf("%s:%d", c.host, s.Port)
	ri.Logger().Infof("connecting to %s", svr)

	client, err := connect(svr, cc)
	if err != nil {
		return files, err
	}
	defer client.Close()

	target := filepath.Join(ri.Dir(), "out", s.Output)
	f, err := os.Create(target)
	if err != nil {
		return files, err
	}
	defer f.Close()

	ri.Logger().Infof("fetching file %s", c.path)
	err = client.CopyFromRemote(context.Background(), f, c.path)
	if err != nil {
		return files, err
	}

	ri.Logger().Infof("fetching file finished")
	files = append(files, target)

	return files, nil
}

func scpExec(s sshscp, c *connector, ri *reusable.RequestInfo) ([]string, error) {

	var (
		err    error
		cc     ssh.ClientConfig
		isCert bool
	)

	files := []string{}

	cc, isCert, err = generateAuth(s.Files, s.Auth, c.user, ri)
	if err != nil {
		return files, err
	}

	if len(c.path) == 0 {
		return files, fmt.Errorf("no path specified for scp")
	}

	svr := fmt.Sprintf("%s:%d", c.host, s.Port)
	ri.Logger().Infof("connecting to %s", svr)

	for a := range s.Files {
		f := s.Files[a]

		client, err := connect(svr, cc)
		if err != nil {
			return files, err
		}
		defer client.Close()

		// don't copy the certificate
		if isCert && f.Name == s.Auth {
			continue
		}

		r, err := f.AsFile(ri, 0)
		if err != nil {
			return files, err
		}
		defer func(ff *os.File) {
			r.Close()
			os.Remove(r.Name())
		}(r)

		fi, err := r.Stat()
		if err != nil {
			return files, err
		}

		smode := strconv.FormatInt(int64(fi.Mode()), 8)

		path := path.Join(c.path, f.Name)
		ri.Logger().Infof("copying %v", path)

		err = client.CopyFromFile(context.Background(), *r, path, fmt.Sprintf("0%s", smode))

		if err != nil && s.Continue {
			ri.Logger().Infof("could not copy file %s: %v", path, err)
			files = append(files, fmt.Sprintf("error %v", path))
		} else if err != nil {
			ri.Logger().Infof("copying failed: %v", err)
			return files, err
		} else {
			files = append(files, path)
		}

	}

	return files, nil
}

func sshExec(s sshscp, c *connector, silent bool, ri *reusable.RequestInfo) ([]map[string]interface{}, error) {

	var (
		err    error
		cc     ssh.ClientConfig
		isCert bool
	)

	var ret []map[string]interface{}

	cc, isCert, err = generateAuth(s.Files, s.Auth, c.user, ri)
	if err != nil {
		return ret, err
	}

	svr := fmt.Sprintf("%s:%d", c.host, s.Port)

	ri.Logger().Infof("connecting to %s", svr)

	conn, err := ssh.Dial("tcp", svr, &cc)
	if err != nil {
		return ret, err
	}
	defer conn.Close()

	for a := range s.Files {
		f := s.Files[a]

		// don't process cert
		if isCert && f.Name == s.Auth {
			continue
		}

		if f.Name == "" {
			f.Name = uuid.New().String()
		}

		result := make(map[string]interface{})

		result["script"] = f.Name

		sess, err := conn.NewSession()
		if err != nil {
			return ret, err
		}
		defer sess.Close()

		ri.Logger().Infof("running %v", f.Name)

		start, err := f.AsFile(ri, 0755)
		if err != nil {
			return ret, err
		}
		defer func() {
			start.Close()
			os.Remove(start.Name())
		}()

		b, err := io.ReadAll(start)
		if err != nil {
			return ret, err
		}

		execute := string(b)
		isScript := false

		if strings.HasPrefix(string(b)[0:2], "#!") {

			// if it is a script we scp it over and delete it later
			start.Seek(0, 0)

			client, err := scp.NewClientBySSH(conn)
			if err != nil {
				return ret, err
			}
			err = client.CopyFromFile(context.Background(), *start, fmt.Sprintf("/tmp/%s", f.Name), "0755")
			if err != nil {
				return ret, err
			}
			execute = fmt.Sprintf("/tmp/%s", f.Name)
			isScript = true
		}

		if len(s.Args) > 0 {
			execute = fmt.Sprintf("%s %s", execute, strings.Join(s.Args, " "))
		}

		ri.Logger().Infof("executing command %d", a)

		var bout bytes.Buffer
		mw := io.MultiWriter(&bout, os.Stdout)
		if !silent {
			mw = io.MultiWriter(&bout, ri.LogWriter(), os.Stdout)
		}

		sess.Stdout = mw
		sess.Stderr = mw

		err = sess.Run(execute)

		result["success"] = true
		if err != nil {
			result["success"] = false
			result["error"] = err.Error()
		}

		if !silent {
			result["stdout"] = bout.String()
		}

		ret = append(ret, result)

		// even with errors we are trying to remove the file
		if isScript {
			sess.Run(fmt.Sprintf("rm -rf %s", execute))
		}

		// load output
		if s.Output != "" {
			catSess, err := conn.NewSession()
			if err != nil {
				ri.Logger().Infof("can not read output: %v", err)
			}
			defer catSess.Close()
			b, err = catSess.Output(fmt.Sprintf("cat %s; rm -Rf %s", s.Output, s.Output))
			if err != nil {
				ri.Logger().Infof("can not read output: %v", err)
			}
			result["output"] = toJSON(ri, string(b))
		}

		if err != nil && s.Continue {
			ri.Logger().Infof("could not execute script %s: %v", execute, err)
		} else if err != nil {
			return ret, err
		}

	}

	return ret, nil
}

func toJSON(ri *reusable.RequestInfo, str string) interface{} {

	str = strings.TrimSpace(str)
	str = stripansi.Strip(str)

	var js json.RawMessage
	err := json.Unmarshal([]byte(str), &js)
	if err != nil {
		// it is a string
		return str
	}

	return json.RawMessage(str)

}

func sshscpHandler(w http.ResponseWriter, r *http.Request, ri *reusable.RequestInfo) {

	var (
		err    error
		files  []string
		c      *connector
		result []map[string]interface{}
	)

	obj := new(requestInput)
	err = reusable.Unmarshal(obj, true, r)
	if err != nil {
		reusable.ReportError(w, reusable.UnmarshallError, err)
		return
	}

	scps := make(map[string][]string)
	scpsFrom := make(map[string][]string)
	sshs := make(map[string][]map[string]interface{})

	ret := make(map[string]interface{})

	for i := range obj.SSHSCP {
		s := obj.SSHSCP[i]

		c, err = parseURL(s.Host, s.Type)
		if err != nil {
			reusable.ReportError(w, errForCode("file"), err)
			return
		}

		// set standard port
		if s.Port == 0 {
			s.Port = 22
		}

		switch s.Type {
		case sshType:
			result, err = sshExec(s, c, obj.Silent, ri)

			// if there is an error we ned to create the success
			// layout and add the error. flows can test for 'failed'
			if err != nil {
				result = []map[string]interface{}{
					make(map[string]interface{}),
				}
				result[0]["failed"] = err
			}

			sshs[c.host] = result
		case scpType:
			files, err = scpExec(s, c, ri)
			scps[c.host] = files
		case scpFromType:
			files, err = scpFromExec(s, c, ri)
			scpsFrom[c.host] = files
		default:
			reusable.ReportError(w, errForCode("sshscp"), fmt.Errorf("an action has to be ssh, scp or scpfrom"))
			return
		}

		if err != nil && !obj.Continue {
			reusable.ReportError(w, errForCode("sshscp"), err)
			return
		} else if err != nil {
			ri.Logger().Infof("could not %s to %s", s.Type, c.host)
		}

	}

	ret[scpType] = scps
	ret[sshType] = sshs
	ret[scpFromType] = scpsFrom

	reusable.ReportResult(w, ret)

}

func main() {
	os.Mkdir("/tmp", 0755)
	reusable.StartServer(sshscpHandler, nil)
}

func errForCode(errCode string) string {
	return fmt.Sprintf("com.sshscp.%s.error", errCode)
}
