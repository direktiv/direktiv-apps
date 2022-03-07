package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path"
	"time"

	"github.com/direktiv/direktiv-apps/pkg/reusable"
	"github.com/docker/cli/cli/config/configfile"
	"github.com/docker/cli/cli/config/types"
)

// if push, this are the credentials
type registry struct {
	Registry string `json:"registry"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type requestInput struct {

	// upload tar instead of git pull
	Tar reusable.File `json:"tar"`

	// pass a brand new dockerfile in
	DockerfileArg reusable.File `json:"dockerfile-arg"`

	// the "url" of the build, e.g. git url
	Context string `json:"context"`

	// name for -t flag
	Tag string `json:"tag"`

	// other args, e.g. --pull
	AdditionalArgs []string `json:"args"`

	// enable build kit
	Buildkit bool `json:"buildkit"`

	// push it
	NoPush bool `json:"no-push"`

	// if push, this are the credentials
	Registries []registry `json:"registries"`

	ri     *reusable.RequestInfo `json:"-"`
	doneCh chan bool             `json:"-"`
	errCh  chan error            `json:"-"`
}

var (
	jobs chan *requestInput
	max  int
)

const timeout = 120

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

func checkForDocker(b chan bool) {

	for {
		err := dockerInfo()
		if err == nil {
			b <- true
			break
		}
	}

}

func waitForDocker() error {

	b := make(chan bool)
	go checkForDocker(b)

	select {
	case <-b:
		break
	case <-time.After(120 * time.Second):
		return fmt.Errorf("docker startup timeoud")
	}

	return nil

}

func managerHandler(w http.ResponseWriter, r *http.Request, ri *reusable.RequestInfo) {

	ri.Logger().Infof("running manager")

	obj := new(requestInput)
	err := reusable.Unmarshal(obj, true, r)
	if err != nil {
		reusable.ReportError(w, reusable.UnmarshallError, err)
		return
	}

	if max == len(jobs) {
		ri.Logger().Infof("build queue is full")
		reusable.ReportError(w, errForCode("queue"), err)
		return
	}

	// waiting till docker is up
	// should only happen first time
	ri.Logger().Infof("waiting for docker")
	err = waitForDocker()
	if err != nil {
		reusable.ReportError(w, errForCode("timeout"), err)
		return
	}

	obj.ri = ri
	obj.doneCh = make(chan bool)
	obj.errCh = make(chan error)

	// send job
	jobs <- obj

	select {
	case err = <-obj.errCh:
		ri.Logger().Infof("error building image: %v", err)
		reusable.ReportError(w, errForCode("build"), err)
		return
	case <-obj.doneCh:
		break
	case <-time.After(timeout * time.Minute):
		reusable.ReportError(w, errForCode("nobuild"), err)
		return
	}

	ri.Logger().Infof("build finished")

	reusable.ReportResult(w, "")

}

func main() {

	// if this file exists it is the UML server
	_, err := os.Stat(dockerDiskFile)
	if err == nil {
		err = startRunner()
		fmt.Printf("error running docker: %v\n", err)
		return
	}

	ds := flag.Int("disksize", 40, "disk size in GB")
	builds := flag.Int("builds", 3, "parallel builds")
	flag.Parse()

	max = *builds

	// start bess.sock. has to work
	err = startNetwork()
	if err != nil {
		panic(err)
	}

	// starts UML with
	go startLinuxDocker(*ds)

	// wait for api.sock, has to work
	err = waitForAPISock()
	if err != nil {
		panic(err)
	}

	jobs = make(chan *requestInput, max)
	for a := 0; a < max; a++ {
		go runBuildAndPush(jobs)
	}

	reusable.StartServer(managerHandler, nil)

}

func dockerInfo() error {
	cmd := exec.Command("/usr/bin/docker", "-H", "tcp://127.0.0.1:2375", "info")
	return cmd.Run()
}

func errForCode(errCode string) string {
	return fmt.Sprintf("com.docker.%s.error", errCode)
}
