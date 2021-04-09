package main

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/vorteil/direktiv-apps/pkg/direktivapps"

	"github.com/aws/aws-sdk-go/service/lambda"
)

// LambdaInvocation struct to invoke a cloud function on aws.
type LambdaInvocation struct {
	Key      string                 `json:"key"`
	Secret   string                 `json:"secret"`
	Region   string                 `json:"region"`
	Function string                 `json:"function"`
	Body     map[string]interface{} `json:"body"`
}

const code = "com.lambdainvoke.error"

func main() {
	direktivapps.StartServer(Lambda)
}

func Lambda(w http.ResponseWriter, r *http.Request) {
	obj := new(LambdaInvocation)
	_, err := direktivapps.Unmarshal(obj, r)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	// Start new aws session using default authentication handlers
	sess := session.New()

	client := lambda.New(sess, &aws.Config{
		Credentials:      credentials.NewStaticCredentials(obj.Key, obj.Secret, ""),
		S3ForcePathStyle: aws.Bool(true),
		Region:           aws.String(obj.Region),
	})

	var payload []byte
	if obj.Body != nil {
		payload, err = json.Marshal(obj.Body)
		if err != nil {
			direktivapps.RespondWithError(w, code, err.Error())
			return
		}
	}

	result, err := client.Invoke(&lambda.InvokeInput{FunctionName: aws.String(obj.Function), Payload: payload})
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	direktivapps.Respond(w, []byte(result.Payload))
}
