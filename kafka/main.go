package main

import (
	"context"
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

func main() {
	g := direktivapps.ActionError{
		ErrorCode:    "com.kafka.error",
		ErrorMessage: "",
	}

	obj := new(KafkaMsgInput)
	direktivapps.ReadIn(obj, g)

	conn, err := kafka.DialLeader(context.Background(), "tcp", obj.Address, obj.Topic, obj.Partition)
	if err != nil {
		g.ErrorMessage = err.Error()
		direktivapps.WriteError(g)
	}

	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	_, err = conn.WriteMessages(
		kafka.Message{Value: []byte(obj.Message)},
	)
	if err != nil {
		g.ErrorMessage = err.Error()
		direktivapps.WriteError(g)
	}

	if err := conn.Close(); err != nil {
		g.ErrorMessage = err.Error()
		direktivapps.WriteError(g)
	}

	direktivapps.WriteOut([]byte{}, g)
}
