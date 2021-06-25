package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/vorteil/direktiv-apps/pkg/direktivapps"
)

// TerraformInput takes different arguments for each authentication service.
type TerraformInput struct {
	Action                 string                 `json:"action"`         // plan, validate, apply, destroy
	TFVars                 map[string]interface{} `json:"variables"`      // the tf variables from the input.
	AdditionalArgs         []string               `json:"args-on-init"`   // additional arguments on init command
	AdditionalArgsOnAction []string               `json:"args-on-action"` // additional arguments on action command
}

// OutputResponse only gets used if a tfstate wasn't provided as the 'state-name' variable.
type OutputResponse struct {
	Output  map[string]interface{} `json:"output"`
	TFState map[string]interface{} `json:"tfstate"`
}

var code = "com.terraform.%s.error"
var terraformBin = "/terraform"

var runningTF map[string]string

// CMDWriter allows us to log to the action id provided.
type CMDWriter struct {
	Aid string `json:"aid"`
}

// Write writes to the new api for logging direktiv apps
func (c *CMDWriter) Write(p []byte) (n int, err error) {
	direktivapps.Log(c.Aid, string(p))
	return len(p), nil
}

func main() {
	runningTF = make(map[string]string)
	direktivapps.StartServer(TerraformHandler)
}

func httpBackend() {
	// mux := http.NewServeMux()
	r := mux.NewRouter()
	r.HandleFunc("/{id}", TFStateHandler)
	r.HandleFunc("/", TFStateHandler)

	srv := &http.Server{
		Addr:    ":8001",
		Handler: r,
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM)

	go func() {
		<-sigs
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		srv.Shutdown(ctx)
	}()

	srv.ListenAndServe()
}

func TFStateHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	url := fmt.Sprintf("http://localhost:8889/var?aid=%s&scope=workflow&key=%s", runningTF[vars["id"]], vars["id"])

	switch r.Method {
	case http.MethodGet:
		direktivapps.Log(runningTF[vars["id"]], "Fetching tfstate variable...")
		resp, err := http.Get(url)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		defer resp.Body.Close()
		data, _ := ioutil.ReadAll(resp.Body)
		if len(data) == 0 {
			// no state currently exists
			data = []byte(`{
				"version": 4
			}`)
		}
		w.Write(data)
	case http.MethodPost:
		direktivapps.Log(runningTF[vars["id"]], "Saving new tfstate variable...")
		req, err := http.NewRequest("POST", url, r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()
		return
	case http.MethodDelete:
		req, err := http.NewRequest("POST", url, bytes.NewReader([]byte("")))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()
		return
	}
}

func checkBackendIsAlive() bool {
	timeout := 1 * time.Second
	_, err := net.DialTimeout("tcp", "localhost:8001", timeout)
	if err != nil {
		return false
	}
	return true
}

