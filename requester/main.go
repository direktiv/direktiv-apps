package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

func main() {

	eb := &EndBody{}
	// Arguments provided are more than the program
	if len(os.Args) > 2 {
		// use first arg for file name
		filepathArg := os.Args[1]
		outFileArg := os.Args[2]

		requester := Manager{}
		err := requester.Initialize(filepathArg)
		if err != nil {
			eb.Error = err.Error()
			finishRunning(outFileArg, eb)
		}

		err = requester.Create()
		if err != nil {
			eb.Error = err.Error()
			finishRunning(outFileArg, eb)
		}

		endBody, err := requester.Send()
		if err != nil {
			eb.Error = err.Error()
			finishRunning(outFileArg, eb)
		}

		finishRunning(outFileArg, endBody)
	} else {
		eb.Error = "No filepath provided as argument to read json"
		fmt.Printf("%+v\n", eb)
		return
	}

}

// finishRunning will write to a file and or print the json body to stdout and exits
func finishRunning(path string, eb *EndBody) {
	ms, _ := json.Marshal(eb)
	_ = ioutil.WriteFile(path, ms, 0644)
	os.Exit(0)
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

// EndBody is the response of this library after a request
type EndBody struct {
	Error         string `json:"error"`
	Response      string `json:"response"`
	Status        int    `json:"statusCode"`
	StatusMessage string `json:"status"`
}

// Initialize reads the file unmarshal json into appropriate struct
func (m *Manager) Initialize(path string) error {

	// Open file and read its contents attempt to unmarshal it from json
	jsonFile, err := os.Open(path)
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	bv, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bv, &m.Request)
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
	m.client = &http.Client{}
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
/*
	{
		error: "",
		response: "",
		status: ""
	}
*/
func (m *Manager) Send() (*EndBody, error) {
	eb := &EndBody{}

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

	eb.Status = resp.StatusCode
	eb.StatusMessage = resp.Status
	eb.Response = string(body)
	eb.Error = ""

	return eb, nil
}
