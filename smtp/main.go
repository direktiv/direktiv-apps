package main

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/vorteil/direktiv-apps/pkg/direktivapps"
	gomail "gopkg.in/mail.v2"
)

var code = "com.smtp.error"

// SMTPEmail is the object to control emailing
type SMTPEmail struct {
	Debug    bool    `json:"debug"`
	From     string  `json:"from"`
	To       string  `json:"to"`
	Subject  string  `json:"subject"`
	Message  string  `json:"message"`
	Server   string  `json:"server"`
	Port     float64 `json:"port"`
	Password string  `json:"password"`
}

func SMTPEmailHandler(w http.ResponseWriter, r *http.Request) {
	tm := new(SMTPEmail)
	aid, err := direktivapps.Unmarshal(tm, r)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	if tm.Debug {
		direktivapps.Log(aid, fmt.Sprintf("Creating new message"))
		direktivapps.Log(aid, fmt.Sprintf("From: %s", tm.From))
		direktivapps.Log(aid, fmt.Sprintf("To: %s", tm.To))
		direktivapps.Log(aid, fmt.Sprintf("Subject: %s", tm.Subject))
		direktivapps.Log(aid, fmt.Sprintf("Body: %s", tm.Message))
	}

	m := gomail.NewMessage()
	// Set E-Mail sender
	m.SetHeader("From", tm.From)

	// Set E-mail receivers
	m.SetHeader("To", tm.To)

	// Set E-mail subject
	m.SetHeader("Subject", tm.Subject)

	// Set E-mail body
	m.SetBody("text/html", tm.Message)

	// Settings for SMTP server
	d := gomail.NewDialer(tm.Server, int(tm.Port), tm.From, tm.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	direktivapps.Respond(w, []byte{})

}

func main() {
	direktivapps.StartServer(SMTPEmailHandler)
}
