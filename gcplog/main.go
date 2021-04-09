package main

import (
	"context"
	"net/http"

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

const code = "com.gcplog.error"

func main() {
	direktivapps.StartServer(GCPLog)
}

func GCPLog(w http.ResponseWriter, r *http.Request) {
	obj := new(GCPLogging)
	_, err := direktivapps.Unmarshal(obj, r)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	err = ioutil.WriteFile("/tmp/key.json", []byte(obj.ServiceAccountKey), 0644)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}
	ctx := context.Background()

	// create a gcp logging client
	client, err := logging.NewClient(ctx, obj.ProjectID, option.WithCredentialsFile("/tmp/key.json"))
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}
	logger := client.Logger(obj.LogName).StandardLogger(logging.Info)

	logger.Println(obj.Message)
	client.Close()

	// Write empty to notified its finished
	direktivapps.Respond(w, []byte{})
}
