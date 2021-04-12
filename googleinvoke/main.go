package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/vorteil/direktiv-apps/pkg/direktivapps"
	"google.golang.org/api/idtoken"
	"google.golang.org/api/option"
)

// InputContainerDetails ...
type InputContainerDetails struct {
	ProjectID         string                 `json:"project-id"`
	Region            string                 `json:"region"`
	Function          string                 `json:"function"`
	ServiceAccountKey string                 `json:"serviceAccountKey"`
	Method            string                 `json:"method"`
	Body              map[string]interface{} `json:"body"`
}

const code = "com.googleinvoke.error"

func main() {
	direktivapps.StartServer(GoogleInvoke)
}

func GoogleInvoke(w http.ResponseWriter, r *http.Request) {
	// Read In - Input
	obj := new(InputContainerDetails)
	_, err := direktivapps.Unmarshal(obj, r)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	// Create Authenticated Client
	ctx := context.Background()
	cloudFuncURL := fmt.Sprintf("https://%s-%s.cloudfunctions.net/%s", obj.Region, obj.ProjectID, obj.Function)
	client, err := idtoken.NewClient(ctx, cloudFuncURL, option.WithCredentialsJSON([]byte(obj.ServiceAccountKey)))
	if err != nil {
		direktivapps.RespondWithError(w, code, fmt.Sprintf("Failed to create an Authenticated Client: %s", err.Error()))
		return
	}

	// Read In - Cloud Function Payload
	var cloudFuncPayload []byte
	if obj.Body != nil {
		cloudFuncPayload, err = json.Marshal(obj.Body)
		if err != nil {
			direktivapps.RespondWithError(w, code, err.Error())
			return
		}
	}

	// Create and Send Request
	resp, err := client.Post(cloudFuncURL, "application/json", bytes.NewReader(cloudFuncPayload))
	if err != nil {
		direktivapps.RespondWithError(w, code, fmt.Sprintf("GCP Request failed: %v", err))
		return
	}
	defer resp.Body.Close()

	bv, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		direktivapps.RespondWithError(w, code, fmt.Sprintf("GCP Request, could not read response: %v", err))
		return
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// error more than likely
		direktivapps.RespondWithError(w, code, fmt.Sprintf("Response Message: %s, Response Code: %v \nResponseBody: %s", resp.Status, resp.StatusCode, bv))
		return
	}

	direktivapps.Respond(w, bv)
}
