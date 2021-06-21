package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	nats "github.com/nats-io/nats.go"
	"github.com/vorteil/direktiv-apps/pkg/direktivapps"
)

const DefaultTimeoutMilliseconds = 1000
const DefaultMaxMessagesCount = 10

type NatsInput struct {
	Address     string `json:"address"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	Subject     string `json:"subject"`
	Timeout     int    `json:"timeout"`
	MaxMessages int    `json:"max-messages"`
}

type NatsResponse struct {
	Messages []string `json:"messages"`
}

func failOnError(g direktivapps.ActionError, err error, msg string) {
	if err != nil {
		g.ErrorMessage = fmt.Sprintf("%s: %s", msg, err)
		direktivapps.WriteError(g)
	}
}

var code = "com.nats-consume.error"

func NatsHandler(w http.ResponseWriter, r *http.Request) {
	obj := new(NatsInput)
	aid, err := direktivapps.Unmarshal(obj, r)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	// Set default timeout
	if obj.Timeout <= 0 {
		obj.Timeout = DefaultTimeoutMilliseconds
	}

	if obj.MaxMessages <= 0 {
		obj.MaxMessages = DefaultMaxMessagesCount
	}

	// Connect to a server
	nc, err := nats.Connect(obj.Address, nats.UserInfo(obj.Username, obj.Password))
	if err != nil {
		direktivapps.RespondWithError(w, code, fmt.Sprintf("Failed to connect to NATS service: %v", err))
		return
	}
	defer nc.Close()

	// Create Response
	messageList := NatsResponse{
		Messages: make([]string, 0),
	}

	// Subscribe
	ch := make(chan *nats.Msg, 64)
	sub, err := nc.ChanSubscribe(obj.Subject, ch)
	if err != nil {
		direktivapps.RespondWithError(w, code, fmt.Sprintf("Failed to subscribe to NATS service: %v", err))
		return
	}

	// Set Timeout
	timeoutCh := make(chan bool, 1)
	go func(timeoutMS int) {
		time.Sleep(time.Duration(timeoutMS) * time.Millisecond)
		timeoutCh <- true
	}(obj.Timeout)

	var loopDone bool

	// Read messages
	for !loopDone {
		select {
		case msg := <-ch:
			messageList.Messages = append(messageList.Messages, string(msg.Data))
			if len(messageList.Messages) >= obj.MaxMessages {
				// Break - Max message count recieved
				loopDone = true
				direktivapps.Log(aid, "NATS max messages recieved")
			}
		case loopDone = <-timeoutCh:
			// Break - Timed out
			direktivapps.Log(aid, "NATS timeout reached")
		}
	}

	sub.Unsubscribe()

	response, err := json.Marshal(messageList)
	if err != nil {
		direktivapps.RespondWithError(w, code, "Failed to decode messages")
		return
	}

	direktivapps.Respond(w, response)
}

func main() {
	direktivapps.StartServer(NatsHandler)
}
