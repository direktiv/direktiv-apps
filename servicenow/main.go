package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	da "github.com/vorteil/direktiv-apps/pkg/direktivapps"
)

type input struct {
	User               string                 `json:"user"`
	Password           string                 `json:"password"`
	Instance           string                 `json:"instance"`
	API                APIPrefix              `json:"api"`
	InsecureSkipVerify bool                   `json:"insecureSkipVerify"`
	Payload            map[string]interface{} `json:"payload"`
	Method             string                 `json:"method"`
	SysID              string                 `json:"sys_id"`
}

type APIPrefix string

var (
	API_CMDB_Instance        = APIPrefix("cmdb-instance")
	API_ChangeRequest_Normal = APIPrefix("change-request-normal")
	API_ChangeRequest        = APIPrefix("change-request")
	API_Incident             = APIPrefix("incident")
	API_User                 = APIPrefix("user")
)

var apiPathMap = map[APIPrefix]string{
	API_CMDB_Instance:        "api/now/cmdb/instance/cmdb_ci_storage_server",
	API_ChangeRequest_Normal: "api/sn_chg_rest/change/normal",
	API_ChangeRequest:        "api/sn_chg_rest/change",
	API_Incident:             "api/now/table/incident",
	API_User:                 "api/now/table/sys_user",
}

func main() {
	da.StartServer(coreLogic)
}

func coreLogic(w http.ResponseWriter, r *http.Request) {

	aid := r.Header.Get(da.DirektivActionIDHeader)

	in := new(input)
	err := json.NewDecoder(r.Body).Decode(in)
	if err != nil {
		da.RespondWithError(w, "servicenow.input.parseError", err.Error())
		return
	}

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: in.InsecureSkipVerify}
	c := http.DefaultClient

	apiPath, ok := apiPathMap[in.API]
	if !ok {
		da.RespondWithError(w, "servicenow.api.unsupported", fmt.Sprintf("unsupported api '%s'", in.API))
		return
	}

	var body io.Reader
	if in.Payload != nil {
		if len(in.Payload) != 0 {
			b, err := json.Marshal(in.Payload)
			if err != nil {
				da.RespondWithError(w, "servicenow.input.payloadParseError", "unable to marshal payload")
				return
			}

			body = bytes.NewReader(b)
		}
	}

	u := fmt.Sprintf("%s/%s", in.Instance, apiPath)
	if in.SysID != "" {
		if !strings.HasPrefix(in.SysID, "?") {
			u += "/"
		}
		u = fmt.Sprintf("%s%s", u, in.SysID)
	}

	req, err := http.NewRequest(in.Method, u, body)
	if err != nil {
		da.RespondWithError(w, "servicenow.request.init", fmt.Sprintf("unable to prepare api request: %s", err.Error()))
		return
	}
	if body != nil {
		defer req.Body.Close()
	}

	req.SetBasicAuth(in.User, in.Password)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Do(req)
	if err != nil {
		da.RespondWithError(w, "servicenow.request.send", fmt.Sprintf("failed to send api request: %s", err.Error()))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode < 300 && resp.StatusCode >= 200 {
		_, err = io.Copy(w, resp.Body)
		if err != nil {
			da.RespondWithError(w, "servicenow.response.write", fmt.Sprintf("unable to write response: %s", err.Error()))
			return
		}

		return
	}

	x, _ := ioutil.ReadAll(resp.Body)
	da.Log(aid, fmt.Sprintf("%s\n", x))
	da.RespondWithError(w, "servicenow.response.code", fmt.Sprintf("%d - %s", resp.StatusCode, http.StatusText(resp.StatusCode)))
	return
}
