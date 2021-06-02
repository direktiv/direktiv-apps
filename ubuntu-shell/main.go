package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"

	"github.com/vorteil/direktiv-apps/pkg/direktivapps"
)

const code = "com.ubuntu.error"

type shell struct {
	Script string   `json:"script"`
	Args   []string `json:"args"`
}

func createScript(s string) (*os.File, error) {

	h, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}

	// create temp file and run it
	file, err := ioutil.TempFile("", "sh")
	if err != nil {
		return nil, err
	}

	_, err = file.Write([]byte(h))
	if err != nil {
		os.Remove(file.Name())
		return nil, err
	}
	file.Close()

	err = os.Chmod(file.Name(), 0755)
	if err != nil {
		os.Remove(file.Name())
		return nil, err
	}

	return file, nil

}

func execScript(path string, args []string) ([]byte, error) {

	cmd := exec.Command(path, args...)
	d, err := cmd.Output()
	if err != nil {
		if e, ok := err.(*exec.ExitError); ok {
			d = e.Stderr
		}
		return d, err
	}

	return d, nil

}

func request(w http.ResponseWriter, r *http.Request) {

	var s shell
	aid, err := direktivapps.Unmarshal(&s, r)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	log(aid, "create script")

	// get the shell script
	file, err := createScript(s.Script)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}
	defer os.Remove(file.Name())

	log(aid, fmt.Sprintf("created script %v", file.Name()))

	ret, err := execScript(file.Name(), s.Args)
	if err != nil {
		direktivapps.RespondWithError(w, code, string(ret))
		return
	}

	log(aid, fmt.Sprintf("script return: %v", string(ret)))

	// check if base64
	var j map[string]interface{}
	err = json.Unmarshal(ret, &j)
	if err != nil {
		o := make(map[string]string)
		o["output"] = string(ret)

		b, err := json.Marshal(&o)
		if err != nil {
			direktivapps.RespondWithError(w, code, string(ret))
			return
		}

		direktivapps.Respond(w, b)
	} else {
		direktivapps.Respond(w, ret)
	}

}

func main() {
	direktivapps.StartServer(request)
}

func log(aid, l string) {
	fmt.Println(l)
	direktivapps.Log(aid, l)
}
