package main

import (
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	gomail "gopkg.in/mail.v2"
)

// SMTPEmail is the object to control emailing
type SMTPEmail struct {
	From     string  `json:"from"`
	To       string  `json:"to"`
	Subject  string  `json:"subject"`
	Message  string  `json:"message"`
	Server   string  `json:"server"`
	Port     float64 `json:"port"`
	Password string  `json:"password"`
}

// EndBody is the response of this library after a request
type EndBody struct {
	Error         string `json:"error"`
	Response      string `json:"response"`
	Status        int    `json:"statusCode"`
	StatusMessage string `json:"status"`
}

func main() {
	tm := &SMTPEmail{}
	eb := &EndBody{}

	// read data in
	data, err := ioutil.ReadFile("/direktiv-data/data.in")
	if err != nil {
		eb.Error = err.Error()
		finishRunning(eb)
		return
	}

	err = json.Unmarshal(data, tm)
	if err != nil {
		eb.Error = err.Error()
		finishRunning(eb)
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
		eb.Error = err.Error()
	}
	finishRunning(eb)
}

// finishRunning will write to a file and or print the json body to stdout and exits
func finishRunning(eb *EndBody) {
	ms, _ := json.Marshal(eb)
	log.Printf("EB: %+v", eb)
	_ = ioutil.WriteFile("/direktiv-data/data.out", []byte(ms), 0755)
	os.Exit(0)
}
