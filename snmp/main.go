package main

import (
	"fmt"
	"net/http"

	"github.com/gosnmp/gosnmp"
	"github.com/vorteil/direktiv-apps/pkg/direktivapps"
)

var code = "com.snmp.%s.error"

type LoggingSNMP struct {
	Aid string `json:"aid"`
}

func (lsnmp *LoggingSNMP) Print(v ...interface{}) {
	for _, f := range v {
		direktivapps.LogDouble(lsnmp.Aid, f.(string))
	}
}

func (lsnmp *LoggingSNMP) Printf(format string, v ...interface{}) {
	direktivapps.LogDouble(lsnmp.Aid, fmt.Sprintf(format, v...))
}

type SnmpInput struct {
	URL           string           `json:"url"`
	Port          int              `json:"port"`
	IsInform      bool             `json:"inform-request"`
	SNMPV1Headers TrapHeaders      `json:"snmpv1-headers"`
	Variables     []gosnmp.SnmpPDU `json:"variables"`
}

type TrapHeaders struct {
	Enterprise   string `json:"enterprise"`
	AgentAddress string `json:"agent-address"`
	GenericTrap  int    `json:"generic-trap"`
	SpecificTrap int    `json:"specific-trap"`
	Timestamp    uint   `json:"timestamp"`
}

func SNMPHandler(w http.ResponseWriter, r *http.Request) {
	var obj SnmpInput
	aid, err := direktivapps.Unmarshal(&obj, r)
	if err != nil {
		direktivapps.RespondWithError(w, fmt.Sprintf(code, "unmarshal-input"), err.Error())
		return
	}

	direktivapps.LogDouble(aid, "setup logger...")

	logger := &LoggingSNMP{
		Aid: aid,
	}

	g := gosnmp.Default

	direktivapps.LogDouble(aid, "reading input...")

	// set the url
	g.Target = obj.URL
	g.Port = uint16(obj.Port)
	// set the logger
	g.Logger = gosnmp.NewLogger(logger)

	err = g.Connect()
	if err != nil {
		direktivapps.RespondWithError(w, fmt.Sprintf(code, "connecting"), err.Error())
		return
	}
	defer g.Conn.Close()

	// send the trap
	result, err := g.SendTrap(gosnmp.SnmpTrap{
		IsInform:     obj.IsInform,
		Variables:    obj.Variables,
		Enterprise:   obj.SNMPV1Headers.Enterprise,
		GenericTrap:  obj.SNMPV1Headers.GenericTrap,
		SpecificTrap: obj.SNMPV1Headers.SpecificTrap,
		Timestamp:    obj.SNMPV1Headers.Timestamp,
	})
	if err != nil {
		direktivapps.RespondWithError(w, fmt.Sprintf(code, "send-trap"), err.Error())
		return
	}

	data, err := result.MarshalMsg()
	if err != nil {
		direktivapps.RespondWithError(w, fmt.Sprintf(code, "marshal-response"), err.Error())
		return
	}

	direktivapps.Respond(w, data)
}

func main() {
	direktivapps.StartServer(SNMPHandler)
}
