package main

import (
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

func diskCleaner() {

	log.Println("starting disk cleaner")

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

			// double check after firewall
			if count() {
				fs := syscall.Statfs_t{}
				err := syscall.Statfs("/", &fs)
				if err != nil {
					fw(false)
					time.Sleep(5 * time.Second)
					continue
				}

				size := fs.Blocks * uint64(fs.Bsize)
				used := size - fs.Bfree*uint64(fs.Bsize)
				perc := (float64(used) / float64(size)) * 100

				log.Printf("disk usage: %0.2f, threshold met: %v", perc, perc > float64(90))
				if perc > float64(90) {
					log.Println("docker prune disk")
					cmd := exec.Command("/usr/bin/docker", "system", "prune", "-a", "-f")
					cmd.Stdout = os.Stdout
					cmd.Stderr = os.Stdout
					cmd.Run()
				}
				fw(false)
			}
		}
		time.Sleep(10 * time.Minute)
	}

}

func startDocker() {

	log.Println("starting docker")
	cmd := exec.Command("/usr/bin/dockerd", "-H", "unix:///var/run/docker.sock",
		"-H", "tcp://0.0.0.0:2375", "--tls=false")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout

	err := cmd.Run()
	if err != nil {
		panic(err)
	}

}

func main() {

	go startDocker()

	diskCleaner()
}