func TerraformHandler(w http.ResponseWriter, r *http.Request) {

	obj := new(TerraformInput)

	aid, err := direktivapps.Unmarshal(obj, r)
	if err != nil {
		direktivapps.RespondWithError(w, fmt.Sprintf(code, "unmarshal-input"), err.Error())
		return
	}

	stateName, ok := obj.TFVars["state-name"].(string)
	if !ok {
		direktivapps.Log(aid, "state-name variable was not provided not using http backend")
	} else {
		direktivapps.Log(aid, "adding to the global map to control action ids")
		if runningTF[stateName] != "" {
			direktivapps.RespondWithError(w, fmt.Sprintf(code, "state-name currently in use with a different action"), err.Error())
			return
		}

		runningTF[stateName] = aid
		defer delete(runningTF, stateName)
	}

	direktivapps.Log(aid, "Finding path to call terraform from...")
	terraformPath := r.Header.Get("Direktiv-TempDir")
	direktivapps.Log(aid, fmt.Sprintf("Found '%s'", terraformPath))

	direktivapps.Log(aid, "Checking if tfstate service http backend is alive...")
	alive := checkBackendIsAlive()
	// if backend not alive spawn backend
	if !alive {
		go httpBackend()
	}

	direktivapps.Log(aid, "Wait till backend service is functional")
	for !alive {
		alive = checkBackendIsAlive()
	}

	direktivapps.Log(aid, "Initializing terraform....")

	cmdW := &CMDWriter{
		Aid: aid,
	}
	dirArg := fmt.Sprintf("-chdir=%s", terraformPath)

	direktivapps.Log(aid, fmt.Sprintf("Reading in TFVars.json..."))

	data, err := json.Marshal(obj.TFVars)
	if err != nil {
		direktivapps.RespondWithError(w, fmt.Sprintf(code, "tfvars-marshal"), err.Error())
		return
	}

	// if data exists write the variables parser
	if string(data) != "null" {
		err = ioutil.WriteFile(path.Join(terraformPath, "terraform.tfvars.json"), data, os.ModePerm)
		if err != nil {
			direktivapps.RespondWithError(w, fmt.Sprintf(code, "tfvars"), err.Error())
			return
		}
	}

	var cmdInit []string
	cmdInit = append(cmdInit, dirArg, "init")
	for _, arg := range obj.AdditionalArgs {
		cmdInit = append(cmdInit, arg)
	}

	init := exec.Command(terraformBin, cmdInit...)
	init.Stderr = cmdW
	init.Stdout = cmdW
	err = init.Run()
	if err != nil {
		direktivapps.RespondWithError(w, fmt.Sprintf(code, "init"), err.Error())
		return
	}

	direktivapps.Log(aid, fmt.Sprintf("Executing '%s' for terraform", obj.Action))
	switch obj.Action {
	case "apply":
		fallthrough
	case "destroy":
		var cmdApply []string
		cmdApply = append(cmdApply, dirArg, obj.Action, "-auto-approve")
		for _, arg := range obj.AdditionalArgsOnAction {
			cmdApply = append(cmdApply, arg)
		}
		cmd := exec.Command(terraformBin, cmdApply...)
		cmd.Stdout = cmdW
		cmd.Stderr = cmdW
		err = cmd.Run()
		if err != nil {
			direktivapps.RespondWithError(w, fmt.Sprintf(code, "apply-destroy"), err.Error())
			return
		}
	case "validate":
		fallthrough
	case "plan":
		var cmdPlan []string
		cmdPlan = append(cmdPlan, dirArg, obj.Action)
		for _, arg := range obj.AdditionalArgsOnAction {
			cmdPlan = append(cmdPlan, arg)
		}
		cmd := exec.Command(terraformBin, cmdPlan...)
		cmd.Stdout = cmdW
		cmd.Stderr = cmdW
		err = cmd.Run()
		if err != nil {
			direktivapps.RespondWithError(w, fmt.Sprintf(code, "validate-plan"), err.Error())
			return
		}
	}

	direktivapps.Log(aid, "Sending output back to direktiv...")
	cmd := exec.Command(terraformBin, dirArg, "output", "-json")
	data, err = cmd.Output()
	if err != nil {
		direktivapps.RespondWithError(w, fmt.Sprintf(code, "output"), err.Error())
		return
	}

	var tfstateData []byte
	if stateName != "" {
		resp, err := http.Get(fmt.Sprintf("http://localhost:8001/%s", stateName))
		if err != nil {
			direktivapps.RespondWithError(w, fmt.Sprintf(code, "read-state"), err.Error())
			return
		}
		defer resp.Body.Close()
		tfstateData, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			direktivapps.RespondWithError(w, fmt.Sprintf(code, "state-data-read"), err.Error())
			return
		}

		if len(tfstateData) == 0 {
			tfstateData, err = ioutil.ReadFile(path.Join(terraformPath, "terraform.tfstate"))
			if err != nil {
				direktivapps.RespondWithError(w, fmt.Sprintf(code, "state-data-read"), err.Error())
				return
			}
			direktivapps.Log(aid, fmt.Sprintf("read tfstate data from file: %s", tfstateData))
		}
	} else {
		tfstateData, err = ioutil.ReadFile(path.Join(terraformPath, "terraform.tfstate"))
		if err != nil {
			direktivapps.RespondWithError(w, fmt.Sprintf(code, "state-data-read"), err.Error())
			return
		}
		direktivapps.Log(aid, fmt.Sprintf("read tfstate data from file: %s", tfstateData))

	}

	var output map[string]interface{}
	err = json.Unmarshal(data, &output)
	if err != nil {
		direktivapps.RespondWithError(w, fmt.Sprintf(code, "unmarshal-json-output"), err.Error())
		return
	}

	var tfstate map[string]interface{}
	err = json.Unmarshal(tfstateData, &tfstate)
	if err != nil {
		direktivapps.RespondWithError(w, fmt.Sprintf(code, "unmarshal-tfstate"), err.Error())
		return
	}

	outputobj := OutputResponse{
		Output:  output,
		TFState: tfstate,
	}

	writeBack, err := json.Marshal(outputobj)
	if err != nil {
		direktivapps.RespondWithError(w, fmt.Sprintf(code, "json-marshal-output"), err.Error())
		return
	}

	direktivapps.Respond(w, writeBack)
}
