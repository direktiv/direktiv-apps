package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"

	"github.com/vorteil/direktiv-apps/pkg/direktivapps"
)

// InputContainerDetails ...
type InputContainerDetails struct {
	ServiceAccountKey string   `json:"serviceAccountKey"`
	Command           []string `json:"command"`
}

func main() {

	g := direktivapps.ActionError{
		ErrorCode:    "com.googlecli.error",
		ErrorMessage: "",
	}

	obj := new(InputContainerDetails)
	direktivapps.ReadIn(obj, g)

	err := ioutil.WriteFile("/key.json", []byte(obj.ServiceAccountKey), 0644)
	if err != nil {
		g.ErrorMessage = fmt.Sprintf("could not write key: %s", err)
		direktivapps.WriteError(g)
	}

	cmd := exec.Command("/usr/bin/gcloud", "auth", "activate-service-account", "--key-file", "/key.json")
	resp, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("resp len = %v ERROR = %s\n", len(resp), err)
		g.ErrorMessage = fmt.Sprintf("failed auth: %s", resp)
		direktivapps.WriteError(g)
	}

	cmd = exec.Command("/usr/bin/gcloud", obj.Command...)
	resp, err = cmd.CombinedOutput()
	if err != nil {
		g.ErrorMessage = fmt.Sprintf("%s", resp)
		direktivapps.WriteError(g)
	}

	direktivapps.WriteOut(resp, g)

}
