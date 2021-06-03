package main

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/vorteil/direktiv-apps/service-now/util"

	da "github.com/vorteil/direktiv-apps/pkg/direktivapps"
)

func main() {
	da.StartServer(coreLogic)
}

type input struct {
	util.ClientAuth
}

func coreLogic(w http.ResponseWriter, r *http.Request) {

	in, err := readInput(r.Body)
	if err != nil {
		// handle it
	}

}

func readInput(r io.ReadCloser) (*input, error) {

	in := new(input)
	err := json.NewDecoder(r).Decode(in)
	if err != nil {
		return nil, err
	}

	return in, nil
}
