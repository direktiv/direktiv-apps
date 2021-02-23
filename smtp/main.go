package main

import (
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"log"

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
	tm := &SMTPEmail{}
	g := ActionError{
		ErrorCode:    "com.smtp.error",
		ErrorMessage: "",
	}
	var err error
	var data []byte
	// read data in
	log.Printf("Reading in Data...")
	data, err = ioutil.ReadFile("/direktiv-data/data.in")
	if err != nil {
		g.ErrorMessage = err.Error()
		writeError(g)
		return
	}

	err = json.Unmarshal(data, tm)
	if err != nil {
		g.ErrorMessage = err.Error()
		writeError(g)
		return
	}

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
		writeError(g)
		return
	}

	// Nothing can be returned
	finishRunning([]byte{}, g)
}

// writeError
func writeError(g ActionError) {
	b, _ := json.Marshal(g)
	err := ioutil.WriteFile("/direktiv-data/error.json", b, 0755)
}

// finishRunning will write to a file and or print the json body to stdout and exits
func finishRunning(eb []byte, g ActionError) {
	var err error
	err = ioutil.WriteFile("/direktiv-data/data.out", eb, 0755)
	if err != nil {
		g.ErrorMessage = err.Error()
		writeError(g)
		return
	}
}
