package main

import (
	"github.com/direktiv/direktiv-apps/pkg/reusable"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func deleteTag(repo *git.Repository, c gitCommand,
	ri *reusable.RequestInfo) (interface{}, error) {
	return ifAvail(c.Args, 0), repo.DeleteTag(c.Args[0])
}

func createTag(repo *git.Repository, c gitCommand,
	ri *reusable.RequestInfo) (interface{}, error) {

	var hash plumbing.Hash

	// we tag head if no ref arg is given
	if len(c.Args) == 1 {
		ref, err := repo.Head()
		if err != nil {
			return nil, err
		}
		hash = ref.Hash()
	} else {
		commit, err := repo.CommitObject(plumbing.NewHash(c.Args[1]))
		if err != nil {
			return nil, err
		}
		hash = commit.Hash
	}

	// TODO sign and message
	tag, err := repo.CreateTag(c.Args[0], hash, nil)
	if err != nil {
		return nil, err
	}

	ret := make(map[string]string)
	ret[tag.Name().Short()] = tag.Hash().String()

	return ret, nil

}

func listTags(repo *git.Repository, c gitCommand,
	ri *reusable.RequestInfo) (interface{}, error) {

	tags, err := repo.Tags()
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
