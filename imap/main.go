package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/mail"
	"github.com/vorteil/direktiv-apps/pkg/direktivapps"
)

var code = "com.imap.error"

type request struct {
	Username string `json:"email"`
	Password string `json:"password"`
	IMAP     string `json:"imap-address"` // includes the port
}

type output struct {
	Message string `json:"msg"` // msg of the email body
}

func IMAPHandler(w http.ResponseWriter, r *http.Request) {
	var rr request
	aid, err := direktivapps.Unmarshal(&rr, r)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	direktivapps.Log(aid, "Connecting to IMAP server...")
	// Connect to server
	c, err := client.DialTLS(rr.IMAP, nil)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}
	direktivapps.Log(aid, "Connected")

	// Don't forget to logout
	defer c.Logout()

	// Login
	if err := c.Login(rr.Username, rr.Password); err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}
	direktivapps.Log(aid, "Logged In")

	// Select INBOX
	mbox, err := c.Select("INBOX", false)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	// Get the last message
	if mbox.Messages == 0 {
		direktivapps.RespondWithError(w, code, "No message in mailbox")
		return
	}
	seqSet := new(imap.SeqSet)
	seqSet.AddNum(mbox.Messages)

	// Get the whole message body
	var section imap.BodySectionName
	items := []imap.FetchItem{section.FetchItem()}

	messages := make(chan *imap.Message, 1)
	go func() {
		if err := c.Fetch(seqSet, items, messages); err != nil {
			direktivapps.RespondWithError(w, code, err.Error())
			return
		}
	}()

	msg := <-messages
	if msg == nil {
		direktivapps.RespondWithError(w, code, "Server didn't returned message")
		return
	}

	read := msg.GetBody(&section)
	if read == nil {
		direktivapps.RespondWithError(w, code, "Server didn't returned message body")
		return
	}

	// Create a new mail reader
	mr, err := mail.CreateReader(read)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	// Process each message's part
	for {
		p, err := mr.NextPart()
		if err == io.EOF {
			break
		} else if err != nil {
			direktivapps.RespondWithError(w, code, err.Error())
			return
		}

		switch p.Header.(type) {
		case *mail.InlineHeader:
			// This is the message's text (can be plain-text or HTML)
			b, _ := ioutil.ReadAll(p.Body)
			var o output
			o.Message = string(b)
			data, err := json.Marshal(o)
			if err != nil {
				direktivapps.RespondWithError(w, code, err.Error())
				return
			}
			direktivapps.Respond(w, data)
		}
	}
}

func main() {
	direktivapps.StartServer(IMAPHandler)
}
