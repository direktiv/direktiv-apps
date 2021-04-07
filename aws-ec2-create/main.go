package main

import (
	"fmt"
	"github.com/go-playground/validator"
	"log"
	"os"
	"os/exec"
	"text/template"

	"bytes"
	"github.com/vorteil/direktiv-apps/pkg/direktivapps"
	"strings"
)

// InputInstanceDetails ...
type InputInstanceDetails struct {
	// Auth
	Key    string `json:"access-key" validate:"required"`
	Secret string `json:"access-secret" validate:"required"`
	Region string `json:"region" validate:"required"`

	// Required
	ImageID      string `json:"image-id" validate:"required"`
	InstanceType string `json:"instance-type" validate:"required"`

	// Optional
	KeyName         string            `json:"key-name"`
	SubnetID        string            `json:"subnet-id"`
	SecurityGroupID string            `json:"security-group-ids"`
	Tags            map[string]string `json:"tag-specifications"`

	// Util
	TagsFormatted string
}

// Defaults
const (
	AWS_CLI_TEMPLATE = ` ec2 run-instances 
	--image-id {{.ImageID}} 
	--instance-type {{.InstanceType}}
	{{if .KeyName}} --key-name {{.KeyName}}{{else}}{{end}}
	{{if .SubnetID}} --subnet-id {{.SubnetID}}{{else}}{{end}}
	{{if .SecurityGroupID}} --security-group-ids {{.SecurityGroupID}}{{else}}{{end}}
	{{if .TagsFormatted}} --tag-specifications {{.TagsFormatted}}{{else}}{{end}}	`
)

func main() {
	g := direktivapps.ActionError{
		ErrorCode:    "com.aws-ec2-create.error",
		ErrorMessage: "",
	}

	var err error
	obj := new(InputInstanceDetails)
	direktivapps.ReadIn(obj, g)

	// Validate Input
	v := validator.New()
	errArr := v.Struct(obj)
	if errArr != nil {
		for _, e := range errArr.(validator.ValidationErrors) {
			if fe, ok := e.(validator.FieldError); ok {
				log.Printf("Input Error: %s is %s\n", fe.Field(), fe.Tag())
			} else {
				log.Println(e)
			}
		}

		g.ErrorMessage = "Invalid input"
		direktivapps.WriteError(g)
	}

	if obj.KeyName == "" {
		log.Println("WARNING Field key-name is empty: If you do not specify a key pair, you can't connect to the instance unless you choose an AMI that is configured to allow users another way to log in.")
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
	t, err := template.New("gcp").Parse(AWS_CLI_TEMPLATE)
	if err != nil {
		g.ErrorMessage = fmt.Sprintf("Failed to parse AWS Command: %v", err)
		direktivapps.WriteError(g)
	}

	var cliCommand bytes.Buffer
	err = t.Execute(&cliCommand, obj)
	if err != nil {
		g.ErrorMessage = fmt.Sprintf("Failed to create AWS Command: %v", err)
		direktivapps.WriteError(g)
	}

	// Auth
	os.Setenv("AWS_ACCESS_KEY_ID", obj.Key)
	os.Setenv("AWS_SECRET_ACCESS_KEY", obj.Secret)
	os.Setenv("AWS_DEFAULT_REGION", obj.Region)

	cmd := exec.Command("/usr/bin/aws", strings.Fields(cliCommand.String())...)
	resp, err := cmd.CombinedOutput()
	if err != nil {
		g.ErrorMessage = fmt.Sprintf("%s", resp)
		direktivapps.WriteError(g)
	}

	direktivapps.WriteOut(resp, g)
}
