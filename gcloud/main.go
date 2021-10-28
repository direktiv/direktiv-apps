package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"

	"github.com/direktiv/direktiv-apps/pkg/direktivapps"
)

// InputContainerDetails ...
type InputContainerDetails struct {
	ServiceAccountKey string   `json:"serviceAccountKey"`
	Command           []string `json:"command"`
	Project           string   `json:"project"`
}

const code = "com.googlecli.error"

func main() {
	direktivapps.StartServer(GCloud)
}

func GCloud(w http.ResponseWriter, r *http.Request) {
	obj := new(InputContainerDetails)
	_, err := direktivapps.Unmarshal(obj, r)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	// json format flag
	obj.Command = append(obj.Command, `--format="json"`)

	if obj.Project == "" {
		direktivapps.RespondWithError(w, code, "input project cannot be empty")
		return
	}

	err = ioutil.WriteFile("/key.json", []byte(obj.ServiceAccountKey), 0644)
	if err != nil {
		direktivapps.RespondWithError(w, code, fmt.Sprintf("could not write key: %s", err))
		return
	}

	cmd := exec.Command("/root/google-cloud-sdk/bin/gcloud", "auth", "activate-service-account", "--key-file", "/key.json")
	resp, err := cmd.CombinedOutput()
	if err != nil {
		if len(resp) > 0 {
			direktivapps.RespondWithError(w, code, fmt.Sprintf("failed auth: %s", resp))
		} else {
			direktivapps.RespondWithError(w, code, fmt.Sprintf("failed auth: %s", err.Error()))
		}
		return
	}

	cmd = exec.Command("/root/google-cloud-sdk/bin/gcloud", "config", "set", "project", obj.Project)
	resp, err = cmd.CombinedOutput()
	if err != nil {
		if len(resp) > 0 {
			direktivapps.RespondWithError(w, code, fmt.Sprintf("invalid project: %s", resp))
		} else {
			direktivapps.RespondWithError(w, code, fmt.Sprintf("invalid project: %s", err.Error()))
		}
		return
	}

	cmd = exec.Command("/root/google-cloud-sdk/bin/gcloud", obj.Command...)
	resp, err = cmd.CombinedOutput()
	if err != nil {
		if len(resp) > 0 {
			direktivapps.RespondWithError(w, code, fmt.Sprintf("%s", resp))
		} else {
			direktivapps.RespondWithError(w, code, fmt.Sprintf("%s", err.Error()))
		}
		return
	}

	direktivapps.Respond(w, resp)
}
