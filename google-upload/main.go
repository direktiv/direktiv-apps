package main

import (
	"context"
	"encoding/base64"
	"io"
	"io/ioutil"
	"net/http"
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
const code = "com.google-upload.error"

func main() {
	direktivapps.StartServer(GoogleUpload)
}

func GoogleUpload(w http.ResponseWriter, r *http.Request) {

	obj := new(GoogleUploadBlob)
	_, err := direktivapps.Unmarshal(obj, r)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}
	ctx := context.Background()

	err = ioutil.WriteFile(credFile, []byte(obj.ServiceAccountKey), 0777)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	client, err := storage.NewClient(ctx, option.WithCredentialsFile(credFile))
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	defer client.Close()

	decoded, err := base64.StdEncoding.DecodeString(obj.Data)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	err = ioutil.WriteFile(obj.Object, decoded, 0700)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	f, err := os.Open(obj.Object)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	defer f.Close()
	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	wc := client.Bucket(obj.Bucket).Object(obj.Object).NewWriter(ctx)
	if _, err = io.Copy(wc, f); err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}
	if err := wc.Close(); err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	direktivapps.Respond(w, []byte{})
}
