package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/acarl005/stripansi"
	"github.com/direktiv/direktiv-apps/pkg/reusable"
	pss "github.com/direktiv/go-powershell"
	"github.com/direktiv/go-powershell/backend"
	"github.com/google/uuid"
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
	Output  interface{} `json:"output,omitempty"`
}

type requestInput struct {
	Scripts  []Script `json:"scripts"`
	Continue bool     `json:"continue"`

	ri *reusable.RequestInfo `json:"-"`

	errCh    chan error                 `json:"-"`
	resultCh chan map[string]*CmdOutput `json:"-"`
	shell    pss.Shell                  `json:"-"`
}

const shells = 5

var jobs chan *requestInput

func runScripts(obj *requestInput, shell pss.Shell, be *backend.Local) (map[string]*CmdOutput, error) {

	be.Writer = obj.ri.LogWriter()

	shell.Execute(fmt.Sprintf("Set-Location -Path %s", obj.ri.Dir()))

	obj.ri.Logger().Infof("getting shell")

	object := make(map[string]*CmdOutput)

	for i := range obj.Scripts {

		co := &CmdOutput{
			Success: true,
		}
		stdout, stderr, err := executeScript(obj.ri, shell, obj.Scripts[i])

		if err != nil {
			obj.ri.Logger().Errorf("error executing script: %v", err)

			// stop if one script fails
			if !obj.Continue {
				return nil, err
			}
			co.Success = false
		}

		// add to result
		co.Stderr = stderr
		co.Stdout = toJSON(obj.ri, stdout)

		// adding output file if configured
		if co.Success && len(obj.Scripts[i].Output) > 0 {

			obj.ri.Logger().Infof("getting output file")

			uf := fmt.Sprintf("%s/%s", obj.ri.Dir(), obj.Scripts[i].Output)
			defer os.Remove(uf)

			o, err := os.ReadFile(uf)
			if err != nil {
				if !obj.Continue {
					return nil, err
				}
				co.Success = false
			}

			co.Output = toJSON(obj.ri, string(o))
		}

		object[obj.Scripts[i].Script.Name] = co
	}

	return object, nil
}

func runPwsh(input chan *requestInput) {

	log.Println("creating shell")

	var err error
	be, shell, err := getShell()
	if err != nil {
		log.Printf("can not create shell: %v", err)
	}

	for {

		log.Println("worker waiting for request")

		obj := <-input

		// we add this so we can timeout in the script later
		obj.shell = shell

		log.Println("worker got request")

		// get a new shell if required
		if shell == nil {
			log.Println("shell nil, creating new one")
			be, shell, err = getShell()
			if err != nil {
				obj.ri.Logger().Errorf("can not create shell: %v", err)
				continue
			}
		}

		result, err := runScripts(obj, shell, be)
		if err != nil {
			obj.ri.Logger().Errorf("can not execute shell: %v", err)
			obj.errCh <- err
		} else {
			obj.resultCh <- result
		}

		// trying to get a new shell
		shell.Exit()
		be, shell, err = getShell()
		if err != nil {
			log.Printf("can not create shell: %v", err)
		}

	}

}

func powercliHandler(w http.ResponseWriter, r *http.Request, ri *reusable.RequestInfo) {

	obj := new(requestInput)
	err := reusable.Unmarshal(obj, true, r)
	if err != nil {
		reusable.ReportError(w, reusable.UnmarshallError, err)
		return
	}

	obj.ri = ri
	obj.errCh = make(chan error, 1)
	obj.resultCh = make(chan map[string]*CmdOutput, 1)

	jobs <- obj

	select {
	case err = <-obj.errCh:
		reusable.ReportError(w, "com.powershell.run.error", err)
	case result := <-obj.resultCh:
		reusable.ReportResult(w, result)
	case <-time.After(60 * time.Minute):
		reusable.ReportError(w, "com.powershell.timeout.error", fmt.Errorf("timed out"))
	}

}

func executeScript(ri *reusable.RequestInfo, shell pss.Shell, script Script) (string, string, error) {

	// create file
	f, err := script.Script.AsFile(ri, 0755)
	if err != nil {
		ri.Logger().Errorf("error creating script: %v", err)
		return "", "", err
	}

	name := f.Name()
	if name == "" {
		name = uuid.New().String()
	}
	ps1 := fmt.Sprintf("%s.ps1", name)
	os.Rename(name, ps1)
	defer os.Remove(ps1)

	full := fmt.Sprintf("%s %s", ps1, strings.Join(script.Args, " "))
	ri.Logger().Infof("running %v", full)

	stdout, stderr, err := shell.Execute(full)
	if err != nil {
		ri.Logger().Errorf("error executing script: %v", err)
		return stdout, stderr, err
	}

	return stdout, stderr, err
}

func getShell() (*backend.Local, pss.Shell, error) {

	back := &backend.Local{}

	s, err := pss.New(back)

	return back, s, err

}

func main() {

	jobs = make(chan *requestInput, shells)
	for a := 0; a < shells; a++ {
		go runPwsh(jobs)
	}

	reusable.StartServer(powercliHandler, nil)
}

func toJSON(ri *reusable.RequestInfo, str string) interface{} {

	str = strings.TrimSpace(str)
	str = stripansi.Strip(str)

	var js json.RawMessage
	err := json.Unmarshal([]byte(str), &js)
	if err != nil {
		return str
	}

	return json.RawMessage(str)

}
