package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/direktiv/direktiv-apps/pkg/direktivapps"
)

// govcinput to run certain commands on govc
type govcinput struct {
	Username string   `json:"username"`
	Password string   `json:"password"`
	Commands []string `json:"commands"`
	URL      string   `json:"url"`
}

const govc = "govc"
const GOVC_URL = "GOVC_URL"
const GOVC_USERNAME = "GOVC_USERNAME"
const GOVC_PASSWORD = "GOVC_PASSWORD"

const code = "com.%s-govc.error"

func GoVCHandler(w http.ResponseWriter, r *http.Request) {
	obj := new(govcinput)
	aid, err := direktivapps.Unmarshal(obj, r)
	if err != nil {
		direktivapps.RespondWithError(w, fmt.Sprintf(code, "unmarshal"), err.Error())
		return
	}

	// direktivapps.LogDouble(aid, "Setting environment variables...")

	// err = setEnvironmentVariables(obj.URL, obj.Username, obj.Password)
	// if err != nil {
	// 	direktivapps.RespondWithError(w, fmt.Sprintf(code, "set-envvars"), err.Error())
	// 	return
	// }

	direktivapps.LogDouble(aid, "Executing commands...")

	x := make([]string, 0)
	for _, command := range obj.Commands {
		data, err := executeCommand(command, obj.Username, obj.Password, obj.URL)
		if err != nil {
			if len(data) > 0 {
				direktivapps.RespondWithError(w, fmt.Sprintf(code, "execute"), string(data))
			} else {
				direktivapps.RespondWithError(w, fmt.Sprintf(code, "execute"), err.Error())
			}
			return
		}
		x = append(x, string(data))
	}

	data, err := json.Marshal(&x)
	if err != nil {
		direktivapps.RespondWithError(w, fmt.Sprintf(code, "marshal"), err.Error())
		return
	}

	direktivapps.Respond(w, data)
}

func executeCommand(command string, username string, password string, url string) ([]byte, error) {
	args := strings.Fields(command)
	cmd := exec.Command(govc, args...)
	envs := os.Environ()
	envs = append(envs, fmt.Sprintf("%s=%s", GOVC_USERNAME, username))
	envs = append(envs, fmt.Sprintf("%s=%s", GOVC_PASSWORD, password))
	envs = append(envs, fmt.Sprintf("%s=%s", GOVC_URL, url))
	cmd.Env = envs
	return cmd.CombinedOutput()
}

// func setEnvironmentVariables(url, username, password string) error {
// 	err := os.Setenv(GOVC_URL, url)
// 	if err != nil {
// 		return err
// 	}
// 	err = os.Setenv(GOVC_USERNAME, username)
// 	if err != nil {
// 		return err
// 	}
// 	err = os.Setenv(GOVC_PASSWORD, password)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

func main() {
	direktivapps.StartServer(GoVCHandler)
}
