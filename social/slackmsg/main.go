package main

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/direktiv/direktiv-apps/pkg/reusable"
)

const code = "com.slack.error"

type requestInput struct {
	URL     string          `json:"url"`
	Message json.RawMessage `json:"message"`
}

func slackHandler(w http.ResponseWriter, r *http.Request, ri *reusable.RequestInfo) {

	obj := new(requestInput)
	err := reusable.Unmarshal(obj, true, r)
	if err != nil {
		reusable.ReportError(w, reusable.UnmarshallError, err)
		return
	}

	ri.Logger().Infof("sending slack message")
	_, err = http.Post(obj.URL, "application/json", bytes.NewBuffer(obj.Message))
	if err != nil {
		reusable.ReportError(w, code, err)
		return
	}

	reusable.ReportResult(w, "")

}

func main() {
	reusable.StartServer(slackHandler, nil)
}
