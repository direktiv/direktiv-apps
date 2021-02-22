package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/jwt"
	sheets "google.golang.org/api/sheets/v4"
)

type ActionError struct {
	ErrorCode    string `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
}

// GoogleServiceAccount is a struct mimicing the service account key json file
type GoogleServiceAccount struct {
	Type        string `json:"type"`
	PrivateKey  string `json:"private_key"`
	ClientEmail string `json:"client_email"`
	TokenURI    string `json:"token_uri"`
}

// GoogleInput takes the data required to talk to the sheets API
type GoogleInput struct {
	Authentication GoogleServiceAccount `json:"authentication"`
	SpreadsheetID  string               `json:"spreadSheetID"`
	Range          string               `json:"range"`
	Values         []interface{}        `json:"values"`
}

// AuthURL the api used to grant authentication
const AuthURL = "https://www.googleapis.com/auth/spreadsheets"

func main() {
	gi := &GoogleInput{}

	g := ActionError{
		ErrorCode:    "com.googlepusher.error",
		ErrorMessage: "",
	}

	data, err := ioutil.ReadFile("/direktiv-data/data.in")
	if err != nil {
		g.ErrorMessage = err.Error()
		writeError(g)
		return
	}

	err = json.Unmarshal(data, gi)
	if err != nil {
		g.ErrorMessage = err.Error()
		writeError(g)
		return
	}

	fmt.Printf("KEY: \n%s\n", gi.Authentication.PrivateKey)

	conf := &jwt.Config{
		Email:      gi.Authentication.ClientEmail,
		PrivateKey: []byte(gi.Authentication.PrivateKey),
		TokenURL:   gi.Authentication.TokenURI,
		Scopes: []string{
			AuthURL,
		},
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	ctx := context.Background()
	sslcli := &http.Client{Transport: tr}
	ctx = context.WithValue(ctx, oauth2.HTTPClient, sslcli)

	client := conf.Client(ctx)

	service, err := sheets.New(client)
	if err != nil {
		g.ErrorMessage = err.Error()
		writeError(g)
		return
	}

	var vr sheets.ValueRange
	writeRange := gi.Range
	vr.Values = append(vr.Values, gi.Values)

	_, err = service.Spreadsheets.Values.Append(gi.SpreadsheetID, writeRange, &vr).ValueInputOption("RAW").Do()
	if err != nil {
		g.ErrorMessage = err.Error()
		writeError(g)
		return
	}

	finishRunning([]byte{})
}

// finishRunning will write to a file and or print the json body to stdout and exits
func finishRunning(eb []byte) {
	var err error
	err = ioutil.WriteFile("/direktiv-data/data.out", eb, 0755)
	if err != nil {
		log.Fatal("can not write out data")
		return
	}
}

// writeError
func writeError(g ActionError) {
	b, _ := json.Marshal(g)
	err := ioutil.WriteFile("/direktiv-data/error.json", b, 0755)
	if err != nil {
		log.Fatal("can not write json error")
		return
	}
}
