package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"

	"github.com/vorteil/direktiv-apps/pkg/direktivapps"
)

const credsPath = "credentials"
const confAws = `[default]
aws_access_key_id=%s
aws_secret_access_key=%s
`

// AwsInputDetails ...
type AwsInputDetails struct {
	Key     string   `json:"access-key"`
	Secret  string   `json:"access-secret"`
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

	usr, err := user.Current()
	if err != nil {
		g.ErrorMessage = err.Error()
		direktivapps.WriteError(g)
	}

	err = os.MkdirAll(filepath.Join(usr.HomeDir, ".aws"), 0700)
	if err != nil {
		g.ErrorMessage = err.Error()
		direktivapps.WriteError(g)
	}

	err = ioutil.WriteFile(filepath.Join(usr.HomeDir, ".aws", "credentials"), []byte(fmt.Sprintf(confAws, obj.Key, obj.Secret)), 0700)
	if err != nil {
		g.ErrorMessage = err.Error()
		direktivapps.WriteError(g)
	}

	cmd := exec.Command("/usr/bin/aws", obj.Command...)
	resp, err := cmd.CombinedOutput()
	if err != nil {
		g.ErrorMessage = fmt.Sprintf("%s", resp)
	}

	direktivapps.WriteOut(resp, g)
}
