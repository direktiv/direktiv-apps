package main

import (
	"golang.org/x/net/publicsuffix"

	"bytes"
	"crypto/tls"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"

	"github.com/vorteil/direktiv-apps/pkg/direktivapps"
)

const code = "com.%s.error"

// request the input object for the requester container
type request struct {
	Method             string                 `json:"method"`
	URL                string                 `json:"url"`
	Body               interface{}            `json:"body"`
	Headers            map[string]interface{} `json:"headers"`
	Params             map[string]interface{} `json:"params"`
	Username           string                 `json:"username"`
	Password           string                 `json:"password"`
	InsecureSkipVerify bool                   `json:"insecureSkipVerify"`
}

// output for the requester container
type output struct {
	Body       interface{} `json:"body,omitempty"` // when the response is able to be unmarshalled
	Headers    http.Header `json:"headers"`
	StatusCode int         `json:"status-code"`
	Status     string      `json:"status"`
}

func Request(w http.ResponseWriter, r *http.Request) {

	obj := new(request)
	aid, err := direktivapps.Unmarshal(obj, r)
	if err != nil {
		direktivapps.RespondWithError(w, fmt.Sprintf(code, "unmarshal-input"), err.Error())
		return
	}

	var b []byte

	direktivapps.Log(aid, "Creating cookie jar")

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
				InsecureSkipVerify: obj.InsecureSkipVerify,
			},
		},
	}

	direktivapps.Log(aid, "Creating new request")

	if obj.Body != nil {
		switch v := obj.Body.(type) {
		case string:
			direktivapps.Log(aid, "Body is a string ignore marshal.")
			b = []byte(obj.Body.(string))
		default:
			direktivapps.Log(aid, fmt.Sprintf("Body is of type %v", v))
			b, err = json.Marshal(obj.Body)
			if err != nil {
				direktivapps.RespondWithError(w, fmt.Sprintf(code, "marshal-body"), err.Error())
				return
			}
		}

		direktivapps.Log(aid, "Body exists, attaching to the request")
	}

	direktivapps.Log(aid, "Creating URL...")
	u, err := url.Parse(obj.URL)
	if err != nil {
		direktivapps.RespondWithError(w, fmt.Sprintf(code, "url-parse"), err.Error())
		return
	}

	q := u.Query()
	for k, v := range obj.Params {
		var actualVal string
		// Handle other types provided and convert to string automatically
		switch t := v.(type) {
		case bool:
			actualVal = strconv.FormatBool(t)
		case float64:
			actualVal = strconv.FormatFloat(t, 'f', 6, 64)
		case string:
			actualVal = t
		}
		direktivapps.Log(aid, fmt.Sprintf("Adding param %s=%s", k, actualVal))
		q.Set(k, actualVal)
	}

	u.RawQuery = q.Encode()

	req, err := http.NewRequest(obj.Method, u.String(), bytes.NewReader(b))
	if err != nil {
		direktivapps.RespondWithError(w, fmt.Sprintf(code, "create-request"), err.Error())
		return
	}

	for k, v := range obj.Headers {
		var actualVal string
		// Handle other types provided and convert to string automatically
		switch t := v.(type) {
		case bool:
			actualVal = strconv.FormatBool(t)
		case float64:
			actualVal = strconv.FormatFloat(t, 'f', 6, 64)
		case string:
			actualVal = t
		}

		// Adding a header requires it to be a string so might as well sprintf
		req.Header.Add(k, actualVal)
	}

	if obj.Username != "" && obj.Password != "" {
		sEnc := b64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", obj.Username, obj.Password)))
		req.Header.Add("Authorization", fmt.Sprintf("Basic %s", sEnc))
		direktivapps.Log(aid, "Adding Basic authorization header")
	}

	direktivapps.Log(aid, "Sending request...")

	resp, err := client.Do(req)
	if err != nil {
		direktivapps.RespondWithError(w, fmt.Sprintf(code, "send-request"), err.Error())
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		direktivapps.RespondWithError(w, fmt.Sprintf(code, "read-resp-body"), err.Error())
		return
	}

	var mapBody map[string]interface{}
	var dataBody interface{}
	var responding output
	responding.Status = resp.Status
	responding.StatusCode = resp.StatusCode
	responding.Headers = resp.Header

	// if body is unable to be marshalled treat as a byte array
	err = json.Unmarshal(body, &mapBody)
	if err != nil {
		err = json.Unmarshal(body, &dataBody)
		if err != nil {
			direktivapps.RespondWithError(w, fmt.Sprintf(code, "unmarshalling"), err.Error())
			return
		}
		responding.Body = dataBody
	} else {
		responding.Body = mapBody
	}

	data, err := json.Marshal(responding)
	if err != nil {
		direktivapps.RespondWithError(w, fmt.Sprintf(code, "marshal-output"), err.Error())
		return
	}

	direktivapps.Respond(w, data)
}

func main() {
	direktivapps.StartServer(Request)
}
