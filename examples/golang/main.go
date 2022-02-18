package main

import (
	"context"
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

	// greeter := new(Greeter)
	aid := r.Header.Get(DirektivActionIDHeader)

	log(aid, "Reading Input")
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respondWithErr(w, fmt.Sprintf(code, "readdata"), err.Error())
		return
	}

	log(aid, fmt.Sprintf("DATA %s", string(data)))
	log(aid, fmt.Sprintf("DATA2 %s", "helloworldhelloworldhelloworldhelloworldhelloworldhelloworldhelloworld12345"))
	log(aid, fmt.Sprintf("TEMPDIR %s", r.Header.Get("Direktiv-TempDir")))

	// POST http://localhost:8889/var?aid=<EXAMPLE>&scope=instance&key=myFiles

	// Body: <VARIABLE DATA>

	// url := fmt.Sprintf("http://localhost:8889/var?aid=%s&scope=workflow&key=HELLO", aid)

	// var jsonStr = []byte(`{"title":"Buy cheese and bread for breakfast."}`)
	// req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	// if err != nil {
	// 	log(aid, fmt.Sprintf("ERR %s", err.Error()))
	// 	return
	// }

	// client := &http.Client{}
	// resp, err := client.Do(req)
	// if err != nil {
	// 	log(aid, fmt.Sprintf("ERR %s", err.Error()))
	// 	return
	// }
	// defer resp.Body.Close()

	// files, err := ioutil.ReadDir(fmt.Sprintf("%s/out", r.Header.Get("Direktiv-TempDir")))
	// if err != nil {
	// 	log(aid, fmt.Sprintf("ERR %s", err.Error()))
	// }

	// for _, f := range files {
	// 	log(aid, f.Name())
	// }

	// file := fmt.Sprintf("%s/out/namespace/testme", r.Header.Get("Direktiv-TempDir"))

	// f, err := os.Create(file)
	// if err != nil {
	// 	log(aid, fmt.Sprintf("ERR %s", err.Error()))
	// 	return
	// }

	// f.Write([]byte("HELLO2!!!"))

	// // registries1.png

	// gg, err := ioutil.ReadFile("/registries1.png")
	// if err != nil {
	// 	log(aid, fmt.Sprintf("ERR %s", err.Error()))
	// 	return
	// }

	// file2 := fmt.Sprintf("%s/out/namespace/whatever", r.Header.Get("Direktiv-TempDir"))
	// f2, err := os.Create(file2)
	// if err != nil {
	// 	log(aid, fmt.Sprintf("ERR %s", err.Error()))
	// 	return
	// }
	// f2.Write(gg)
	// err = os.Rename("/registries1.png", fmt.Sprintf("%s/out/workflow/whatever", r.Header.Get("Direktiv-TempDir")))
	// if err != nil {
	// 	log(aid, fmt.Sprintf("ERR %s", err.Error()))
	// 	return
	// }

	// /mnt/shared/example/out/workflow

	// rdr := bytes.NewReader(data)
	// dec := json.NewDecoder(rdr)

	// dec.DisallowUnknownFields()

	// log(aid, "Decoding Input")
	// err = dec.Decode(greeter)
	// if err != nil {
	// 	respondWithErr(w, fmt.Sprintf(code, "decode"), err.Error())
	// 	return
	// }

	// var output ReturnMessage
	// output.Greeting = fmt.Sprintf("Welcome to Direktiv2, %s!", greeter.Name)

	// marshalBytes, err := json.Marshal(output)
	// if err != nil {
	// 	respondWithErr(w, fmt.Sprintf(code, "marshal"), err.Error())
	// 	return
	// }

	// log(aid, "Writing Output")
	// respond(w, marshalBytes)

	respond(w, []byte("{}"))
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
