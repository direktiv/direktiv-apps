package main

import (
	"fmt"

	"github.com/streadway/amqp"
	"github.com/vorteil/direktiv-apps/pkg/direktivapps"
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

func main() {
	g := direktivapps.ActionError{
		ErrorCode:    "com.rabbitmq.error",
		ErrorMessage: "",
	}

	obj := new(RabbitMQInput)
	direktivapps.ReadIn(obj, g)

	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s/", obj.Username, obj.Password, obj.Address))
	failOnError(g, err, "Failed to connect to RabbitMq")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(g, err, "Failed to open a channel")

	q, err := ch.QueueDeclare(
		obj.Queue, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       //arguments
	)
	failOnError(g, err, "Failed to declare a queue")

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
	failOnError(g, err, "Failed to publish a message")

	direktivapps.WriteOut([]byte{}, g)
}
