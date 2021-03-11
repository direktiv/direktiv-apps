package main

import (
	"encoding/json"
	"fmt"

	"github.com/vorteil/direktiv-apps/pkg/direktivapps"
)

type Greeter struct {
	Name string `json:"name"`
}

type ReturnGreeting struct {
	Greeting string `json:"greeting"`
}

func main() {
	g := direktivapps.ActionError{
		ErrorCode:    "com.greeting.error",
		ErrorMessage: "",
	}

	var err error
	obj := new(Greeter)
	direktivapps.ReadIn(obj, g)

	var rg ReturnGreeting
	rg.Greeting = fmt.Sprintf("Welcome to Direktiv, %s!", obj.Name)

	bv, err := json.Marshal(rg)
	if err != nil {
		g.ErrorMessage = err.Error()
		direktivapps.WriteError(g)
	}

	direktivapps.WriteOut(bv, g)
}
