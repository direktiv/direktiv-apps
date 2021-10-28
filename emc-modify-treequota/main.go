package main

import (
	"bytes"
	"crypto/tls"
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"

	bytefmt "code.cloudfoundry.org/bytefmt"
	"github.com/direktiv/direktiv-apps/pkg/direktivapps"
	"golang.org/x/net/publicsuffix"
)

// EmcModifyInput the input object for the requester container
type EmcModifyInput struct {
	URL         string `json:"url"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	Path        string `json:"path"`
	SoftLimit   string `json:"softlimit"`
	HardLimit   string `json:"hardlimit"`
	Description string `json:"description"`
}

// TreeQuotaResponse the response from EMC for the tree quota entries
type TreeQuotaResponse struct {
	Base    string    `json:"@base"`
	Updated time.Time `json:"updated"`
	Links   []struct {
		Rel  string `json:"rel"`
		Href string `json:"href"`
	} `json:"links"`
	Entries []struct {
		Content struct {
			ID                   string `json:"id"`
			State                int    `json:"state"`
			Path                 string `json:"path"`
			Description          string `json:"description"`
			HardLimit            int64  `json:"hardLimit"`
			SoftLimit            int64  `json:"softLimit"`
			RemainingGracePeriod int    `json:"remainingGracePeriod"`
			SizeUsed             int64  `json:"sizeUsed"`
			Filesystem           struct {
				ID string `json:"id"`
			} `json:"filesystem"`
			QuotaConfig struct {
				ID string `json:"id"`
			} `json:"quotaConfig"`
		} `json:"content"`
	} `json:"entries"`
}

// EMCModifyTreeQuotaResponse returns the new size of the quota and description.
type EMCModifyTreeQuotaResponse struct {
	Hardlimit   string `json:"hardlimit"`
	Softlimit   string `json:"softlimit"`
	Description string `json:"description"`
}

// code return error
const code = "com.emc-modify-tree-quota.%s.error"

// EMCModifyTreeQuota request handler
func EMCModifyTreeQuota(w http.ResponseWriter, r *http.Request) {
	obj := new(EmcModifyInput)

	aid, err := direktivapps.Unmarshal(obj, r)
	if err != nil {
		direktivapps.RespondWithError(w, fmt.Sprintf(code, "unmarshal-input"), err.Error())
		return
	}

	jar, err := cookiejar.New(&cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	})
	if err != nil {
		direktivapps.RespondWithError(w, fmt.Sprintf(code, "cookie-jar"), err.Error())
		return
	}

	client := &http.Client{
		Jar: jar,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	direktivapps.Log(aid, fmt.Sprintf("fetch treeQuota id from path '%s'", obj.Path))
	id, hsize, ssize, err := getTreeQuotaId(client, obj.URL, aid, obj.Path, obj.Username, obj.Password)
	if err != nil {
		direktivapps.RespondWithError(w, fmt.Sprintf(code, "get-treequota-id"), err.Error())
		return
	}

	direktivapps.Log(aid, fmt.Sprintf("modify '%s' treeQuota", id))
	// work out the new size were changing to
	updateH, err := bytefmt.ToBytes(obj.HardLimit)
	if err != nil {
		direktivapps.RespondWithError(w, fmt.Sprintf(code, "convert-to-bytes"), err.Error())
		return
	}

	updateS, err := bytefmt.ToBytes(obj.SoftLimit)
	if err != nil {
		direktivapps.RespondWithError(w, fmt.Sprintf(code, "convert-to-bytes"), err.Error())
		return
	}

	newHSize := hsize + int64(updateH)
	newSSize := ssize + int64(updateS)

	direktivapps.Log(aid, fmt.Sprintf("adjusting softLimit from '%v' to '%v'", ssize, newSSize))
	direktivapps.Log(aid, fmt.Sprintf("adjusting hardlimit from '%v' to '%v'", hsize, newHSize))

	err = modifyHardLimit(client, obj.URL, aid, id, newHSize, newSSize, obj.Description, obj.Username, obj.Password)
	if err != nil {
		direktivapps.RespondWithError(w, fmt.Sprintf(code, "modify-treequota"), err.Error())
		return
	}

	var respondingOutput EMCModifyTreeQuotaResponse
	respondingOutput.Hardlimit = bytefmt.ByteSize(uint64(newHSize))
	respondingOutput.Softlimit = bytefmt.ByteSize(uint64(newSSize))
	if obj.Description != "" {
		respondingOutput.Description = obj.Description
	}

	returnData, err := json.Marshal(respondingOutput)
	if err != nil {
		direktivapps.RespondWithError(w, fmt.Sprintf(code, "unmarshal-response"), err.Error())
		return
	}

	direktivapps.Respond(w, returnData)
}

func modifyHardLimit(client *http.Client, urlpath string, aid, id string, newHSize int64, newSSize int64, description, username, password string) error {
	u, err := url.Parse(fmt.Sprintf("%s//api/instances/treeQuota/%s/action/modify", urlpath, id))
	if err != nil {
		return err
	}

	var modifyArgs string
	if description != "" {
		modifyArgs = fmt.Sprintf(`{
			"hardLimit": %v,
			"softLimit": %v,
			"description": "%s"
		}`, newHSize, newSSize, description)
	} else {
		modifyArgs = fmt.Sprintf(`{
			"hardLimit": %v,
			"softLimit": %v
		}`, newHSize, newSSize)
	}

	req, err := http.NewRequest("POST", u.String(), bytes.NewReader([]byte(modifyArgs)))
	if err != nil {
		return err
	}

	direktivapps.Log(aid, "Adding required headers")

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-EMC-REST-CLIENT", "true")

	direktivapps.Log(aid, "Adding authorization header")
	sEnc := b64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", username, password)))
	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", sEnc))

	direktivapps.Log(aid, "Sending request")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func getTreeQuotaId(client *http.Client, urlpath string, aid string, path string, username, password string) (string, int64, int64, error) {
	u, err := url.Parse(fmt.Sprintf("%s/api/types/treeQuota/instances", urlpath))
	if err != nil {
		return "", 0, 0, err
	}

	q := u.Query()
	q.Set("fields", "id,filesystem,quotaConfig,path,description,state,hardLimit,softLimit,remainingGracePeriod,sizeUsed")
	q.Set("compact", "true")

	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return "", 0, 0, err
	}

	direktivapps.Log(aid, "adding required headers")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-EMC-REST-CLIENT", "true")

	sEnc := b64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", username, password)))
	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", sEnc))

	resp, err := client.Do(req)
	if err != nil {
		return "", 0, 0, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", 0, 0, err
	}

	var tq TreeQuotaResponse
	err = json.Unmarshal(data, &tq)
	if err != nil {
		return "", 0, 0, err
	}

	for _, t := range tq.Entries {
		if path == t.Content.Path {
			return t.Content.ID, t.Content.HardLimit, t.Content.SoftLimit, nil
		}
	}

	return "", 0, 0, errors.New("unable to find tree quota id")
}

func main() {
	direktivapps.StartServer(EMCModifyTreeQuota)
}
