package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/vorteil/direktiv-apps/pkg/direktivapps"
	"golang.org/x/text/language"
	"google.golang.org/api/option"

	"cloud.google.com/go/translate"
)

const credFile = "/tmp/creds"

type InputTranslatorGoogle struct {
	TargetLanguage    string `json:"target-language"`
	Message           string `json:"message"`
	ServiceAccountKey string `json:"serviceAccountKey"`
}

type OutputMessage struct {
	Message string `json:"message"`
}

func main() {

	g := direktivapps.ActionError{
		ErrorCode:    "com.google-translator.error",
		ErrorMessage: "",
	}

	obj := new(InputTranslatorGoogle)
	direktivapps.ReadIn(obj, g)

	ctx := context.Background()

	err := ioutil.WriteFile(credFile, []byte(obj.ServiceAccountKey), 0777)
	if err != nil {
		g.ErrorMessage = err.Error()
		direktivapps.WriteError(g)
	}

	lang, err := language.Parse(obj.TargetLanguage)
	if err != nil {
		g.ErrorMessage = err.Error()
		direktivapps.WriteError(g)
	}

	client, err := translate.NewClient(ctx, option.WithCredentialsFile(credFile))
	if err != nil {
		g.ErrorMessage = err.Error()
		direktivapps.WriteError(g)
	}
	defer client.Close()

	resp, err := client.Translate(ctx, []string{obj.Message}, lang, &translate.Options{
		Format: "text",
	})
	if err != nil {
		g.ErrorMessage = err.Error()
		direktivapps.WriteError(g)
	}
	if len(resp) == 0 {
		g.ErrorMessage = "Translate returned empty response."
		direktivapps.WriteError(g)
	}

	var output OutputMessage
	fmt.Printf("Translated: %v\n", resp[0].Text)
	output.Message = resp[0].Text

	data, err := json.Marshal(output)
	if err != nil {
		g.ErrorMessage = err.Error()
		direktivapps.WriteError(g)
	}

	direktivapps.WriteOut(data, g)
}
