package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
)

func main() {

	setupLog()

	m, err := getData()
	if err != nil {
		log.Fatalf(err.Error())
	}

	log.Printf("data in: %v", m)

	if _, ok := m["error"]; ok {

		log.Printf("error data in payload")

		error := `{
				"code": "badInput",
				"message": "don't send me error messages"
		}`

		err := ioutil.WriteFile("/direktiv-data/error.json", []byte(error), 0644)
		if err != nil {
			log.Fatalf(err.Error())
		}

		os.Create("/direktiv-data/done")

	}

	// prepare response
	resp := make(map[string]string)
	resp["response"] = "Hello World"

	b, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf(err.Error())
	}

	err = ioutil.WriteFile("/direktiv-data/out.log", b, 0644)
	if err != nil {
		log.Fatalf(err.Error())
	}

	os.Create("/direktiv-data/done")

}

func setupLog() {

	lf, err := os.OpenFile("/direktiv-data/out.log",
		os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		log.Fatalf(err.Error())
	}

	mw := io.MultiWriter(os.Stdout, lf)
	log.SetOutput(mw)

}

func getData() (map[string]interface{}, error) {

	var data map[string]interface{}

	in, err := os.Open("/direktiv-data/input.json")
	if err != nil {
		log.Printf("error open input: %v", err)
		return data, err
	}

	b, err := ioutil.ReadAll(in)
	if err != nil {
		log.Printf("error reading input: %v", err)
		return data, err
	}

	err = json.Unmarshal(b, &data)
	if err != nil {
		log.Printf("error marshal input: %v", err)
		return data, err
	}

	return data, nil
}
