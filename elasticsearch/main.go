package main

import (
	"io/ioutil"
	"net/http"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esutil"
	"github.com/direktiv/direktiv-apps/pkg/direktivapps"
)

type ElasticsearchInput struct {
	Method    string                 `json:"method"`    // 'store' or 'query' to either index or return results
	Index     string                 `json:"index"`     // index to search or create records on
	Addresses []string               `json:"addresses"` // array of addresses to hit the elasticsearch host
	Record    map[string]interface{} `json:"record"`    // record that you wish to store note.. not needed if method is query
}

const code = "com.elasticsearch.error"

func main() {
	direktivapps.StartServer(ElasticSearch)
}

func ElasticSearch(w http.ResponseWriter, r *http.Request) {
	obj := new(ElasticsearchInput)
	_, err := direktivapps.Unmarshal(obj, r)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	cfg := elasticsearch.Config{
		Addresses: obj.Addresses,
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	var data []byte

	switch obj.Method {
	case "query":
		res, err := es.Search(es.Search.WithIndex(obj.Index))
		if err != nil {
			direktivapps.RespondWithError(w, code, err.Error())
			return
		}
		defer res.Body.Close()

		data, err = ioutil.ReadAll(res.Body)
		if err != nil {
			direktivapps.RespondWithError(w, code, err.Error())
			return
		}

	case "store":

		res, err := es.Index(obj.Index, esutil.NewJSONReader(&obj.Record))
		if err != nil {
			direktivapps.RespondWithError(w, code, err.Error())
			return
		}
		defer res.Body.Close()

		data, err = ioutil.ReadAll(res.Body)
		if err != nil {
			direktivapps.RespondWithError(w, code, err.Error())
			return
		}

	}

	direktivapps.Respond(w, data)
}
