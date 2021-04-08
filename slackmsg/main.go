package main

import (
	"net/http"

	"github.com/vorteil/direktiv-apps/pkg/direktivapps"
	"github.com/vorteil/direktiv-apps/pkg/requester"
)

const code = "com.slack.error"

// Information is the input struct needed to post a message to a channel
type Information struct {
	URL     string `json:"url"`
	Message string `json:"message"`
}

func SlackMsg(w http.ResponseWriter, r *http.Request) {
	obj := new(Information)
	aid, err := direktivapps.Unmarshal(obj, r)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	mgr := requester.Manager{
		Request: &requester.Request{
			Method: "POST",
			URL:    obj.URL,
			Body: map[string]interface{}{
				"text": obj.Message,
			},
			Headers: map[string]interface{}{
				"Content-Type": "application/json",
			},
		},
	}

	err = mgr.Create(aid)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	resp, err := mgr.Send(aid)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	direktivapps.Respond(w, resp)
}

func main() {
	direktivapps.StartServer(SlackMsg)
}
