package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path"

	"github.com/codeclysm/extract/v3"
	"github.com/direktiv/direktiv-apps/pkg/reusable"
	"github.com/go-git/go-git/v5"
)

type cloneCommand struct {
	Repo     string `json:"repo"`
	Depth    int    `json:"depth"`
	Ref      string `json:"ref"`
	User     string `json:"user"`
	Password string `json:"pwd"`
	Scope    string `json:"scope"`
	Name     string `json:"name"`

	InitScope string `json:"init-scope"`
	InitName  string `json:"init-name"`
}

type requestInput struct {
	Clone cloneCommand `json:"clone"`
	Cmds  []gitCommand `json:"cmds"`
}

type gitCommand struct {
	Cmd      string        `json:"cmd"`
	Args     []string      `json:"args"`
	Script   reusable.File `json:"script"`
	Continue bool          `json:"continue"`

	// needed to push it down to the push command
	user     string `json:"-"`
	password string `json:"-"`
}

var (
	cmdList = []gitExecutor{
		newGitExecutorImpl("list-tags", 0, listTags),
		newGitExecutorImpl("create-tag", 1, createTag),
		newGitExecutorImpl("delete-tag", 1, deleteTag),
		newGitExecutorImpl("delete-push-tag", 1, deletePushTag),
		newGitExecutorImpl("push-tags", 0, pushTags),
		newGitExecutorImpl("list-branches", 0, listBranch),
		newGitExecutorImpl("delete-branch", 1, deleteBranch),
		newGitExecutorImpl("create-branch", 1, createBranch),
		newGitExecutorImpl("checkout", 2, doCheckout),

		newGitExecutorImpl("diff", 2, doDiff),
		newGitExecutorImpl("script", 0, doScript),
		newGitExecutorImpl("status", 0, doStatus),

		newGitExecutorImpl("add", 1, add),
		newGitExecutorImpl("commit", 3, commit),
		newGitExecutorImpl("list-commits", 0, listCommits),
		newGitExecutorImpl("push", 1, push),
		newGitExecutorImpl("logs", 0, doLogs),

		newGitExecutorImpl("get-files", 0, getFiles),
		newGitExecutorImpl("get-file", 1, getFile),
	}
	cmds = make(map[string]gitExecutor)
)

func init() {

	for a := range cmdList {
		c := cmdList[a]
		cmds[c.name()] = c
	}

}

func gitHandler(w http.ResponseWriter, r *http.Request, ri *reusable.RequestInfo) {

	obj := new(requestInput)
	err := reusable.Unmarshal(obj, true, r)
	if err != nil {
		reusable.ReportError(w, reusable.UnmarshallError, err)
		return
	}

	var repo *git.Repository
	if obj.Clone.InitName != "" && obj.Clone.InitScope != "" {

		r, _, err := ri.ReadVar(obj.Clone.InitScope, obj.Clone.InitName)
		if err != nil {
			reusable.ReportError(w, errForCode("clone"), err)
			return
		}

		dir := path.Join(ri.Dir(), "clone")
		err = extract.Gz(context.Background(), r, dir, nil)
		if err != nil {
			reusable.ReportError(w, errForCode("clone"), err)
			return
		}

		// load repo
		repo, err = git.PlainOpen(dir)
		if err != nil {
			reusable.ReportError(w, errForCode("clone"), err)
			return
		}

	} else {
		repo, err = clone(obj.Clone.Repo, obj.Clone.User, obj.Clone.Password, obj.Clone.Depth, ri)
		if err != nil {
			reusable.ReportError(w, errForCode("clone"), err)
			return
		}
	}

	ret := make(map[int]interface{})
	for i := range obj.Cmds {
		c := obj.Cmds[i]

		cmd, ok := cmds[c.Cmd]
		if !ok {
			reusable.ReportError(w, errForCode("execute"), fmt.Errorf("command %s unkown", c.Cmd))
			return
		}

		if len(c.Args) < cmd.requiredArgs() {
			reusable.ReportError(w, errForCode("execute"),
				fmt.Errorf("not enough args for command %s unkown, has %d but needs %d",
					cmd.name(), len(c.Args), cmd.requiredArgs()))
			return
		}

		if cmd.name() == "push" {
			c.user = obj.Clone.User
			c.password = obj.Clone.Password
		}

		ri.Logger().Infof("running %v, args %v", cmd.name(), c.Args)
		result, err := cmd.execute(repo, c, ri)

		if err != nil && c.Continue {
			ri.Logger().Infof("error executing command %s: %v, continuing", cmd.name(), err)
			result = fmt.Sprintf("error: %s", err.Error())
		} else if err != nil {
			reusable.ReportError(w, errForCode("execute"), err)
			return
		}

		ret[i] = result

	}

	// store as variable if set
	if obj.Clone.Scope != "" && obj.Clone.Name != "" {

		ri.Logger().Infof("storing repo in %s/%s", obj.Clone.Scope, obj.Clone.Name)
		outFile := path.Join(ri.Dir(), "out", obj.Clone.Scope, obj.Clone.Name)
		srcFile := path.Join(ri.Dir(), "clone")
		err = os.Rename(srcFile, outFile)
		if err != nil {
			reusable.ReportError(w, errForCode("store"), err)
			return
		}

	}

	reusable.ReportResult(w, ret)

}

func main() {
	reusable.StartServer(gitHandler, nil)
}

func errForCode(errCode string) string {
	return fmt.Sprintf("com.git.%s.error", errCode)
}

func ifAvail(d []string, a int) string {
	if a >= len(d) {
		return ""
	}
	return d[a]
}
