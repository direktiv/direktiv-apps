package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/direktiv/direktiv-apps/pkg/reusable"
)

type git struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

// if push, this are the credentials
type registry struct {
	Registry string `json:"registry"`
	User     string `json:"user"`
	Password string `json:"password"`
}

const (
	buildLog = "build.out"
	pushLog  = "pushLog.out"
	errOut   = "error.out"
)

type requestInput struct {

	// upload tar instead of git pull
	Tar reusable.File `json:"tar"`

	// pass a brand new dockerfile in
	DockerfileArg reusable.File `json:"dockerfile-arg"`

	// sets the dockerfile path within the context
	Dockerfile string `json:"dockerfile"`

	// the "url" of the build, e.g. git url
	Context string `json:"context"`

	// name for -t flag
	Name string `json:"name"`

	// other args, e.g. --pull
	AdditionalArgs []string `json:"args"`

	// push it
	NoPush bool `json:"no-push"`
	NoTLS  bool `json:"no-tls"`

	// if push, this are the credentials
	Registries []registry `json:"registries"`

	Git git `json:"git"`

	// aux
	ri     *reusable.RequestInfo `json:"-"`
	doneCh chan bool             `json:"-"`
	errCh  chan error            `json:"-"`
}

var jobs chan *requestInput

func managerHandler(w http.ResponseWriter, r *http.Request, ri *reusable.RequestInfo) {

	ri.Logger().Infof("running manager")

	obj := new(requestInput)
	err := reusable.Unmarshal(obj, true, r)
	if err != nil {
		reusable.ReportError(w, reusable.UnmarshallError, err)
		return
	}

	obj.ri = ri
	obj.doneCh = make(chan bool, 1)
	obj.errCh = make(chan error, 1)

	go func() {
		jobs <- obj
	}()

	select {
	case <-obj.doneCh:
		ri.Logger().Infof("build finished")
		break
	case err := <-obj.errCh:
		ri.Logger().Infof("build failed: %v", err)
		reusable.ReportError(w, errForCode("timeout"), err)
		return
	case <-time.After(3 * time.Hour):
		// let's assume that failed :)
		ri.Logger().Infof("build failed, timed out")
		reusable.ReportError(w, errForCode("timeout"), fmt.Errorf("build timed out"))
		return
	}

}

func main() {

	manager := flag.Bool("manager", true, "runs as manager")
	i := flag.Int("builds", 3, "parallel builds")
	dir := flag.String("dir", "", "directory with build data (runner only)")
	ds := flag.Int("disksize", 1, "disk size in GB")

	flag.Parse()

	if *manager {
		jobs = make(chan *requestInput, *i)
		for a := 0; a < *i; a++ {
			go startRunner(a, *ds, jobs)
		}
		reusable.StartServer(managerHandler, nil)
	} else {
		runnerHandler(*dir)
	}

}

func errForCode(errCode string) string {
	return fmt.Sprintf("com.container.%s.error", errCode)
}
