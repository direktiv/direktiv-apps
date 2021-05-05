package requester

import (
	"bytes"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/vorteil/direktiv-apps/pkg/direktivapps"
)

// Manager is used to maintain the request object
type Manager struct {
	Request *Request
	client  *http.Client
	req     *http.Request
}

// Request is the struct we unmarshal the JSON input
type Request struct {
	Method   string                 `json:"method"`
	URL      string                 `json:"url"`
	Debug    bool                   `json:"debug"`
	Body     map[string]interface{} `json:"body"`
	Headers  map[string]interface{} `json:"headers"`
	Username string                 `json:"username"`
	Password string                 `json:"password"`
}

// Create initializes the http client
func (m *Manager) Create(aid string) error {
	var err error

	// Generate the body from the interface provided
	var b []byte

	if m.Request.Body != nil {
		b, err = json.Marshal(m.Request.Body)
		if err != nil {
			return err
		}
		if m.Request.Debug {
			direktivapps.Log(aid, fmt.Sprintf("Body Provided: %s", b))
		}
	}

	// Initialize client and the request
	m.client = &http.Client{}

	if m.Request.Debug {
		direktivapps.Log(aid, fmt.Sprintf("Method: %s, Sending to %s", m.Request.Method, m.Request.URL))
	}

	m.req, err = http.NewRequest(m.Request.Method, m.Request.URL, bytes.NewReader(b))
	if err != nil {
		return err
	}

	// range the header map and attach to the request if required
	for k, v := range m.Request.Headers {
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
		if m.Request.Debug {
			direktivapps.Log(aid, fmt.Sprintf("Adding %s=%s", k, actualVal))
		}
		// Adding a header requires it to be a string so might as well sprintf
		m.req.Header.Add(k, actualVal)
	}

	// Handle basic authentication.
	if m.Request.Username != "" && m.Request.Password != "" {
		sEnc := b64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", m.Request.Username, m.Request.Password)))
		m.req.Header.Add("Authorization", fmt.Sprintf("Basic %s", sEnc))
	}

	return nil
}

// Send sends the http request to provided host and responds
func (m *Manager) Send(aid string) ([]byte, error) {

	// Perform the request with the client we're using
	resp, err := m.client.Do(m.req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if m.Request.Debug {
		direktivapps.Log(aid, fmt.Sprintf("Response body: %s", body))
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// error more than likely
		return nil, fmt.Errorf("Response Message: %s, Response Code: %v \nResponseBody: %s", resp.Status, resp.StatusCode, body)
	}

	return body, nil
}
