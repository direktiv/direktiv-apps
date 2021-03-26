package main

import (
	"encoding/base64"
	"io/ioutil"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/vorteil/direktiv-apps/pkg/direktivapps"
)

type AmazonUpload struct {
	Bucket string `json:"bucket"`
	Region string `json:"region"`
	Name   string `json:"upload-name"`
	Key    string `json:"key"`
	Secret string `json:"secret"`
	Data   string `json:"data"`
}

func main() {

	g := direktivapps.ActionError{
		ErrorCode:    "com.amazon-upload.error",
		ErrorMessage: "",
	}

	obj := new(AmazonUpload)
	direktivapps.ReadIn(obj, g)

	sess := session.New(&aws.Config{
		Region: aws.String(obj.Region),
	})

	decoded, err := base64.StdEncoding.DecodeString(obj.Data)
	if err != nil {
		g.ErrorMessage = err.Error()
		direktivapps.WriteError(g)
	}
	err = ioutil.WriteFile(obj.Name, decoded, 0700)
	if err != nil {
		g.ErrorMessage = err.Error()
		direktivapps.WriteError(g)
	}

	file, err := os.Open(obj.Name)
	if err != nil {
		g.ErrorMessage = err.Error()
		direktivapps.WriteError(g)
	}

	defer file.Close()

	uploader := s3manager.NewUploader(sess)

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(obj.Bucket),
		Key:    aws.String(obj.Name),
		Body:   file,
	})
	if err != nil {
		g.ErrorMessage = err.Error()
		direktivapps.WriteError(g)
	}

	direktivapps.WriteOut([]byte{}, g)
}
