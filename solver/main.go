package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Knetic/govaluate"
)

const (
	dataIn   = "/direktiv-data/data.in"
	dataOut  = "/direktiv-data/data.out"
	errorOut = "/direktiv-data/error.json"
)

func main() {

	f, err := os.Open(dataIn)
	if err != nil {
		writeError("error.disk.read", err.Error())
	}
	defer f.Close()

	input := make(map[string]interface{})

	dec := json.NewDecoder(f)
	err = dec.Decode(&input)
	if err != nil {
		writeError("error.input.parse", err.Error())
	}

	var expressionString string

	if e, ok := input["x"]; ok {

		// e2, ok := e.([]interface{})
		// if !ok {
		// 	writeError("error.input.parse", fmt.Sprintf("expressions field was of type %s", reflect.TypeOf(e)))
		// }

		// for _, exp := range e2 {
		expressionString, ok = exp.(string)
		if !ok {
			writeError("error.input.parse", fmt.Sprintf("index of expressions array was not of type string"))
		}

		// expressions = append(expressions, e3)
		// }

	} else {
		writeError("error.input.parse", err.Error())
	}

	// out := make([]string, 0)
	// for _, exp := range expressions {
	expression, err := govaluate.NewEvaluableExpression(expressionString)
	if err != nil {
		writeError("error.math.invalid", err.Error())
	}

	res, err := expression.Evaluate(nil)
	if err != nil {
		writeError("error.math.eval", err.Error())
	}

	// out = append(out, fmt.Sprintf("%v", res))
	// }

	writeOut(fmt.Sprintf("%v", res))
}

func writeOut(x interface{}) {

	f, err := os.OpenFile(dataOut, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		writeError("error.disk.write", err.Error())
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	err = enc.Encode(x)
	if err != nil {
		writeError("error.output.encode", err.Error())
	}

}

func writeError(code, msg string) {

	out := make(map[string]string)
	out["errorCode"] = code
	out["errorMessage"] = msg

	f, err := os.OpenFile(errorOut, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	err = enc.Encode(out)
	if err != nil {
		panic(err)
	}

	os.Exit(1)
}
