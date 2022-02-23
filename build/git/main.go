package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/direktiv/direktiv-apps/pkg/reusable"
	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/cache"
	"github.com/go-git/go-git/v5/plumbing/object"
	httpgit "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/filesystem"
	"github.com/google/uuid"
)

type gitCommand struct {
	Cmd      string        `json:"cmd"`
	Args     []string      `json:"args"`
	Script   reusable.File `json:"script"`
	Continue bool          `json:"continue"`
}

type cloneCommand struct {
	Repo     string `json:"repo"`
	Depth    int    `json:"depth"`
	Ref      string `json:"ref"`
	User     string `json:"user"`
	Password string `json:"pwd"`
}

type requestInput struct {
	Clone cloneCommand `json:"clone"`
	Cmds  []gitCommand `json:"cmds"`
}

const (
	ltag = "list-tags"
	dtag = "delete-tag"
	ctag = "create-tag"

	lbranch = "list-branches"
	dbranch = "delete-branch"
	cbranch = "create-branch"

	checkout = "checkout"

	gfile = "get-file"
	tfile = "get-files"

	lcommits = "list-commits"
	script   = "script"
	acommit  = "add"
	ccommit  = "commit"
	push     = "push"

	status = "status"

	diff = "diff"
	logs = "logs"
)

func createTree(in string, repo *git.Repository) (*object.Tree, error) {

	var ph plumbing.Hash
	t1, err := repo.Tag(in)
	if err != nil {
		commit, err := repo.CommitObject(plumbing.NewHash(in))
		if err != nil {
			return nil, err
		}
		ph = commit.Hash
	} else {
		ph = t1.Hash()
	}

	co, err := repo.CommitObject(ph)
	if err != nil {
		return nil, err
	}

	return co.Tree()
}

