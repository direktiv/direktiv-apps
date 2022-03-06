package main

import (
	"fmt"

	"github.com/direktiv/direktiv-apps/pkg/reusable"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func doCheckout(repo *git.Repository, c gitCommand, ri *reusable.RequestInfo) (interface{}, error) {

	wt, err := repo.Worktree()
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
			if err != nil {
				return nil, err
			}
			co = &git.CheckoutOptions{
				Hash: commit.Hash,
			}
		}
	default:
		return nil, fmt.Errorf("arg for checkout unknown: %s", c.Args[0])
	}

	return c.Args[1], wt.Checkout(co)

}
