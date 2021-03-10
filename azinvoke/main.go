package main

import (
	"fmt"

	"github.com/vorteil/direktiv-apps/pkg/direktivapps"
	"github.com/vorteil/direktiv-apps/pkg/requester"
)

const funcURL = `https://%s.azurewebsites.net/api/%s?code=%s`

// AzureFuncTriggerDetails ...
type AzureFuncTriggerDetails struct {
	FunctionApp string                 `json:"function-app"`
	Function    string                 `json:"function-name"`
	Key         string                 `json:"function-key"`
	Body        map[string]interface{} `json:"body"`
}

func main() {
	g := direktivapps.ActionError{
		ErrorCode:    "com.azureinvoke.error",
		ErrorMessage: "",
	}
	var err error

	obj := new(AzureFuncTriggerDetails)
	direktivapps.ReadIn(obj, g)

	mgr := requester.Manager{
		Request: &requester.Request{
			Method: "POST",
			URL:    fmt.Sprintf(funcURL, obj.FunctionApp, obj.Function, obj.Key),
			Body:   obj.Body,
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
