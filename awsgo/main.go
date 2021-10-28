package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"

	"github.com/direktiv/direktiv-apps/pkg/direktivapps"
)

// AwsInputDetails ...
type AwsInputDetails struct {
	Key     string   `json:"access-key"`
	Secret  string   `json:"access-secret"`
	Region  string   `json:"region"`
	Command []string `json:"command"`
}

const code = "com.awsgo.error"

func main() {
	direktivapps.StartServer(AWSGo)
}

func AWSGo(w http.ResponseWriter, r *http.Request) {
	var err error
	obj := new(AwsInputDetails)
	_, err = direktivapps.Unmarshal(obj, r)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	os.Setenv("AWS_ACCESS_KEY_ID", obj.Key)
	os.Setenv("AWS_SECRET_ACCESS_KEY", obj.Secret)
	os.Setenv("AWS_DEFAULT_REGION", obj.Region)

	cmd := exec.Command("/usr/bin/aws", obj.Command...)
	resp, err := cmd.CombinedOutput()
	if err != nil {
		direktivapps.RespondWithError(w, code, fmt.Sprintf("%s", resp))
		return
	}

	direktivapps.Respond(w, resp)
}
