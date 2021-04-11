package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/vorteil/direktiv-apps/pkg/direktivapps"
	"google.golang.org/api/option"

	language "cloud.google.com/go/language/apiv1"
	languagepb "google.golang.org/genproto/googleapis/cloud/language/v1"
)

const credFile = "/creds"

type InputSentimentAnalysis struct {
	ServiceAccountKey string `json:"serviceAccountKey"`
	Message           string `json:"message"`
}
type SentimentAnalysis struct {
	Feeling   string  `json:"feeling"`
	Score     float32 `json:"score"`
	Magnitude float32 `json:"magnitude"`
}

const code = "com.google-sentiment-check.error"

func main() {
	direktivapps.StartServer(GoogleSentimentCheck)
}

func GoogleSentimentCheck(w http.ResponseWriter, r *http.Request) {
	obj := new(InputSentimentAnalysis)
	aid, err := direktivapps.Unmarshal(obj, r)
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

	client, err := language.NewClient(ctx, option.WithCredentialsFile(credFile))
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	// Detects the sentiment of the text.
	sentiment, err := client.AnalyzeSentiment(ctx, &languagepb.AnalyzeSentimentRequest{
		Document: &languagepb.Document{
			Source: &languagepb.Document_Content{
				Content: obj.Message,
			},
			Type: languagepb.Document_PLAIN_TEXT,
		},
		EncodingType: languagepb.EncodingType_UTF8,
	})
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	direktivapps.Log(aid, fmt.Sprintf("Text being checked: %v\n", obj.Message))

	var outputSentiment SentimentAnalysis
	outputSentiment.Magnitude = sentiment.DocumentSentiment.Magnitude
	outputSentiment.Score = sentiment.DocumentSentiment.Score

	if sentiment.DocumentSentiment.Score >= 0 && sentiment.DocumentSentiment.Score < 0.5 {
		outputSentiment.Feeling = "Somewhat Positive"
		if sentiment.DocumentSentiment.Magnitude == 0 {
			outputSentiment.Feeling = "Somewhat Positive/Neutral"
		}
	} else if sentiment.DocumentSentiment.Score >= 0.6 {
		outputSentiment.Feeling = "Positive"
	} else if sentiment.DocumentSentiment.Score < 0 {
		if sentiment.DocumentSentiment.Score >= -0.5 {
			outputSentiment.Feeling = "Somewhat Negative"
			if sentiment.DocumentSentiment.Magnitude == 0 {
				outputSentiment.Feeling = "Somewhat Negative/Neutral"
			}
		} else {
			outputSentiment.Feeling = "Negative"
		}
	}

	data, err := json.Marshal(outputSentiment)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	direktivapps.Respond(w, data)
}
