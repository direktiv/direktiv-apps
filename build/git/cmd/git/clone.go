package main

import (
	"io"
	"os"
	"path"

	"github.com/direktiv/direktiv-apps/pkg/reusable"
	"github.com/go-git/go-git/v5"
	httpgit "github.com/go-git/go-git/v5/plumbing/transport/http"
)

func clone(url, user, pwd string, depth int, ri *reusable.RequestInfo) (*git.Repository, error) {

	ri.Logger().Infof("cloning %s", url)

	mw := io.MultiWriter(os.Stdout, ri.LogWriter())
	cloneOptions := &git.CloneOptions{
		URL:      url,
		Progress: mw,
	}

	if depth > 0 {
		cloneOptions.Depth = depth
	}

	if user != "" && pwd != "" {
		ri.Logger().Infof("authenticating with %s", url)

		cloneOptions.Auth = &httpgit.BasicAuth{
			Username: user,
			Password: pwd,
		}
	}

	dir := path.Join(ri.Dir(), "clone")
	return git.PlainClone(dir, false, cloneOptions)

}
