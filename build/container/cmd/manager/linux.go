package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"time"
)

const (
	dockerDiskFile = "/disk/containers.img"
)

var api = "/api.sock"

func dockerDisk(size int) error {

	cmds := [][]string{
		{
			"/usr/bin/dd",
			"if=/dev/zero",
			fmt.Sprintf("of=%s", dockerDiskFile),
			"bs=1",
			"count=0",
			fmt.Sprintf("seek=%dG", size),
		},
		{
			"mkfs.ext4",
			dockerDiskFile,
		},
	}

	for a := range cmds {
		c := cmds[a]
		cmd := exec.Command(c[0], c[1:]...)
		log.Printf("executing: %v\n", cmd)
		err := cmd.Run()
		if err != nil {
			return err
		}
	}

	return nil
}

func startNetwork() error {

	fmt.Println("setting up bess network")

	cmd := exec.Command("/slirp4netns", "-a", api, "--target-type=bess", "/bess.sock")
	cmd.Start()

	for {
		_, err := os.Stat("/bess.sock")
		if err == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}

	log.Printf("socket up")

	return nil

}

func waitForAPISock() error {

	for {
		_, err := os.Stat(api)
		if err == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}

	c, err := net.Dial("unix", api)
	if err != nil {
		return err
	}

	log.Printf("api socket up\n")
	json := `
	{"execute": "add_hostfwd",
		"arguments": {"proto": "tcp", "host_addr": "0.0.0.0", "host_port": 2375,
		"guest_addr": "10.0.2.100", "guest_port": 2375}}`

	_, err = c.Write([]byte(json))
	return err
}

func startLinuxDocker(size int) error {

	err := dockerDisk(size)
	if err != nil {
		return err
	}

	linuxArgs := []string{"rootfstype=hostfs", "rw",
		"mem=4G", "quiet", "vec0:transport=bess,dst=/bess.sock,depth=128,gro=1",
		"init=/usr/local/bin/manager",
		"cgroup_enable=cpuset",
		"cgroup_enable=memory",
		"cgroup_memory=1",
		"swapaccount=1",
		"quiet"}

	cmd := exec.Command("linux", linuxArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// runs a server inside so it should work fine
	return cmd.Run()

}
