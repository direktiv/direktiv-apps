package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/vorteil/direktiv-apps/pkg/direktivapps"
)

type AWSSnsInput struct {
	Key      string `json:"key"`
	Secret   string `json:"secret"`
	Region   string `json:"region"`
	TopicARN string `json:"topic-arn"`
	Message  string `json:"message"`
}

func main() {
	g := direktivapps.ActionError{
		ErrorCode:    "com.amazon-sns.error",
		ErrorMessage: "",
	}

	obj := new(AWSSnsInput)
	direktivapps.ReadIn(obj, g)

	sess := session.New()

	client := sns.New(sess, &aws.Config{
		Credentials: credentials.NewStaticCredentials(obj.Key, obj.Secret, ""),
		Region:      aws.String(obj.Region),
	})

	_, err := client.Publish(&sns.PublishInput{
		TopicArn: aws.String(obj.TopicARN),
		Message:  aws.String(obj.Message),
	})
	if err != nil {
		g.ErrorMessage = err.Error()
		direktivapps.WriteError(g)
	}

	direktivapps.WriteOut([]byte{}, g)
}
