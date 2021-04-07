package main

import (
	"io/ioutil"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esutil"
	"github.com/vorteil/direktiv-apps/pkg/direktivapps"
)

type ElasticsearchInput struct {
	Method    string                 `json:"method"`    // 'store' or 'query' to either index or return results
	Index     string                 `json:"index"`     // index to search or create records on
	Addresses []string               `json:"addresses"` // array of addresses to hit the elasticsearch host
	Record    map[string]interface{} `json:"record"`    // record that you wish to store note.. not needed if method is query
}

func main() {
	g := direktivapps.ActionError{
		ErrorCode:    "com.elasticsearch.error",
		ErrorMessage: "",
	}

	obj := new(ElasticsearchInput)
	direktivapps.ReadIn(obj, g)

	cfg := elasticsearch.Config{
		Addresses: obj.Addresses,
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		g.ErrorMessage = err.Error()
		direktivapps.WriteError(g)
	}

	var data []byte

	switch obj.Method {
	case "query":
		res, err := es.Search(es.Search.WithIndex(obj.Index))
		if err != nil {
			g.ErrorMessage = err.Error()
			direktivapps.WriteError(g)
		}
		defer res.Body.Close()

		data, err = ioutil.ReadAll(res.Body)
		if err != nil {
			g.ErrorMessage = err.Error()
			direktivapps.WriteError(g)
		}

	case "store":

		res, err := es.Index(obj.Index, esutil.NewJSONReader(&obj.Record))
		if err != nil {
			g.ErrorMessage = err.Error()
			direktivapps.WriteError(g)
		}
		defer res.Body.Close()

		data, err = ioutil.ReadAll(res.Body)
		if err != nil {
			g.ErrorMessage = err.Error()
			direktivapps.WriteError(g)
		}

	}

	direktivapps.WriteOut(data, g)
}
