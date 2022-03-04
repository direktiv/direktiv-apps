package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"time"

	"github.com/codeclysm/extract/v3"
	"github.com/direktiv/direktiv-apps/pkg/reusable"
	"github.com/nxadm/tail"
)

func pullLogs(path string, ri *reusable.RequestInfo) {

	t, err := tail.TailFile(path,
		tail.Config{MustExist: true, Follow: true,
			Logger: tail.DiscardingLogger, Poll: true, ReOpen: false})
	if err != nil {
		ri.Logger().Infof("can not tail build log: %v", err)
		return
	}

	for line := range t.Lines {
		fmt.Printf("%v: %s\n", line.Time.Format(time.RFC822), line.Text)
	}

	ri.Logger().Infof("tailing %s finished", path)

}

func prepareBuildFolder(id, size int) string {

	// crate base directory
	baseDir := path.Join(os.TempDir(), fmt.Sprintf("%d", id))
	log.Printf("base directory %s", baseDir)

	err := os.RemoveAll(baseDir)
	if err != nil {
		log.Printf("can not remove %s", baseDir)
	}

	err = os.MkdirAll(baseDir, 0755)
	if err != nil {
		log.Fatalf("can not create base directory: %v", err)
	}

	// create mount disk in dir
	err = createDisk(baseDir, size)
	if err != nil {
		log.Fatalf("error creating base disk: %v", err)
	}

	return baseDir
}

func createDisk(dir string, size int) error {

	fmt.Printf("creatgin disk in %s\n", dir)

	cmds := [][]string{
		{
			"/usr/bin/dd",
			"if=/dev/zero",
			fmt.Sprintf("of=%s/containers.img", dir),
			"bs=1",
			"count=0",
			fmt.Sprintf("seek=%dG", size),
		},
		{
			"mkfs.ext4",
			fmt.Sprintf("%s/containers.img", dir),
		},
	}

	for a := range cmds {
		c := cmds[a]
		cmd := exec.Command(c[0], c[1:]...)

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		fmt.Printf("executing: %v\n", cmd)
		err := cmd.Run()
		if err != nil {
			return err
		}
	}

	return nil
}

func startRunner(id, size int, jobs <-chan *requestInput) {

	cont := func(dir, msg string, req *requestInput, err error) {
		os.RemoveAll(dir)
		req.ri.Logger().Infof(fmt.Sprintf("%s: %v", msg, err))
		req.errCh <- err
	}

	baseDir := prepareBuildFolder(id, size)

	for req := range jobs {
		req.ri.Logger().Infof("worker %d started", id)

		prep, err := prepWorkDir(req)
		if err != nil {
			req.ri.Logger().Infof("error preparing dir: %v", err)
			req.errCh <- err
			continue
		}

		// start pulling logs
		go pullLogs(path.Join(prep.dir, buildLog), req.ri)
		go pullLogs(path.Join(prep.dir, pushLog), req.ri)

		req.ri.Logger().Infof("using dir %s", prep.dir)

		s := genShell(id, prep, req)
		initPath := path.Join(prep.dir, "run.sh")
		run, err := os.Create(initPath)
		if err != nil {
			cont(prep.dir, "can not create run.sh file", req, err)
			continue
		}
		run.WriteString(s)
		run.Close()
		os.Chmod(initPath, 0755)

		err = genCmds(id, prep, req)
		if err != nil {
			cont(prep.dir, "can not create commands for build", req, err)
			continue
		}

		req.ri.Logger().Infof("running uml in %v", baseDir)

		linuxArgs := []string{"rootfstype=hostfs", "rw",
			"mem=1024m", "quiet", "eth0=slirp,,/usr/bin/slirp-fullbolt",
			fmt.Sprintf("init=%s", path.Join(prep.dir, "run.sh"))}

		fw := func(lookup string) {
			v, ok := os.LookupEnv(lookup)
			if ok {
				linuxArgs = append(linuxArgs,
					fmt.Sprintf("%s=%s", lookup, v))
			}
		}

		// forward proxy settings
		fw("HTTP_PROXY")
		fw("HTTPS_PROXY")
		fw("NO_PROXY")

		cmd := exec.Command("linux", linuxArgs...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		// this is an error in any case because the buildah process ends
		// which leads to a kernel panic in UML
		cmd.Run()

		// we need to check the error file if there was an issue instead
		errTxt, err := os.ReadFile(path.Join(prep.dir, errOut))

		// error file exists, something happened
		if err == nil {
			os.RemoveAll(prep.dir)
			err = fmt.Errorf("error during build: %v", string(errTxt))
			req.ri.Logger().Infof(err.Error())
			req.errCh <- err
			continue
		}

		os.RemoveAll(prep.dir)
		req.doneCh <- true

		req.ri.Logger().Infof("worker %d finished", id)
	}
}

func prepWorkDir(req *requestInput) (*prep, error) {

	prep := &prep{}

	// prep linux build
	dir, err := os.MkdirTemp("", req.ri.ActionID())
	if err != nil {
		return nil, err
	}
	prep.dir = dir

	if req.Tar.Data != "" {

		req.ri.Logger().Infof("tar file attached")

		r, err := req.Tar.AsReader()
		if err != nil {
			return nil, err
		}

		err = extract.Gz(context.Background(), r, dir, nil)
		if err != nil {
			return nil, err
		}

		prep.tar = path.Join(dir, req.Tar.Name)

	}

	if req.DockerfileArg.Data != "" {
		req.ri.Logger().Infof("dockerfile file attached")

		f, err := req.DockerfileArg.AsFile(0644)
		if err != nil {
			return nil, err
		}
		f.Close()
		req.ri.Logger().Infof("copying dockerfile to %s", path.Join(dir, dfName))
		err = os.Rename(f.Name(), path.Join(dir, dfName))
		if err != nil {
			return nil, err
		}

		prep.dockerfile = path.Join(dir, dfName)
	}

	// create the log files so they can be tailed
	pl, err := os.Create(path.Join(dir, pushLog))
	if err != nil {
		return nil, err
	}
	defer pl.Close()

	bl, err := os.Create(path.Join(dir, buildLog))
	if err != nil {
		return nil, err
	}
	defer bl.Close()

	return prep, err
}
