package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/direktiv/direktiv-apps/pkg/direktivapps"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

// TwitterDetails is a struct that takes the input.
type TwitterDetails struct {
	ConsumerKey    string `json:"consumerKey"`
	ConsumerSecret string `json:"consumerSecret"`
	TokenKey       string `json:"tokenKey"`
	TokenSecret    string `json:"tokenSecret"`
	Message        string `json:"message"`
}

var code = "com.tweet.error"

func Tweet(w http.ResponseWriter, r *http.Request) {
	obj := new(TwitterDetails)

	_, err := direktivapps.Unmarshal(obj, r)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	config := oauth1.NewConfig(obj.ConsumerKey, obj.ConsumerSecret)
	token := oauth1.NewToken(obj.TokenKey, obj.TokenSecret)

	httpClient := config.Client(oauth1.NoContext, token)

	client := twitter.NewClient(httpClient)

	resp, tweet, err := client.Statuses.Update(obj.Message, nil)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	if tweet.StatusCode < 200 || tweet.StatusCode >= 300 {
		direktivapps.RespondWithError(w, code, fmt.Sprintf(fmt.Sprintf("Response Message: %s, Response Code: %v", tweet.Status, tweet.StatusCode)))
		return
	}

	bv, err := json.Marshal(resp)
	if err != nil {
		direktivapps.RespondWithError(w, code, fmt.Sprintf("unable to unmarshal: %v", err))
		return
	}

	direktivapps.Respond(w, bv)
}

func main() {
	direktivapps.StartServer(Tweet)
}
