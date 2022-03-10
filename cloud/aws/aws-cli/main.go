package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/direktiv/direktiv-apps/pkg/reusable"
	"github.com/go-playground/validator"
)

const (
	defaultRegion = "us-east-1"
	awsErr        = "com.aws.execute"
)

type requestInput struct {
	AccessKey string `json:"access-key" validate:"required" `
	SecretKey string `json:"secret-key" validate:"required"`
	Region    string `json:"region"`

	Commands []string `json:"commands"`
	Continue bool     `json:"continue"`
	Print    bool     `json:"print"`

	Script reusable.File `json:"script"`
}

type result struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Error   string      `json:"error"`
}

func dummyHandler(w http.ResponseWriter, r *http.Request, ri *reusable.RequestInfo) {

	obj := new(requestInput)
	err := reusable.Unmarshal(obj, true, r)
	if err != nil {
		reusable.ReportError(w, reusable.UnmarshallError, err)
		return
	}

	validate := validator.New()
	err = validate.Struct(obj)

	if err != nil {
		reusable.ReportError(w, reusable.UnmarshallError, err)
		return
	}

	if obj.Region == "" {
		obj.Region = defaultRegion
	}

	awsEnvs := []string{
		fmt.Sprintf("AWS_ACCESS_KEY_ID=%s", obj.AccessKey),
		fmt.Sprintf("AWS_SECRET_ACCESS_KEY=%s", obj.SecretKey),
		fmt.Sprintf("AWS_DEFAULT_REGION=%s", obj.Region),
		"AWS_DEFAULT_OUTPUT=json",
	}

	envs := os.Environ()
	envs = append(envs, awsEnvs...)

	results := []result{}

	for a := range obj.Commands {

		c := obj.Commands[a]

		r := result{
			Success: true,
		}

		// split at whitespace
		args := strings.Fields(c)

		cmd := exec.Command("aws", args...)

		var b bytes.Buffer
		mw := io.MultiWriter(ri.LogWriter(), os.Stdout, &b)

		if !obj.Print {
			cmd.Stdout = mw
		}
		cmd.Stderr = mw

		cmd.Env = envs
		cmd.Dir = ri.Dir()

		ri.Logger().Infof("executing %s...", cmd.String()[:60])

		err := cmd.Run()
		if err != nil && !obj.Continue {
			reusable.ReportError(w, awsErr,
				fmt.Errorf("error running command: %v: %v", b.String(), err))
			return
		} else if err != nil {
			r.Success = false
			r.Error = err.Error()
		} else {
			r.Data = reusable.ToJSON(b.String())
		}

		results = append(results, r)

	}

	if obj.Script.Data != "" {

		ri.Logger().Infof("running script")
		var script strings.Builder

		script.WriteString("#!/bin/bash\n")

		s, err := obj.Script.AsString(ri)
		if err != nil {
			reusable.ReportError(w, awsErr, err)
			return
		}
		script.WriteString(s)

		scriptPath := path.Join(ri.Dir(), "script.sh")
		err = os.WriteFile(scriptPath, []byte(script.String()), 0755)
		if err != nil {
			reusable.ReportError(w, awsErr, err)
			return
		}

		var b bytes.Buffer
		mw := io.MultiWriter(ri.LogWriter(), os.Stdout, &b)

		cmd := exec.Command(scriptPath)
		cmd.Env = envs
		cmd.Stdout = mw
		cmd.Stderr = mw
		cmd.Dir = ri.Dir()

		err = cmd.Run()
		if err != nil {
			reusable.ReportError(w, awsErr, err)
			return
		}

		// doesn't have to be there, we ignore errors
		o, _ := reusable.GetFileAsJson(path.Join(ri.Dir(), "output.json"))

		r := result{
			Success: true,
			Data:    o,
		}

		results = append(results, r)

	}

	reusable.ReportResult(w, results)

}

func main() {
	reusable.StartServer(dummyHandler, nil)
}
