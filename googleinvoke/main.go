package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/vorteil/direktiv-apps/pkg/direktivapps"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/jwt"
)

// AuthURL to fetch a new token from for executing cloud functions
const AuthURL = "https://www.googleapis.com/auth/cloudfunctions"

// InputContainerDetails ...
type InputContainerDetails struct {
	Region            string                 `json:"region"`
	Function          string                 `json:"function"`
	ServiceAccountKey string                 `json:"serviceAccountKey"`
	Method            string                 `json:"method"`
	Body              map[string]interface{} `json:"body"`
}

// Authentication is the struct to unmarshal the service account key into
type Authentication struct {
	Type        string `json:"type"`
	ProjectID   string `json:"project_id"`
	PrivateKey  string `json:"private_key"`
	ClientEmail string `json:"client_email"`
	TokenURI    string `json:"token_uri"`
}

const code = "com.googleinvoke.error"

func main() {
	direktivapps.StartServer(GoogleInvoke)
}

func GoogleInvoke(w http.ResponseWriter, r *http.Request) {
	obj := new(InputContainerDetails)
	_, err := direktivapps.Unmarshal(obj, r)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	authentication := &Authentication{}
	// unmarshal into another struct
	err = json.Unmarshal([]byte(obj.ServiceAccountKey), authentication)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	conf := &jwt.Config{
		Email:      authentication.ClientEmail,
		PrivateKey: []byte(authentication.PrivateKey),
		TokenURL:   authentication.TokenURI,
		Scopes: []string{
			AuthURL,
		},
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	ctx := context.Background()
	cli := &http.Client{Transport: tr}
	ctx = context.WithValue(ctx, oauth2.HTTPClient, cli)

	client := conf.Client(ctx)

	var payload []byte
	if obj.Body != nil {
		payload, err = json.Marshal(obj.Body)
		if err != nil {
			direktivapps.RespondWithError(w, code, err.Error())
			return
		}
	}

	req, err := http.NewRequest(obj.Method, fmt.Sprintf("https://%s-%s.cloudfunctions.net/%s", obj.Region, authentication.ProjectID, obj.Function), bytes.NewReader(payload))
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	defer resp.Body.Close()

	bv, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// error more than likely
		direktivapps.RespondWithError(w, code, fmt.Sprintf("Response Message: %s, Response Code: %v \nResponseBody: %s", resp.Status, resp.StatusCode, bv))
		return
	}

	direktivapps.Respond(w, bv)
}
