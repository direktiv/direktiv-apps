package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/vorteil/direktiv-apps/pkg/direktivapps"
	"golang.org/x/text/language"
	"google.golang.org/api/option"

	"cloud.google.com/go/translate"
)

const credFile = "/creds"

type InputTranslatorGoogle struct {
	TargetLanguage    string `json:"target-language"`
	Message           string `json:"message"`
	ServiceAccountKey string `json:"serviceAccountKey"`
}

type OutputMessage struct {
	Message string `json:"message"`
}

const code = "com.google-translator.error"

func main() {
	direktivapps.StartServer(GoogleTranslate)
}

func GoogleTranslate(w http.ResponseWriter, r *http.Request) {

	obj := new(InputTranslatorGoogle)
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

	lang, err := language.Parse(obj.TargetLanguage)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	client, err := translate.NewClient(ctx, option.WithCredentialsFile(credFile))
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}
	defer client.Close()

	resp, err := client.Translate(ctx, []string{obj.Message}, lang, &translate.Options{
		Format: "text",
	})
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}
	if len(resp) == 0 {
		direktivapps.RespondWithError(w, code, "Translate returned empty response.")
		return
	}

	var output OutputMessage
	direktivapps.Log(aid, fmt.Sprintf("Translated: %v\n", resp[0].Text))
	output.Message = resp[0].Text

	data, err := json.Marshal(output)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	direktivapps.Respond(w, data)
}
