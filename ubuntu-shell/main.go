package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"

	"github.com/direktiv/direktiv-apps/pkg/reusable"
)

type requestInput struct {
	Script reusable.File `json:"script"`
	Args   []string      `json:"args"`
	Envs   []string      `json:"envs"`
	Output string        `json:"output"`
}

func ubuntuHandler(w http.ResponseWriter, r *http.Request, ri *reusable.RequestInfo) {

	obj := new(requestInput)
	err := reusable.Unmarshal(obj, false, r)
	if err != nil {
		reusable.ReportError(w, reusable.UnmarshallError, err)
		return
	}

	file, err := obj.Script.AsFile(ri, 0755)
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

	cmd.Env = append(os.Environ(), obj.Envs...)
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

func main() {
	reusable.StartServer(ubuntuHandler, nil)
}

func errForCode(errCode string) string {
	return fmt.Sprintf("com.ubuntu.%s.error", errCode)
}
