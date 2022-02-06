package main

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/direktiv/direktiv-apps/pkg/reusable"
	mail "github.com/xhit/go-simple-mail/v2"
)

type requestInput struct {
	Body        string          `json:"body"`
	CC          []string        `json:"cc"`
	Bcc         []string        `json:"bcc"`
	Subject     string          `json:"subject"`
	TLS         bool            `json:"tls"`
	From        string          `json:"from"`
	User        string          `json:"user"`
	Password    string          `json:"password"`
	To          []string        `json:"to"`
	Host        string          `json:"host"`
	Attachments []reusable.File `json:"attachments"`
	Port        float64         `json:"port"`
	Individual  bool            `json:"individual"`
}

func smtpHandler(w http.ResponseWriter, r *http.Request, ri *reusable.RequestInfo) {

	obj := new(requestInput)
	err := reusable.Unmarshal(obj, false, r)
	if err != nil {
		reusable.ReportError(w, reusable.UnmarshallError, err)
		return
	}

	if len(obj.To) == 0 && len(obj.Bcc) == 0 {
		reusable.ReportError(w, errForCode("connect"), fmt.Errorf("to and bcc are empty"))
		return
	}

	if len(obj.Host) == 0 || obj.Port == 0 {
		reusable.ReportError(w, errForCode("connect"), fmt.Errorf("host and port not set"))
		return
	}

	server := mail.NewSMTPClient()
	server.Host = obj.Host
	server.Port = int(obj.Port)

	server.Authentication = mail.AuthNone
	if len(obj.Password) > 0 {
		ri.Logger().Infof("using password authentication")
		server.Username = obj.User
		server.Password = obj.Password
		server.Authentication = mail.AuthPlain
	}

	server.SendTimeout = 120 * time.Second
	server.KeepAlive = false
	if obj.Individual {
		server.KeepAlive = true
	}
	server.ConnectTimeout = 30 * time.Second

	if obj.TLS {
		ri.Logger().Infof("using tls")
		server.Encryption = mail.EncryptionTLS
	}

	// SMTP client
	ri.Logger().Infof("smtp connect")
	ee, err := server.Connect()
	if err != nil {
		reusable.ReportError(w, errForCode("connect"), err)
		return
	}

	if obj.Individual {

		ri.Logger().Infof("smtp send individual")
		for i := range obj.To {
			ri.Logger().Infof("sending to %v", obj.To[i])
			email := generateEmail([]string{obj.To[i]}, obj, ri)
			err = email.Send(ee)
			if err != nil {
				ri.Logger().Infof("error sending to %v: %v", obj.To[i], err)
			}

		}

	} else {

		ri.Logger().Infof("sending email to %v", obj.To)
		email := generateEmail(obj.To, obj, ri)
		err = email.Send(ee)
		if err != nil {
			reusable.ReportError(w, errForCode("send"), err)
			return
		}

	}

	reusable.ReportResult(w, []byte("{}"))

}

func generateEmail(to []string, obj *requestInput, ri *reusable.RequestInfo) *mail.Email {

	email := mail.NewMSG()
	email.SetFrom(obj.From)
	email.AddBcc(obj.Bcc...)
	email.AddCc(obj.CC...)
	email.SetSubject(obj.Subject)
	email.SetBody(mail.TextHTML, obj.Body)
	email.AddTo(to...)
	attachFiles(email, reusable.NewFileIterator(obj.Attachments, ri))

	return email

}

func attachFiles(email *mail.Email, fi *reusable.FileIterator) error {

	for {
		f, err := fi.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		if err != nil {
			return err
		}

		c, err := f.AsBase64()
		if err != nil {
			return err
		}

		attachFile := &mail.File{
			B64Data:  c,
			Name:     f.Name,
			MimeType: f.ContentType,
		}

		email.Attach(attachFile)

	}

	return nil

}

func errForCode(errCode string) string {
	return fmt.Sprintf("com.smtp-bare.%s.error", errCode)
}

func main() {
	reusable.StartServer(smtpHandler, nil)
}

// import (
// 	"fmt"
// 	"io"
// 	"io/ioutil"
// 	"net/http"
// 	"os"
// 	"path/filepath"
// 	"strings"
// 	"time"

// 	"github.com/direktiv/direktiv-apps/pkg/direktivapps"
// 	mail "github.com/xhit/go-simple-mail/v2"
// )

// var code = "com.smtp-bare.%s.error"

// type Attachment struct {
// 	Name string `json:"name"`
// 	Data string `json:"data"`
// 	Type string `json:"type"`
// }
// type RequestInput struct {
// 	Body        string       `json:"body"`
// 	CC          []string     `json:"cc"`
// 	Bcc         []string     `json:"bcc"`
// 	Subject     string       `json:"subject"`
// 	TLS         bool         `json:"tls"`
// 	From        string       `json:"from"`
// 	Password    string       `json:"password"`
// 	To          []string     `json:"to"`
// 	Host        string       `json:"address"`
// 	Attachments []Attachment `json:"attachments"`
// 	Port        float64      `json:"port"`
// }