func gitHandler(w http.ResponseWriter, r *http.Request, ri *reusable.RequestInfo) {

	obj := new(requestInput)
	err := reusable.Unmarshal(obj, true, r)
	if err != nil {
		reusable.ReportError(w, reusable.UnmarshallError, err)
		return
	}

	rdir := ri.Dir()
	// rdir := "/tmp"

	dir := fmt.Sprintf("%s/clone", rdir)
	d := osfs.New(dir)
	dot, err := d.Chroot("/git")
	if err != nil {
		reusable.ReportError(w, errForCode("clone"), err)
		return
	}

	storage := filesystem.NewStorage(d, cache.NewObjectLRUDefault())

	mw := io.MultiWriter(os.Stdout, ri.LogWriter())
	cloneOptions := &git.CloneOptions{
		URL:      obj.Clone.Repo,
		Progress: mw,
	}

	if obj.Clone.Depth > 0 {
		cloneOptions.Depth = obj.Clone.Depth
	}

	if obj.Clone.User != "" && obj.Clone.Password != "" {
		cloneOptions.Auth = &httpgit.BasicAuth{
			Username: obj.Clone.User,
			Password: obj.Clone.Password,
		}
	}

	repo, err := git.Clone(storage, dot, cloneOptions)
	if err != nil {
		reusable.ReportError(w, errForCode("clone"), err)
		return
	}

	// repo, err := git.Open(storage, dot)
	// if err != nil {
	// 	reusable.ReportError(w, errForCode("clone"), err)
	// 	return
	// }

	ret := make(map[int]interface{})

	for i := range obj.Cmds {
		c := obj.Cmds[i]

		switch c.Cmd {
		case push:

			if len(c.Args) < 1 {
				reusable.ReportError(w, errForCode("git"), fmt.Errorf("not enough arguments for push"))
				return
			}

			ri.Logger().Infof("pushing to %s", c.Args[0])

			mw := io.MultiWriter(os.Stdout, ri.LogWriter())
			localRef := plumbing.NewBranchReferenceName(c.Args[0])

			spec := fmt.Sprintf("%s:%s", localRef.String(), localRef.String())

			pushOptions := &git.PushOptions{
				Progress: mw,
				RefSpecs: []config.RefSpec{config.RefSpec(spec)},
			}

			if obj.Clone.User != "" && obj.Clone.Password != "" {
				pushOptions.Auth = &httpgit.BasicAuth{
					Username: obj.Clone.User,
					Password: obj.Clone.Password,
				}
			}

			err := repo.Push(pushOptions)
			err = handleErr(c, err, ri.Logger())
			if err != nil {
				reusable.ReportError(w, errForCode("git"), err)
				return
			}

			ret[i] = err

		case diff:

			ri.Logger().Infof("running diff between %s and %s", c.Args[0], c.Args[1])

			t1, err := createTree(c.Args[0], repo)
			err = handleErr(c, err, ri.Logger())
			if err != nil {
				reusable.ReportError(w, errForCode("git"), err)
				return
			}

			t2, err := createTree(c.Args[1], repo)
			err = handleErr(c, err, ri.Logger())
			if err != nil {
				reusable.ReportError(w, errForCode("git"), err)
				return
			}

			changes, err := object.DiffTree(t1, t2)
			err = handleErr(c, err, ri.Logger())
			o, err := changes.Patch()
			if err != nil {
				reusable.ReportError(w, errForCode("git"), err)
				return
			}

			// repsonse as base64
			var b bytes.Buffer
			o.Encode(&b)
			ret[i] = base64.StdEncoding.EncodeToString(b.Bytes())

		case logs:

			options := &git.LogOptions{
				Order: git.LogOrderCommitterTime,
			}

			if len(c.Args) > 1 {

				switch c.Args[0] {
				case "tag":
					t, err := repo.Tag(c.Args[1])
					err = handleErr(c, err, ri.Logger())
					if err != nil {
						reusable.ReportError(w, errForCode("git"), err)
						return
					}
					options.From = t.Hash()

				case "ref":
					commit, err := repo.CommitObject(plumbing.NewHash(c.Args[1]))
					err = handleErr(c, err, ri.Logger())
					if err != nil {
						reusable.ReportError(w, errForCode("git"), err)
						return
					}
					options.From = commit.Hash
				}
			}

			logs, err := repo.Log(options)
			err = handleErr(c, err, ri.Logger())
			if err != nil {
				reusable.ReportError(w, errForCode("git"), err)
				return
			}

			var ll []interface{}

			logs.ForEach(func(comm *object.Commit) error {
				m := make(map[string]interface{})
				m["author"] = comm.Author
				m["committer"] = comm.Committer
				m["message"] = comm.Message
				m["ref"] = comm.Hash.String()
				ll = append(ll, m)
				return nil
			})

			ret[i] = ll

		case status:
			wt, err := repo.Worktree()
			err = handleErr(c, err, ri.Logger())
			if err != nil {
				reusable.ReportError(w, errForCode("git"), err)
				return
			}
			st, err := wt.Status()
			err = handleErr(c, err, ri.Logger())
			if err != nil {
				reusable.ReportError(w, errForCode("git"), err)
				return
			}
			ret[i] = st
		case script:
			{
				ri.Logger().Infof("running script %v", c.Script.Name)
				runScript(&c.Script, ri, dot.Root())
			}
		case gfile, tfile:
			result, err := handleFiles(repo, dot, c, ri.Logger())
			if err != nil {
				reusable.ReportError(w, errForCode("git"), err)
				return
			}
			ret[i] = result
		case lcommits, acommit, ccommit:
			result, err := handleCommit(repo, c, ri.Logger())
			if err != nil {
				reusable.ReportError(w, errForCode("git"), err)
				return
			}
			ret[i] = result
		case checkout:
			result, err := handleCheckout(repo, c, ri.Logger())
			if err != nil {
				reusable.ReportError(w, errForCode("git"), err)
				return
			}
			ret[i] = result
		case ltag, dtag, ctag:
			result, err := handleTag(repo, c, ri.Logger())
			if err != nil {
				reusable.ReportError(w, errForCode("git"), err)
				return
			}
			ret[i] = result
		case lbranch, dbranch, cbranch:
			result, err := handleBranch(repo, c, ri.Logger())
			if err != nil {
				reusable.ReportError(w, errForCode("git"), err)
				return
			}
			ret[i] = result
		default:
			reusable.ReportError(w, errForCode("git"), fmt.Errorf("unknown command: %s", c.Cmd))
			return
		}

	}

	reusable.ReportResult(w, ret)

}

