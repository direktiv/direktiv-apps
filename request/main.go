package main

import (
	"fmt"
	"net/http"

	"github.com/vorteil/direktiv-apps/pkg/direktivapps"
	"github.com/vorteil/direktiv-apps/pkg/requester"
)

func Request(w http.ResponseWriter, r *http.Request) {

	obj := new(requester.Request)
	err := direktivapps.Unmarshal(obj, r)
	if err != nil {
		fmt.Println(err)
		return
	}

	mgr := requester.Manager{
		Request: obj,
	}

	if mgr.Request.Debug {
		fmt.Println("Requester has been initialized.")
	}

	err = mgr.Create()
	if err != nil {
		fmt.Println(err)
		return
	}

	resp, err := mgr.Send()
	if err != nil {
		fmt.Println(err)
		return
	}

	direktivapps.Respond(w, resp)
}

func main() {
	// g := direktivapps.ActionError{
	// 	ErrorCode:    "com.request.error",
	// 	ErrorMessage: "",
	// }
	// var err error

	// Gather data for running request application
	// obj := new(requester.Request)
	// direktivapps.ReadIn(obj, g)

	// mgr := requester.Manager{
	// 	Request: obj,
	// }

	// if mgr.Request.Debug {
	// 	log.Printf("Requester has been initialized")
	// }

	// err = mgr.Create()
	// if err != nil {
	// 	g.ErrorMessage = err.Error()
	// 	direktivapps.WriteError(g)
	// }

	// resp, err := mgr.Send()
	// if err != nil {
	// 	g.ErrorMessage = err.Error()
	// 	direktivapps.WriteError(g)
	// }

	// direktivapps.WriteOut(resp, g)
	direktivapps.StartServer(Request)
}
