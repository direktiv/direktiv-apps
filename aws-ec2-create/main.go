package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"text/template"

	"bytes"
	"strings"

	"github.com/vorteil/direktiv-apps/pkg/direktivapps"
	"github.com/vorteil/direktiv-apps/pkg/validator"
)

// InputInstanceDetails ...
type InputInstanceDetails struct {
	// Auth
	Key    string `json:"access-key" validate:"required"`
	Secret string `json:"access-secret" validate:"required"`

	// Required
	Region       string `json:"region" validate:"required"`
	ImageID      string `json:"image-id" validate:"required"`
	InstanceType string `json:"instance-type" validate:"required"`

	// Optional
	Name            string            `json:"name"`
	KeyName         string            `json:"key-name"`
	SubnetID        string            `json:"subnet-id"`
	SecurityGroupID string            `json:"security-group-ids"`
	Tags            map[string]string `json:"tag-specifications"`

	// Util
	TagsFormatted string
}

// Defaults
const (
	code             = "com.aws-ec2-create.error"
	AWS_CLI_TEMPLATE = ` ec2 run-instances 
	--image-id {{.ImageID}} 
	--instance-type {{.InstanceType}}
	{{if .Name}} --name {{.Name}}{{else}}{{end}}
	{{if .KeyName}} --key-name {{.KeyName}}{{else}}{{end}}
	{{if .SubnetID}} --subnet-id {{.SubnetID}}{{else}}{{end}}
	{{if .SecurityGroupID}} --security-group-ids {{.SecurityGroupID}}{{else}}{{end}}
	{{if .TagsFormatted}} --tag-specifications {{.TagsFormatted}}{{else}}{{end}}	`
)

func main() {
	direktivapps.StartServer(AWSInstanceCreate)
}

func AWSInstanceCreate(w http.ResponseWriter, r *http.Request) {
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
			direktivapps.Log(aid, fmt.Sprintf("Input Error: %s is required\n", mf))
		}

		direktivapps.RespondWithError(w, code, fmt.Sprintf("Invalid input: Fields [%s] are required", strings.Join(missingFields, ",")))
		return
	}

	if obj.KeyName == "" {
		direktivapps.Log(aid, "WARNING Field key-name is empty: If you do not specify a key pair, you can't connect to the instance unless you choose an AMI that is configured to allow users another way to log in.")
	}

	// Format Tags
	tagsJSON := make([]string, 0)
	for tagKey, tagValue := range obj.Tags {
		tagsJSON = append(tagsJSON, fmt.Sprintf("{Key=%s,Value=%s}", tagKey, tagValue))
	}

	if len(tagsJSON) > 0 {
		obj.TagsFormatted = fmt.Sprintf("'ResourceType=instance,Tags=[%s]'", strings.Join(tagsJSON, ","))
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
