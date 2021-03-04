package main

import (
	"github.com/vorteil/direktiv-apps/pkg/direktivapps"
	"github.com/vorteil/direktiv-apps/pkg/requester"
)

// GoogleChatInfo is the struct
type GoogleChatInfo struct {
	Message string `json:"message"`
	URL     string `json:"url"`
}

func main() {
	g := direktivapps.ActionError{
		ErrorCode:    "com.googlechat.error",
		ErrorMessage: "",
	}
	var err error

	obj := new(GoogleChatInfo)
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

	err = mgr.Create()
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
