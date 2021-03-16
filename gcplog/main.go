package main

import (
	"context"

	"io/ioutil"

	"cloud.google.com/go/logging"
	"github.com/vorteil/direktiv-apps/pkg/direktivapps"
	"google.golang.org/api/option"
)

type GCPLogging struct {
	ServiceAccountKey string `json:"serviceAccountKey"`
	Message           string `json:"message"`
	ProjectID         string `json:"project-id"`
	LogName           string `json:"log-name"`
}

func main() {
	g := direktivapps.ActionError{
		ErrorCode:    "com.gcplog.error",
		ErrorMessage: "",
	}

	obj := new(GCPLogging)
	direktivapps.ReadIn(obj, g)

	err := ioutil.WriteFile("/tmp/key.json", []byte(obj.ServiceAccountKey), 0644)
	if err != nil {
		g.ErrorMessage = err.Error()
		direktivapps.WriteError(g)
	}
	ctx := context.Background()

	// create a gcp logging client
	client, err := logging.NewClient(ctx, obj.ProjectID, option.WithCredentialsFile("/tmp/key.json"))
	if err != nil {
		g.ErrorMessage = err.Error()
		direktivapps.WriteError(g)
	}
	logger := client.Logger(obj.LogName).StandardLogger(logging.Info)

	logger.Println(obj.Message)
	client.Close()

	// Write empty to notified its finished
	direktivapps.WriteOut([]byte{}, g)
}
