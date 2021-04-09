package main

import (
	"context"
	"net/http"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/vorteil/direktiv-apps/pkg/direktivapps"
)

type KafkaMsgInput struct {
	Topic     string `json:"topic"`
	Partition int    `json:"partition"`
	Message   string `json:"message"`
	Address   string `json:"address"`
}

var code = "com.kafka.error"

func KafkaHandler(w http.ResponseWriter, r *http.Request) {
	obj := new(KafkaMsgInput)
	_, err := direktivapps.Unmarshal(obj, r)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	conn, err := kafka.DialLeader(context.Background(), "tcp", obj.Address, obj.Topic, obj.Partition)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	_, err = conn.WriteMessages(
		kafka.Message{Value: []byte(obj.Message)},
	)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	if err := conn.Close(); err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	direktivapps.Respond(w, []byte{})
}

func main() {
	direktivapps.StartServer(KafkaHandler)
}
