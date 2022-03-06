package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/codeclysm/extract"
)

func runBuildAndPush(input chan *requestInput) {

	var ad, tar, dfName string
	var err error

	report := func(objin *requestInput, err error) {

		os.RemoveAll(ad)
		os.Remove(tar)
		os.Remove(dfName)
		if err != nil {
			objin.errCh <- err
			return
		}
		objin.doneCh <- true

	}

	for {

		obj := <-input
		ri := obj.ri

		ri.Logger().Infof("worker picked up task")

		// build args
		nArgs := []string{"build"}
		env := os.Environ()

		ad, err = os.MkdirTemp("", ri.ActionID())
		if err != nil {
			report(obj, err)
			continue
		}

		ri.Logger().Infof("creating auth files in %s", ad)
		err = createConfigJsonFile(ad, obj.Registries)
		if err != nil {
			report(obj, err)
			continue
		}

		// add config dir to process
		env = append(env, fmt.Sprintf("DOCKER_CONFIG=%s", ad))

		for a := range obj.AdditionalArgs {

			arg := obj.AdditionalArgs[a]
			if strings.HasPrefix(arg, "-f=") &&
				obj.DockerfileArg.Data != "" {
				ri.Logger().Infof("ignoring -f flag: %v", arg)
				continue
			}
			nArgs = append(nArgs, arg)

		}

		if obj.Tag == "" && !obj.NoPush {
			report(obj, fmt.Errorf("tag is required for push"))
			continue
		}

		nArgs = append(nArgs, "-t", obj.Tag)

		//  if there is a tar or dockerfile we need to put it somewhere
		if obj.Tar.Data != "" {

			ri.Logger().Infof("tar attached, extracting")

			r, err := obj.Tar.AsReader()
			if err != nil {
				report(obj, err)
				continue
			}

			td, err := os.MkdirTemp("", ri.ActionID())
			if err != nil {
				report(obj, err)
				continue
			}

			err = extract.Gz(context.TODO(), r, td, nil)
			if err != nil {
				report(obj, err)
				continue
			}

			obj.Context = path.Join(td, obj.Tar.Name)
			nArgs = append(nArgs, obj.Context)

		} else if obj.Context != "" {
			nArgs = append(nArgs, obj.Context)
		} else {
			report(obj, fmt.Errorf("either tar file or context has to be provided"))
			continue
		}

		if obj.Buildkit {
			nArgs = append(nArgs, "--network", "host")
			env = append(env, "DOCKER_BUILDKIT=1")
		}

		// extract extra docker file
		if obj.DockerfileArg.Data != "" {

			// this only works for tar

			df, err := obj.DockerfileArg.AsFile(0644)
			if err != nil {
				report(obj, err)
				continue
			}

			dfName = path.Join(obj.Context, "Dockerfile.direktiv")

			err = os.Rename(df.Name(), dfName)
			if err != nil {
				report(obj, err)
				continue
			}
			nArgs = append(nArgs, "-f", dfName)

		}

		mw := io.MultiWriter(os.Stdout, ri.LogWriter())
		cmd := exec.Command("docker", nArgs...)
		cmd.Env = env
		cmd.Stdout = mw
		cmd.Stderr = mw

		ri.Logger().Infof("running %v", cmd)

		err = cmd.Run()
		if err != nil {
			report(obj, err)
			continue
		}

		if !obj.NoPush {
			err = runPush(obj.Tag, env, mw)
			if err != nil {
				report(obj, err)
				continue
			}
		}

		report(obj, nil)

	}

}

func runPush(image string, env []string, mw io.Writer) error {

	args := []string{"push", image}

	cmd := exec.Command("docker", args...)
	cmd.Env = env
	cmd.Stdout = mw
	cmd.Stderr = mw

	mw.Write([]byte(fmt.Sprintf("running %v\n", cmd)))

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error pushing image: %v", err)
	}

	return nil
}
