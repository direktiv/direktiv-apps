package main

import (
	"fmt"
	"net/http"
	"os/exec"

	"github.com/direktiv/direktiv-apps/pkg/direktivapps"
)

// AzureDetails are the details provided to authenticate and execute a command using the azcli
type AzureDetails struct {
	Name     string   `json:"name"`
	Password string   `json:"password"`
	Tenant   string   `json:"tenant"`
	Command  []string `json:"command"`
}

const code = "com.azcli.error"

func main() {
	direktivapps.StartServer(AzureGo)
}

func AzureGo(w http.ResponseWriter, r *http.Request) {
	var err error

	obj := new(AzureDetails)
	_, err = direktivapps.Unmarshal(obj, r)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	// Authenticate with the azcli using a service principal
	cmd := exec.Command("/usr/bin/az", "login", "--service-principal", "-u", obj.Name, "-p", obj.Password, "--tenant", obj.Tenant)
	resp, err := cmd.CombinedOutput()
	if err != nil {
		direktivapps.RespondWithError(w, code, fmt.Sprintf("output: %s", resp))
		return
	}
	// Allow azure to update itself without a tty
	cmd = exec.Command("/usr/bin/az", "config", "set", "extension.use_dynamic_install=no")
	resp, err = cmd.CombinedOutput()
	if err != nil {
		direktivapps.RespondWithError(w, code, fmt.Sprintf("output: %s", resp))
		return
	}

	// Execute command provided via the input of container
	cmd = exec.Command("/usr/bin/az", obj.Command...)
	resp, err = cmd.CombinedOutput()
	if err != nil {
		direktivapps.RespondWithError(w, code, fmt.Sprintf("output: %s", resp))
		return
	}

	direktivapps.Respond(w, resp)
}
