package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

func setupCmds() error {

	err := os.MkdirAll("/var/lib/docker", 0655)
	if err != nil {
		return err
	}

	cmds := [][]string{
		{
			"/usr/bin/mknod",
			"/dev/anon", "c", "1", "10",
		},
		{
			"/usr/bin/chmod",
			"666", "/dev/anon",
		},
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

		log.Printf("%v", string(b))
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

func diskCleaner() {

	// write the counter file
	shell := `#!/bin/bash
	netstat -anp 2>/dev/null | grep :2375 | grep ESTABLISHED | wc -l
	`
	f, _ := os.OpenFile("/count.sh", os.O_CREATE|os.O_RDWR, 0755)
	f.WriteString(shell)
	f.Close()

	// method which use count function
	count := func() bool {
		cmd := exec.Command("/count.sh")
		b, err := cmd.CombinedOutput()
		if err != nil {
			return false
		}

		log.Printf("%s clients connected", strings.TrimSpace(string(b)))
		return strings.TrimSpace(string(b)) == "0"
	}

	// enable / disable firewall
	fw := func(add bool) {

		c := "-A"
		if !add {
			c = "-D"
		}
		log.Printf("clean action %v\n", c)
		cmd := exec.Command("/usr/sbin/iptables", c, "INPUT", "-p", "tcp", "--destination-port", "2375", "-j", "DROP")
		cmd.Run()
	}

	for {

		if count() {
			fw(true)
			// after the trafffic is blocked, lets check again
			if count() {
				// check disk size, if above threshold we wipe it
				fs := syscall.Statfs_t{}
				err := syscall.Statfs("/var/lib/docker", &fs)
				if err != nil {
					fw(false)
					continue
				}

				size := fs.Blocks * uint64(fs.Bsize)
				used := size - fs.Bfree*uint64(fs.Bsize)
				perc := float64(used) / float64(size)

				log.Printf("disk usage: %0.2f", perc)

				if perc > 80.0 {

				}
				cmd := exec.Command("/usr/bin/docker", "system", "prune", "-a", "-f")
				cmd.Run()

			}
			fw(false)
		}

		time.Sleep(5 * time.Second)
	}

}

func startRunner() error {

	go diskCleaner()

	err := setupCmds()
	if err != nil {
		return err
	}

	return startDocker()

}
