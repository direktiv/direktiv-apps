package main

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"path"

	"github.com/direktiv/direktiv-apps/pkg/reusable"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/google/uuid"
)

func doStatus(repo *git.Repository, c gitCommand, ri *reusable.RequestInfo) (interface{}, error) {

	wt, err := repo.Worktree()
	if err != nil {
		return nil, err
	}

	return wt.Status()

}

func doLogs(repo *git.Repository, c gitCommand, ri *reusable.RequestInfo) (interface{}, error) {

	options := &git.LogOptions{
		Order: git.LogOrderCommitterTime,
	}

	if len(c.Args) > 1 {

		switch c.Args[0] {
		case "tag":
			t, err := repo.Tag(c.Args[1])
			if err != nil {
				return nil, err
			}
			options.From = t.Hash()

		case "ref":
			commit, err := repo.CommitObject(plumbing.NewHash(c.Args[1]))
			if err != nil {
				return nil, err
			}
			options.From = commit.Hash
		}
	}

	logs, err := repo.Log(options)
	if err != nil {
		return nil, err
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

	return ll, nil
}

func doScript(repo *git.Repository, c gitCommand, ri *reusable.RequestInfo) (interface{}, error) {

	if len(c.Script.Data) == 0 {
		return nil, nil
	}

	if len(c.Script.Name) == 0 {
		c.Script.Name = uuid.New().String()
	}

	file, err := c.Script.AsFile(0755)
	if err != nil {
		return nil, err
	}
	file.Close()

	defer func() {
		os.Remove(file.Name())
	}()

	cmd := exec.Command(file.Name())
	ri.Logger().Infof("executing %v", cmd)

	var b bytes.Buffer
	mw := io.MultiWriter(&b, os.Stdout, ri.LogWriter())

	cmd.Stderr = mw
	cmd.Stdout = mw
	cmd.Dir = path.Join(ri.Dir(), "clone")

	err = cmd.Run()

	return b.String(), err

}
