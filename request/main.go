package main

import (
	"log"

	"github.com/vorteil/direktiv-apps/pkg/direktivapps"
	"github.com/vorteil/direktiv-apps/pkg/requester"
)

func main() {
	g := direktivapps.ActionError{
		ErrorCode:    "com.request.error",
		ErrorMessage: "",
	}
	var err error

	// Gather data for running request application
	obj := new(requester.Request)
	direktivapps.ReadIn(obj, g)

	mgr := requester.Manager{
		Request: obj,
	}

	if mgr.Request.Debug {
		log.Printf("Requester has been initialized")
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
