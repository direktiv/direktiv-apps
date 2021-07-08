package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"github.com/mattn/go-shellwords"
	"github.com/vorteil/direktiv-apps/pkg/direktivapps"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os/exec"
	"time"
)

const ps = "/bin/pwsh"
const code = "com.vmware-power-cli.%s.error"

// PowerShell struct
type PowerShell struct {
	powerShell string
	aid        string
}

// New create new session
func New(aid string) *PowerShell {
	ps, _ := exec.LookPath("pwsh")
	return &PowerShell{
		powerShell: ps,
		aid:        aid,
	}
}

func (p *PowerShell) execute(args ...string) (stdOut string, stdErr string, err error) {
	args = append([]string{"-NoProfile", "-NonInteractive", "-c"}, args...)
	hash := NewSHA1Hash()
	args = append(args, fmt.Sprintf("1>%s", hash))
	cmd := exec.Command(ps, args...)
	direktivapps.LogDouble(p.aid, fmt.Sprintf("executing '%v", cmd.Args))
	// var stdout bytes.Buffer
	var stderr bytes.Buffer
	// cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()

	data, err := ioutil.ReadFile(hash)
	if err != nil {
		return "", stdErr, err
	}
	stdOut = string(data)
	// stdOut, stdErr = stdout.String(), stderr.String()
	// direktivapps.LogDouble(p.aid, fmt.Sprintf("stdout: %s\nstderr: %s", stdOut, stdErr))
	return stdOut, stdErr, err
}

type VMWarePowerCLIInput struct {
	Run []string `json:"run"`
}

func VMWarePowerCLIHandler(w http.ResponseWriter, r *http.Request) {
	var obj VMWarePowerCLIInput
	aid, err := direktivapps.Unmarshal(&obj, r)
	if err != nil {
		direktivapps.RespondWithError(w, fmt.Sprintf(code, "unmarshal-input"), err.Error())
		return
	}
	posh := New(aid)

	direktivapps.LogDouble(aid, "reading input...")

	object := make(map[string]interface{})
	for _, r := range obj.Run {
		args, err := shellwords.Parse(r)
		if err != nil {
			direktivapps.RespondWithError(w, fmt.Sprintf(code, "shellwords"), err.Error())
			return
		}
		o, e, err := posh.execute(args...)
		if err != nil {
			direktivapps.RespondWithError(w, fmt.Sprintf(code, "execution"), fmt.Sprintf("%v: %s", err, string(o)))
			return
		}
		direktivapps.LogDouble(aid, fmt.Sprintf("%s%s", o, e))

		object[r] = string(o + e)
	}

	direktivapps.LogDouble(aid, "fetching output...")

	data, err := json.Marshal(&object)
	if err != nil {
		direktivapps.RespondWithError(w, fmt.Sprintf(code, "marshal-response"), err.Error())
		return
	}

	direktivapps.Respond(w, data)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	direktivapps.StartServer(VMWarePowerCLIHandler)
}

// NewSHA1Hash generates a new SHA1 hash based on
// a random number of characters.
func NewSHA1Hash(n ...int) string {
	noRandomCharacters := 32

	if len(n) > 0 {
		noRandomCharacters = n[0]
	}

	randString := RandomString(noRandomCharacters)

	hash := sha1.New()
	hash.Write([]byte(randString))
	bs := hash.Sum(nil)

	return fmt.Sprintf("%x", bs)
}

var characterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

// RandomString generates a random string of n length
func RandomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = characterRunes[rand.Intn(len(characterRunes))]
	}
	return string(b)
}
