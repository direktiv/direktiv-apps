package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	da "github.com/direktiv/direktiv-apps/pkg/direktivapps"
)

func main() {
	da.StartServer(coreLogic)
}

type input struct {
	HostKeyChecking bool     `json:"hostKeyChecking"`
	Playbook        string   `json:"playbook"`
	PrivateKey      string   `json:"privateKey"`
	Args            []string `json:"args"`
	Envs            []string `json:"envs"`
}

const (
	errCode    = "com.ansible.%s"
	ansibleBin = "ansible-playbook"
)

func reportError(w http.ResponseWriter, code string, err error) {
	da.RespondWithError(w, fmt.Sprintf(errCode, code), err.Error())
}

func coreLogic(w http.ResponseWriter, r *http.Request) {

	obj := new(input)
	aid, err := da.Unmarshal(obj, r)
	if err != nil {
		reportError(w, "inputUnmarshal", err)
		return
	}

	mntPath := r.Header.Get("Direktiv-TempDir")
	if len(mntPath) == 0 {
		reportError(w, "tmpDir", fmt.Errorf("direktiv tmp directory not provided in Direktiv-TempDir header"))
		return
	}

	if obj.Playbook == "" {
		reportError(w, "badRequest", fmt.Errorf("name of playbook file must be provided in input data"))
		return
	}

	if obj.PrivateKey == "" {
		reportError(w, "privateKey", fmt.Errorf("name of private key file must be provided in input data"))
		return
	}

	pkPath := filepath.Join(mntPath, obj.PrivateKey)
	err = os.Chmod(pkPath, 0600)
	if err != nil {
		reportError(w, "chmod", err)
		return
	}

	if obj.Args == nil {
		obj.Args = make([]string, 0)
	}

	obj.Args = append(obj.Args, "--private-key", pkPath, filepath.Join(mntPath, obj.Playbook))

	cmd := exec.Command(ansibleBin, obj.Args...)
	if !obj.HostKeyChecking {
		cmd.Env = append(cmd.Env, "ANSIBLE_HOST_KEY_CHECKING=False")
	}
	osenv := os.Environ()
	osenv = append(osenv, obj.Envs...)
	da.LogDouble(aid, fmt.Sprintf("attaching envs: %v", osenv))
	cmd.Env = osenv

	out, err := cmd.CombinedOutput()
	if err != nil {
		if out != nil {
			err = fmt.Errorf("%w: %s", err, out)
		}
		reportError(w, "cmd", err)
		return
	}

	out = []byte(strings.TrimSpace(string(out)))
	da.LogDouble(aid, fmt.Sprintf("cmd output: %v", string(out)))

	if !json.Valid(out) {
		o := make(map[string]string)
		o["output"] = string(out)

		b, err := json.Marshal(o)
		if err != nil {
			reportError(w, "marshalMap", err)
			return
		}

		da.Respond(w, b)
	} else {
		da.Respond(w, out)
	}

}
