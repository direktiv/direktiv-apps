package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func setupCmds() error {

	err := os.MkdirAll("/var/lib/docker", 0655)
	if err != nil {
		return err
	}

	cmds := [][]string{
		{
			"/usr/sbin/ip",
			"link", "set", "dev", "lo", "up",
		},
		{
			"/usr/sbin/ip",
			"link", "set", "dev", "vec0", "up",
		},
		{
			"/usr/sbin/ip",
			"addr", "add", "10.0.2.100/24", "dev", "vec0",
		},
		{
			"/usr/sbin/ip",
			"route", "add", "default", "via", "10.0.2.2",
		},
		{
			"/usr/bin/mount",
			"-t", "proc", "proc", "/proc/",
		},
		{
			"/usr/bin/mount",
			"-t", "sysfs", "sys", "/sys/",
		},
		{
			"/usr/bin/mount",
			"-t", "cgroup2", "none", "/sys/fs/cgroup",
		},
		{
			"/usr/bin/mount",
			"-t", "ext4", "/disk/containers.img", "/var/lib/docker",
		},
	}

	for a := range cmds {
		c := cmds[a]
		cmd := exec.Command(c[0], c[1:]...)
		log.Printf("executing: %v\n", cmd)
		b, err := cmd.CombinedOutput()

		fmt.Printf("%v", string(b))
		if err != nil {
			return err
		}
	}

	return nil
}

func startDocker() error {

	cmd := exec.Command("/usr/bin/dockerd", "-H", "unix:///var/run/docker.sock",
		"-H", "tcp://0.0.0.0", "--tls=false")
	cmd.Stderr = os.Stdout
	cmd.Stdout = os.Stdout

	path := fmt.Sprintf("PATH=/usr/bin:/usr/sbin:%s", os.Getenv("PATH"))
	path = strings.TrimSuffix(path, ":")
	cmd.Env = append(cmd.Env, path)

	return cmd.Run()

}

func startRunner() error {

	err := setupCmds()
	if err != nil {
		return err
	}

	return startDocker()

}
