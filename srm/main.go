package main

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/direktiv/direktiv-apps/pkg/direktivapps"
)

var code = "com.srm.%s.error"

const lbAddress = "srm-lb-address"
const lbPort = "srm-lb-port"
const timestamp = "timestamp"
const variable = "variable"
const value = "value"
const group = "group"

type SRMConn struct {
	conn net.Conn
}

func (src SRMConn) Write(content string) (int, error) {
	writer := bufio.NewWriter(src.conn)
	number, err := writer.WriteString(content)
	if err == nil {
		err = writer.Flush()
	}
	return number, err
}

func SRMHandler(w http.ResponseWriter, r *http.Request) {
	obj := make([]map[string]interface{}, 0)
	aid, err := direktivapps.Unmarshal(&obj, r)
	if err != nil {
		direktivapps.RespondWithError(w, fmt.Sprintf(code, "unmarshal"), err.Error())
		return
	}

	direktivapps.LogDouble(aid, "validating input...")

	var newObjects []map[string]interface{}
	for i, metric := range obj {
		// check to see if the address exists
		err = checkForAddress(i, metric)
		if err != nil {
			direktivapps.RespondWithError(w, fmt.Sprintf(code, "address-missing"), err.Error())
			return
		}

		err = checkForValue(i, metric)
		if err != nil {
			direktivapps.RespondWithError(w, fmt.Sprintf(code, "value-missing"), err.Error())
			return
		}

		err = checkForValue(i, metric)
		if err != nil {
			direktivapps.RespondWithError(w, fmt.Sprintf(code, "variable-missing"), err.Error())
			return
		}

		checkForPort(metric)
		checkForTime(metric)
		checkForGroup(metric)

		newObjects = append(newObjects, metric)
	}

	direktivapps.LogDouble(aid, "finished validating...")

	direktivapps.LogDouble(aid, "sending data...")
	for _, obj := range newObjects {
		address := obj[lbAddress].(string)
		port := obj[lbPort].(string)

		delete(obj, lbAddress)
		delete(obj, lbPort)

		raw := fmt.Sprintf("+r\t%v\t", int(obj[timestamp].(int64)))
		// add group, variable and value
		raw += fmt.Sprintf("%s\t%s\t%v", obj[group], obj[variable], obj[value])
		delete(obj, "group")
		delete(obj, "value")
		delete(obj, "variable")
		delete(obj, "timestamp")
		for key, element := range obj {
			raw += fmt.Sprintf("\t%s=%v", key, element)
		}
		raw += "\n"
		direktivapps.LogDouble(aid, "%s >> %s:%s", raw, address, port)

		tcpAddr, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf("%s:%s", address, port))
		if err != nil {
			direktivapps.RespondWithError(w, fmt.Sprintf(code, "invalid-address"), err.Error())
			return
		}

		conn, err := net.DialTCP("tcp", nil, tcpAddr)
		if err != nil {
			direktivapps.RespondWithError(w, fmt.Sprintf(code, "dial-tcp"), err.Error())
			return
		}

		srcConn := &SRMConn{
			conn: conn,
		}

		_, err = srcConn.Write(raw)
		// _, err = io.Copy(conn, strings.NewReader(raw))
		if err != nil {
			direktivapps.RespondWithError(w, fmt.Sprintf(code, "write-tcp"), err.Error())
			return
		}

		conn.Close()
	}

	direktivapps.Respond(w, []byte{})
}

func checkForVariable(i int, metric map[string]interface{}) error {
	if metric[variable] != nil {
		// check if address is provided in the objects
		if str, ok := metric[lbAddress].(string); ok {
			if str == "" {
				return errors.New(fmt.Sprintf("result '%v' is missing variable", i))
			}
		} else {
			return errors.New(fmt.Sprintf("result '%v' is not a valid variable", i))
		}
	} else {
		return errors.New(fmt.Sprintf("result '%v' is missing variable", i))
	}
	return nil
}

func checkForValue(i int, metric map[string]interface{}) error {
	if metric[value] != nil {
		// check if address is provided in the objects
		if str, ok := metric[lbAddress].(string); ok {
			if str == "" {
				return errors.New(fmt.Sprintf("result '%v' is missing value", i))
			}
		} else {
			return errors.New(fmt.Sprintf("result '%v' is not a valid value", i))
		}
	} else {
		return errors.New(fmt.Sprintf("result '%v' is missing value", i))
	}
	return nil
}

func checkForGroup(metric map[string]interface{}) {
	if metric[group] == nil {
		metric[group] = group
	}
}

func checkForTime(metric map[string]interface{}) {
	// if not provided set it to default
	if metric[timestamp] == nil {
		metric[timestamp] = time.Now().Unix()
	} else {
		if ts, ok := metric[timestamp].(float64); ok {
			metric[timestamp] = int64(ts)
		}
	}
}

func checkForPort(metric map[string]interface{}) {
	// if not provided it set it to 2020
	if metric[lbPort] == nil {
		metric[lbPort] = "2020"
	}
}

func checkForAddress(i int, metric map[string]interface{}) error {
	if metric[lbAddress] != nil {
		// check if address is provided in the objects
		if str, ok := metric[lbAddress].(string); ok {
			if str == "" {
				return errors.New(fmt.Sprintf("result '%v' is missing address", i))
			}
		} else {
			return errors.New(fmt.Sprintf("result '%v' is not a valid address", i))
		}
	} else {
		return errors.New(fmt.Sprintf("result '%v' is missing address", i))
	}
	return nil
}

func main() {
	direktivapps.StartServer(SRMHandler)
}
