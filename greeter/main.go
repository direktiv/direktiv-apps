package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

type Greeter struct {
	Name string `json:"name"`
}

type ReturnGreeting struct {
	Greeting string `json:"greeting"`
}

type ActionError struct {
	ErrorCode    string `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
}

func main() {
	var gr Greeter
	g := ActionError{
		ErrorCode:    "com.greeting.error",
		ErrorMessage: "",
	}

	var err error
	var data []byte

	log.Printf("Reading in Data...")

	data, err = ioutil.ReadFile("/direktiv-data/data.in")
	if err != nil {
		g.ErrorMessage = err.Error()
		writeError(g)
		return
	}

	err = json.Unmarshal(data, &gr)
	if err != nil {
		g.ErrorMessage = err.Error()
		writeError(g)
		return
	}

	var rg ReturnGreeting
	rg.Greeting = fmt.Sprintf("Welcome to Direktiv, %s!", gr.Name)

	bv, err := json.Marshal(rg)
	if err != nil {
		g.ErrorMessage = err.Error()
		writeError(g)
		return
	}

	finishRunning(bv, g)
}

// finishRunning
func finishRunning(eb []byte, g ActionError) {
	var err error
	err = ioutil.WriteFile("/direktiv-data/data.out", eb, 0755)
	if err != nil {
		g.ErrorMessage = err.Error()
		writeError(g)
		return
	}
}

// writeError
func writeError(g ActionError) {
	b, _ := json.Marshal(g)
	ioutil.WriteFile("/direktiv-data/error.json", b, 0755)
}
