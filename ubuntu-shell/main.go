package main

import (
	"encoding/json"
	"fmt"
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

func execScript(path string, args []string, envs []string) ([]byte, error) {

	cmd := exec.Command(path, args...)
	cmd.Env = envs

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

	if s.Script == "" {
		direktivapps.RespondWithError(w, code, "no script provided")
		return
	}

	direktivapps.LogDouble(aid, "getting script")

	dtmp := r.Header.Get("Direktiv-TempDir")
	if len(dtmp) == 0 {
		direktivapps.RespondWithError(w, code, fmt.Sprintf("direktiv tmp directory not provided"))
		return
	}

	f := path.Join(dtmp, s.Script)

	err = os.Chmod(f, 0755)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	direktivapps.LogDouble(aid, fmt.Sprintf("found script %v", f))

	envs := []string{fmt.Sprintf("Direktiv_TempDir=%s",
		r.Header.Get("Direktiv-TempDir"))}

	ret, err := execScript(f, s.Args, envs)
	if err != nil {
		direktivapps.RespondWithError(w, code, fmt.Sprintf("%v: %s", err, string(ret)))
		return
	}

	direktivapps.LogDouble(aid, fmt.Sprintf("script return: %v", string(ret)))

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
