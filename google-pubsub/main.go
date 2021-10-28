package main

import (
	"context"
	"io/ioutil"
	"net/http"

	"cloud.google.com/go/pubsub"
	"github.com/direktiv/direktiv-apps/pkg/direktivapps"
	"google.golang.org/api/option"
)

type PubSubInput struct {
	ServiceAccountKey string            `json:"serviceAccountKey"`
	Message           string            `json:"message"`
	ProjectID         string            `json:"project-id"`
	TopicID           string            `json:"topic-id"`
	Attributes        map[string]string `json:"attributes"`
}

const code = "com.google-pubsub.error"

func main() {
	direktivapps.StartServer(GooglePubSub)
}

func GooglePubSub(w http.ResponseWriter, r *http.Request) {
	obj := new(PubSubInput)
	_, err := direktivapps.Unmarshal(obj, r)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	err = ioutil.WriteFile("/key.json", []byte(obj.ServiceAccountKey), 0700)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, obj.ProjectID, option.WithCredentialsFile("/key.json"))
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	topic := client.Topic(obj.TopicID)
	topic.Publish(ctx, &pubsub.Message{
		Data:       []byte(obj.Message),
		Attributes: obj.Attributes,
	})

	topic.Stop()
	client.Close()

	direktivapps.Respond(w, []byte{})
}
