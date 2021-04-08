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


// Respond writes out to the responsewriter the json marshalled data
func Respond(w http.ResponseWriter, data []byte) {
	w.Write(data)
}

// Unmarshal reads the req body and unmarshals the data
func Unmarshal(obj interface{}, r *http.Request) (error) {

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	rdr := bytes.NewReader(data)
	dec := json.NewDecoder(rdr)

	dec.DisallowUnknownFields()

	err = dec.Decode(obj)
	if err != nil {
		return err
	}
}

// StartServer starts a new server
func StartServer(f func(w http.ResponseWriter, r *http.Request)) *http.Server {
	
	fmt.Println("Starting server")

	mux := http.NewServeMux()
	mux.HandleFunc("/", f)

	srv := &http.Server{
		Addr: ":8080",
		Handler: mux,
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM)

	go func(){
		<-sigs
		ShutDown(srv)
	}()

	srv.ListenAndServe()
}

// Shutdown turns off the server
func ShutDown(srv *http.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}

// Log sends a string to log via kubernetes
func Log(aid, l string) {
	http.Post(fmt.Sprintf("http://localhost:8889/log?aid=%s", aid), "plain/text", strings.NewReader(l))
}


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
