package main

import (
	"context"
	"io/ioutil"

	"cloud.google.com/go/pubsub"
	"github.com/vorteil/direktiv-apps/pkg/direktivapps"
	"google.golang.org/api/option"
)

type PubSubInput struct {
	ServiceAccountKey string            `json:"serviceAccountKey"`
	Message           string            `json:"message"`
	ProjectID         string            `json:"project-id"`
	TopicID           string            `json:"topic-id"`
	Attributes        map[string]string `json:"attributes"`
}

func main() {

	g := direktivapps.ActionError{
		ErrorCode:    "com.google-pubsub.error",
		ErrorMessage: "",
	}

	obj := new(PubSubInput)
	direktivapps.ReadIn(obj, g)

	err := ioutil.WriteFile("/tmp/key.json", []byte(obj.ServiceAccountKey), 0700)
	if err != nil {
		g.ErrorMessage = err.Error()
		direktivapps.WriteError(g)
	}

	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, obj.ProjectID, option.WithCredentialsFile("/tmp/key.json"))
	if err != nil {
		g.ErrorMessage = err.Error()
		direktivapps.WriteError(g)
	}

	topic := client.Topic(obj.TopicID)
	topic.Publish(ctx, &pubsub.Message{
		Data:       []byte(obj.Message),
		Attributes: obj.Attributes,
	})

	topic.Stop()
	client.Close()

	direktivapps.WriteOut([]byte{}, g)
}
