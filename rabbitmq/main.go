package main

import (
	"fmt"
	"net/http"

	"github.com/streadway/amqp"
	"github.com/direktiv/direktiv-apps/pkg/direktivapps"
)

type RabbitMQInput struct {
	Address  string `json:"address"`
	Username string `json:"username"`
	Password string `json:"password"`
	Queue    string `json:"queue"`
	Message  string `json:"message"`
}

func failOnError(g direktivapps.ActionError, err error, msg string) {
	if err != nil {
		g.ErrorMessage = fmt.Sprintf("%s: %s", msg, err)
		direktivapps.WriteError(g)
	}
}

var code = "com.rabbitmq.error"

func RabbitMQHandler(w http.ResponseWriter, r *http.Request) {
	obj := new(RabbitMQInput)
	_, err := direktivapps.Unmarshal(obj, r)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
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

	err = ch.Publish(
		"",     // exchange
		q.Name, //routing key,
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(obj.Message),
		},
	)
	if err != nil {
		direktivapps.RespondWithError(w, code, "Failed to publish a message")
		return
	}

	direktivapps.Respond(w, []byte{})
}

func main() {
	direktivapps.StartServer(RabbitMQHandler)
}
