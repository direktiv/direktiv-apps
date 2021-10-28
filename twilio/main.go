package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/direktiv/direktiv-apps/pkg/direktivapps"
)

// TwilioMessage input struct to send an sms or email
type TwilioMessage struct {
	Debug       bool   `json:""debug`
	TypeOf      string `json:"typeof"` // Email, sms
	Sid         string `json:"sid"`
	Token       string `json:"token"`
	Subject     string `json:"subject"`     // subject header of email
	Message     string `json:"message"`     // contents of email/sms
	HTMLMessage string `json:"htmlMessage"` // contents if you want to display in html
	From        string `json:"from"`        // who we sending from
	To          string `json:"to"`          // who we sending to
}

var code = "com.twilio.error"

func TwilioMessageHandler(w http.ResponseWriter, r *http.Request) {
	tm := new(TwilioMessage)
	aid, err := direktivapps.Unmarshal(tm, r)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	var response []byte

	switch tm.TypeOf {
	case "email":
		if tm.Debug {
			direktivapps.Log(aid, "Sending Email")
		}
		response, err = SendEmail(tm, aid)
		if err != nil {
			direktivapps.RespondWithError(w, code, err.Error())
			return
		}
	case "sms":
		if tm.Debug {
			direktivapps.Log(aid, "Sending SMS")
		}
		response, err = SendSMS(tm, aid)
		if err != nil {
			direktivapps.RespondWithError(w, code, err.Error())
			return
		}
	default:
		errs := fmt.Errorf("'%s' is not a valid type to use the twilio application", tm.TypeOf).Error()
		direktivapps.RespondWithError(w, code, errs)
		return
	}

	direktivapps.Respond(w, response)
}

func main() {
	direktivapps.StartServer(TwilioMessageHandler)
}

// SendEmail sends a message to the provided email from the input json
func SendEmail(tm *TwilioMessage, aid string) ([]byte, error) {

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

	if tm.Debug {
		direktivapps.Log(aid, fmt.Sprintf("Send to %s\n From %s\n  Body %s", tm.To, tm.From, tm.Message))
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

	if tm.Debug {
		direktivapps.Log(aid, fmt.Sprintf("Response Body: %s", br))
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// error more than likely
		return nil, fmt.Errorf("Response Message: %s, Response Code: %v \nResponseBody: %s", resp.Status, resp.StatusCode, br)
	}

	return br, nil
}

// SendSMS sends a sms to the provided mobile number from the input json
func SendSMS(tm *TwilioMessage, aid string) ([]byte, error) {
	urlStr := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", tm.Sid)
	msgData := url.Values{}
	msgData.Set("To", tm.To)
	msgData.Set("From", tm.From)
	msgData.Set("Body", tm.Message)

	if tm.Debug {
		direktivapps.Log(aid, fmt.Sprintf("Send to %s\nFrom %s\n Body: %s", tm.To, tm.From, tm.Message))
	}

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

	if tm.Debug {
		log.Printf("Response Body: %s", body)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// error more than likely
		return nil, fmt.Errorf("Response Message: %s, Response Code: %v \nResponseBody: %s", resp.Status, resp.StatusCode, body)
	}
	return body, nil
}
