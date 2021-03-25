package direktivapps

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// ActionError is a struct Direktiv uses to report application errors.
type ActionError struct {
	ErrorCode    string `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
}

const outPath = "/direktiv-data/data.out"
const dataInPath = "/direktiv-data/data.in"
const errorPath = "/direktiv-data/error.json"

// ReadIn reads data from dataInPath and returns struct provided with json fields
func ReadIn(obj interface{}, g ActionError) {
	f, err := os.Open(dataInPath)
	if err != nil {
		g.ErrorMessage = err.Error()
		WriteError(g)
	}

	defer f.Close()

	dec := json.NewDecoder(f)
	dec.DisallowUnknownFields()

	err = dec.Decode(obj)
	if err != nil {
		g.ErrorMessage = err.Error()
		WriteError(g)
	}
}

// WriteError writes an error to errorPath
func WriteError(g ActionError) {
	b, _ := json.Marshal(g)
	ioutil.WriteFile(errorPath, b, 0755)
	os.Exit(0)
}

// WriteOut writes out data to outPath
func WriteOut(by []byte, g ActionError) {
	var err error
	err = ioutil.WriteFile(outPath, by, 0755)
	if err != nil {
		g.ErrorMessage = err.Error()
		WriteError(g)
	}
	os.Exit(0)
}
