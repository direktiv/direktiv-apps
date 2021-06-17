package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/streadway/amqp"
	"github.com/vorteil/direktiv-apps/pkg/direktivapps"
)

const DefaultWaitTimeMilliseconds = 1000

type RabbitMQInput struct {
	Address  string `json:"address"`
	Username string `json:"username"`
	Password string `json:"password"`
	Queue    string `json:"queue"`
	WaitTime int    `json:"waitTime`
}

type RabbitResponseMsg struct {
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

func failOnError(g direktivapps.ActionError, err error, msg string) {
	if err != nil {
		g.ErrorMessage = fmt.Sprintf("%s: %s", msg, err)
		direktivapps.WriteError(g)
	}
}

var code = "com.rabbitmq-consumer.error"

func RabbitMQHandler(w http.ResponseWriter, r *http.Request) {
	obj := new(RabbitMQInput)
	_, err := direktivapps.Unmarshal(obj, r)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	// Set default wait time
	if obj.WaitTime == 0 {
		obj.WaitTime = DefaultWaitTimeMilliseconds
	}

	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s/", obj.Username, obj.Password, obj.Address))
	if err != nil {
		direktivapps.RespondWithError(w, code, "Failed to connect to RabbitMq")
		return
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		direktivapps.RespondWithError(w, code, "Failed to open a channel")
		return
	}

	q, err := ch.QueueDeclare(
		obj.Queue, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       //arguments
	)
	if err != nil {
		direktivapps.RespondWithError(w, code, "Failed to declare a queue")
		return
	}

	msgs, err := ch.Consume(
		q.Name, //queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		direktivapps.RespondWithError(w, code, "Failed to consume a messages")
		return
	}

	messageList := make([]RabbitResponseMsg, 0)
	readMessage := true
	go func() {
		for d := range msgs {
			if readMessage {
				log.Printf("Received a message")
				messageList = append(messageList, RabbitResponseMsg{
					Message:   string(d.Body),
					Timestamp: d.Timestamp.String(),
				})
			}
		}
	}()

	// Wait for messages
	time.Sleep(time.Duration(obj.WaitTime) * time.Millisecond)
	readMessage = false

	b, err := json.Marshal(messageList)
	if err != nil {
		direktivapps.RespondWithError(w, code, "Failed to decode messages")
		return
	}

	direktivapps.Respond(w, b)
}

func main() {
	direktivapps.StartServer(RabbitMQHandler)
}
