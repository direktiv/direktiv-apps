package main

import (
	"bytes"
	"encoding/base64"

	"github.com/direktiv/direktiv-apps/pkg/reusable"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
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

func doDiff(repo *git.Repository, c gitCommand,
	ri *reusable.RequestInfo) (interface{}, error) {

	ri.Logger().Infof("running diff between %s and %s", c.Args[0], c.Args[1])

	t1, err := createTree(c.Args[0], repo)
	if err != nil {
		return nil, err
	}

	t2, err := createTree(c.Args[1], repo)
	if err != nil {
		return nil, err
	}

	changes, err := object.DiffTree(t1, t2)
	if err != nil {
		return nil, err
	}

	o, err := changes.Patch()
	if err != nil {
		return nil, err
	}

	// repsonse as base64
	var b bytes.Buffer
	o.Encode(&b)
	return base64.StdEncoding.EncodeToString(b.Bytes()), nil

}
