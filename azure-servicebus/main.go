package main

import (
	"github.com/michaelbironneau/asbclient"
	"github.com/direktiv/direktiv-apps/pkg/direktivapps"
	"net/http"
)

type ServiceBusInput struct {
	Policy    string `json:"policy"`
	Key       string `json:"key"`
	Namespace string `json:"namespace"`
	Message   string `json:"message"`
	Queue     string `json:"queue"`
}

const code = "com-azure-servicebus.error"

func main() {
	direktivapps.StartServer(AzureServiceBus)
}

func AzureServiceBus(w http.ResponseWriter, r *http.Request) {
	obj := new(ServiceBusInput)
	_, err := direktivapps.Unmarshal(obj, r)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	client := asbclient.New(asbclient.Topic, obj.Namespace, obj.Policy, obj.Key)
	err = client.Send(obj.Queue, &asbclient.Message{
		Body: []byte(obj.Message),
	})
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	direktivapps.Respond(w, []byte{})
}
