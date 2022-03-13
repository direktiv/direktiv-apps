package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/direktiv/direktiv-apps/pkg/reusable"
	"github.com/google/uuid"
)

type keptnPayload struct {
	Project  string            `json:"project"`
	Service  string            `json:"service"`
	Stage    string            `json:"stage"`
	Sequence string            `json:"sequence,omitempty"`
	Lables   map[string]string `json:"labels"`
}

type requestInput struct {
	Keptn string `json:"keptn"`
	Token string `json:"token"`

	TLS        bool `json:"tls"`
	SkipVerify bool `json:"skip-verify"`

	Payload keptnPayload `json:"data"`
}

const (
	keptnAPIPath  = "/api/v1/event"
	direktivLabel = "triggered_by"
	direktivValue = "direktiv"

	keptnContext = "keptnContext"
)

func keptnHandler(w http.ResponseWriter, r *http.Request, ri *reusable.RequestInfo) {

	obj := new(requestInput)
	err := reusable.Unmarshal(obj, true, r)
	if err != nil {
		reusable.ReportError(w, reusable.UnmarshallError, err)
		return
	}

	ri.Logger().Infof("go keptn request for %s", obj.Keptn)

	scheme := "http"
	if obj.TLS {
		scheme = "https"
	}

	// generate URL
	url := url.URL{
		Scheme: scheme,
		Host:   obj.Keptn,
		Path:   keptnAPIPath,
	}

	je, err := json.Marshal(generateEvent(obj))
	if err != nil {
		reusable.ReportError(w, "com.keptn.error", err)
		return
	}

	req, err := http.NewRequest("POST", url.String(),
		bytes.NewBuffer(je))
	if err != nil {
		reusable.ReportError(w, "com.keptn.error", err)
		return
	}

	req.Header.Add("x-token", obj.Token)
	req.Header.Add("Accept", "application/cloudevents+json")
	req.Header.Add("Content-Type", "application/cloudevents+json")

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: obj.SkipVerify},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	if err != nil {
		reusable.ReportError(w, "com.keptn.error", err)
		return
	}

	rm, err := handleResponse(resp.Body)
	if err != nil {
		reusable.ReportError(w, "com.keptn.error", err)
		return
	}

	reusable.ReportResult(w, rm)

}

func handleResponse(r io.ReadCloser) (map[string]interface{}, error) {

	rm := make(map[string]interface{})

	b, err := io.ReadAll(r)
	if err != nil {
		return rm, err
	}
	defer r.Close()

	err = json.Unmarshal(b, &rm)
	if err != nil {
		return rm, err
	}

	if _, ok := rm[keptnContext]; !ok {
		return nil, fmt.Errorf("error pushing keptn request: %v", rm)
	}

	return rm, nil

}

func generateEvent(obj *requestInput) event.Event {

	event := cloudevents.NewEvent()
	event.SetSource("direktiv")
	event.SetID(uuid.New().String())
	event.SetType(fmt.Sprintf("sh.keptn.event.%s.%s.triggered", obj.Payload.Stage, obj.Payload.Sequence))

	// wipe the sequence so it does snot show up in json payload for keptn
	obj.Payload.Sequence = ""

	// add direktiv label
	if obj.Payload.Lables == nil {
		obj.Payload.Lables = make(map[string]string)
	}
	obj.Payload.Lables[direktivLabel] = direktivValue

	event.SetData(cloudevents.ApplicationJSON, obj.Payload)

	return event

}

func main() {
	reusable.StartServer(keptnHandler, nil)
}
