package main

import (
	"bytes"
	"crypto/tls"
	b64 "encoding/base64"
	"fmt"
	"html/template"
	"net/http"
	"net/smtp"
	"strings"
	"path/filepath"

	"github.com/vorteil/direktiv-apps/pkg/direktivapps"
	gomail "gopkg.in/mail.v2"
)

var code = "com.smtp.%s.error"

// SMTPEmail is the object to control emailing
type SMTPEmail struct {
	From     string                 `json:"from"`
	To       []string               `json:"to"`
	Subject  string                 `json:"subject"`
	Message  string                 `json:"message"`
	Base64   bool                   `json:"template"`
	Server   string                 `json:"server"`
	Port     float64                `json:"port"`
	Password string                 `json:"password"`
	Args     map[string]interface{} `json:"args"` // optional
	Images []string `json:"images"`
}

func SMTPEmailHandler(w http.ResponseWriter, r *http.Request) {
	tm := new(SMTPEmail)
	var err error
	aid, err := direktivapps.Unmarshal(tm, r)
	if err != nil {
		direktivapps.RespondWithError(w, fmt.Sprintf(code, "unmarshal"), err.Error())
		return
	}

	direktivapps.Log(aid, "Decoding Message Input")

	var message []byte
	if tm.Base64 {
		direktivapps.Log(aid, "decoding from base 64 string")
		message, err = b64.StdEncoding.DecodeString(tm.Message)
		if err != nil {
			message = []byte(tm.Message)
		}
	} else {
		direktivapps.Log(aid, "using message as plain string")
		message = []byte(tm.Message)
	}

	direktivapps.Log(aid, "Creating Template")

	t := template.New("email.html")
	t, err = t.Parse(fmt.Sprintf("%s", message))
	if err != nil {
		direktivapps.RespondWithError(w, fmt.Sprintf(code, "parse-template"), err.Error())
		return
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, tm.Args); err != nil {
		direktivapps.RespondWithError(w, fmt.Sprintf(code, "template-execute"), err.Error())
		return
	}

	direktivapps.Log(aid, "Creating Message")

	m := gomail.NewMessage()
	// Set E-Mail sender
	m.SetHeader("From", tm.From)

	// Set E-mail receivers
	m.SetHeader("To", tm.To...)

	// Set E-mail subject
	m.SetHeader("Subject", tm.Subject)

	for _, imageBody := range tm.Images {
		direktivapps.Log(aid, fmt.Sprintf("%s",filepath.Join(r.Header.Get("Direktiv-TempDir"), imageBody)))
		m.Embed(filepath.Join(r.Header.Get("Direktiv-TempDir"), imageBody))
	}

	// Set E-mail body
	m.SetBody("text/html", tpl.String())

	// if not auth is provided use smtp send
	if tm.Password == "" {
		c, err := smtp.Dial(fmt.Sprintf("%s:%v", tm.Server, int(tm.Port)))
		if err != nil {
			direktivapps.RespondWithError(w, fmt.Sprintf(code, "smtp-dial"), err.Error())
			return
		}
		defer c.Close()

		if err = c.Mail(tm.From); err != nil {
			direktivapps.RespondWithError(w, fmt.Sprintf(code, "smtp-from-address"), err.Error())
			return
		}

		for _, to := range tm.To {
			if err = c.Rcpt(to); err != nil {
				direktivapps.RespondWithError(w, fmt.Sprintf(code, "smtp-to-address"), err.Error())
				return
			}
		}

		wf, err := c.Data()
		if err != nil {
			direktivapps.RespondWithError(w, fmt.Sprintf(code, "smtp-data"), err.Error())
			return
		}

		msg := fmt.Sprintf("To: %s\r\nFrom: %s\r\nSubject: %s\r\nContent-Type: text/html; charset=\"UTF-8\"\r\n\r\n%s", strings.Join(tm.To, ","), tm.From, tm.Subject, tpl.String())
		_, err = wf.Write([]byte(msg))
		if err != nil {
			direktivapps.RespondWithError(w, fmt.Sprintf(code, "smtp-write-msg"), err.Error())
			return
		}
		err = wf.Close()
		if err != nil {
			direktivapps.RespondWithError(w, fmt.Sprintf(code, "close-writer"), err.Error())
			return
		}

		c.Quit()
	} else {
		// Settings for SMTP server
		d := gomail.NewDialer(tm.Server, int(tm.Port), tm.From, tm.Password)
		d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

		direktivapps.Log(aid, "Sending Message")
		if err := d.DialAndSend(m); err != nil {
			direktivapps.RespondWithError(w, fmt.Sprintf(code, "send-message"), err.Error())
			return
		}

	}
	direktivapps.Respond(w, []byte{})

}

func main() {
	direktivapps.StartServer(SMTPEmailHandler)
}
