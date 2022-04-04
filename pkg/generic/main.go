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

type cmd struct {
	Cmd      string   `json:"cmd"`
	Env      []string `json:"env"`
	Continue bool     `json:"continue"`
	Silent   bool     `json:"silent"`
	Output   string   `json:"output"`
}

type requestInput struct {
	Commands []cmd `json:"cmds"`
}

const (
	cmdError = "io.direktiv.generic.error"
)

func genericHandler(w http.ResponseWriter, r *http.Request, ri *reusable.RequestInfo) {

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

		args, err := shellwords.Parse(c.Cmd)
		if err != nil {
			reusable.ReportError(w, cmdError, err)
			return
		}

		ir := make(map[string]interface{})

		cmd := exec.Command(args[0], args[1:]...)

		var o bytes.Buffer
		writer := []io.Writer{os.Stdout, &o}
		if !c.Silent {
			writer = append(writer, ri.LogWriter())
		}

		mw := io.MultiWriter(writer...)
		cmd.Stdout = mw
		cmd.Stderr = mw
		cmd.Dir = ri.Dir()
		cmd.Env = append(os.Environ(), c.Env...)

		success := true
		err = cmd.Run()
		if err != nil {
			ri.Logger().Errorf("error running command %d: %v", a, err)
			if !c.Continue {
				reusable.ReportError(w, cmdError, err)
				return
			}
			success = false
			result[a] = ir
		}

		rj := reusable.ToJSON(o.String())
		if c.Output != "" {
			b, err := os.ReadFile(c.Output)
			if err != nil {
				reusable.ReportError(w, cmdError, fmt.Errorf("error cmd %d: %s", a, err))
				return
			}
			rj = reusable.ToJSON(string(b))
		}

		ir["success"] = success
		ir["result"] = rj
		result[a] = ir

	}

	reusable.ReportResult(w, result)
}

func main() {
	reusable.StartServer(genericHandler, nil)
}
