package main

import (
	"fmt"
	"net/http"

	nats "github.com/nats-io/nats.go"
	"github.com/vorteil/direktiv-apps/pkg/direktivapps"
)

type NatsInput struct {
	Address  string `json:"address"`
	Username string `json:"username"`
	Password string `json:"password"`
	Subject  string `json:"subject"`
	Message  string `json:"message"`
}

func failOnError(g direktivapps.ActionError, err error, msg string) {
	if err != nil {
		g.ErrorMessage = fmt.Sprintf("%s: %s", msg, err)
		direktivapps.WriteError(g)
	}
}

var code = "com.nats.error"

func NatsHandler(w http.ResponseWriter, r *http.Request) {
	obj := new(NatsInput)
	_, err := direktivapps.Unmarshal(obj, r)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	// Connect to a server
	nc, err := nats.Connect(obj.Address, nats.UserInfo(obj.Username, obj.Password))
	if err != nil {
		direktivapps.RespondWithError(w, code, fmt.Sprintf("Failed to connect to NATS service: %v", err))
		return
	}
	defer nc.Close()

	if err := nc.Publish(obj.Subject, []byte(obj.Message)); err != nil {
		direktivapps.RespondWithError(w, code, fmt.Sprintf("Failed to write message to NATS service: %v", err))
		return
	}

	direktivapps.Respond(w, []byte{})
}

func main() {
	direktivapps.StartServer(NatsHandler)
}
