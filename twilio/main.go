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
	"strings"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type ActionError struct {
	ErrorCode    string `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
}

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

func main() {
	var response []byte

	tm := TwilioMessage{}
	g := ActionError{
		ErrorCode:    "com.twilio.error",
		ErrorMessage: "",
	}

	// read data in
	data, err := ioutil.ReadFile("/direktiv-data/data.in")
	if err != nil {
		g.ErrorMessage = err.Error()
		writeError(g)
		return
	}

	err = json.Unmarshal(data, &tm)
	if err != nil {
		g.ErrorMessage = err.Error()
		writeError(g)
		return
	}

	switch tm.TypeOf {
	case "email":
		response, err = SendEmail(&tm)
		if err != nil {
			g.ErrorMessage = err.Error()
			writeError(g)
			return
		}
	case "sms":
		response, err = SendSMS(&tm)
		if err != nil {
			g.ErrorMessage = err.Error()
			writeError(g)
			return
		}
	default:
		g.ErrorMessage = fmt.Errorf("'%s' is not a valid type to use the twilio application", tm.TypeOf).Error()
		writeError(g)
		return
	}

	finishRunning(response)
}

// SendEmail sends a message to the provided email from the input json
func SendEmail(tm *TwilioMessage) ([]byte, error) {

	from := mail.NewEmail("", tm.From)
	subject := tm.Subject
	to := mail.NewEmail("", tm.To)

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
		return nil, err
	}
	defer resp.Body.Close()

	br, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		// error more than likely
		return nil, fmt.Errorf("Response Message: %s, Response Code: %v \nResponseBody: %s", resp.Status, resp.StatusCode, body)
	}

	return br, nil
}

// SendSMS sends a sms to the provided mobile number from the input json
func SendSMS(tm *TwilioMessage) ([]byte, error) {
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

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		// error more than likely
		return nil, fmt.Errorf("Response Message: %s, Response Code: %v \nResponseBody: %s", resp.Status, resp.StatusCode, body)
	}
	return body, nil
}

// finishRunning will write to a file and or print the json body to stdout and exits
func finishRunning(eb []byte) {
	var err error
	err = ioutil.WriteFile("/direktiv-data/data.out", eb, 0755)
	if err != nil {
		log.Fatal("can not write out data")
		return
	}
}

// writeError
func writeError(g ActionError) {
	b, _ := json.Marshal(g)
	err := ioutil.WriteFile("/direktiv-data/error.json", b, 0755)
	if err != nil {
		log.Fatal("can not write json error")
		return
	}
}
