package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"
)

type prep struct {
	dir        string
	tar        string
	dockerfile string
}

const (
	dfName = "Dockerfile.add"

	authTempl = `{
		"auths": 
			%s	
		}`
)

var (
	shellTemplate = `#!/bin/bash
	
ip link set dev lo up
ip link set dev eth0 up
route add default dev eth0
ifconfig eth0 10.0.2.15

export PATH=$PATH:/usr/bin:/usr/local/bin

mount -t proc proc /proc/
mount -t sysfs sys /sys/

mkdir -p /var/lib/containersWORKERID/
mount -t ext4 /buildtmp/WORKERID/containers.img /var/lib/containersWORKERID

manager --manager=false --dir=RUNDIRECTORY
`
)

func genCmds(id int, prep *prep, req *requestInput) error {

	c := []string{"bud", "--root", fmt.Sprintf("/var/lib/containers%d", id)}
	ps := []string{"push", "--root", fmt.Sprintf("/var/lib/containers%d", id)}

	if prep.dockerfile != "" {
		c = append(c, []string{"-f", prep.dockerfile}...)
	}

	if req.Name != "" {
		c = append(c, []string{"-t", req.Name}...)
	}

	var regs = make(map[string]map[string]string)
	if len(req.Registries) > 0 {
		for i := range req.Registries {
			r := req.Registries[i]

			authString := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s",
				r.User, r.Password)))

			// add registry under name key
			reg := make(map[string]string)
			reg["auth"] = authString
			regs[r.Registry] = reg
		}

		authPath := path.Join(prep.dir, "auth.config")

		b, err := json.Marshal(regs)
		if err != nil {
			return err
		}

		ap, err := os.Create(authPath)
		if err != nil {
			return err
		}

		a := fmt.Sprintf(authTempl, string(b))
		_, err = ap.WriteString(a)
		if err != nil {
			return err
		}

		c = append(c, []string{"--authfile", authPath}...)
		ps = append(ps, []string{"--authfile", authPath}...)
	}

	for a := range req.AdditionalArgs {
		arg := req.AdditionalArgs[a]
		c = append(c, arg)
	}

	if prep.tar != "" {
		c = append(c, prep.tar)
	} else {
		c = append(c, req.Context)
	}

	if !req.NoPush {

		// need a name like docker.io/hello/world
		if req.Name == "" {
			return fmt.Errorf("can not push without a name")
		}

		if req.NoTLS {
			ps = append(ps, "--tls-verify=false")
		}
		ps = append(ps, req.Name)
	} else {
		ps = []string{}
	}

	cmds := make(map[string][]string)
	cmds["build"] = c
	cmds["push"] = ps

	f, err := os.Create(path.Join(prep.dir, "commands.json"))
	if err != nil {
		return err
	}

	return json.NewEncoder(f).Encode(cmds)

}

func genShell(id int, prep *prep, req *requestInput) string {

	s := strings.Replace(shellTemplate, "RUNDIRECTORY", prep.dir, -1)
	s = strings.Replace(s, "WORKERID", fmt.Sprintf("%d", id), -1)

	return s
}
