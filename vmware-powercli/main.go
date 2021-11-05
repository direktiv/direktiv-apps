package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/direktiv/direktiv-apps/pkg/direktivapps"
	pss "github.com/direktiv/go-powershell"
	"github.com/direktiv/go-powershell/backend"
)

const ps = "/bin/pwsh"
const code = "com.vmware-power-cli.%s.error"

// PowerShell struct
type PowerShell struct {
	powerShell string
	aid        string
}

type CmdOutput struct {
	Result string      `json:"result"`
	Output interface{} `json:"output"`
}

type Script struct {
	Name string `json:"script"`
}

type VMWarePowerCLIInput struct {
	Host    string   `json:"host"`
	User    string   `json:"user"`
	Pwd     string   `json:"password"`
	Run     []string `json:"run"`
	OnError string   `json:"on-error"`
	Full    bool     `json:"full-command"`
	Print   bool     `json:"print"`
	Scripts []struct {
		Name string   `json:"name"`
		Args []string `json:"args"`
	} `json:"scripts"`
}

func printOut(aid, o, e string) {
	if len(o) > 0 {
		direktivapps.LogDouble(aid, fmt.Sprintf("shell stdout: %v", o))
	}
	if len(e) > 0 {
		direktivapps.LogDouble(aid, fmt.Sprintf("shell stderr: %v", e))
	}
}

func runShell(aid, server, user, password string) (pss.Shell, error) {

	back := &backend.Local{}
	shell, err := pss.New(back)
	if err != nil {
		return nil, err
	}

	if server != "none" {
		o, e, err := shell.Execute(fmt.Sprintf("Connect-VIServer -Server %s -User %s -Password %s", server, user, password))
		if err != nil {
			return nil, err
		}
		printOut(aid, o, e)
	}

	_, _, err = shell.Execute("$WarningPreference = 'SilentlyContinue'")
	if err != nil {
		return nil, err
	}

	return shell, nil

}

func execute(shell pss.Shell, psCmd string) (string, error) {

	stdout, _, err := shell.Execute(psCmd)
	if err != nil {
		return "", err
	}

	return stdout, nil
}

func executeFile(aid string, shell pss.Shell, file string, args []string) (string, string, error) {

	cmd := []string{
		ps,
		"-f",
		file,
	}

	cmd = append(cmd, args...)
	fullCommand := strings.Join(cmd, " ")
	stdout, stderr, err := shell.Execute(fullCommand)
	if err != nil {
		return "", "", err
	}

	return stdout, stderr, nil
}

func toJSON(aid, str string) interface{} {

	var js json.RawMessage
	err := json.Unmarshal([]byte(str), &js)
	if err != nil {
		direktivapps.LogDouble(aid, fmt.Sprintf("output is string value: %v", err))
		return str
	}

	direktivapps.LogDouble(aid, fmt.Sprintf("output is json"))
	return json.RawMessage(str)

}

func createShell(aid string, w http.ResponseWriter, host, user, pwd string) (pss.Shell, error) {

	direktivapps.LogDouble(aid, fmt.Sprintf("getting powershell for host %v", host))
	return runShell(aid, host, user, pwd)

}

func VMWarePowerCLIHandler(w http.ResponseWriter, r *http.Request) {

	var obj VMWarePowerCLIInput
	aid, err := direktivapps.Unmarshal(&obj, r)
	if err != nil {
		direktivapps.RespondWithError(w, fmt.Sprintf(code, "unmarshal-input"), err.Error())
		return
	}

	// get shell for commands if there are any
	var fsh, sh pss.Shell
	if len(obj.Run) > 0 {
		sh, err = createShell(aid, w, obj.Host, obj.User, obj.Pwd)
		if err != nil {
			direktivapps.RespondWithError(w, fmt.Sprintf(code, "connecting"), err.Error())
			return
		}
		defer sh.Exit()
	}

	if len(obj.Scripts) > 0 {
		fsh, err = createShell(aid, w, "none", "", "")
		if err != nil {
			direktivapps.RespondWithError(w, fmt.Sprintf(code, "connecting"), err.Error())
			return
		}
		defer fsh.Exit()
	}

	object := make(map[string]interface{})

	// run commands
	for i := range obj.Run {
		if obj.Full {
			direktivapps.LogDouble(aid, fmt.Sprintf("running command: %s", obj.Run[i]))
		} else {
			direktivapps.LogDouble(aid, fmt.Sprintf("running command: %s...", obj.Run[i][0:15]))
		}

		o, err := execute(sh, obj.Run[i])
		if err != nil {
			object[fmt.Sprintf("%d", i)] = &CmdOutput{
				Result: "error",
				Output: err.Error(),
			}
			direktivapps.LogDouble(aid, fmt.Sprintf("command error: %v", err))

			if obj.OnError == "stop" {
				direktivapps.RespondWithError(w, fmt.Sprintf(code, "execute"), err.Error())
				return
			}

			continue
		}

		direktivapps.LogDouble(aid, fmt.Sprintf("command output: %v", o))
		object[fmt.Sprintf("%d", i)] = &CmdOutput{
			Result: "success",
			Output: toJSON(aid, o),
		}

	}

	// run scripts
	for i := range obj.Scripts {
		tmpDir := r.Header.Get("Direktiv-TempDir")
		direktivapps.LogDouble(aid, fmt.Sprintf("running script %v", obj.Scripts[i].Name))
		o, e, err := executeFile(aid, fsh, fmt.Sprintf("%s/%s", tmpDir, obj.Scripts[i].Name), obj.Scripts[i].Args)

		if obj.Print {
			printOut(aid, o, e)
		}

		if err != nil {
			object[fmt.Sprintf("script-%d", i)] = &CmdOutput{
				Result: "error",
				Output: err.Error(),
			}
			direktivapps.LogDouble(aid, fmt.Sprintf("command error: %v", err))

			if obj.OnError == "stop" {
				direktivapps.RespondWithError(w, fmt.Sprintf(code, "execute"), err.Error())
				return
			}
			continue
		}
		object[fmt.Sprintf("script-%d", i)] = &CmdOutput{
			Result: "success",
			Output: toJSON(aid, o),
		}
	}

	data, err := json.Marshal(&object)
	if err != nil {
		direktivapps.RespondWithError(w, fmt.Sprintf(code, "marshal-response"), err.Error())
		return
	}
	direktivapps.Respond(w, data)

}

func main() {
	direktivapps.StartServer(VMWarePowerCLIHandler)
}
