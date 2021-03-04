package main

import (
	"encoding/json"
	"fmt"

	"github.com/vorteil/direktiv-apps/pkg/direktivapps"

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

func main() {
	g := direktivapps.ActionError{
		ErrorCode:    "com.tweet.error",
		ErrorMessage: "",
	}
	obj := new(TwitterDetails)
	direktivapps.ReadIn(obj, g)

	config := oauth1.NewConfig(obj.ConsumerKey, obj.ConsumerSecret)
	token := oauth1.NewToken(obj.TokenKey, obj.TokenSecret)

	httpClient := config.Client(oauth1.NoContext, token)

	client := twitter.NewClient(httpClient)

	resp, tweet, err := client.Statuses.Update(obj.Message, nil)
	if err != nil {
		g.ErrorMessage = err.Error()
		direktivapps.WriteError(g)
	}

	if tweet.StatusCode < 200 || tweet.StatusCode >= 300 {
		g.ErrorMessage = fmt.Sprintf("Response Message: %s, Response Code: %v", tweet.Status, tweet.StatusCode)
		direktivapps.WriteError(g)
	}

	bv, err := json.Marshal(resp)
	if err != nil {
		g.ErrorMessage = fmt.Sprintf("unable to unmarshal: %v", err)
		direktivapps.WriteError(g)
	}

	direktivapps.WriteOut(bv, g)
}
