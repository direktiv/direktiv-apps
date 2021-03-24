package main

import (
	"github.com/michaelbironneau/asbclient"
	"github.com/vorteil/direktiv-apps/pkg/direktivapps"
)

type ServiceBusInput struct {
	Policy    string `json:"policy"`
	Key       string `json:"key"`
	Namespace string `json:"namespace"`
	Message   string `json:"message"`
	Queue     string `json:"queue"`
}

func main() {
	g := direktivapps.ActionError{
		ErrorCode:    "com-azure-servicebus.error",
		ErrorMessage: "",
	}

	obj := new(ServiceBusInput)
	direktivapps.ReadIn(obj, g)

	client := asbclient.New(asbclient.Topic, obj.Namespace, obj.Policy, obj.Key)
	err := client.Send(obj.Queue, &asbclient.Message{
		Body: []byte(obj.Message),
	})
	if err != nil {
		g.ErrorMessage = err.Error()
		direktivapps.WriteError(g)
	}

	direktivapps.WriteOut([]byte{}, g)
}
