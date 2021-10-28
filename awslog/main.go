package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/direktiv/direktiv-apps/pkg/direktivapps"
)

type AWSLogDetails struct {
	Key       string `json:"key"`
	Secret    string `json:"secret"`
	Region    string `json:"region"`
	LogGroup  string `json:"log-group"`
	LogStream string `json:"log-stream"`
	Message   string `json:"message"`
}

const code = "com.awslog.error"

func main() {
	direktivapps.StartServer(AWSLog)
}

func AWSLog(w http.ResponseWriter, r *http.Request) {
	obj := new(AWSLogDetails)
	_, err := direktivapps.Unmarshal(obj, r)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

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
		direktivapps.RespondWithError(w, code, err.Error())
		return
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
		direktivapps.RespondWithError(w, code, fmt.Sprintf("log stream: '%s' not found", obj.LogStream))
		return
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
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	if peo.RejectedLogEventsInfo != nil {
		direktivapps.RespondWithError(w, code, peo.RejectedLogEventsInfo.GoString())
		return
	}

	// write empty to notify when its finished
	direktivapps.Respond(w, []byte{})
}
