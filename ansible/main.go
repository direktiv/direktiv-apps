package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/direktiv/direktiv-apps/pkg/direktivapps"
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
	Collections     []string `json:"collections"`
	ShowPlaybook    bool     `json:"show"`
}

const (
	errCode       = "com.ansible.%s"
	ansibleBin    = "ansible-playbook"
	ansibleGalaxy = "ansible-galaxy"
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

	lw, _ := direktivapps.NewDirektivLogWriter(aid)

	// install collections, e.g ansible-galaxy collection install devsec.hardening
	for i := range obj.Collections {
		galaxyCmd := exec.Command(ansibleGalaxy, "collection", "install", obj.Collections[i])

		galaxyCmd.Stderr = lw
		galaxyCmd.Stdout = lw
		err = galaxyCmd.Run()
		if err != nil {
			reportError(w, "galaxy", err)
			return
		}
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
	osenv = append(osenv, "ANSIBLE_CALLBACKS_ENABLED=json",
		"ANSIBLE_STDOUT_CALLBACK=json")
	osenv = append(osenv, obj.Envs...)

	cmd.Env = osenv

	if obj.ShowPlaybook {
		b, err := ioutil.ReadFile(filepath.Join(mntPath, obj.Playbook))
		if err != nil {
			da.LogDouble(aid, fmt.Sprintf("can not find playbook: %v", err))
			reportError(w, "playbook", err)
			return
		}
		da.LogDouble(aid, fmt.Sprintf("running playbook %v", obj.Playbook))
		da.LogDouble(aid, fmt.Sprintf("%v", string(b)))
	} else {
		da.LogDouble(aid, fmt.Sprintf("running playbook %v", obj.Playbook))
	}

	outb := &bytes.Buffer{}
	errb := &bytes.Buffer{}

	cmd.Stdout = outb
	cmd.Stderr = errb

	err = cmd.Run()
	if err != nil {
		da.LogDouble(aid, fmt.Sprintf("running playbook failed: %v", err))
		da.LogDouble(aid, fmt.Sprintf("%v", string(errb.Bytes())))
		reportError(w, "cmd", err)
		return
	}

	direktivapps.LogDouble(aid, fmt.Sprintf("Output: %v", string(outb.Bytes())))

	if json.Valid(outb.Bytes()) {
		o := make(map[string]interface{})

		o["output"] = json.RawMessage(string(outb.Bytes()))

		b, err := json.Marshal(o)
		if err != nil {
			reportError(w, "marshalMap", err)
			return
		}

		da.Respond(w, b)
	} else {
		da.Respond(w, outb.Bytes())
	}

}
