package main

import (
	"github.com/vorteil/direktiv-apps/pkg/direktivapps"
	"github.com/vorteil/direktiv-apps/pkg/requester"
)

// Information is the input struct needed to post a message to a channel
type Information struct {
	URL     string `json:"url"`
	Message string `json:"message"`
}

func main() {

	g := direktivapps.ActionError{
		ErrorCode:    "com.slack.error",
		ErrorMessage: "",
	}

	obj := new(Information)
	direktivapps.ReadIn(obj, g)

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

	err := mgr.Create()
	if err != nil {
		g.ErrorMessage = err.Error()
		direktivapps.WriteError(g)
	}

	resp, err := mgr.Send()
	if err != nil {
		g.ErrorMessage = err.Error()
		direktivapps.WriteError(g)
	}

	direktivapps.WriteOut(resp, g)
}
