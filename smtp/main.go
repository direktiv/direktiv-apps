package main

import (
	"encoding/json"
	"io/ioutil"
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
	if len(os.Args) > 2 {
		tm := &SMTPEmail{}
		eb := &EndBody{}

		inputFile := os.Args[1]
		outputFile := os.Args[2]

		// Read input file and unmarshal into struct
		b, err := ioutil.ReadFile(inputFile)
		if err != nil {
			eb.Error = err.Error()
			finishRunning(outputFile, eb)
		}

		err = json.Unmarshal(b, tm)
		if err != nil {
			eb.Error = err.Error()
			finishRunning(outputFile, eb)
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

		if err := d.DialAndSend(m); err != nil {
			eb.Error = err.Error()
		}
		finishRunning(outputFile, eb)
	}
}

// finishRunning will write to a file and or print the json body to stdout and exits
func finishRunning(path string, eb *EndBody) {
	ms, _ := json.Marshal(eb)
	_ = ioutil.WriteFile(path, ms, 0644)
	os.Exit(0)
}
