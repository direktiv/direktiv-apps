package main

import (
	"github.com/direktiv/direktiv-apps/pkg/reusable"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
)

func deleteBranch(repo *git.Repository, c gitCommand, ri *reusable.RequestInfo) (interface{}, error) {
	return ifAvail(c.Args, 0), repo.DeleteBranch(c.Args[0])
}

func createBranch(repo *git.Repository, c gitCommand, ri *reusable.RequestInfo) (interface{}, error) {

	branchName := c.Args[0]
	localRef := plumbing.NewBranchReferenceName(branchName)

	bc := &config.Branch{
		Name:   branchName,
		Remote: "origin",
		Merge:  localRef,
	}

	err := repo.CreateBranch(bc)
	if err != nil {
		return nil, err
	}

	headRef, err := repo.Head()
	if err != nil {
		return nil, err
	}
	ref := plumbing.NewHashReference(localRef, headRef.Hash())

	// The created reference is saved in the storage.
	return ifAvail(c.Args, 0), repo.Storer.SetReference(ref)
}

func listBranch(repo *git.Repository, c gitCommand, ri *reusable.RequestInfo) (interface{}, error) {

	branches, err := repo.Branches()
	if err != nil {
		return nil, err
	}

	ret := make(map[string]string)
	branches.ForEach(func(ref *plumbing.Reference) error {
		ret[ref.Name().Short()] = ref.Hash().String()
		return nil
	})

	return ret, nil

}
