package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"

	"github.com/direktiv/direktiv-apps/pkg/reusable"
)

type requestInput struct {
	Script reusable.File   `json:"script"`
	Args   []string        `json:"args"`
	Envs   []string        `json:"envs"`
	Output string          `json:"output"`
	Files  []reusable.File `json:"files"`
}

func ubuntuHandler(w http.ResponseWriter, r *http.Request, ri *reusable.RequestInfo) {

	obj := new(requestInput)
	err := reusable.Unmarshal(obj, false, r)
	if err != nil {
		reusable.ReportError(w, reusable.UnmarshallError, err)
		return
	}

	err = writeFiles(ri.Dir(), reusable.NewFileIterator(obj.Files, ri))
	if err != nil {
		reusable.ReportError(w, errForCode("files"), err)
		return
	}

	file, err := obj.Script.AsFile(0755)
	if err != nil {
		reusable.ReportError(w, errForCode("execute"), err)
		return
	}
	defer os.Remove(file.Name())
	file.Close()

	cmd := exec.Command(file.Name(), obj.Args...)
	obj.Envs = append(obj.Envs, fmt.Sprintf("Direktiv_TempDir=%s", ri.Dir()))

	cmd.Stderr = ri.LogWriter()
	cmd.Stdout = ri.LogWriter()

	cmd.Env = obj.Envs
	cmd.Dir = ri.Dir()

	ri.Logger().Infof("running %s", cmd)

	if err := cmd.Start(); err != nil {
		reusable.ReportError(w, errForCode("execute"), err)
		return
	}
	if err := cmd.Wait(); err != nil {
		reusable.ReportError(w, errForCode("execute"), err)
		return
	}

	// read output file and use it as return
	if len(obj.Output) > 0 {
		ri.Logger().Infof("preparing output %s", obj.Output)

		uf := fmt.Sprintf("%s/%s", ri.Dir(), obj.Output)

		b, err := os.ReadFile(uf)
		if err != nil {
			reusable.ReportError(w, errForCode("json"), err)
			return
		}

		ret := make(map[string]interface{})
		err = json.Unmarshal(b, &ret)
		if err != nil {
			reusable.ReportError(w, errForCode("json"), err)
			return
		}
		defer os.Remove(uf)
		reusable.ReportResult(w, ret)
	}

}

func writeFiles(dir string, fi *reusable.FileIterator) error {

	for {
		f, err := fi.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		file, err := os.Create(fmt.Sprintf("%s/%s", dir, f.Name))
		if err != nil {
			return err
		}
		defer file.Close()

		err = os.Chmod(file.Name(), 0755)
		if err != nil {
			return err
		}

		script, err := f.AsReader()
		if err != nil {
			return err
		}
		defer script.Close()

		_, err = io.Copy(file, script)
		if err != nil {
			return err
		}

	}

	return nil

}

func main() {
	reusable.StartServer(ubuntuHandler, nil)
}

func errForCode(errCode string) string {
	return fmt.Sprintf("com.ubuntu.%s.error", errCode)
}

// import (
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"os"
// 	"os/exec"
// 	"path"
// 	"strings"

// 	"github.com/direktiv/direktiv-apps/pkg/direktivapps"
// )

// const code = "com.ubuntu.error"

// type shell struct {
// 	Script string   `json:"script"`
// 	Args   []string `json:"args"`
// }

// func execScript(path string, args []string, envs []string) ([]byte, error) {

// 	cmd := exec.Command(path, args...)
// 	cmd.Env = envs

// 	d, err := cmd.Output()
// 	if err != nil {
// 		if e, ok := err.(*exec.ExitError); ok {
// 			d = e.Stderr
// 		}
// 		return d, err
// 	}

// 	return d, nil

// }

// func request(w http.ResponseWriter, r *http.Request) {

// 	var s shell
// 	aid, err := direktivapps.Unmarshal(&s, r)
// 	if err != nil {
// 		direktivapps.RespondWithError(w, code, err.Error())
// 		return
// 	}

// 	if s.Script == "" {
// 		direktivapps.RespondWithError(w, code, "no script provided")
// 		return
// 	}

// 	direktivapps.LogDouble(aid, "getting script")

// 	dtmp := r.Header.Get("Direktiv-TempDir")
// 	if len(dtmp) == 0 {
// 		direktivapps.RespondWithError(w, code, fmt.Sprintf("direktiv tmp directory not provided"))
// 		return
// 	}

// 	f := path.Join(dtmp, s.Script)

// 	err = os.Chmod(f, 0755)
// 	if err != nil {
// 		direktivapps.RespondWithError(w, code, err.Error())
// 		return
// 	}

// 	direktivapps.LogDouble(aid, fmt.Sprintf("found script %v", f))

// 	envs := []string{fmt.Sprintf("Direktiv_TempDir=%s",
// 		r.Header.Get("Direktiv-TempDir"))}

// 	ret, err := execScript(f, s.Args, envs)
// 	if err != nil {
// 		direktivapps.RespondWithError(w, code, fmt.Sprintf("%v: %s", err, string(ret)))
// 		return
// 	}

// 	ret = []byte(strings.TrimSpace(string(ret)))
// 	direktivapps.LogDouble(aid, fmt.Sprintf("script return: %v", string(ret)))

// 	if !json.Valid(ret) {
// 		o := make(map[string]string)
// 		o["output"] = string(ret)

// 		b, err := json.Marshal(&o)
// 		if err != nil {
// 			direktivapps.RespondWithError(w, code, string(ret))
// 			return
// 		}

// 		direktivapps.Respond(w, b)
// 	} else {
// 		direktivapps.Respond(w, ret)
// 	}

// }

// func main() {
// 	direktivapps.StartServer(request)
// }
