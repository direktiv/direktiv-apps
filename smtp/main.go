package main

import (
	"bytes"
	"crypto/tls"
	b64 "encoding/base64"
	"fmt"
	"html/template"
	"net/http"

	"github.com/vorteil/direktiv-apps/pkg/direktivapps"
	gomail "gopkg.in/mail.v2"
)

var code = "com.smtp.error"

// SMTPEmail is the object to control emailing
type SMTPEmail struct {
	From     string                 `json:"from"`
	To       string                 `json:"to"`
	Subject  string                 `json:"subject"`
	Message  string                 `json:"message"`
	Server   string                 `json:"server"`
	Port     float64                `json:"port"`
	Password string                 `json:"password"`
	Args     map[string]interface{} `json:"args"` // optional
}

func SMTPEmailHandler(w http.ResponseWriter, r *http.Request) {
	tm := new(SMTPEmail)
	var err error
	aid, err := direktivapps.Unmarshal(tm, r)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	direktivapps.Log(aid, "Decoding Message Input")

	var message []byte
	message, err = b64.StdEncoding.DecodeString(tm.Message)
	if err != nil {
		direktivapps.Log(aid, "Message Input is not a base64 string defaulting to normal string content.")
		message = []byte(tm.Message)
	}

	direktivapps.Log(aid, "Creating Template")

	t := template.New("email.html")
	t, err = t.Parse(fmt.Sprintf("%s", message))
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, tm.Args); err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	direktivapps.Log(aid, "Creating Message")

	m := gomail.NewMessage()
	// Set E-Mail sender
	m.SetHeader("From", tm.From)

	// Set E-mail receivers
	m.SetHeader("To", tm.To)

	// Set E-mail subject
	m.SetHeader("Subject", tm.Subject)

	// Set E-mail body
	m.SetBody("text/html", tpl.String())

	// Settings for SMTP server
	d := gomail.NewDialer(tm.Server, int(tm.Port), tm.From, tm.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	direktivapps.Log(aid, "Sending Message")
	if err := d.DialAndSend(m); err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	direktivapps.Respond(w, []byte{})
}

func main() {
	direktivapps.StartServer(SMTPEmailHandler)
}