func handleErr(c gitCommand, err error, logger *reusable.DirektivLogger) error {

	if c.Continue && err != nil {
		logger.Infof("error running %s %v: %v", c.Cmd, c.Args, err)
		return nil
	}

	return err
}

func handleFiles(repo *git.Repository, fs billy.Filesystem, c gitCommand, logger *reusable.DirektivLogger) (interface{}, error) {

	switch c.Cmd {
	case gfile:
		if len(c.Args) < 1 {
			return nil, fmt.Errorf("not enough args for file command")
		}
		f, err := fs.Open(c.Args[0])
		err = handleErr(c, err, logger)
		if err != nil {
			return nil, err
		}

		b, err := io.ReadAll(f)
		err = handleErr(c, err, logger)
		if err != nil {
			return nil, err
		}

		return base64.StdEncoding.EncodeToString(b), nil

	case tfile:
		var files []string
		h, _ := repo.Head()
		commit, _ := repo.CommitObject(h.Hash())
		tree, _ := commit.Tree()

		for _, entry := range tree.Entries {
			files = append(files, entry.Name)
		}
		return files, nil
	}

	return nil, fmt.Errorf("unknown file command: %s", c.Cmd)

}

func handleCommit(repo *git.Repository, c gitCommand, logger *reusable.DirektivLogger) (interface{}, error) {

	wt, err := repo.Worktree()
	err = handleErr(c, err, logger)
	if err != nil {
		return nil, err
	}

	switch c.Cmd {
	case ccommit:
		if len(c.Args) < 3 {
			return nil, fmt.Errorf("not enough args for add command")
		}

		name := c.Args[0]
		email := c.Args[1]
		msg := c.Args[2]

		commit, err := wt.Commit(msg, &git.CommitOptions{
			Author: &object.Signature{
				Name:  name,
				Email: email,
				When:  time.Now(),
			},
		})
		err = handleErr(c, err, logger)
		return commit.String(), err

	case acommit:
		if len(c.Args) < 1 {
			return nil, fmt.Errorf("not enough args for add command")
		}

		h, err := wt.Add(c.Args[0])
		err = handleErr(c, err, logger)

		return h.String(), err

	case lcommits:
		{
			logger.Infof("listing commits")

			commits, err := repo.CommitObjects()
			err = handleErr(c, err, logger)
			if err != nil {
				return nil, err
			}

			var ret []map[string]interface{}

			commits.ForEach(func(ref *object.Commit) error {
				m := make(map[string]interface{})
				m["author"] = ref.Author
				m["committer"] = ref.Committer
				m["message"] = ref.Message
				m["ref"] = ref.Hash.String()
				return nil
			})

			return ret, nil
		}
	}

	return nil, fmt.Errorf("unknown command %s", c.Cmd)

}

func handleCheckout(repo *git.Repository, c gitCommand, logger *reusable.DirektivLogger) (interface{}, error) {

	if len(c.Args) < 2 {
		return nil, fmt.Errorf("not enough args for checkout command")
	}

	wt, err := repo.Worktree()
	err = handleErr(c, err, logger)
	if err != nil {
		return nil, err
	}

	var co *git.CheckoutOptions

	switch c.Args[0] {
	case "branch":
		{
			co = &git.CheckoutOptions{
				Branch: plumbing.NewBranchReferenceName(c.Args[1]),
			}
		}
	case "tag":
		{
			t, err := repo.Tag(c.Args[1])
			err = handleErr(c, err, logger)
			if err != nil {
				return nil, err
			}
			co = &git.CheckoutOptions{
				Hash: t.Hash(),
			}
		}
	case "ref":
		{
			commit, err := repo.CommitObject(plumbing.NewHash(c.Args[1]))
			err = handleErr(c, err, logger)
			if err != nil {
				return nil, err
			}
			co = &git.CheckoutOptions{
				Hash: commit.Hash,
			}
		}
	default:
		err = handleErr(c, fmt.Errorf("command for checkout unknown: %s", c.Args[0]), logger)
		return c.Args[1], err
	}

	logger.Infof("checking out %s", c.Args[1])
	err = handleErr(c, wt.Checkout(co), logger)

	return c.Args[1], err

}

