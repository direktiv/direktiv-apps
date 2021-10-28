package main

import (
	"fmt"
	"net/http"

	"github.com/direktiv/direktiv-apps/pkg/direktivapps"
	"github.com/direktiv/direktiv-apps/pkg/requester"
)

const funcURL = `https://%s.azurewebsites.net/api/%s?code=%s`
const code = "com.azureinvoke.error"

// AzureFuncTriggerDetails ...
type AzureFuncTriggerDetails struct {
	FunctionApp string                 `json:"function-app"`
	Function    string                 `json:"function-name"`
	Key         string                 `json:"function-key"`
	Body        map[string]interface{} `json:"body"`
}

func main() {
	direktivapps.StartServer(AzureInvoke)
}

func AzureInvoke(w http.ResponseWriter, r *http.Request) {
	var err error

	obj := new(AzureFuncTriggerDetails)
	aid, err := direktivapps.Unmarshal(obj, r)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	mgr := requester.Manager{
		Request: &requester.Request{
			Method: "POST",
			URL:    fmt.Sprintf(funcURL, obj.FunctionApp, obj.Function, obj.Key),
			Body:   obj.Body,
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
