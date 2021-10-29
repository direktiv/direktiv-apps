package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/direktiv/direktiv-apps/pkg/direktivapps"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/jwt"
	sheets "google.golang.org/api/sheets/v4"
)

// GoogleServiceAccount is a struct mimicing the service account key json file
type GoogleServiceAccount struct {
	Type        string `json:"type"`
	PrivateKey  string `json:"private_key"`
	ClientEmail string `json:"client_email"`
	TokenURI    string `json:"token_uri"`
}

const code = "com.store.error"

// GoogleInput takes the data required to talk to the sheets API
type GoogleInput struct {
	Authentication GoogleServiceAccount `json:"authentication"`
	Debug          bool                 `json:"debug"`
	SpreadsheetID  string               `json:"spreadSheetID"`
	Range          string               `json:"range"`
	Values         []interface{}        `json:"values"`
}

// AuthURL the api used to grant authentication
const AuthURL = "https://www.googleapis.com/auth/spreadsheets"

func WriteToSpreadsheet(w http.ResponseWriter, r *http.Request) {
	obj := new(GoogleInput)

	aid, err := direktivapps.Unmarshal(obj, r)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	conf := &jwt.Config{
		Email:      obj.Authentication.ClientEmail,
		PrivateKey: []byte(obj.Authentication.PrivateKey),
		TokenURL:   obj.Authentication.TokenURI,
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
	if obj.Debug {
		direktivapps.Log(aid, "JWT has been created and verified")
	}

	service, err := sheets.New(client)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	if obj.Debug {
		direktivapps.Log(aid, "Create a new sheets service")
	}

	var vr sheets.ValueRange
	writeRange := obj.Range
	vr.Values = append(vr.Values, obj.Values)

	if obj.Debug {
		direktivapps.Log(aid, "Appending new sheet values")
		direktivapps.Log(aid, fmt.Sprintf("Writing %v", vr.Values))
	}

	_, err = service.Spreadsheets.Values.Append(obj.SpreadsheetID, writeRange, &vr).ValueInputOption("RAW").Do()
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	direktivapps.Respond(w, []byte{})

}

func main() {
	direktivapps.StartServer(WriteToSpreadsheet)
}
