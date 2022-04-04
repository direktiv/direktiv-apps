package main

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/direktiv/direktiv-apps/pkg/reusable"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	httpgit "github.com/go-git/go-git/v5/plumbing/transport/http"
)

func commit(repo *git.Repository, c gitCommand, ri *reusable.RequestInfo) (interface{}, error) {

	wt, err := repo.Worktree()
	if err != nil {
		return nil, err
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
	return commit.String(), err

}

func pushTags(repo *git.Repository, c gitCommand, ri *reusable.RequestInfo) (interface{}, error) {

	mw := io.MultiWriter(os.Stdout, ri.LogWriter())

	pushOptions := &git.PushOptions{
		Progress: mw,
		RefSpecs: []config.RefSpec{config.RefSpec("refs/tags/*:refs/tags/*")},
	}

	if c.user != "" && c.password != "" {
		pushOptions.Auth = &httpgit.BasicAuth{
			Username: c.user,
			Password: c.password,
		}
	}

	return "", repo.Push(pushOptions)
}

func push(repo *git.Repository, c gitCommand, ri *reusable.RequestInfo) (interface{}, error) {

	mw := io.MultiWriter(os.Stdout, ri.LogWriter())
	localRef := plumbing.NewBranchReferenceName(c.Args[0])
	spec := fmt.Sprintf("%s:%s", localRef.String(), localRef.String())

	pushOptions := &git.PushOptions{
		Progress: mw,
		RefSpecs: []config.RefSpec{config.RefSpec(spec)},
	}

	if c.user != "" && c.password != "" {
		pushOptions.Auth = &httpgit.BasicAuth{
			Username: c.user,
			Password: c.password,
		}
	}

	return "", repo.Push(pushOptions)
}

func add(repo *git.Repository, c gitCommand, ri *reusable.RequestInfo) (interface{}, error) {

	wt, err := repo.Worktree()
	if err != nil {
		return nil, err
	}
	h, err := wt.Add(c.Args[0])
	return h.String(), err

}

func listCommits(repo *git.Repository, c gitCommand, ri *reusable.RequestInfo) (interface{}, error) {

	commits, err := repo.CommitObjects()
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
