package main

import (
	"fmt"
	"net/http"

	"github.com/direktiv/direktiv-apps/pkg/reusable"
	"github.com/gosnmp/gosnmp"
)

type requestInput struct {
	Server string `json:"server"`
	Port   int    `json:"port"`

	Version   int    `json:"version"`
	Transport string `json:"transport"`

	Community string `json:"community"`

	Inform bool `json:"inform"`

	// v1
	GenericTrap  int    `json:"generic"`
	SpecificTrap int    `json:"specific"`
	Enterprise   string `json:"enterprise"`

	Variables []gosnmp.SnmpPDU `json:"variables"`
}

type SNMPLogWriter struct {
	w *reusable.DirektivLogger
}

func (lw *SNMPLogWriter) Print(v ...interface{}) {
	lw.w.Infof("%v", v)
}

func (lw *SNMPLogWriter) Printf(format string, v ...interface{}) {
	lw.w.Infof(format, v)
}

const snmpError = "direktiv.snmp.error"

func snmpHandler(w http.ResponseWriter, r *http.Request, ri *reusable.RequestInfo) {

	obj := new(requestInput)
	err := reusable.Unmarshal(obj, true, r)
	if err != nil {
		reusable.ReportError(w, reusable.UnmarshallError, err)
		return
	}

	if obj.Server == "" {
		reusable.ReportError(w, snmpError,
			fmt.Errorf("server required for snmp"))
		return
	}

	if obj.Port == 0 {
		obj.Port = 161
	}

	gsnmp := gosnmp.Default

	lw := &SNMPLogWriter{
		w: ri.Logger(),
	}

	gosnmp.Default.Logger = gosnmp.NewLogger(lw)

	if obj.Transport == "tcp" {
		gsnmp.Transport = obj.Transport
	}

	gsnmp.Target = obj.Server
	gsnmp.Port = uint16(obj.Port)

	ri.Logger().Infof("sending trap to %s:%d via %s",
		gsnmp.Target, gsnmp.Port, gsnmp.Transport)

	if obj.Community != "" {
		gsnmp.Community = obj.Community
	}

	ri.Logger().Infof("using community %s", gsnmp.Community)

	switch obj.Version {
	case 1:
		gsnmp.Version = gosnmp.Version1
	case 3:
		gsnmp.Version = gosnmp.Version3
	}

	ri.Logger().Infof("sending trap version %v", gsnmp.Version)

	trap := gosnmp.SnmpTrap{
		IsInform:     obj.Inform,
		Variables:    obj.Variables,
		Enterprise:   obj.Enterprise,
		AgentAddress: obj.Server,
		GenericTrap:  1,
		SpecificTrap: obj.SpecificTrap,
	}

	err = gsnmp.Connect()
	if err != nil {
		reusable.ReportError(w, snmpError,
			fmt.Errorf("can not connect to snmp agent"))
		return
	}
	defer gsnmp.Conn.Close()

	_, err = gsnmp.SendTrap(trap)
	if err != nil {
		reusable.ReportError(w, snmpError, err)
		return
	}

	reusable.ReportResult(w, obj.Variables)
}

func main() {
	reusable.StartServer(snmpHandler, nil)
}
