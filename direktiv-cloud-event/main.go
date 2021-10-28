package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/direktiv/direktiv-apps/pkg/direktivapps"
	"github.com/direktiv/direktiv-apps/pkg/requester"
)

type CloudEventSend struct {
	Type        string                 `json:"type"`         // the type of event sent
	Data        map[string]interface{} `json:"body"`         // the data of the cloud event
	Source      string                 `json:"source"`       // where did this event come from?
	AccessToken string                 `json:"access_token"` // the access token from direktiv with correct permissions
	Namespace   string                 `json:"namespace"`    // the namespace to send the request to
}

const DirektivURL = "https://playground.direktiv.io/api/namespaces/%s/event"
const code = "com.cloud-event.%s.error"

func CloudEventHandler(w http.ResponseWriter, r *http.Request) {

	ces := new(CloudEventSend)
	aid, err := direktivapps.Unmarshal(ces, r)
	if err != nil {
		direktivapps.RespondWithError(w, fmt.Sprintf(code, "unmarshal"), err.Error())
		return
	}

	event := cloudevents.NewEvent()
	event.SetSource(ces.Source)
	event.SetType(ces.Type)
	event.SetData(cloudevents.ApplicationJSON, ces.Data)

	bytes, err := json.Marshal(event)
	if err != nil {
		direktivapps.RespondWithError(w, fmt.Sprintf(code, "marshal-event"), err.Error())
		return
	}

	var newEvent map[string]interface{}
	err = json.Unmarshal(bytes, &newEvent)
	if err != nil {
		direktivapps.RespondWithError(w, fmt.Sprintf(code, "unmarshal-event"), err.Error())
		return
	}

	mgr := requester.Manager{
		Request: &requester.Request{
			Method: "POST",
			URL:    fmt.Sprintf(DirektivURL, ces.Namespace),
			Body:   newEvent,
			Headers: map[string]interface{}{
				"Content-Type":  "application/json",
				"Authorization": fmt.Sprintf("Bearer %s", ces.AccessToken),
			},
		},
	}

	err = mgr.Create(aid)
	if err != nil {
		direktivapps.RespondWithError(w, fmt.Sprintf(code, "request-create"), err.Error())
		return
	}

	resp, err := mgr.Send(aid)
	if err != nil {
		direktivapps.RespondWithError(w, fmt.Sprintf(code, "request-sent"), err.Error())
		return
	}

	direktivapps.Respond(w, resp)
}

func main() {
	direktivapps.StartServer(CloudEventHandler)
}
