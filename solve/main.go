package main

import (
	"fmt"
	"net/http"

	"github.com/Knetic/govaluate"
	"github.com/direktiv/direktiv-apps/pkg/direktivapps"
)

var code = "com.solve.error"

const (
	dataIn   = "/direktiv-data/data.in"
	dataOut  = "/direktiv-data/data.out"
	errorOut = "/direktiv-data/error.json"
)

type SolveExpression struct {
	Expression string `json:"x"`
}

func Solve(w http.ResponseWriter, r *http.Request) {
	obj := new(SolveExpression)
	_, err := direktivapps.Unmarshal(obj, r)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	expression, err := govaluate.NewEvaluableExpression(obj.Expression)
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
