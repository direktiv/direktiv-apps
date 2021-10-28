package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/direktiv/direktiv-apps/pkg/direktivapps"
)

var code = "com.teams.%s.error"

type TeamsInput struct {
	Body map[string]interface{} `json:"body"`
	URL  string                 `json:"url"`
}

func TeamsHandler(w http.ResponseWriter, r *http.Request) {
	var obj TeamsInput
	aid, err := direktivapps.Unmarshal(&obj, r)
	if err != nil {
		direktivapps.RespondWithError(w, fmt.Sprintf(code, "unmarshal-input"), err.Error())
		return
	}

	direktivapps.LogDouble(aid, "reading input...")

	data, err := json.Marshal(obj.Body)
	if err != nil {
		direktivapps.RespondWithError(w, fmt.Sprintf(code, "marshal-body"), err.Error())
		return
	}

	direktivapps.LogDouble(aid, "create request...")

	req, err := http.NewRequest("POST", obj.URL, bytes.NewReader(data))
	if err != nil {
		direktivapps.RespondWithError(w, fmt.Sprintf(code, "create-request"), err.Error())
		return
	}

	req.Header.Add("Content-Type", "application/json")

	direktivapps.LogDouble(aid, "send request...")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		direktivapps.RespondWithError(w, fmt.Sprintf(code, "send-request"), err.Error())
		return
	}
	defer resp.Body.Close()

	direktivapps.LogDouble(aid, "read response...")

	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		direktivapps.RespondWithError(w, fmt.Sprintf(code, "read-response"), err.Error())
		return
	}

	direktivapps.Respond(w, respData)
}

func main() {
	direktivapps.StartServer(TeamsHandler)
}
