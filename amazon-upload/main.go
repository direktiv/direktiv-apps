package main

import (
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"os"
	"path"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/direktiv/direktiv-apps/pkg/direktivapps"
)

type AmazonUpload struct {
	Bucket   string `json:"bucket"`
	Region   string `json:"region"`
	Name     string `json:"upload-name"`
	Key      string `json:"key"`
	Secret   string `json:"secret"`
	Data     string `json:"data"`
	FileName string `json:"filename"` // filename gets used pulling from the temp dir
}

const code = "com.amazon-upload.error"

func main() {
	direktivapps.StartServer(Upload)
}

func Upload(w http.ResponseWriter, r *http.Request) {
	obj := new(AmazonUpload)
	_, err := direktivapps.Unmarshal(obj, r)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	os.Setenv("AWS_ACCESS_KEY_ID", obj.Key)
	os.Setenv("AWS_SECRET_ACCESS_KEY", obj.Secret)

	var file *os.File
	if obj.Data != "" {
		decoded, err := base64.StdEncoding.DecodeString(obj.Data)
		if err != nil {
			direktivapps.RespondWithError(w, code, err.Error())
			return
		}
		err = ioutil.WriteFile(obj.Name, decoded, 0700)
		if err != nil {
			direktivapps.RespondWithError(w, code, err.Error())
			return
		}

		file, err = os.Open(obj.Name)
		if err != nil {
			direktivapps.RespondWithError(w, code, err.Error())
			return
		}
	} else {
		// check a variable
		filePath := path.Join(r.Header.Get("Direktiv-TempDir"), obj.FileName)
		file, err = os.Open(filePath)
		if err != nil {
			direktivapps.RespondWithError(w, code, err.Error())
			return
		}
	}
	defer file.Close()

	sess := session.New(&aws.Config{
		Region: aws.String(obj.Region),
	})

	uploader := s3manager.NewUploader(sess)

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(obj.Bucket),
		Key:    aws.String(obj.Name),
		Body:   file,
	})
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	direktivapps.Respond(w, []byte{})
}
