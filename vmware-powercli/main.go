package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/acarl005/stripansi"
	"github.com/direktiv/direktiv-apps/pkg/reusable"
	pss "github.com/direktiv/go-powershell"
	"github.com/direktiv/go-powershell/backend"
)

type Script struct {
	Script reusable.File `json:"script"`
	Args   []string      `json:"args"`
	Output string        `json:"output"`
}

type CmdOutput struct {
	Success bool        `json:"result"`
	Stdout  interface{} `json:"stdout"`
	Stderr  interface{} `json:"stderr"`
	Output  interface{} `json:"output"`
}

type requestInput struct {
	Scripts  []Script `json:"scripts"`
	Continue bool     `json:"continue"`
	Host     string   `json:"host"`
	User     string   `json:"user"`
	Pwd      string   `json:"password"`
}

func powercliHandler(w http.ResponseWriter, r *http.Request, ri *reusable.RequestInfo) {

	obj := new(requestInput)
	err := reusable.Unmarshal(obj, true, r)
	if err != nil {
		reusable.ReportError(w, reusable.UnmarshallError, err)
		return
	}

	// create the shell
	sh, err := getShell(ri, obj.Host, obj.User, obj.Pwd)
	if err != nil {
		reusable.ReportError(w, errForCode("shell"), err)
		return
	}
	defer sh.Exit()

	object := make(map[string]*CmdOutput)

	for i := range obj.Scripts {
		co := &CmdOutput{
			Success: true,
		}
		stdout, stderr, err := executeScript(ri, sh, obj.Scripts[i])

		if err != nil {
			ri.Logger().Errorf("error executing script: %v", err)

			// stop if one script fails
			if !obj.Continue {
				reusable.ReportError(w, reusable.UnmarshallError, err)
				return
			}
			co.Success = false
		}

		// add to result
		co.Stderr = stderr
		co.Stdout = toJSON(ri, stdout)

		if co.Success && len(obj.Scripts[i].Output) > 0 {
			ri.Logger().Infof("getting output file")

			uf := fmt.Sprintf("%s/%s", ri.Dir(), obj.Scripts[i].Output)

			b, err := os.ReadFile(uf)
			if err != nil {
				reusable.ReportError(w, errForCode("output"), err)
				return
			}

			ret := make(map[string]interface{})
			err = json.Unmarshal(b, &ret)
			if err != nil {
				reusable.ReportError(w, errForCode("json"), err)
				return
			}
			defer os.Remove(uf)
			co.Output = ret
		}

		object[obj.Scripts[i].Script.Name] = co

	}

	reusable.ReportResult(w, object)

}

func executeScript(ri *reusable.RequestInfo, shell pss.Shell, script Script) (string, string, error) {

	// create file
	f, err := script.Script.AsFile(0755)
	if err != nil {
		ri.Logger().Errorf("error creating script: %v", err)
		return "", "", err
	}

	ps1 := fmt.Sprintf("%s.ps1", f.Name())

	os.Rename(f.Name(), ps1)
	defer os.Remove(ps1)

	full := fmt.Sprintf("%s %s", ps1, strings.Join(script.Args, " "))
	ri.Logger().Infof("running %v", full)

	stdout, stderr, err := shell.Execute(full)
	if err != nil {
		ri.Logger().Errorf("error executing script: %v", err)
		return stdout, stderr, err
	}
	printOut(ri, stdout, stderr)

	return stdout, stderr, err
}

func getShell(ri *reusable.RequestInfo, server, user, password string) (pss.Shell, error) {

	ri.Logger().Infof("creating shell")

	back := &backend.Local{
		Dir: ri.Dir(),
	}
	shell, err := pss.New(back)
	if err != nil {
		ri.Logger().Errorf("error creating new shell: %v", err)
		return nil, err
	}

	if server != "" {
		ri.Logger().Infof("connecting to %s", server)
		o, e, err := shell.Execute(fmt.Sprintf("Connect-VIServer -Server %s -User %s -Password %s", server, user, password))
		if err != nil {
			ri.Logger().Errorf("error connecting to %s: %v", server, err)
			return nil, err
		}
		printOut(ri, o, e)
	}

	_, _, err = shell.Execute("$WarningPreference = 'SilentlyContinue'")
	if err != nil {
		ri.Logger().Errorf("error executing SilentlyContinue: %v", err)
		return nil, err
	}

	return shell, nil

}

func main() {
	reusable.StartServer(powercliHandler, nil)
}

func errForCode(errCode string) string {
	return fmt.Sprintf("com.powershell.%s.error", errCode)
}

func printOut(ri *reusable.RequestInfo, o, e string) {
	if len(o) > 0 {
		ri.Logger().Infof(fmt.Sprintf("shell stdout: %v", o))
	}
	if len(e) > 0 {
		ri.Logger().Infof(fmt.Sprintf("shell stderr: %v", e))
	}
}

func toJSON(ri *reusable.RequestInfo, str string) interface{} {

	str = strings.TrimSpace(str)
	str = stripansi.Strip(str)

	var js json.RawMessage
	err := json.Unmarshal([]byte(str), &js)
	if err != nil {

		if e, ok := err.(*json.SyntaxError); ok {
			ri.Logger().Infof("syntax error at byte offset %d", e.Offset)
		}

		ri.Logger().Infof("response is string: %v", err)
		return str
	}

	ri.Logger().Infof("response is json")
	return json.RawMessage(str)

}
