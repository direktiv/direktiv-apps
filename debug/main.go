package main

import (
	"encoding/json"

	"github.com/vorteil/direktiv-apps/pkg/direktivapps"
)

func main() {

	g := direktivapps.ActionError{
		ErrorCode:    "com.debug.error",
		ErrorMessage: "",
	}

	obj := new(map[string]interface{})
	direktivapps.ReadIn(obj, g)

	data, err := json.Marshal(obj)
	if err != nil {
		g.ErrorMessage = err.Error()
		direktivapps.WriteError(g)
	}

	direktivapps.WriteOut(data, g)
}
