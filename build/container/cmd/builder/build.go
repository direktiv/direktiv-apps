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
	"github.com/docker/cli/cli/config/configfile"
	"github.com/docker/cli/cli/config/types"
)

func serveBuild() {

	var ad, tar, dfName string
	var err error

	report := func(objin *requestInput, err error) {

		os.RemoveAll(ad)
		os.RemoveAll(tar)
		os.Remove(dfName)
		if err != nil {
			objin.errCh <- err
			return
		}
		objin.doneCh <- true

	}

	for {

		obj := <-jobs
		ri := obj.ri
		ri.Logger().Infof("worker picked up task")

		// build args
		nArgs := []string{"build", "--progress", "plain"}
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

		if obj.Tag != "" {
			nArgs = append(nArgs, "-t", obj.Tag)
		}

		//  if there is a tar or dockerfile we need to put it somewhere
		if obj.Tar.Data != "" {

			ri.Logger().Infof("tar attached")

			r, err := obj.Tar.AsReader(ri)
			if err != nil {
				report(obj, err)
				continue
			}

			tar, err = os.MkdirTemp("", ri.ActionID())
			if err != nil {
				report(obj, err)
				continue
			}

			err = extract.Gz(context.Background(), r, tar, nil)
			if err != nil {
				report(obj, err)
				continue
			}

			obj.Context = path.Join(tar, obj.Tar.Name)
			nArgs = append(nArgs, obj.Context)

		} else if obj.Context != "" {
			nArgs = append(nArgs, obj.Context)
		} else {
			report(obj, fmt.Errorf("either tar file or context has to be provided"))
			continue
		}

		if obj.Buildkit {
			env = append(env, "DOCKER_BUILDKIT=1")
		}

		// extract extra docker file
		if obj.DockerfileArg.Data != "" {

			// this only works for tar, not for git
			df, err := obj.DockerfileArg.AsFile(ri, 0644)
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

func createConfigJsonFile(dir string, rs []registry) error {

	cf := configfile.ConfigFile{
		AuthConfigs: make(map[string]types.AuthConfig),
		Proxies:     make(map[string]configfile.ProxyConfig),
	}

	for i := range rs {
		r := rs[i]
		authConfig := types.AuthConfig{
			Username: r.User,
			Password: r.Password,
		}
		cf.AuthConfigs[r.Registry] = authConfig
	}

	var proxyConfig configfile.ProxyConfig
	hp, ok := os.LookupEnv("HTTPS_PROXY")
	if ok {
		proxyConfig.HTTPSProxy = hp
	}
	hp, ok = os.LookupEnv("HTTP_PROXY")
	if ok {
		proxyConfig.HTTPProxy = hp
	}
	hp, ok = os.LookupEnv("NO_PROXY")
	if ok {
		proxyConfig.NoProxy = hp
	}

	cf.Proxies["default"] = proxyConfig

	authPath := path.Join(dir, "config.json")
	ap, err := os.Create(authPath)
	if err != nil {
		return err
	}

	return cf.SaveToWriter(ap)

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
