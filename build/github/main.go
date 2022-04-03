package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"

	"github.com/direktiv/direktiv-apps/pkg/reusable"
	shellwords "github.com/mattn/go-shellwords"
)

const (
	cmdError = "io.direktiv.github.cmd.error"
)

type requestInput struct {
	Token    string   `json:"token"`
	Commands []string `json:"cmds"`
	Continue bool     `json:"continue"`
}

func ghHandler(w http.ResponseWriter, r *http.Request, ri *reusable.RequestInfo) {

	obj := new(requestInput)
	err := reusable.Unmarshal(obj, true, r)
	if err != nil {
		reusable.ReportError(w, reusable.UnmarshallError, err)
		return
	}

	result := make(map[int]map[string]interface{})

	for a := range obj.Commands {

		ri.Logger().Infof("running command %d", a)

		c := obj.Commands[a]

		args, err := shellwords.Parse(c)
		if err != nil {
			reusable.ReportError(w, cmdError, err)
			return
		}

		ir := make(map[string]interface{})

		cmd := exec.Command("/usr/bin/gh", args...)
		cmd.Env = append(os.Environ(), fmt.Sprintf("GITHUB_TOKEN=%s", obj.Token))

		var o bytes.Buffer

		mw := io.MultiWriter(os.Stdout, ri.LogWriter(), &o)
		cmd.Stdout = mw
		cmd.Stderr = mw
		cmd.Dir = ri.Dir()

		success := true
		err = cmd.Run()
		if err != nil {
			ri.Logger().Errorf("error running command %d: %v", a, err)
			if !obj.Continue {
				reusable.ReportError(w, cmdError, err)
				return
			}
			success = false
			result[a] = ir
		}

		rj := reusable.ToJSON(o.String())

		ir["success"] = success
		ir["result"] = rj

		result[a] = ir

		reusable.ReportResult(w, result)
	}

}

func main() {
	reusable.StartServer(ghHandler, nil)
}
