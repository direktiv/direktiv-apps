package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/direktiv/direktiv-apps/pkg/direktivapps"
	"github.com/direktiv/direktiv-apps/pkg/template"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

var code = "com.sendgrid.%s.error"

type SendGridError struct {
	Errors []struct {
		Message string      `json:"message"`
		Field   string      `json:"field"`
		Help    interface{} `json:"help"`
	} `json:"errors"`
}

// SendGrid is the object to control emailing via sendgrid
type SendGrid struct {
	SenderName  string `json:"sender-name"`
	SenderEmail string `json:"sender-email"`
	Subject     string `json:"subject"`
	RecvName    string `json:"recv-name"`
	RecvEmail   string `json:"recv-email"`
	APIKey      string `json:"apikey"`

	PlainMessage string      `json:"message"`
	HTMLMessage  string      `json:"html-message"`
	Template     string      `json:"template"`
	TemplateData interface{} `json:"template-data"`
}

func SendGridHandler(w http.ResponseWriter, r *http.Request) {

	tm := new(SendGrid)
	var err error
	aid, err := direktivapps.Unmarshal(tm, r)
	if err != nil {
		direktivapps.Log(aid, fmt.Sprintf("unmarshalling failed: %s", err.Error()))
		direktivapps.RespondWithError(w, fmt.Sprintf(code, "unmarshal"), err.Error())
		return
	}

	direktivapps.Log(aid, fmt.Sprintf("sending message to %s", tm.RecvEmail))

	from := mail.NewEmail(tm.SenderName, tm.SenderEmail)
	to := mail.NewEmail(tm.RecvName, tm.RecvEmail)

	var message *mail.SGMailV3
	if tm.Template != "" {
		msg, err := template.Render(tm.Template, tm.TemplateData)
		if err != nil {
			direktivapps.Log(aid, fmt.Sprintf("template failed: %s", err.Error()))
			direktivapps.RespondWithError(w, fmt.Sprintf(code, "template"), err.Error())
			return
		}
		message = mail.NewSingleEmail(from, tm.Subject, to, msg, msg)
	} else {
		message = mail.NewSingleEmail(from, tm.Subject, to, tm.PlainMessage, tm.HTMLMessage)
	}

	client := sendgrid.NewSendClient(tm.APIKey)
	response, err := client.Send(message)
	if err != nil {
		direktivapps.Log(aid, fmt.Sprintf("sending failed: %s", err.Error()))
		direktivapps.RespondWithError(w, fmt.Sprintf(code, "send"), err.Error())
		return
	}

	if response.StatusCode >= 300 {
		// we only handle the first error
		var e SendGridError
		err := json.Unmarshal([]byte(response.Body), &e)
		if err != nil {
			direktivapps.RespondWithError(w, fmt.Sprintf(code, "handleresponse"), err.Error())
		}
		errMsg := fmt.Sprintf("%s, field: %s", e.Errors[0].Message, e.Errors[0].Field)
		direktivapps.Log(aid, fmt.Sprintf("sending failed: %s", errMsg))
		direktivapps.RespondWithError(w, fmt.Sprintf(code, "response"), errMsg)
		return
	}

	direktivapps.Log(aid, fmt.Sprintf("sent message to %s", tm.RecvEmail))
	direktivapps.Respond(w, []byte{})

}

func main() {
	direktivapps.StartServer(SendGridHandler)
}
