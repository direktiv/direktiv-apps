package main

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/vorteil/direktiv-apps/pkg/direktivapps"
)

type AWSLogDetails struct {
	Key       string `json:"key"`
	Secret    string `json:"secret"`
	Region    string `json:"region"`
	LogGroup  string `json:"log-group"`
	LogStream string `json:"log-stream"`
	Message   string `json:"message"`
}

func main() {

	g := direktivapps.ActionError{
		ErrorCode:    "com.awslog.error",
		ErrorMessage: "",
	}

	obj := new(AWSLogDetails)
	direktivapps.ReadIn(obj, g)

	sess := session.New()

	client := cloudwatchlogs.New(sess, &aws.Config{
		Credentials: credentials.NewStaticCredentials(obj.Key, obj.Secret, ""),
		Region:      aws.String(obj.Region),
	})

	// describe log stream to get sequence token
	streams, err := client.DescribeLogStreams(&cloudwatchlogs.DescribeLogStreamsInput{
		LogGroupName: aws.String(obj.LogGroup),
	})
	if err != nil {
		g.ErrorMessage = err.Error()
		direktivapps.WriteError(g)
	}

	// get the upload sequence token as we're only writing one log
	var uploadSequenceToken *string
	var found bool
	for _, stream := range streams.LogStreams {
		if *stream.LogStreamName == obj.LogStream {
			uploadSequenceToken = stream.UploadSequenceToken
			found = true
			break
		}
	}

	if !found {
		g.ErrorMessage = fmt.Sprintf("log stream: '%s' not found", obj.LogStream)
		direktivapps.WriteError(g)
	}

	logs := make([]*cloudwatchlogs.InputLogEvent, 0)
	logs = append(logs, &cloudwatchlogs.InputLogEvent{
		Message:   aws.String(obj.Message),
		Timestamp: aws.Int64(time.Now().Unix() * 1000), // Needs to be microseconds
	})

	peo, err := client.PutLogEvents(&cloudwatchlogs.PutLogEventsInput{
		LogGroupName:  aws.String(obj.LogGroup),
		LogStreamName: aws.String(obj.LogStream),
		SequenceToken: uploadSequenceToken,
		LogEvents:     logs,
	})
	if err != nil {
		g.ErrorMessage = err.Error()
		direktivapps.WriteError(g)
	}

	if peo.RejectedLogEventsInfo != nil {
		g.ErrorMessage = peo.RejectedLogEventsInfo.GoString()
		direktivapps.WriteError(g)
	}

	// write empty to notify when its finished
	direktivapps.WriteOut([]byte{}, g)
}
