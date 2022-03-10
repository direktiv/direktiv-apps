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
	PreRun  reusable.File `json:"pre-run"`
	PostRun reusable.File `json:"post-run"`
	Script  reusable.File `json:"script"`
	Args    []string      `json:"args"`
	Envs    []string      `json:"envs"`
	Output  string        `json:"output"`
}

func pythonHandler(w http.ResponseWriter, r *http.Request, ri *reusable.RequestInfo) {

	obj := new(requestInput)
	err := reusable.Unmarshal(obj, true, r)
	if err != nil {
		reusable.ReportError(w, reusable.UnmarshallError, err)
		return
	}

	ri.Logger().Infof("starting pre-scripts")

	// run pre-script
	err = runScript(&obj.PreRun, obj.Envs, ri)
	if err != nil {
		reusable.ReportError(w, errForCode("pre-run"), err)
		return
	}

	file, err := obj.Script.AsFile(0755)
	if err != nil {
		reusable.ReportError(w, errForCode("execute"), err)
		return
	}
	defer os.Remove(file.Name())

	args := []string{
		file.Name(),
	}
	args = append(args, obj.Args...)

	cmd := exec.Command("python3", args...)

	ri.Logger().Infof("executing %v", cmd)

	// printing everything to the screen
	cmd.Stderr = ri.LogWriter()
	cmd.Stdout = ri.LogWriter()

	ri.Logger().Infof("running %s", cmd)

	ri.Logger().Infof("using %d envs: %v", len(obj.Envs), obj.Envs)
	obj.Envs = append(obj.Envs, fmt.Sprintf("Direktiv_TempDir=%s", ri.Dir()))
	cmd.Env = obj.Envs
	cmd.Dir = ri.Dir()

	if err := cmd.Start(); err != nil {
		reusable.ReportError(w, errForCode("execute"), err)
		return
	}
	if err := cmd.Wait(); err != nil {
		reusable.ReportError(w, errForCode("execute"), err)
		return
	}

	// run post-script
	err = runScript(&obj.PostRun, obj.Envs, ri)
	if err != nil {
		reusable.ReportError(w, errForCode("post-run"), err)
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

func errForCode(errCode string) string {
	return fmt.Sprintf("com.python.%s.error", errCode)
}

func runScript(f *reusable.File, envs []string, ri *reusable.RequestInfo) error {

	if len(f.Data) == 0 {
		return nil
	}

	file, err := f.AsFile(0755)
	if err != nil {
		return err
	}
	file.Close()
	defer os.Remove(file.Name())

	cmd := exec.Command(file.Name())

	ri.Logger().Infof("executing %v", cmd)

	cmd.Stderr = ri.LogWriter()
	cmd.Stdout = ri.LogWriter()
	cmd.Env = envs
	cmd.Dir = ri.Dir()

	if err := cmd.Start(); err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		return err
	}

	return nil

}

func main() {
	reusable.StartServer(pythonHandler, nil)
}
