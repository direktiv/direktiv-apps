package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/vorteil/direktiv-apps/pkg/direktivapps"
)

const code = "com.git.error"

type cmdIn struct {
	Cmds []string `json:"cmds"`
}

type item struct {
	Cmd    string      `json:"cmd"`
	Output interface{} `json:"output"`
}

var mtx sync.Mutex

func request(w http.ResponseWriter, r *http.Request) {

	var cmds cmdIn

	aid, err := direktivapps.Unmarshal(&cmds, r)
	if err != nil {
		fmt.Printf("error unmarshalling: %v", err)
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	if len(cmds.Cmds) == 0 {
		fmt.Printf("no git commands provided")
		direktivapps.RespondWithError(w, code, "no git commands provided")
		return
	}

	itemList := make(map[string]*item)

	// we lock the run of all commands
	mtx.Lock()
	defer mtx.Unlock()

	for i, c := range cmds.Cmds {

		// we add the output anyway to see the error in the frontend
		cmdStr := c
		if strings.Contains(c, "@") {
			cmdStr = "..."
		}

		logMsg := fmt.Sprintf("running command %d '%s'", i, cmdStr)
		fmt.Println(logMsg)
		direktivapps.Log(aid, logMsg)

		// cut git command in front
		if strings.HasPrefix(c, "git") {
			c = c[3:len(c)]
			c = strings.TrimSpace(c)
		}

		d, isJSON, err := runGitCmd(c)

		key := fmt.Sprintf("cmd%d", i)

		// we add the output anyway to see the error in the frontend
		if strings.Contains(c, "@") {
			c = "..."
		}

		itemList[key] = &item{
			Cmd: c,
		}

		if isJSON {
			itemList[key].Output = d
		} else {
			itemList[key].Output = string(d.(string))
		}

		if err != nil {
			logMsg := fmt.Sprintf("error running command '%v'", err)
			fmt.Println(logMsg)
			direktivapps.Log(aid, logMsg)
			break
		}

	}

	b, err := json.Marshal(itemList)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	direktivapps.Respond(w, b)

}

func runGitCmd(cmd string) (interface{}, bool, error) {

	var clonedDir string

	// clone and log commands need to run before the actual git command
	f := strings.Fields(cmd)

	if strings.Contains(cmd, "log ") || strings.Contains(cmd, " log") {

		f = append(f, "--pretty=format:{%n  \"commit\": \"%H\",%n  \"author\": \"%aN <%aE>\",%n  \"date\": \"%ad\",%n  \"message\": \"%f\"%n},")

	} else if strings.Contains(cmd, "clone ") {

		// clone directory is getting removed with every clone
		cloneDir := func(f []string) string {
			if len(f) > 2 {
				// if the last two don't have - as prefix, folder is provided
				if !strings.HasPrefix(f[len(f)-1], "-") &&
					!strings.HasPrefix(f[len(f)-2], "-") {
					os.RemoveAll(f[len(f)-1])
					return f[len(f)-1]
				}
			}

			base := path.Base(f[len(f)-1])
			dir := strings.TrimSuffix(base, filepath.Ext(base))

			os.RemoveAll(dir)

			return dir
		}

		clonedDir = cloneDir(f)

	}

	git := exec.Command("git", f...)
	d, err := git.Output()
	if err != nil {
		if e, ok := err.(*exec.ExitError); ok {
			d = e.Stderr
		}
		return string(d), false, err
	}

	// check for log, tag, trying to json the output
	if strings.Contains(cmd, "log ") || strings.Contains(cmd, " log") {

		if len(d) > 1 {
			// remove trailing comma and make array
			in := "[" + string(d[0:len(string(d))-1]) + "]"
			var msgs []map[string]interface{}
			err = json.Unmarshal([]byte(in), &msgs)
			if err != nil {
				return err.Error(), false, err
			}
			return msgs, true, err

		}
		return string(d), false, nil

	} else if strings.Contains(cmd, "clone ") {

		d = []byte(clonedDir)

	} else if strings.Contains(cmd, "tag ") || strings.Contains(cmd, " tag") {

		tags := strings.Split(string(d), "\n")
		return tags[0 : len(tags)-1], true, err
	}

	return string(d), false, nil

}

func main() {
	direktivapps.StartServer(request)
}
