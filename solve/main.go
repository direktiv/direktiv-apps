package main

import (
	"fmt"
	"net/http"

	"github.com/Knetic/govaluate"
	"github.com/vorteil/direktiv-apps/pkg/direktivapps"
)

var code = "com.solve.error"

const (
	dataIn   = "/direktiv-data/data.in"
	dataOut  = "/direktiv-data/data.out"
	errorOut = "/direktiv-data/error.json"
)

func Solve(w http.ResponseWriter, r *http.Request) {
	obj := make(map[string]interface{})
	_, err := direktivapps.Unmarshal(obj, r)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	var expressionString string

	if e, ok := obj["x"]; ok {
		expressionString, ok = e.(string)
		if !ok {
			direktivapps.RespondWithError(w, code, err.Error())
			return
		}
	} else {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	expression, err := govaluate.NewEvaluableExpression(expressionString)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	res, err := expression.Evaluate(nil)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	direktivapps.Respond(w, []byte(fmt.Sprintf("%v", res)))
}

func main() {
	direktivapps.StartServer(Solve)
}
