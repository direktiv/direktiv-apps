package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"

	"github.com/direktiv/direktiv-apps/pkg/direktivapps"
)

var code = "com.go.error"

// CMDWriter allows us to log to the action id provided.
type CMDWriter struct {
	Aid string `json:"aid"`
}

// Write writes to the new api for logging direktiv apps
func (c *CMDWriter) Write(p []byte) (n int, err error) {
	direktivapps.LogDouble(c.Aid, string(p))
	return len(p), nil
}

type goinput struct {
	Arguments       []string `json:"args"`
	ExecutionFolder string   `json:"execution-folder"`
	Variable        string   `json:"variable"`      // the name the variable will be saved as and the file from a go build
	VariableType    string   `json:"variable-type"` // instance, workflow or namespace
}

func GoHandler(w http.ResponseWriter, r *http.Request) {
	var goi goinput
	aid, err := direktivapps.Unmarshal(&goi, r)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}
	cmdW := &CMDWriter{
		Aid: aid,
	}

	executionPath := path.Join(r.Header.Get("Direktiv-TempDir"), goi.ExecutionFolder)

	direktivapps.LogDouble(aid, fmt.Sprintf("changing directory to '%s'", executionPath))

	err = os.Chdir(executionPath)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	direktivapps.LogDouble(aid, "Executing go command...")

	cmd := exec.Command("go", goi.Arguments...)
	cmd.Stderr = cmdW
	cmd.Stdout = cmdW

	err = cmd.Run()
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	switch goi.Arguments[0] {
	case "build":
		f, err := os.Open(path.Join(executionPath, goi.Variable))
		if err != nil {
			direktivapps.RespondWithError(w, code, err.Error())
			return
		}
		defer f.Close()
		err = uploadVariable(aid, f, goi.Variable, goi.VariableType)
		if err != nil {
			direktivapps.RespondWithError(w, code, err.Error())
			return
		}
	}

	direktivapps.Respond(w, []byte(""))
}

func uploadVariable(aid string, f *os.File, name string, vartype string) error {
	fi, err := os.Stat(f.Name())
	if err != nil {
		return err
	}
	// get the size
	size := fi.Size()
	direktivapps.LogDouble(aid, fmt.Sprintf("uploading '%s', size: '%v'", f.Name(), size))

	req, err := http.NewRequest("POST", fmt.Sprintf("http://localhost:8889/var?aid=%s&scope=%s&key=%s", aid, vartype, name), f)
	if err != nil {
		return err
	}
	req.ContentLength = size

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return errors.New(string(data))
	}

	return nil
}

func main() {
	direktivapps.StartServer(GoHandler)
}