func handleBranch(repo *git.Repository, c gitCommand, logger *reusable.DirektivLogger) (interface{}, error) {

	switch c.Cmd {
	case dbranch:

		if len(c.Args) < 1 {
			return nil, fmt.Errorf("not enough args for create branch command")
		}
		logger.Infof("deleting branch %s", c.Args[0])
		return c.Args[0], handleErr(c, repo.DeleteBranch(c.Args[0]), logger)

	case lbranch:

		logger.Infof("listing branches")
		branches, err := repo.Branches()
		err = handleErr(c, err, logger)
		if err != nil {
			return nil, err
		}

		ret := make(map[string]string)

		branches.ForEach(func(ref *plumbing.Reference) error {
			ret[ref.Name().Short()] = ref.Hash().String()
			return nil
		})

		return ret, nil

	case cbranch:

		if len(c.Args) < 1 {
			return nil, fmt.Errorf("not enough args for create branch command")
		}

		branchName := c.Args[0]
		localRef := plumbing.NewBranchReferenceName(branchName)

		logger.Infof("creating branch %s: %s", branchName, localRef)

		bc := &config.Branch{
			Name:   branchName,
			Remote: "origin",
			Merge:  localRef,
		}

		err := handleErr(c, repo.CreateBranch(bc), logger)
		if err != nil {
			return nil, err
		}

		headRef, err := repo.Head()
		if err != nil {
			return nil, err
		}
		ref := plumbing.NewHashReference(localRef, headRef.Hash())

		// The created reference is saved in the storage.
		err = repo.Storer.SetReference(ref)

		return c.Args[0], nil

	}

	return nil, fmt.Errorf("unknown branch command %s", c.Cmd)
}

func handleTag(repo *git.Repository, c gitCommand,
	logger *reusable.DirektivLogger) (interface{}, error) {

	switch c.Cmd {
	case dtag:
		{
			if len(c.Args) < 1 {
				return nil, fmt.Errorf("not enough args for delete tag command")
			}

			logger.Infof("deleting tag %s", c.Args[0])

			err := handleErr(c, repo.DeleteTag(c.Args[0]), logger)

			return nil, err
		}
	case ctag:
		{
			if len(c.Args) < 1 {
				return nil, fmt.Errorf("not enough args for create tag command")
			}

			logger.Infof("creating tag %s", c.Args[0])

			ref, err := repo.Head()
			err = handleErr(c, err, logger)
			if err != nil {
				return nil, err
			}

			h := ref.Hash()
			if len(c.Args) > 1 {
				commit, err := repo.CommitObject(plumbing.NewHash(c.Args[1]))
				err = handleErr(c, err, logger)
				if err != nil {
					return nil, err
				}
				h = commit.Hash
			}

			// TODO sign and message
			tag, err := repo.CreateTag(c.Args[0], h, nil)
			err = handleErr(c, err, logger)
			if err != nil {
				return nil, err
			}

			ret := make(map[string]string)
			ret[tag.Name().Short()] = tag.Hash().String()

			return ret, nil
		}
	case ltag:
		{
			logger.Infof("listing tags")

			tags, err := repo.Tags()
			err = handleErr(c, err, logger)
			if err != nil {
				return nil, err
			}

			ret := make(map[string]string)

			tags.ForEach(func(ref *plumbing.Reference) error {
				ret[ref.Name().Short()] = ref.Hash().String()
				return nil
			})

			return ret, nil
		}
	}

	return nil, fmt.Errorf("unknown tag command %s", c.Cmd)

}

func main() {
	reusable.StartServer(gitHandler, nil)
}

func errForCode(errCode string) string {
	return fmt.Sprintf("com.git.%s.error", errCode)
}

func runScript(f *reusable.File, ri *reusable.RequestInfo, wd string) error {

	if len(f.Data) == 0 {
		return nil
	}

	if len(f.Name) == 0 {
		f.Name = uuid.New().String()
	}

	file, err := f.AsFile(0755)
	if err != nil {
		return err
	}
	file.Close()
	defer os.Remove(file.Name())
	cmd := exec.Command(file.Name())

	ri.Logger().Infof("executing %v %s", cmd, wd)

	mw := io.MultiWriter(os.Stdout, ri.LogWriter())

	cmd.Stderr = mw
	cmd.Stdout = mw
	cmd.Dir = wd

	return cmd.Run()
}
