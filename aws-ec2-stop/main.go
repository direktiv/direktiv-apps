package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"text/template"

	"bytes"
	"strings"

	"github.com/direktiv/direktiv-apps/pkg/direktivapps"
	"github.com/direktiv/direktiv-apps/pkg/validator"
)

// InputInstanceDetails ...
type InputInstanceDetails struct {
	// Auth
	Key    string `json:"access-key" validate:"required"`
	Secret string `json:"access-secret" validate:"required"`

	// Required
	Region     string `json:"region" validate:"required"`
	InstanceID string `json:"instance-id" validate:"required"`
}

// Defaults
const (
	code             = "com.aws-ec2-stop.error"
	AWS_CLI_TEMPLATE = ` ec2 stop-instances 
	--instance-ids {{.InstanceID}}`
)

func main() {
	direktivapps.StartServer(AWSInstanceStop)
}

func AWSInstanceStop(w http.ResponseWriter, r *http.Request) {
	var err error
	obj := new(InputInstanceDetails)
	aid, err := direktivapps.Unmarshal(obj, r)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	// Validate Input
	v := validator.CreateValidator()

	if missingFields := v.ValidateRequired(obj); len(missingFields) > 0 {
		for _, mf := range missingFields {
			direktivapps.LogDouble(aid, "Input Error: %s is required\n", mf)
		}

		direktivapps.RespondWithError(w, code, fmt.Sprintf("Invalid input: Fields [%s] are required", strings.Join(missingFields, ",")))
		return
	}

	// Create cli command from template
	t, err := template.New("aws").Parse(AWS_CLI_TEMPLATE)
	if err != nil {
		direktivapps.RespondWithError(w, code, fmt.Sprintf("Failed to parse AWS Command: %v", err))
		return
	}

	var cliCommand bytes.Buffer
	err = t.Execute(&cliCommand, obj)
	if err != nil {
		direktivapps.RespondWithError(w, code, fmt.Sprintf("Failed to create AWS Command: %v", err))
		return
	}

	// Auth
	os.Setenv("AWS_ACCESS_KEY_ID", obj.Key)
	os.Setenv("AWS_SECRET_ACCESS_KEY", obj.Secret)
	os.Setenv("AWS_DEFAULT_REGION", obj.Region)

	cmd := exec.Command("/usr/bin/aws", strings.Fields(cliCommand.String())...)
	resp, err := cmd.CombinedOutput()
	if err != nil {
		direktivapps.RespondWithError(w, code, fmt.Sprintf("%s", resp))
		return
	}

	direktivapps.Respond(w, resp)
}
