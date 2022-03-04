package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
)

func createRunnerError(dir, phase string, err error) {
	errMsg := fmt.Sprintf("error in %s: %v", phase, err)
	fmt.Println(errMsg)
	f, err := os.Create(path.Join(dir, errOut))
	if err != nil {
		// can not do much
		fmt.Printf("could not write error file: %v\n", err)
		return
	}
	f.WriteString(errMsg)
	f.Close()
}

func runPhase(dir, phase, log string, cmds []string) error {

	delFile := func(d *os.File) {
		d.Close()
		os.Remove(d.Name())
	}

	fout, err := os.OpenFile(log, os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		createRunnerError(dir, phase, err)
		return err
	}
	defer delFile(fout)

	fout.WriteString(fmt.Sprintf("starting phase %s\n", phase))

	var errBytes bytes.Buffer
	mw := io.MultiWriter(fout, &errBytes)

	cmd := exec.Command("/usr/bin/buildah", cmds...)
	fout.WriteString(fmt.Sprintf("running %v", cmd))

	// we catch errors in case
	cmd.Stdout = os.Stdout
	cmd.Stderr = mw

	err = cmd.Run()
	if err != nil {
		createRunnerError(dir, phase, fmt.Errorf("%v: %v", err.Error(), errBytes.String()))
		return err
	}

	return nil

}

func runnerHandler(dir string) {

	c, err := os.ReadFile(path.Join(dir, "commands.json"))
	if err != nil {
		createRunnerError(dir, "build", err)
		return
	}

	var cmds map[string][]string

	err = json.Unmarshal(c, &cmds)
	if err != nil {
		createRunnerError(dir, "build", err)
		return
	}

	err = runPhase(dir, "build", path.Join(dir, buildLog), cmds["build"])
	if err != nil {
		return
	}

	if len(cmds["push"]) > 0 {
		runPhase(dir, "push", path.Join(dir, pushLog), cmds["push"])
	}

}
