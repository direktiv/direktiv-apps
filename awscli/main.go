package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/vorteil/direktiv-apps/pkg/direktivapps"
)

// AwsInputDetails ...
type AwsInputDetails struct {
	Key     string   `json:"access-key"`
	Secret  string   `json:"access-secret"`
	Region  string   `json:"region"`
	Command []string `json:"command"`
}

func main() {
	g := direktivapps.ActionError{
		ErrorCode:    "com.awsgo.error",
		ErrorMessage: "",
	}

	var err error
	obj := new(AwsInputDetails)
	direktivapps.ReadIn(obj, g)

	os.Setenv("AWS_ACCESS_KEY_ID", obj.Key)
	os.Setenv("AWS_SECRET_ACCESS_KEY", obj.Secret)
	os.Setenv("AWS_DEFAULT_REGION", obj.Region)

	cmd := exec.Command("/usr/bin/aws", obj.Command...)
	resp, err := cmd.CombinedOutput()
	if err != nil {
		g.ErrorMessage = fmt.Sprintf("%s", resp)
	}

	direktivapps.WriteOut(resp, g)
}
