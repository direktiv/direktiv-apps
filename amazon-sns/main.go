package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/direktiv/direktiv-apps/pkg/direktivapps"
	"net/http"
)

type AWSSnsInput struct {
	Key      string `json:"key"`
	Secret   string `json:"secret"`
	Region   string `json:"region"`
	TopicARN string `json:"topic-arn"`
	Message  string `json:"message"`
}

const code = "com.amazon-sns.error"

func main() {
	direktivapps.StartServer(AmazonSNS)
}

func AmazonSNS(w http.ResponseWriter, r *http.Request) {
	obj := new(AWSSnsInput)
	_, err := direktivapps.Unmarshal(obj, r)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	sess := session.New()

	client := sns.New(sess, &aws.Config{
		Credentials: credentials.NewStaticCredentials(obj.Key, obj.Secret, ""),
		Region:      aws.String(obj.Region),
	})

	_, err = client.Publish(&sns.PublishInput{
		TopicArn: aws.String(obj.TopicARN),
		Message:  aws.String(obj.Message),
	})
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	direktivapps.Respond(w, []byte{})
}
