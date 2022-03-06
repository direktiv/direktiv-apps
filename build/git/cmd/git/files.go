package main

import (
	"encoding/base64"
	"io"
	"os"
	"path"

	"github.com/direktiv/direktiv-apps/pkg/reusable"
	"github.com/go-git/go-git/v5"
)

func getFile(repo *git.Repository, c gitCommand, ri *reusable.RequestInfo) (interface{}, error) {

	ff := path.Join(ri.Dir(), "clone", c.Args[0])

	f, err := os.Open(ff)
	if err != nil {
		return nil, err
	}

	b, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return base64.StdEncoding.EncodeToString(b), nil
}

func getFiles(repo *git.Repository, c gitCommand, ri *reusable.RequestInfo) (interface{}, error) {

	var files []string
	h, err := repo.Head()
	if err != nil {
		return nil, err
	}

	commit, err := repo.CommitObject(h.Hash())
	if err != nil {
		return nil, err
	}

	tree, err := commit.Tree()
	if err != nil {
		return nil, err
	}

	for _, entry := range tree.Entries {
		files = append(files, entry.Name)
	}
	return files, nil

}
