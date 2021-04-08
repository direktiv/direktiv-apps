package main

import (
	"encoding/json"
	"net/http"

	"github.com/vorteil/direktiv-apps/pkg/direktivapps"
)

const code = "com.debug.error"

func Debug(w http.ResponseWriter, r *http.Request) {
	obj := new(map[string]interface{})
	_, err := direktivapps.Unmarshal(obj, r)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	data, err := json.Marshal(obj)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	direktivapps.Respond(w, data)
}

func main() {
	direktivapps.StartServer(Debug)
}
