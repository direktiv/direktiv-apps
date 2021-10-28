package main

import (
	"bytes"
	"crypto/tls"
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"

	"github.com/direktiv/direktiv-apps/pkg/direktivapps"
	"golang.org/x/net/publicsuffix"
)

const code = "com.emc.modify-alert.%s.error"

type AlertResponse struct {
	Base    string    `json:"@base"`
	Updated time.Time `json:"updated"`
	Links   []struct {
		Rel  string `json:"rel"`
		Href string `json:"href"`
	} `json:"links"`
	Entries []struct {
		Content struct {
			ID      string `json:"id"`
			State   int    `json:"state"`
			Message string `json:"message"`
		} `json:"content"`
	} `json:"entries"`
}

// EMCModifyAlert the input object for the requester container
type EMCModifyAlert struct {
	URL      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
	Message  string `json:"message"`
}

type responseid struct {
	ID string `json:"id"`
}

func Request(w http.ResponseWriter, r *http.Request) {
	obj := new(EMCModifyAlert)

	aid, err := direktivapps.Unmarshal(obj, r)
	if err != nil {
		direktivapps.RespondWithError(w, fmt.Sprintf(code, "unmarshal-input"), err.Error())
		return
	}

	jar, err := cookiejar.New(&cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	})
	if err != nil {
		direktivapps.RespondWithError(w, fmt.Sprintf(code, "cookie-jar"), err.Error())
		return
	}

	client := &http.Client{
		Jar: jar,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	id, err := getAlertId(client, obj.URL, obj.Message, obj.Username, obj.Password, aid)
	if err != nil {
		direktivapps.RespondWithError(w, fmt.Sprintf(code, "getAlertId"), err.Error())
		return
	}

	direktivapps.Log(aid, fmt.Sprintf("acknowledge '%s' alert", id))
	u, err := url.Parse(fmt.Sprintf("%s//api/instances/alert/%s/action/modify", obj.URL, id))
	if err != nil {
		direktivapps.RespondWithError(w, fmt.Sprintf(code, "url-parse"), err.Error())
		return
	}

	modifyArgs := fmt.Sprintf(`{
		"isAcknowledged": true
	}`)

	req, err := http.NewRequest("POST", u.String(), bytes.NewReader([]byte(modifyArgs)))
	if err != nil {
		direktivapps.RespondWithError(w, fmt.Sprintf(code, "create-request"), err.Error())
		return
	}

	direktivapps.Log(aid, "Adding required headers")

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-EMC-REST-CLIENT", "true")

	direktivapps.Log(aid, "Adding authorization header")
	sEnc := b64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", obj.Username, obj.Password)))
	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", sEnc))

	direktivapps.Log(aid, "Sending request")
	resp, err := client.Do(req)
	if err != nil {
		direktivapps.RespondWithError(w, fmt.Sprintf(code, "send-response"), err.Error())
		return
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		direktivapps.RespondWithError(w, fmt.Sprintf(code, "read-response"), err.Error())
		return
	}

	direktivapps.Log(aid, fmt.Sprintf("%s", data))

	var response responseid
	response.ID = id

	body, err := json.Marshal(response)
	if err != nil {
		direktivapps.RespondWithError(w, fmt.Sprintf(code, "marshal"), err.Error())
		return
	}

	direktivapps.Respond(w, body)
}

func getAlertId(client *http.Client, urlpath, message, username, password, aid string) (string, error) {
	u, err := url.Parse(fmt.Sprintf("%s/api/types/alert/instances", urlpath))
	if err != nil {
		return "", err
	}

	q := u.Query()

	direktivapps.Log(aid, "setting fields and compact params")
	q.Set("fields", "message,state")
	q.Set("compact", "true")

	u.RawQuery = q.Encode()

	direktivapps.Log(aid, fmt.Sprintf("requesting '%s'", u.String()))
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return "", err
	}

	direktivapps.Log(aid, "Adding required headers")

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-EMC-REST-CLIENT", "true")

	direktivapps.Log(aid, "Adding authorization header")
	sEnc := b64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", username, password)))

	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", sEnc))

	direktivapps.Log(aid, "Sending request")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var results AlertResponse
	err = json.Unmarshal(data, &results)
	if err != nil {
		return "", err
	}

	for _, r := range results.Entries {
		if r.Content.Message == message && r.Content.State == 0 {
			return r.Content.ID, nil
		}
	}
	return "", errors.New("unable to find alert id")
}

func main() {
	direktivapps.StartServer(Request)
}
