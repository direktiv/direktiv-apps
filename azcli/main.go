package main

import (
	"os/exec"

	"github.com/vorteil/direktiv-apps/pkg/direktivapps"
)

// AzureDetails are the details provided to authenticate and execute a command using the azcli
type AzureDetails struct {
	Name     string   `json:"name"`
	Password string   `json:"password"`
	Tenant   string   `json:"tenant"`
	Command  []string `json:"command"`
}

func main() {
	g := direktivapps.ActionError{
		ErrorCode:    "com.azcli.error",
		ErrorMessage: "",
	}

	var err error

	obj := new(AzureDetails)
	direktivapps.ReadIn(obj, g)

	// Authenticate with the azcli using a service principal
	cmd := exec.Command("/usr/bin/az", "--service-principal", "login", "-u", obj.Name, "-p", obj.Password, "--tenant", obj.Tenant)
	err = cmd.Run()
	if err != nil {
		g.ErrorMessage = err.Error()
		direktivapps.WriteError(g)
	}

	// Execute command provided via the input of container
	cmd = exec.Command("/usr/bin/az", obj.Command...)

	resp, err := cmd.CombinedOutput()
	if err != nil {
		g.ErrorMessage = err.Error()
		direktivapps.WriteError(g)
	}

	direktivapps.WriteOut(resp, g)
}
