package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type ActionError struct {
	ErrorCode    string `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
}

func main() {
	g := ActionError{
		ErrorCode:    "com.request.error",
		ErrorMessage: "",
	}
	var err error
	var data []byte

	// read data in
	data, err = ioutil.ReadFile("/direktiv-data/data.in")
	if err != nil {
		g.ErrorMessage = err.Error()
		writeError(g)
		return
	}

	requester := Manager{}
	err = requester.Initialize(data)
	if err != nil {
		g.ErrorMessage = err.Error()
		writeError(g)
		return
	}

	err = requester.Create()
	if err != nil {
		g.ErrorMessage = err.Error()
		writeError(g)
		return
	}

	resp, err := requester.Send()
	if err != nil {
		g.ErrorMessage = err.Error()
		writeError(g)
		return
	}

	finishRunning(resp)

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

// finishRunning will write to a file and or print the json body to stdout and exits
func finishRunning(eb []byte) {
	var err error
	err = ioutil.WriteFile("/direktiv-data/data.out", eb, 0755)
	if err != nil {
		log.Fatal("can not write out data")
		return
	}
}

// Manager is used to maintain the request object
type Manager struct {
	Request Request
	client  *http.Client
	req     *http.Request
}

// Request is the struct we unmarshal the JSON input
type Request struct {
	Method  string                 `json:"method"`
	Host    string                 `json:"host"`
	Body    map[string]interface{} `json:"body"`
	Headers map[string]interface{} `json:"headers"`
}

// Initialize reads the file unmarshal json into appropriate struct
func (m *Manager) Initialize(bv []byte) error {
	err := json.Unmarshal(bv, &m.Request)
	if err != nil {
		return err
	}

	return nil
}

// Create initializes the http client
func (m *Manager) Create() error {
	var err error

	// Generate the body from the interface provided
	bvMap, err := json.Marshal(m.Request.Body)
	if err != nil {
		return err
	}
	b := bytes.NewBuffer([]byte(bvMap))

	// Initialize client and the request
	m.client = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	m.req, err = http.NewRequest(m.Request.Method, m.Request.Host, b)
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
		// Adding a header requires it to be a string so might as well sprintf
		m.req.Header.Add(k, actualVal)
	}

	return nil
}

// Send sends the http request to provided host and responds
func (m *Manager) Send() ([]byte, error) {

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

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		// error more than likely
		return nil, fmt.Errorf("Response Message: %s, Response Code: %v \nResponseBody: %s", resp.Status, resp.StatusCode, body)
	}

	return body, nil
}
