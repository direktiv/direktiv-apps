package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sethvargo/go-password/password"
	"github.com/direktiv/direktiv-apps/pkg/direktivapps"
)

type GenPWInput struct {
	Length           float64 `json:"length"`
	Digits           float64 `json:"digits"`
	Symbols          float64 `json:"symbols"`
	Uppercase        bool    `json:"uppercase"`
	RepeatCharacters bool    `json:"repeat-characters"`
}

type GenPWOutput struct {
	Password string `json:"password"`
}

var code = "com.genpw.error"

func GenPWHandler(w http.ResponseWriter, r *http.Request) {
	var gi GenPWInput
	aid, err := direktivapps.Unmarshal(&gi, r)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}
	direktivapps.LogDouble(aid, fmt.Sprintf("generating password, length: %v, digits: %v and symbols: %v", int(gi.Length), int(gi.Digits), int(gi.Symbols)))
	res, err := password.Generate(int(gi.Length), int(gi.Digits), int(gi.Symbols), gi.Uppercase, gi.RepeatCharacters)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	var got GenPWOutput
	got.Password = res

	data, err := json.Marshal(got)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	direktivapps.Respond(w, data)
}

func main() {
	direktivapps.StartServer(GenPWHandler)
}
