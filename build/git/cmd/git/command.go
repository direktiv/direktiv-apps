package main

import (
	"github.com/direktiv/direktiv-apps/pkg/reusable"
	"github.com/go-git/go-git/v5"
)

type gitExecutor interface {
	name() string
	requiredArgs() int
	execute(repo *git.Repository, c gitCommand,
		ri *reusable.RequestInfo) (interface{}, error)
}

type gitExecutorImpl struct {
	cmdName string
	cmdArgs int
	cmdFunc func(repo *git.Repository, c gitCommand,
		ri *reusable.RequestInfo) (interface{}, error)
}

func newGitExecutorImpl(name string, reqArgs int, fn func(repo *git.Repository, c gitCommand,
	ri *reusable.RequestInfo) (interface{}, error)) gitExecutor {

	return &gitExecutorImpl{
		cmdName: name,
		cmdArgs: reqArgs,
		cmdFunc: fn,
	}

}

func (gei *gitExecutorImpl) name() string {
	return gei.cmdName
}

func (gei *gitExecutorImpl) requiredArgs() int {
	return gei.cmdArgs
}

func (gei *gitExecutorImpl) execute(repo *git.Repository, c gitCommand,
	ri *reusable.RequestInfo) (interface{}, error) {
	return gei.cmdFunc(repo, c, ri)
}
