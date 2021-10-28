package main

import (
	"net/http"

	"github.com/direktiv/direktiv-apps/pkg/direktivapps"
	"github.com/direktiv/direktiv-apps/pkg/requester"
)

const code = "com.discordmsg.error"

// DiscordInformation is the struct provided to the webhook request
type DiscordInformation struct {
	Message string `json:"message"`
	URL     string `json:"url"`
	TTS     bool   `json:"tts"`
}

func DiscordMsg(w http.ResponseWriter, r *http.Request) {
	obj := new(DiscordInformation)
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
				"content": obj.Message,
				"tts":     obj.TTS,
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
	direktivapps.StartServer(DiscordMsg)
}
