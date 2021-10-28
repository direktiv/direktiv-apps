package main

import (
	"context"
	"fmt"
	"net/http"
	"text/template"
	"time"

	"bytes"
	"io/ioutil"
	"strings"

	"github.com/direktiv/direktiv-apps/pkg/direktivapps"
	"github.com/direktiv/direktiv-apps/pkg/validator"
	"golang.org/x/oauth2/google"
)

// InputInstanceDetails ...
type InputInstanceDetails struct {
	ServiceAccountKey string   `json:"serviceAccountKey" validate:"required"`
	Project           string   `json:"project" validate:"required"`
	ImageProject      string   `json:"imageProject"` // if left unset, defaults to Project
	Name              string   `json:"name" validate:"required"`
	MachineType       string   `json:"machineType" validate:"required"`
	Zone              string   `json:"zone" validate:"required"`
	Region            string   // Extrapolate from Zone
	DiskType          string   `json:"diskType"` // Defaults to GCP_DiskType
	Image             string   `json:"image" validate:"required"`
	ExternalIP        bool     `json:"externalIP"`
	Tags              []string `json:"tags"`
	TagsFormatted     string
}

// Defaults
const (
	GCP_DiskType       = "pd-balanced"
	GCP_AuthURL        = "https://www.googleapis.com/auth/compute"
	GCP_RequestTimeout = 80 * time.Second
	GCP_RequestBody    = `{
		"kind": "compute#instance",
		"name": "{{ .Name}}",
		"zone": "projects/{{ .Project}}/zones/{{ .Zone}}",
		"machineType": "projects/{{ .Project}}/zones/{{ .Zone}}/machineTypes/{{ .MachineType}}",
		"displayDevice": {
		  "enableDisplay": false
		},
		"metadata": {
		  "kind": "compute#metadata",
		  "items": []
		},
		"tags": {
		  "items": [{{ .TagsFormatted}}]
		},
		"disks": [
		  {
			"kind": "compute#attachedDisk",
			"type": "PERSISTENT",
			"boot": true,
			"mode": "READ_WRITE",
			"autoDelete": true,
			"deviceName": "{{ .Name}}",
			"initializeParams": {
			  "sourceImage": "projects/{{ .ImageProject}}/global/images/{{ .Image}}",
			  "diskType": "projects/{{ .Project}}/zones/{{ .Zone}}/diskTypes/{{ .DiskType}}",
			  "diskSizeGb": "10",
			  "labels": {}
			},
			"diskEncryptionKey": {}
		  }
		],
		"canIpForward": false,
		"networkInterfaces": [
		  {
			"kind": "compute#networkInterface",
			"subnetwork": "projects/{{ .Project}}/regions/{{ .Region}}/subnetworks/default",
			{{if .ExternalIP}} "accessConfigs": [
				{
				  "kind": "compute#accessConfig",
				  "name": "External NAT",
				  "type": "ONE_TO_ONE_NAT",
				  "networkTier": "PREMIUM"
				}
			  ],
			  {{else}}{{end}}"aliasIpRanges": []
		  }
		],
		"description": "",
		"labels": {},
		"scheduling": {
		  "preemptible": false,
		  "onHostMaintenance": "MIGRATE",
		  "automaticRestart": true,
		  "nodeAffinities": []
		},
		"deletionProtection": false,
		"reservationAffinity": {
		  "consumeReservationType": "ANY_RESERVATION"
		},
		"confidentialInstanceConfig": {
		  "enableConfidentialCompute": false
		}
	  }`
)

const code = "com.gcloud-instance-create.error"

func main() {
	direktivapps.StartServer(GCPComputeCreate)
}

func GCPComputeCreate(w http.ResponseWriter, r *http.Request) {
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

	// Load Defaults
	if obj.ImageProject == "" {
		obj.ImageProject = obj.Project
	}

	if obj.DiskType == "" {
		obj.DiskType = GCP_DiskType
	}

	// Get Region
	if i := strings.LastIndex(obj.Zone, "-"); i > 0 {
		obj.Region = obj.Zone[:i]
	} else {
		direktivapps.RespondWithError(w, code, fmt.Sprintf("Potentially invalid Zone: could not extract region from inputted zone \"%s\"", obj.Zone))
		return
	}

	// Format Tags
	for i, tag := range obj.Tags {
		obj.TagsFormatted += `"` + tag + `"`

		if i < len(obj.Tags)-1 {
			obj.TagsFormatted += `, `
		}
	}

	// Create request body from template
	t, err := template.New("gcp").Parse(GCP_RequestBody)
	if err != nil {
		direktivapps.RespondWithError(w, code, fmt.Sprintf("Failed to parse GCP Request: %v", err))
		return
	}

	var reqBody bytes.Buffer
	err = t.Execute(&reqBody, obj)
	if err != nil {
		direktivapps.RespondWithError(w, code, fmt.Sprintf("Failed to create GCP Request: %v", err))
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
	requestURL := fmt.Sprintf("https://www.googleapis.com/compute/v1/projects/%s/zones/%s/instances", obj.Project, obj.Zone)
	resp, err := client.Post(requestURL, "application/json", bytes.NewReader(reqBody.Bytes()))
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
