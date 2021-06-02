package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"

	"github.com/vorteil/direktiv-apps/pkg/direktivapps"
)

const code = "com.ubuntu.error"

type shell struct {
	Script string   `json:"script"`
	Args   []string `json:"args"`
}

func copyFile(path string) (*os.File, error) {

	in, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	out, err := ioutil.TempFile("", "exe")
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(out, in)
	if err != nil {
		os.RemoveAll(out.Name())
		return nil, err
	}

	out.Sync()

	err = os.Chmod(out.Name(), 0755)
	if err != nil {
		os.RemoveAll(out.Name())
		return nil, err
	}

	out.Close()

	return out, nil
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

	log(aid, "getting script")

	dtmp := r.Header.Get("Direktiv-TempDir")
	if len(dtmp) == 0 {
		direktivapps.RespondWithError(w, code, fmt.Sprintf("direktiv tmp direktory not provided"))
		return
	}

	f := path.Join(dtmp, s.Script)
	rf, err := copyFile(f)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}
	defer os.Remove(rf.Name())

	log(aid, fmt.Sprintf("found script %v -> %v", f, rf.Name()))

	ret, err := execScript(rf.Name(), s.Args)
	if err != nil {
		direktivapps.RespondWithError(w, code, fmt.Sprintf("%v: %s", err, string(ret)))
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
