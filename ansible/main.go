package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"

	da "github.com/vorteil/direktiv-apps/pkg/direktivapps"
)

func main() {
	da.StartServer(handler)
}

type inputData struct {
	User       string `json:"user"`
	Host       string `json:"host"`
	PrivateKey string `json:"privateKey"`
	Playbook   string `json:"playbook"`
}

const errCode = "com.ansible.%s"

func reportError(w http.ResponseWriter, code string, err error) {
	da.RespondWithError(w, fmt.Sprintf(errCode, code), err.Error())
}

func handler(w http.ResponseWriter, r *http.Request) {

	obj := new(inputData)
	aid, err := da.Unmarshal(obj, r)
	if err != nil {
		reportError(w, "inputUnmarshal", err)
		return
	}

	// write base64-decoded privatekey to filesystem
	pkFile, err := os.OpenFile("pkFile", os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		reportError(w, "openFile", err)
		return
	}
	defer pkFile.Close()

	pkBytes, err := base64.StdEncoding.DecodeString(obj.PrivateKey)
	if err != nil {
		reportError(w, "b64decode", err)
		return
	}

	_, err = io.Copy(pkFile, bytes.NewReader(pkBytes))
	if err != nil {
		reportError(w, "writePrivateKeyData", err)
		return
	}

	err = pkFile.Close()
	if err != nil {
		reportError(w, "closeFile", err)
		return
	}

	// write base64-decoded playbook to filesystem
	pbFile, err := os.OpenFile("pbFile", os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		reportError(w, "openFile", err)
		return
	}
	defer pbFile.Close()

	pbBytes, err := base64.StdEncoding.DecodeString(obj.Playbook)
	if err != nil {
		reportError(w, "b64decode", err)
		return
	}

	_, err = io.Copy(pbFile, bytes.NewReader(pbBytes))
	if err != nil {
		reportError(w, "writePlaybookData", err)
		return
	}

	err = pbFile.Close()
	if err != nil {
		reportError(w, "closeFile", err)
		return
	}

	// write to ssh_config
	f, err := os.OpenFile("/etc/ssh/ssh_config", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		reportError(w, "openSSHConfig", err)
		return
	}
	defer f.Close()

	_, err = io.Copy(f, strings.NewReader(fmt.Sprintf("\n\nHost %s\n", strings.TrimSuffix(obj.Host, ","))))
	if err != nil {
		reportError(w, "modifySSHConfig", err)
		return
	}

	err = f.Close()
	if err != nil {
		reportError(w, "closeSSHConfig", err)
		return
	}

	// exec command to run ansible
	cmd := exec.Command("ansible-playbook", "-i", obj.Host, "-u", obj.User, "--private-key", pkFile.Name(), pbFile.Name(), "-vvvv")
	cmd.Env = append(os.Environ(), "ANSIBLE_HOST_KEY_CHECKING=False")
	b, err := cmd.CombinedOutput()
	if err != nil {
		if b != nil {
			da.Log(aid, string(b))
		}
		reportError(w, "exec", err)
		return
	}

	da.Log(aid, string(b))

}
