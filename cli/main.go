package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

type DockerResponse struct {
	Count    int            `json:"count"`
	Next     interface{}    `json:"next"`
	Previous interface{}    `json:"previous"`
	Results  []DockerResult `json:"results"`
}

type DockerResult struct {
	Creator int         `json:"creator"`
	ID      int         `json:"id"`
	ImageID interface{} `json:"image_id"`
	Images  []struct {
		Architecture string      `json:"architecture"`
		Features     string      `json:"features"`
		Variant      interface{} `json:"variant"`
		Digest       string      `json:"digest"`
		Os           string      `json:"os"`
		OsFeatures   string      `json:"os_features"`
		OsVersion    interface{} `json:"os_version"`
		Size         int         `json:"size"`
		Status       string      `json:"status"`
		LastPulled   time.Time   `json:"last_pulled"`
		LastPushed   time.Time   `json:"last_pushed"`
	} `json:"images"`
	LastUpdated         time.Time `json:"last_updated"`
	LastUpdater         int       `json:"last_updater"`
	LastUpdaterUsername string    `json:"last_updater_username"`
	Name                string    `json:"name"`
	Repository          int       `json:"repository"`
	FullSize            int       `json:"full_size"`
	V2                  bool      `json:"v2"`
	TagStatus           string    `json:"tag_status"`
	TagLastPulled       time.Time `json:"tag_last_pulled"`
	TagLastPushed       time.Time `json:"tag_last_pushed"`
}

type ByLastUpdated []DockerResult

func (a ByLastUpdated) Len() int { return len(a) }

var r, _ = regexp.Compile("v\\d")

func (a ByLastUpdated) Less(i, j int) bool {
	if r.Match([]byte(a[i].Name)) && r.Match([]byte(a[j].Name)) {
		ivers := strings.Split(a[i].Name, "v")[1][:1]
		jvers := strings.Split(a[j].Name, "v")[1][:1]
		return ivers > jvers
	}
	return false
}
func (a ByLastUpdated) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

func main() {
	repository := os.Getenv("REPOSITORY")
	container := os.Getenv("CONTAINER")

	resp, err := http.Get(fmt.Sprintf("https://registry.hub.docker.com/v2/repositories/%s/%s/tags?page_size=1024", repository, container))
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	var dr DockerResponse
	err = json.Unmarshal(data, &dr)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	sort.Sort(ByLastUpdated(dr.Results))

	if len(dr.Results) >= 2 {
		f, err := strconv.Atoi(strings.Split(dr.Results[1].Name, "v")[1])
		if err != nil {
			fmt.Printf("v1")
			return
		}
		fmt.Printf("v%v", f+1)
		return
	}

	// default to v1 if unable to work out
	fmt.Printf("v1")
}
