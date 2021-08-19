package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	regex "regexp"

	"github.com/vorteil/direktiv-apps/pkg/direktivapps"
)

type request struct {
	Msg   string `json:"msg"`
	Regex string `json:"regex"`
}

type output struct {
	Results []string `json:"results"`
}

var code = "com.regex.error"

func RegexHandler(w http.ResponseWriter, r *http.Request) {
	var o request

	_, err := direktivapps.Unmarshal(&o, r)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	re, err := regex.Compile(fmt.Sprintf(`%s`, o.Regex))
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	var ot output
	ot.Results = re.FindAllString(o.Msg, -1)

	data, err := json.Marshal(&ot)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	direktivapps.Respond(w, data)
}

func main() {
	direktivapps.StartServer(RegexHandler)
}
