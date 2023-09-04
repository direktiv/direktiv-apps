package main

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/direktiv/direktiv-apps/pkg/reusable"
	"github.com/go-ping/ping"
)

type requestInput struct {
	App      string   `json:"app"`
	Targets  []string `json:"targets"`
	Count    int      `json:"count"`
	Interval int      `json:"interval"`
	Result   string   `json:"result"` // can be stat / detail or somethign else for basic
}

type pingResult struct {
	Host    string           `json:"host"`
	Packets []*ping.Packet   `json:"packets,omitempty"`
	Stats   *ping.Statistics `json:"statistic,omitempty"`
	Success bool             `json:"success"`
}

type lookupResult struct {
	Host    string   `json:"host"`
	Addrs   []string `json:"address,omitempty"`
	Fails   int      `json:"fails"`
	Runs    int      `json:"runs"`
	Reason  string   `json:"reason,omitempty"`
	Success bool     `json:"success"`
}

func networkHandler(w http.ResponseWriter, r *http.Request, ri *reusable.RequestInfo) {

	obj := new(requestInput)
	err := reusable.Unmarshal(obj, false, r)
	if err != nil {
		reusable.ReportError(w, reusable.UnmarshallError, err)
		return
	}

	if len(obj.Targets) == 0 {
		reusable.ReportError(w, reusable.UnmarshallError, fmt.Errorf("no targets defined"))
		return
	}

	switch obj.App {
	case "ping":
		r, err := pingTarget(obj, w, ri)
		if err != nil {
			reusable.ReportError(w, errForCode("ping"), err)
			return
		}
		reusable.ReportResult(w, r)
	case "lookup":
		r, err := lookupTarget(obj, w, ri)
		if err != nil {
			reusable.ReportError(w, errForCode("ping"), err)
			return
		}
		reusable.ReportResult(w, r)
	case "lookupaddr":
		r, err := lookupAddr(obj, w, ri)
		if err != nil {
			reusable.ReportError(w, errForCode("ping"), err)
			return
		}
		reusable.ReportResult(w, r)
	default:
		reusable.ReportError(w, errForCode("app"), fmt.Errorf("application %s does not exist", obj.App))
		return
	}

}

func lookupAddr(obj *requestInput, w http.ResponseWriter, ri *reusable.RequestInfo) (interface{}, error) {

	hostList := make([]*lookupResult, 0)

	if obj.Interval <= 0 {
		obj.Interval = 1
	}
	if obj.Count <= 0 {
		obj.Count = 1
	}

	for i := range obj.Targets {
		t := obj.Targets[i]

		if len(t) == 0 {
			continue
		}

		lr := &lookupResult{
			Host:    t,
			Success: false,
		}

		for i := 0; i < obj.Count; i++ {

			ri.Logger().Infof("lookup address %s", t)

			lr.Runs++
			addrs, err := net.LookupAddr(t)
			if err != nil {
				lr.Fails++
				lr.Reason = err.Error()
				time.Sleep(time.Duration(obj.Interval) * time.Second)
				continue
			}
			ri.Logger().Infof("address %v", addrs)
			if obj.Result == "detail" {
				lr.Addrs = addrs
			}
			time.Sleep(time.Duration(obj.Interval) * time.Second)
		}

		if lr.Fails == 0 {
			lr.Success = true
		}

		hostList = append(hostList, lr)
	}

	return hostList, nil
}

func lookupTarget(obj *requestInput, w http.ResponseWriter, ri *reusable.RequestInfo) (interface{}, error) {

	hostList := make([]*lookupResult, 0)

	if obj.Interval <= 0 {
		obj.Interval = 1
	}
	if obj.Count <= 0 {
		obj.Count = 1
	}

	for i := range obj.Targets {
		t := obj.Targets[i]

		if len(t) == 0 {
			continue
		}

		lr := &lookupResult{
			Host:    t,
			Success: false,
		}

		for i := 0; i < obj.Count; i++ {

			ri.Logger().Infof("lookup address %s", t)

			lr.Runs++
			addrs, err := net.LookupHost(t)
			if err != nil {
				lr.Fails++
				lr.Reason = err.Error()
				time.Sleep(time.Duration(obj.Interval) * time.Second)
				continue
			}
			ri.Logger().Infof("address %v", addrs)
			if obj.Result == "detail" {
				lr.Addrs = addrs
			}
			time.Sleep(time.Duration(obj.Interval) * time.Second)
		}

		if lr.Fails == 0 {
			lr.Success = true
		}

		hostList = append(hostList, lr)
	}

	return hostList, nil
}

func pingTarget(obj *requestInput, w http.ResponseWriter, ri *reusable.RequestInfo) (interface{}, error) {

	hostList := make([]*pingResult, 0)

	for i := range obj.Targets {

		t := obj.Targets[i]

		if len(t) == 0 {
			continue
		}

		ri.Logger().Infof("pinging %s", t)

		pr := &pingResult{
			Host:    t,
			Packets: make([]*ping.Packet, 0),
			Success: false,
		}

		pinger, err := ping.NewPinger(t)
		if err != nil {
			return nil, err
		}
		pinger.SetPrivileged(true)

		pinger.OnRecv = func(pkt *ping.Packet) {
			ri.Logger().Infof("%d bytes from %s: icmp_seq=%d time=%v ttl=%v",
				pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt, pkt.Ttl)
			if obj.Result == "detail" {
				pr.Packets = append(pr.Packets, pkt)
			}
		}

		pinger.OnFinish = func(stats *ping.Statistics) {
			if obj.Result == "stats" {
				pr.Stats = stats
			}
			if stats.PacketLoss == 0 {
				pr.Success = true
			}
		}

		if obj.Interval <= 0 {
			obj.Interval = 1
		}
		if obj.Count <= 0 {
			obj.Count = 3
		}

		pinger.Count = obj.Count
		pinger.Interval = time.Duration(obj.Interval) * time.Second
		pinger.Timeout = time.Second * 3

		err = pinger.Run()
		if err != nil {
			return nil, err
		}

		hostList = append(hostList, pr)

		ri.Logger().Infof("done")

	}

	return hostList, nil

}

func main() {

	reusable.StartServer(networkHandler, nil)

}

func errForCode(errCode string) string {
	return fmt.Sprintf("com.network.%s.error", errCode)
}
