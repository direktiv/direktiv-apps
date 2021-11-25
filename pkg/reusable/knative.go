package reusable

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
)

const (
	DirektivActionIDHeader     = "Direktiv-ActionID"
	DirektivErrorCodeHeader    = "Direktiv-ErrorCode"
	DirektivErrorMessageHeader = "Direktiv-ErrorMessage"
)

const (
	UnmarshallError = "io.direktiv.unmarshal"
	MarshallError   = "io.direktiv.marshal"
)

type ActionError struct {
	ErrorCode    string `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
}

var sm sync.Map

func aid(r *http.Request) (string, error) {
	aid := r.Header.Get(DirektivActionIDHeader)
	if aid == "" {
		return "", fmt.Errorf("no Direktiv-ActionID header set")
	}
	return aid, nil
}

// Start Server
func StartServer(f func(w http.ResponseWriter, r *http.Request, ri *RequestInfo)) {

	logger := getZeroLogger(nil)
	logger.Info().Msg("starting server")

	r := mux.NewRouter()
	r.HandleFunc("/", cancelHandler).Methods(http.MethodDelete)
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		aid, err := aid(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Direktiv-ActionID header is required to cancel the instance"))
			return
		}

		ctx, cancel := context.WithCancel(r.Context())
		r = r.WithContext(ctx)
		sm.Store(aid, cancel)
		f(w, r, newRequestInfo(aid))
		sm.Delete(aid)
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	handleShutdown(srv)
	log.Fatal(srv.ListenAndServe())

}

func handleShutdown(srv *http.Server) {

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM)

	go func() {
		<-sigs
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		srv.Shutdown(ctx)
	}()

}

func cancelHandler(w http.ResponseWriter, r *http.Request) {

	aid := r.Header.Get(DirektivActionIDHeader)
	if aid == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Direktiv-ActionID header is required to cancel the instance"))
		return
	}

	cancel, ok := sm.Load(aid)
	if !ok {
		return
	}

	cf, ok := cancel.(context.CancelFunc)
	if !ok {
		return
	}

	cf()
	sm.Delete(aid)
}

func ReportResult(w http.ResponseWriter, data interface{}) {

	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		ReportError(w, MarshallError, err)
	}
	w.Write(b)

}

func ReportError(w http.ResponseWriter, errcode string,
	err error) {

	w.Header().Set(DirektivErrorCodeHeader, errcode)
	w.Header().Set(DirektivErrorMessageHeader, err.Error())

	ae := &ActionError{
		ErrorCode:    errcode,
		ErrorMessage: err.Error(),
	}

	b, _ := json.MarshalIndent(ae, "", "  ")
	w.Write(b)

}

// Unmarshal reads the req body and unmarshals the data
func Unmarshal(obj interface{}, strict bool, r *http.Request) error {

	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	rdr := bytes.NewReader(data)
	dec := json.NewDecoder(rdr)

	if strict {
		dec.DisallowUnknownFields()
	}

	err = dec.Decode(obj)
	if err != nil {
		return err
	}

	return nil
}