// func main() {
// 	direktivapps.StartServer(SMTPHandler)
// }

// func SMTPHandler(w http.ResponseWriter, r *http.Request) {
// 	var ri RequestInput

// 	aid, err := direktivapps.Unmarshal(&ri, r)
// 	if err != nil {
// 		direktivapps.RespondWithError(w, fmt.Sprintf(code, "unmarshal"), err.Error())
// 		return
// 	}

// 	server := mail.NewSMTPClient()
// 	server.Host = ri.Host

// 	server.Port = int(ri.Port)

// 	// check for authentication
// 	if ri.From != "" && ri.Password != "" {
// 		direktivapps.LogDouble(aid, "using username and password")
// 		server.Username = ri.From
// 		server.Password = ri.Password
// 	}
// 	// check for TLS
// 	if ri.TLS {
// 		direktivapps.LogDouble(aid, "using encryption")
// 		server.Encryption = mail.EncryptionTLS
// 	}

// 	if ri.Password == "" {
// 		direktivapps.LogDouble(aid, "using auth none")
// 		server.Authentication = mail.AuthNone
// 	} else {
// 		direktivapps.LogDouble(aid, "using auth plain")
// 		server.Authentication = mail.AuthPlain
// 	}

// 	// Variable to keep alive connection
// 	server.KeepAlive = false

// 	// Timeout for connect to SMTP Server
// 	server.ConnectTimeout = 30 * time.Second

// 	// Timeout for send the data and wait respond
// 	server.SendTimeout = 120 * time.Second

// 	// SMTP client
// 	ee, err := server.Connect()
// 	if err != nil {
// 		direktivapps.LogDouble(aid, "could not connect: %v", err.Error())
// 		direktivapps.RespondWithError(w, fmt.Sprintf(code, "smtp-connect"), err.Error())
// 		return
// 	}

// 	email := mail.NewMSG()
// 	email.SetFrom(ri.From)
// 	email.AddTo(ri.To...)
// 	email.AddBcc(ri.Bcc...)
// 	email.AddCc(ri.CC...)
// 	email.SetSubject(ri.Subject)
// 	email.SetBody(mail.TextHTML, ri.Body)

// 	direktivapps.LogDouble(aid, fmt.Sprintf("%v\n", ri.Attachments))
// 	for _, attach := range ri.Attachments {
// 		direktivapps.LogDouble(aid, fmt.Sprintf("TYPE: '%#v'\n", attach.Type))
// 		switch attach.Type {
// 		case "base64":
// 			direktivapps.LogDouble(aid, "adding base64 attachment")

// 			if attach.Data == "" {
// 				direktivapps.LogDouble(aid, "pulling from temp variable")
// 				data, err := ioutil.ReadFile(filepath.Join(r.Header.Get("Direktiv-TempDir"), attach.Name))
// 				if err != nil {
// 					direktivapps.RespondWithError(w, fmt.Sprintf(code, "read-file"), err.Error())
// 					return
// 				}
// 				attach.Data = strings.TrimSuffix(strings.TrimPrefix(string(data), "\""), "\"")

// 			}
// 			email.AddAttachmentBase64(attach.Data, attach.Name)
// 		default:
// 			var f *os.File
// 			if attach.Data == "" {
// 				direktivapps.LogDouble(aid, "pulling from temp variable")
// 				f, err = os.Open(filepath.Join(r.Header.Get("Direktiv-TempDir"), attach.Name))
// 				if err != nil {
// 					direktivapps.RespondWithError(w, fmt.Sprintf(code, "read-file"), err.Error())
// 					return
// 				}
// 				defer f.Close()
// 			} else {
// 				f, err = os.Create(attach.Name)
// 				if err != nil {
// 					direktivapps.RespondWithError(w, fmt.Sprintf(code, "create-file"), err.Error())
// 					return
// 				}
// 				defer f.Close()
// 				_, err = io.Copy(f, strings.NewReader(attach.Data))
// 				if err != nil {
// 					direktivapps.RespondWithError(w, fmt.Sprintf(code, "copy"), err.Error())
// 					return
// 				}
// 			}

// 			email.AddAttachment(f.Name(), filepath.Base(f.Name()))
// 		}
// 	}

// 	err = email.Send(ee)
// 	if err != nil {
// 		direktivapps.RespondWithError(w, fmt.Sprintf(code, "send-mail"), err.Error())
// 		return
// 	}

// 	direktivapps.Respond(w, []byte{})
// }
