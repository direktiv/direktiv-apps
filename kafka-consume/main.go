package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/direktiv/direktiv-apps/pkg/direktivapps"
)

const DefaultTimeoutMilliseconds = 1000

type KafkaMsgInput struct {
	Topic     string `json:"topic"`
	Partition int    `json:"partition"`
	Address   string `json:"address"`
	Timeout   int    `json:"timeout"`
}

type KafkaResponse struct {
	Messages []string `json:"messages"`
}

var code = "com.kafka-read.error"

func KafkaHandler(w http.ResponseWriter, r *http.Request) {
	obj := new(KafkaMsgInput)
	aid, err := direktivapps.Unmarshal(obj, r)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	// Set default timeout
	if obj.Timeout == 0 {
		obj.Timeout = DefaultTimeoutMilliseconds
	}

	conn, err := kafka.DialLeader(context.Background(), "tcp", obj.Address, obj.Topic, obj.Partition)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	messageList := KafkaResponse{
		Messages: make([]string, 0),
	}

	conn.SetReadDeadline(time.Now().Add(time.Duration(obj.Timeout) * time.Millisecond))
	batch := conn.ReadBatch(0, 3e6) // max 3MB of messages
	b := make([]byte, 3e6)

	// Read messages
	for {
		_, err := batch.Read(b)
		if err != nil {
			break
		}

		// trim empty bytes
		messageList.Messages = append(messageList.Messages, string(bytes.Trim(b, "\x00")))
	}

	if err := batch.Close(); err != nil {
		if err == kafka.RequestTimedOut {
			direktivapps.Log(aid, "Warning - kafka request timed out")
		} else {
			direktivapps.RespondWithError(w, code, err.Error())
			return
		}
	}

	if err := conn.Close(); err != nil {
		direktivapps.RespondWithError(w, code, fmt.Sprintf("failed to close connection: %s", err.Error()))
		return
	}

	response, err := json.Marshal(messageList)
	if err != nil {
		direktivapps.RespondWithError(w, code, "Failed to decode messages")
		return
	}

	direktivapps.Respond(w, response)
}

func main() {
	direktivapps.StartServer(KafkaHandler)
}
