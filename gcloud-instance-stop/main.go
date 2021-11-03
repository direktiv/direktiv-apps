package main

import (
	"context"
	"fmt"
	"time"

	"io/ioutil"
	"net/http"
	"strings"

	"github.com/direktiv/direktiv-apps/pkg/direktivapps"
	"github.com/direktiv/direktiv-apps/pkg/validator"
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
	code               = "com.gcloud-instance-stop.error"
	GCP_AuthURL        = "https://www.googleapis.com/auth/compute"
	GCP_RequestTimeout = 80 * time.Second
)

func main() {
	direktivapps.StartServer(GCPComputeStop)
}

func GCPComputeStop(w http.ResponseWriter, r *http.Request) {
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
			direktivapps.LogDouble(aid, fmt.Sprintf("Input Error: %s is required\n", mf))
		}

		direktivapps.RespondWithError(w, code, fmt.Sprintf("Invalid input: Fields [%s] are required", strings.Join(missingFields, ",")))
		return
	}

	// Create client
	conf, err := google.JWTConfigFromJSON([]byte(obj.ServiceAccountKey), GCP_AuthURL)
	if err != nil {
		direktivapps.RespondWithError(w, code, fmt.Sprintf("Failed to create GCP JWT from service account key: %v", err))
		return
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
		direktivapps.RespondWithError(w, code, fmt.Sprintf("GCP Request creation failed: %v", err))
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		direktivapps.RespondWithError(w, code, fmt.Sprintf("GCP Request failed: %v", err))
		return
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		direktivapps.RespondWithError(w, code, fmt.Sprintf("GCP Request, could not read response: %v", err))
		return
	}

	if resp.StatusCode != 200 {
		direktivapps.RespondWithError(w, code, fmt.Sprintf("GCP Request completed with errors. Response:\n%s\n", string(bytes)))
		return
	}

	direktivapps.Respond(w, bytes)
}
