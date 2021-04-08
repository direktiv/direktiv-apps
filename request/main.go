package main

import (
	"net/http"

	"github.com/vorteil/direktiv-apps/pkg/direktivapps"
	"github.com/vorteil/direktiv-apps/pkg/requester"
)

const code = "com.request.error"

func Request(w http.ResponseWriter, r *http.Request) {

	obj := new(requester.Request)
	aid, err := direktivapps.Unmarshal(obj, r)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	mgr := requester.Manager{
		Request: obj,
	}

	if mgr.Request.Debug {
		direktivapps.Log(aid, "Requester has been initialized")
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
	direktivapps.StartServer(Request)
}
