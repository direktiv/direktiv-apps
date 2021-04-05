package main

import (
	"github.com/vorteil/direktiv-apps/pkg/direktivapps"
	"github.com/vorteil/direktiv-apps/pkg/requester"
)

// DiscordInformation is the struct provided to the webhook request
type DiscordInformation struct {
	Message string `json:"message"`
	URL     string `json:"url"`
	TTS     bool   `json:"tts"`
}

func main() {

	g := direktivapps.ActionError{
		ErrorCode:    "com.discord.error",
		ErrorMessage: "",
	}

	obj := new(DiscordInformation)
	direktivapps.ReadIn(obj, g)

	mgr := requester.Manager{
		Request: &requester.Request{
			Method: "POST",
			URL:    obj.URL,
			Body: map[string]interface{}{
				"content": obj.Message,
				"tts":     obj.TTS,
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
