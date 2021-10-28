package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/direktiv/direktiv-apps/pkg/direktivapps"
)

type Greeter struct {
	Name string `json:"name"`
}

type ReturnGreeting struct {
	Greeting string `json:"greeting"`
}

const code = "com.greeting.error"

func GreetingHandler(w http.ResponseWriter, r *http.Request) {
	obj := new(Greeter)
	_, err := direktivapps.Unmarshal(obj, r)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}
	var rg ReturnGreeting
	rg.Greeting = fmt.Sprintf("Welcome to Direktiv, %s!", obj.Name)

	bv, err := json.Marshal(rg)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	direktivapps.Respond(w, bv)
}

func main() {
	direktivapps.StartServer(GreetingHandler)
}
