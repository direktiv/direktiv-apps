package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"

	"github.com/vorteil/direktiv-apps/pkg/direktivapps"
)

type TerraformInput struct {
	Tf             string                 `json:"url"`       // url that wgets the terraform file
	Tfstate        map[string]interface{} `json:"tfstate"`   // the tfstate being parsed into
	Action         string                 `json:"action"`    // "plan", "validate", "apply", "destroy"
	ClientID       string                 `json:"client_id"` //
	ClientSecret   string                 `json:"client_secret"`
	SubscriptionID string                 `json:"subscription_id"`
	TenantID       string                 `json:"tenant_id"`
}

var code = "com.terraform.error"
var terraformBin = "/terraform"

type CMDWriter struct {
	Aid string
}

func (c *CMDWriter) Write(p []byte) (n int, err error) {
	direktivapps.Log(c.Aid, string(p))
	return len(p), nil
}

func main() {
	direktivapps.StartServer(TerraformHandler)
}

func TerraformHandler(w http.ResponseWriter, r *http.Request) {

	obj := new(TerraformInput)
	aid, err := direktivapps.Unmarshal(obj, r)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	// set tf state
	if obj.Tfstate != nil {
		tfstate, err := json.Marshal(obj.Tfstate)
		if err != nil {
			direktivapps.RespondWithError(w, code, err.Error())
			return
		}
		direktivapps.Log(aid, fmt.Sprintf("\n\n\n%+v\n\n\n", tfstate))

		err = ioutil.WriteFile("./terraform.tfstate", tfstate, os.ModePerm)
		if err != nil {
			direktivapps.RespondWithError(w, code, err.Error())
			return
		}
	}

	os.Setenv("TF_LOG", "DEBUG")

	// setup custom writer to send command output to direktiv
	cmdW := &CMDWriter{
		Aid: aid,
	}

	direktivapps.Log(aid, "Fetching TF file externally")

	resp, err := http.Get(obj.Tf)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	direktivapps.Log(aid, "Generating authentication for terraform...")

	// set azure auth
	err = os.Setenv("ARM_CLIENT_ID", obj.ClientID)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	err = os.Setenv("ARM_CLIENT_SECRET", obj.ClientSecret)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	err = os.Setenv("ARM_SUBSCRIPTION_ID", obj.SubscriptionID)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	err = os.Setenv("ARM_TENANT_ID", obj.TenantID)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	direktivapps.Log(aid, "Creating the tf file for terraform to read from...")
	// write terraform file from string
	err = ioutil.WriteFile("/main.tf", []byte(data), os.ModePerm)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	direktivapps.Log(aid, "Initializing terraform...")
	// initialize the terraform get required plugins etc
	init := exec.Command(terraformBin, "init")
	init.Stderr = cmdW
	init.Stdout = cmdW
	err = init.Run()
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	direktivapps.Log(aid, "Executing action for terraform...")
	switch obj.Action {
	case "apply":
		fallthrough
	case "destroy":
		cmd := exec.Command(terraformBin, obj.Action, "-auto-approve")
		cmd.Stdout = cmdW
		cmd.Stderr = cmdW
		err = cmd.Run()
		if err != nil {
			direktivapps.RespondWithError(w, code, err.Error())
			return
		}
	case "validate":
		fallthrough
	case "plan":
		cmd := exec.Command(terraformBin, obj.Action)
		cmd.Stdout = cmdW
		cmd.Stderr = cmdW
		cmd.Run()
		if err != nil {
			direktivapps.RespondWithError(w, code, err.Error())
			return
		}
	}

	state, err := ioutil.ReadFile("terraform.tfstate")
	if err != nil {
		if os.IsNotExist(err) {
		} else {
			direktivapps.RespondWithError(w, code, err.Error())
			return
		}
	}

	if state == nil {
		direktivapps.Respond(w, []byte("{}"))
	}

	direktivapps.Respond(w, state)
}
