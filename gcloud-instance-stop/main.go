package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"io/ioutil"
	"net/http"
	"strings"

	"github.com/go-playground/validator"
	"github.com/vorteil/direktiv-apps/pkg/direktivapps"
	"golang.org/x/oauth2/google"
)

// InputInstanceDetails ...
type InputInstanceDetails struct {
	ServiceAccountKey string `json:"serviceAccountKey" validate:"required"`
	Project           string `json:"project" validate:"required"`
	Zone              string `json:"zone" validate:"required"`
	InstanceID        string `json:"instanceID" validate:"required"`
}

// Defaults
const (
	GCP_AuthURL        = "https://www.googleapis.com/auth/compute"
	GCP_RequestTimeout = 80 * time.Second
)

func main() {
	g := direktivapps.ActionError{
		ErrorCode:    "com.gcloud-instance-stop.error",
		ErrorMessage: "",
	}

	obj := new(InputInstanceDetails)
	direktivapps.ReadIn(obj, g)

	// Validate Input
	v := validator.New()
	errArr := v.Struct(obj)
	if errArr != nil {
		invalidFields := make([]string, 0)
		for _, e := range errArr.(validator.ValidationErrors) {
			if fe, ok := e.(validator.FieldError); ok {
				invalidFields = append(invalidFields, fe.Field())
				log.Printf("Input Error: %s is %s\n", fe.Field(), fe.Tag())
			} else {
				log.Println(e)
			}
		}

		g.ErrorMessage = fmt.Sprintf("Invalid input: Fields [%s] are required", strings.Join(invalidFields, ","))
		direktivapps.WriteError(g)
	}

	// Create client
	conf, err := google.JWTConfigFromJSON([]byte(obj.ServiceAccountKey), GCP_AuthURL)
	if err != nil {
		g.ErrorMessage = fmt.Sprintf("Failed to create GCP JWT from service account key: %v", err)
		direktivapps.WriteError(g)
	}

	ctx := context.Background()
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(GCP_RequestTimeout))
	defer cancel()
	conf.Expires = GCP_RequestTimeout
	client := conf.Client(ctx)

	// Do Request
	requestURL := fmt.Sprintf("https://www.googleapis.com/compute/v1/projects/%s/zones/%s/instances/%s/stop", obj.Project, obj.Zone, obj.InstanceID)
	req, err := http.NewRequest("POST", requestURL, nil)
	if err != nil {
		g.ErrorMessage = fmt.Sprintf("GCP Request creation failed: %v", err)
		direktivapps.WriteError(g)
	}

	resp, err := client.Do(req)
	if err != nil {
		g.ErrorMessage = fmt.Sprintf("GCP Request failed: %v", err)
		direktivapps.WriteError(g)
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		g.ErrorMessage = fmt.Sprintf("GCP Request, could not read response: %v", err)
		direktivapps.WriteError(g)
	}

	if resp.StatusCode != 200 {
		g.ErrorMessage = fmt.Sprintf("GCP Request completed with errors. Response:\n%s\n", string(bytes))
		direktivapps.WriteError(g)
	}

	direktivapps.WriteOut(bytes, g)
}
