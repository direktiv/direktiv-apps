package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type Greeter struct {
	Name string `json:"name"`
}

type ReturnMessage struct {
	Greeting string `json:"greeting"`
}

const (
	DirektivActionIDHeader = "Direktiv-ActionID"

	DirektivErrorCodeHeader    = "Direktiv-ErrorCode"
	DirektivErrorMessageHeader = "Direktiv-ErrorMessage"
)

const code = "com.greeting-%s.error"

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/", GreetingHandler)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM)

	go func() {
		<-sigs
		shutdown(srv)
	}()

	srv.ListenAndServe()
}

func GreetingHandler(w http.ResponseWriter, r *http.Request) {
	greeter := new(Greeter)
	aid := r.Header.Get(DirektivActionIDHeader)

	log(aid, "Reading Input")
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respondWithErr(w, fmt.Sprintf(code, "readdata"), err.Error())
		return
	}

	rdr := bytes.NewReader(data)
	dec := json.NewDecoder(rdr)

	dec.DisallowUnknownFields()

	log(aid, "Decoding Input")
	err = dec.Decode(greeter)
	if err != nil {
		respondWithErr(w, fmt.Sprintf(code, "decode"), err.Error())
		return
	}

	var output ReturnMessage
	output.Greeting = fmt.Sprintf("Welcome to Direktiv, %s!", greeter.Name)

	marshalBytes, err := json.Marshal(output)
	if err != nil {
		respondWithErr(w, fmt.Sprintf(code, "marshal"), err.Error())
		return
	}

	log(aid, "Writing Output")
	respond(w, marshalBytes)
}

func shutdown(srv *http.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}

func log(aid, l string) {
	if aid == "development" || aid == "Development" {
		fmt.Println(l)
	} else {
		http.Post(fmt.Sprintf("http://localhost:8889/log?aid=%s", aid), "plain/text", strings.NewReader(l))
	}
}

func respond(w http.ResponseWriter, data []byte) {
	w.Write(data)
}

func respondWithErr(w http.ResponseWriter, code, err string) {
	w.Header().Set(DirektivErrorCodeHeader, code)
	w.Header().Set(DirektivErrorMessageHeader, err)
}
