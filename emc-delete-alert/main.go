package main

import (
	"crypto/tls"
	b64 "encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"github.com/direktiv/direktiv-apps/pkg/direktivapps"
	"golang.org/x/net/publicsuffix"
)

const code = "com.emc.delete-alert.%s.error"

// request the input object
type request struct {
	URL      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
	ID       string `json:"id"`
}

func authenticate(client *http.Client, urlpath, username, password, aid string) error {
	u, err := url.Parse(fmt.Sprintf("%s/api/types/alert/instances", urlpath))
	if err != nil {
		return err
	}

	q := u.Query()

	direktivapps.LogDouble(aid, "setting fields and compact params")
	q.Set("fields", "message,state")
	q.Set("compact", "true")

	u.RawQuery = q.Encode()

	direktivapps.LogDouble(aid, fmt.Sprintf("requesting '%s'", u.String()))
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return err
	}

	direktivapps.LogDouble(aid, "Adding required headers")

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-EMC-REST-CLIENT", "true")

	direktivapps.LogDouble(aid, "Adding authorization header")
	sEnc := b64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", username, password)))

	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", sEnc))

	direktivapps.LogDouble(aid, "Sending request")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	return nil
}

func Request(w http.ResponseWriter, r *http.Request) {
	obj := new(request)

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

	direktivapps.LogDouble(aid, "authenticate")
	err = authenticate(client, obj.URL, obj.Username, obj.Password, aid)
	if err != nil {
		direktivapps.RespondWithError(w, fmt.Sprintf(code, "authenticate"), err.Error())
		return
	}

	u, err := url.Parse(fmt.Sprintf("%s//api/instances/alert/%s", obj.URL, obj.ID))
	if err != nil {
		direktivapps.RespondWithError(w, fmt.Sprintf(code, "url-parse"), err.Error())
		return
	}

	req, err := http.NewRequest("DELETE", u.String(), nil)
	if err != nil {
		direktivapps.RespondWithError(w, fmt.Sprintf(code, "create-request"), err.Error())
		return
	}

	direktivapps.LogDouble(aid, "Adding required headers")

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-EMC-REST-CLIENT", "true")

	direktivapps.LogDouble(aid, "Adding authorization header")
	sEnc := b64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", obj.Username, obj.Password)))
	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", sEnc))

	direktivapps.LogDouble(aid, "Sending request")
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

	direktivapps.LogDouble(aid, fmt.Sprintf("%s", data))

	direktivapps.Respond(w, []byte{})
}

func main() {
	direktivapps.StartServer(Request)
}
