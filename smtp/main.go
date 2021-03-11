package main

import (
	"crypto/tls"
	"log"

	"github.com/vorteil/direktiv-apps/pkg/direktivapps"
	gomail "gopkg.in/mail.v2"
)

type ActionError struct {
	ErrorCode    string `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
}

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

func main() {
	g := direktivapps.ActionError{
		ErrorCode:    "com.smtp.error",
		ErrorMessage: "",
	}

	tm := new(SMTPEmail)
	direktivapps.ReadIn(tm, g)

	if tm.Debug {
		log.Printf("Creating new message")
		log.Printf("From: %s", tm.From)
		log.Printf("To: %s", tm.To)
		log.Printf("Subject: %s", tm.Subject)
		log.Printf("Body: %s", tm.Message)
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
		g.ErrorMessage = err.Error()
		direktivapps.WriteError(g)
	}

	// Nothing can be returned
	direktivapps.WriteOut([]byte{}, g)
}
