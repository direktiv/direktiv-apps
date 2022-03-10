package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/direktiv/direktiv-apps/pkg/reusable"
	"github.com/mackerelio/go-osstat/memory"
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
	jobs    chan *requestInput
	max     = 5
	timeout = 240
)

func waitForDocker() error {

	for i := 0; i < timeout; i++ {
		err := dockerInfo()
		if err == nil {
			return nil
		}
		time.Sleep(1 * time.Second)
	}

	return fmt.Errorf("docker start timed out")

}

func managerHandler(w http.ResponseWriter, r *http.Request, ri *reusable.RequestInfo) {

	ri.Logger().Infof("running handler")

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
	case <-time.After(time.Duration(timeout) * time.Minute):
		reusable.ReportError(w, errForCode("timeout"), err)
		return
	}

	ri.Logger().Infof("build finished")

	reusable.ReportResult(w, "")

}

func runVM() {

	log.Printf("starting with vm")
	cpus := runtime.NumCPU()
	if cpus > 8 {
		cpus = 8
	}

	mem := uint64(4096)
	memory, err := memory.Get()
	if err == nil {
		mem = memory.Total
		mem /= (1024 * 1024)
		mem = uint64(math.Min(float64(mem/2), 8192))
	}

	log.Printf("starting vm with %d cpus & %dm", cpus, mem)
	cmd := exec.Command("/usr/bin/qemu-system-x86_64", "--accel", "tcg,thread=multi",
		"-machine", "q35",
		"-smp", fmt.Sprintf("%d", cpus),
		"-m", fmt.Sprintf("%d", mem),
		"-serial", "stdio",
		"-display", "none",
		"-device", "virtio-scsi-pci,id=scsi", "-device", "scsi-hd,drive=hd0",
		"-drive", "if=none,file=/base.vmdk,format=vmdk,id=hd0",
		"-netdev", "user,id=network0,hostfwd=tcp::2375-:2375",
		"-device", "virtio-net-pci,netdev=network0,id=virtio0,mac=26:10:05:00:00:0a",
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout

	log.Printf("starting with command: %v", cmd)

	err = cmd.Run()

	// we can not do much, if the VM dies we are done here
	if err != nil {
		panic(err)
	}

}

func main() {

	go runVM()

	jobs = make(chan *requestInput, max)
	for a := 0; a < max; a++ {
		go serveBuild()
	}

	reusable.StartServer(managerHandler, nil)

}

func dockerInfo() error {
	cmd := exec.Command("/usr/bin/docker", "info")
	return cmd.Run()
}

func errForCode(errCode string) string {
	return fmt.Sprintf("com.docker.%s.error", errCode)
}
