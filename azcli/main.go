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

	var generateLoginCommand []string
	var generateExecCommand []string

	generateLoginCommand = append(generateLoginCommand, "-m", "azure.cli", "login", "-u", obj.Name, "-p", obj.Password, "--tenant", obj.Tenant)
	generateExecCommand = append(generateExecCommand, "-m", "azure.cli")
	generateExecCommand = append(generateExecCommand, obj.Command...)

	// Authenticate with the azcli using a service principal
	cmd := exec.Command("/azure-cli/bin/python", generateLoginCommand...)
	err = cmd.Run()
	if err != nil {
		g.ErrorMessage = err.Error()
		direktivapps.WriteError(g)
	}

	// Execute command provided via the input of container
	cmd = exec.Command("/azure-cli/bin/python", generateExecCommand...)

	resp, err := cmd.CombinedOutput()
	if err != nil {
		g.ErrorMessage = err.Error()
		direktivapps.WriteError(g)
	}

	direktivapps.WriteOut(resp, g)
}
