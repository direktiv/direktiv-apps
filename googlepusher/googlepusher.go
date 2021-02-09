package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

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

// GoogleInput takes the data required to talk to the sheets API
type GoogleInput struct {
	Authentication GoogleServiceAccount `json:"authentication"`
	SpreadsheetID  string               `json:"spreadSheetID"`
	Range          string               `json:"range"`
	Values         []interface{}        `json:"values"`
}

// EndBody is the response of this library after a request
type EndBody struct {
	Error         string `json:"error"`
	Response      string `json:"response"`
	Status        int    `json:"statusCode"`
	StatusMessage string `json:"status"`
}

// AuthURL the api used to grant authentication
const AuthURL = "https://www.googleapis.com/auth/spreadsheets"

func main() {
	eb := &EndBody{}
	gi := &GoogleInput{}

	// inputFile := os.Args[1]
	// outputFile := os.Args[2]

	// b, err := ioutil.ReadFile(inputFile)
	// if err != nil {
	// eb.Error = err.Error()
	// finishRunning(outputFile, eb)
	// }

	data, err := ioutil.ReadFile("/direktiv-data/data.in")
	if err != nil {
		eb.Error = err.Error()
		finishRunning(eb)
		return
	}

	err = json.Unmarshal(data, gi)
	if err != nil {
		eb.Error = err.Error()
		finishRunning(eb)
	}

	conf := &jwt.Config{
		Email:      gi.Authentication.ClientEmail,
		PrivateKey: []byte(gi.Authentication.PrivateKey),
		TokenURL:   gi.Authentication.TokenURI,
		Scopes: []string{
			AuthURL,
		},
	}

	client := conf.Client(oauth2.NoContext)

	service, err := sheets.New(client)
	if err != nil {
		eb.Error = err.Error()
		finishRunning(eb)
	}

	var vr sheets.ValueRange
	writeRange := gi.Range
	vr.Values = append(vr.Values, gi.Values)

	_, err = service.Spreadsheets.Values.Append(gi.SpreadsheetID, writeRange, &vr).ValueInputOption("RAW").Do()
	if err != nil {
		eb.Error = err.Error()
	}

	finishRunning(eb)
}

// finishRunning will write to a file and or print the json body to stdout and exits
func finishRunning(eb *EndBody) {
	ms, _ := json.Marshal(eb)
	err := ioutil.WriteFile("/direktiv-data/data.out", []byte(ms), 0755)
	if err != nil {
		log.Fatal("can not write out data")
		return
	}
	os.Exit(0)
}
