package main

import (
	"context"
	"encoding/base64"
	"io"
	"io/ioutil"
	"os"
	"time"

	"github.com/vorteil/direktiv-apps/pkg/direktivapps"
	"google.golang.org/api/option"

	"cloud.google.com/go/storage"
)

type GoogleUploadBlob struct {
	Bucket            string `json:"bucket"`
	Object            string `json:"upload-name"`
	Data              string `json:"data"`
	ServiceAccountKey string `json:"serviceAccountKey"`
}

const credFile = "/tmp/creds"

func main() {
	g := direktivapps.ActionError{
		ErrorCode:    "com.google-upload.error",
		ErrorMessage: "",
	}

	obj := new(GoogleUploadBlob)
	direktivapps.ReadIn(obj, g)
	ctx := context.Background()

	err := ioutil.WriteFile(credFile, []byte(obj.ServiceAccountKey), 0777)
	if err != nil {
		g.ErrorMessage = err.Error()
		direktivapps.WriteError(g)
	}

	client, err := storage.NewClient(ctx, option.WithCredentialsFile(credFile))
	if err != nil {
		g.ErrorMessage = err.Error()
		direktivapps.WriteError(g)
	}

	defer client.Close()

	decoded, err := base64.StdEncoding.DecodeString(obj.Data)
	if err != nil {
		g.ErrorMessage = err.Error()
		direktivapps.WriteError(g)
	}

	err = ioutil.WriteFile(obj.Object, decoded, 0700)
	if err != nil {
		g.ErrorMessage = err.Error()
		direktivapps.WriteError(g)
	}

	f, err := os.Open(obj.Object)
	if err != nil {
		g.ErrorMessage = err.Error()
		direktivapps.WriteError(g)
	}

	defer f.Close()
	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	wc := client.Bucket(obj.Bucket).Object(obj.Object).NewWriter(ctx)
	if _, err = io.Copy(wc, f); err != nil {
		g.ErrorMessage = err.Error()
		direktivapps.WriteError(g)
	}
	if err := wc.Close(); err != nil {
		g.ErrorMessage = err.Error()
		direktivapps.WriteError(g)
	}

	direktivapps.WriteOut([]byte{}, g)
}
