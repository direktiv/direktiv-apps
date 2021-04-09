package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"net/http"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	api "github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/vorteil/direktiv-apps/pkg/direktivapps"
)

// InputFluxDetails
type InputFluxDetails struct {
	URL          string                 `json:"url"`
	Token        string                 `json:"token"`
	Type         string                 `json:"type"`
	Organisation string                 `json:"organisation"`
	Bucket       string                 `json:"bucket"`
	Measurement  string                 `json:"measurement"`
	Tags         map[string]string      `json:"tags,omitempty"`
	Fields       map[string]interface{} `json:"fields,omitempty"`
	Query        string                 `json:"query,omitempty"`
}

var code = "com.influxdb.error"

func InfluxDBHandler(w http.ResponseWriter, r *http.Request) {
	obj := new(InputFluxDetails)
	_, err := direktivapps.Unmarshal(obj, r)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	// create new influxdb client
	client := influxdb2.NewClientWithOptions(obj.URL, obj.Token, influxdb2.DefaultOptions().SetTLSConfig(&tls.Config{
		InsecureSkipVerify: true,
	}))
	defer client.Close()

	var bv []byte

	switch obj.Type {
	case "write":
		wapi := client.WriteAPIBlocking(obj.Organisation, obj.Bucket)
		bv, err = WriteData(wapi, obj)
		if err != nil {
			direktivapps.RespondWithError(w, code, err.Error())
			return
		}
	default:
		qapi := client.QueryAPI(obj.Organisation)
		bv, err = QueryData(qapi, obj.Query)
		if err != nil {
			direktivapps.RespondWithError(w, code, err.Error())
			return
		}
	}

	direktivapps.Respond(w, bv)
}

func main() {
	direktivapps.StartServer(InfluxDBHandler)
}

// WriteData writes a point with measurement to influxdb
func WriteData(wapi api.WriteAPIBlocking, obj *InputFluxDetails) ([]byte, error) {
	p := influxdb2.NewPointWithMeasurement(obj.Measurement)

	// Add tags to point
	for k, v := range obj.Tags {
		p.AddTag(k, v)
	}
	// Add fields to point
	for k, v := range obj.Fields {
		p.AddField(k, v)
	}
	// setTime on insert
	p.SetTime(time.Now())

	err := wapi.WritePoint(context.Background(), p)
	if err != nil {
		return nil, err
	}

	var ws struct {
		Success bool `json:"success"`
	}

	// At this point the write was successful give json feedback back
	ws.Success = true
	data, err := json.Marshal(ws)
	if err != nil {
		return nil, err
	}

	return data, err
}

// QueryData takes a string query and returns the results
func QueryData(qapi api.QueryAPI, line string) ([]byte, error) {

	var records struct {
		Records []interface{} `json:"records"`
	}

	result, err := qapi.Query(context.Background(), line)
	if err != nil {
		return nil, err
	}

	for result.Next() {
		record := result.Record()
		records.Records = append(records.Records, map[string]interface{}{record.Field(): record.Value(), "time": record.Time()})
	}
	if result.Err() != nil {
		return nil, result.Err()
	}

	data, err := json.Marshal(records)
	if err != nil {
		return nil, err
	}

	return data, err
}
