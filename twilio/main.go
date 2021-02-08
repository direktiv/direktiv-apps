package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// TwilioMessage input struct to send an sms or email
type TwilioMessage struct {
	TypeOf      string `json:"typeof"` // Email, sms
	Sid         string `json:"sid"`
	Token       string `json:"token"`
	Subject     string `json:"subject"`     // subject header of email
	Message     string `json:"message"`     // contents of email/sms
	HTMLMessage string `json:"htmlMessage"` // contents if you want to display in html
	From        string `json:"from"`        // who we sending from
	To          string `json:"to"`          // who we sending to
}

// EndBody is the response of this library after a request
type EndBody struct {
	Error         string `json:"error"`
	Response      string `json:"response"`
	Status        int    `json:"statusCode"`
	StatusMessage string `json:"status"`
}

func main() {
	// if len(os.Args) > 2 {
	tm := &TwilioMessage{}
	eb := &EndBody{}

	// inputFile := os.Args[1]
	// outputFile := os.Args[2]

	// read data in
	data, err := ioutil.ReadFile("/direktiv-data/data.in")
	if err != nil {
		eb.Error = err.Error()
		finishRunning(eb)
		return
	}

	log.Printf("in data: %s", string(data))

	err = json.Unmarshal(data, &tm)
	if err != nil {
		eb.Error = err.Error()
		finishRunning(eb)
	}

	switch tm.TypeOf {
	case "email":
		eb, err = SendEmail(tm)
		if err != nil {
			fmt.Println(err, eb)
			eb.Error = err.Error()
			finishRunning(eb)

		}
	case "sms":
		eb, err = SendSMS(tm)
		if err != nil {
			eb.Error = err.Error()
			finishRunning(eb)

		}
	}
	finishRunning(eb)
	// }
}

// SendEmail sends a message to the provided email from the input json
func SendEmail(tm *TwilioMessage) (*EndBody, error) {
	eb := &EndBody{}

	from := mail.NewEmail("", tm.From)
	subject := tm.Subject
	to := mail.NewEmail("", tm.To)

	// from, subject header, send to, content, htmlContent
	message := mail.NewSingleEmail(from, subject, to, tm.Message, tm.HTMLMessage)
	b := bytes.NewReader(mail.GetRequestBody(message))
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	req, _ := http.NewRequest("POST", "https://api.sendgrid.com/v3/mail/send", b)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tm.Token))
	req.Header.Add("User-Agent", fmt.Sprintf("sendgrid/%s;go", sendgrid.Version))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	// client := sendgrid.NewSendClient(tm.Token)
	resp, err := client.Do(req)
	if err != nil {
		return eb, err
	}
	defer resp.Body.Close()

	br, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return eb, err
	}

	eb.Status = resp.StatusCode
	eb.Response = string(br)
	eb.Error = ""

	return eb, nil
}

// SendSMS sends a sms to the provided mobile number from the input json
func SendSMS(tm *TwilioMessage) (*EndBody, error) {
	eb := &EndBody{}
	urlStr := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", tm.Sid)
	msgData := url.Values{}
	msgData.Set("To", tm.To)
	msgData.Set("From", tm.From)
	msgData.Set("Body", tm.Message)
	msgDataReader := *strings.NewReader(msgData.Encode())

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	req, _ := http.NewRequest("POST", urlStr, &msgDataReader)
	req.SetBasicAuth(tm.Sid, tm.Token)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	eb.Status = resp.StatusCode
	eb.StatusMessage = resp.Status
	eb.Response = string(body)
	eb.Error = ""

	return eb, nil
}

// finishRunning will write to a file and or print the json body to stdout and exits
func finishRunning(eb *EndBody) {
	ms, _ := json.Marshal(eb)
	_ = ioutil.WriteFile("/direktiv-data/data.out", []byte(ms), 0755)
	os.Exit(0)
}
