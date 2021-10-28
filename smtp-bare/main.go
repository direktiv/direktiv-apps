package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/direktiv/direktiv-apps/pkg/direktivapps"
	mail "github.com/xhit/go-simple-mail"
)

var code = "com.smtp-bare.%s.error"

type Attachment struct {
	Name string `json:"name"`
	Data string `json:"data"`
	Type string `json:"type"`
}
type RequestInput struct {
	Body        string       `json:"body"`
	CC          []string     `json:"cc"`
	Bcc         []string     `json:"bcc"`
	Subject     string       `json:"subject"`
	TLS         bool         `json:"tls"`
	From        string       `json:"from"`
	Password    string       `json:"password"`
	To          []string     `json:"to"`
	Host        string       `json:"address"`
	Attachments []Attachment `json:"attachments"`
	Port        float64      `json:"port"`
}

func main() {
	direktivapps.StartServer(SMTPHandler)
}

func SMTPHandler(w http.ResponseWriter, r *http.Request) {
	var ri RequestInput

	aid, err := direktivapps.Unmarshal(&ri, r)
	if err != nil {
		direktivapps.RespondWithError(w, fmt.Sprintf(code, "unmarshal"), err.Error())
		return
	}

	server := mail.NewSMTPClient()
	server.Host = ri.Host

	server.Port = int(ri.Port)

	// server.Authentication = mail.AuthPlain

	// check for authentication
	if ri.From != "" && ri.Password != "" {
		server.Username = ri.From
		server.Password = ri.Password
	}
	// check for TLS
	if ri.TLS {
		server.Encryption = mail.EncryptionTLS
	}

	// if ri.Password == "" {
	// 	server.Authentication = mail.EncodingNone
	// } else {
	// 	server.Authentication = mail.AuthPlain
	// }

	// Variable to keep alive connection
	server.KeepAlive = false

	// Timeout for connect to SMTP Server
	server.ConnectTimeout = 10 * time.Second

	// Timeout for send the data and wait respond
	server.SendTimeout = 60 * time.Second

	// SMTP client
	ee, err := server.Connect()
	if err != nil {
		direktivapps.RespondWithError(w, fmt.Sprintf(code, "smtp-connect"), err.Error())
		return
	}

	email := mail.NewMSG()
	email.SetFrom(ri.From)
	email.AddTo(ri.To...)
	email.AddBcc(ri.Bcc...)
	email.AddCc(ri.CC...)
	email.SetSubject(ri.Subject)
	email.SetBody(mail.TextHTML, ri.Body)

	direktivapps.LogDouble(aid, fmt.Sprintf("%v\n", ri.Attachments))
	for _, attach := range ri.Attachments {
		direktivapps.LogDouble(aid, fmt.Sprintf("TYPE: '%#v'\n", attach.Type))
		switch attach.Type {
		case "base64":
			direktivapps.LogDouble(aid, "adding base64 attachment")

			if attach.Data == "" {
				direktivapps.LogDouble(aid, "pulling from temp variable")
				data, err := ioutil.ReadFile(filepath.Join(r.Header.Get("Direktiv-TempDir"), attach.Name))
				if err != nil {
					direktivapps.RespondWithError(w, fmt.Sprintf(code, "read-file"), err.Error())
					return
				}
				attach.Data = strings.TrimSuffix(strings.TrimPrefix(string(data), "\""), "\"")

			}
			email.AddAttachmentBase64(attach.Data, attach.Name)
		default:
			var f *os.File
			if attach.Data == "" {
				direktivapps.LogDouble(aid, "pulling from temp variable")
				f, err = os.Open(filepath.Join(r.Header.Get("Direktiv-TempDir"), attach.Name))
				if err != nil {
					direktivapps.RespondWithError(w, fmt.Sprintf(code, "read-file"), err.Error())
					return
				}
				defer f.Close()
			} else {
				f, err = os.Create(attach.Name)
				if err != nil {
					direktivapps.RespondWithError(w, fmt.Sprintf(code, "create-file"), err.Error())
					return
				}
				defer f.Close()
				_, err = io.Copy(f, strings.NewReader(attach.Data))
				if err != nil {
					direktivapps.RespondWithError(w, fmt.Sprintf(code, "copy"), err.Error())
					return
				}
			}

			email.AddAttachment(f.Name(), filepath.Base(f.Name()))
		}
	}

	err = email.Send(ee)
	if err != nil {
		direktivapps.RespondWithError(w, fmt.Sprintf(code, "send-mail"), err.Error())
		return
	}

	direktivapps.Respond(w, []byte{})
}
