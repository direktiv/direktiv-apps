package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/vorteil/direktiv-apps/pkg/direktivapps"
)

const code = "com.azlog.error"

type AzureLogDetails struct {
	Message     string `json:"message"`
	WorkspaceID string `json:"workspace-id"` // log analytics workspace ID
	Key         string `json:"key"`
	Type        string `json:"type"` // name of the log type you wish to record on azure
}

func main() {
	direktivapps.StartServer(AzureLog)
}

func AzureLog(w http.ResponseWriter, r *http.Request) {
	obj := new(AzureLogDetails)
	_, err := direktivapps.Unmarshal(obj, r)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	var msgData struct {
		Message string `json:"msg"`
	}
	msgData.Message = obj.Message

	// marshal message object
	data, err := json.Marshal(msgData)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	dateString := time.Now().UTC().Format(time.RFC1123)
	dateString = strings.Replace(dateString, "UTC", "GMT", -1)

	hashString := fmt.Sprintf("POST\n%v\napplication/json\nx-ms-date:%s\n/api/logs", strconv.Itoa(len(data)), dateString)
	hashedString, err := BuildSignature(hashString, obj.Key)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	signature := fmt.Sprintf("SharedKey %s:%s", obj.WorkspaceID, hashedString)
	url := fmt.Sprintf("https://%s.ods.opinsights.azure.com/api/logs?api-version=2016-04-01", obj.WorkspaceID)

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	req.Header.Add("Log-Type", obj.Type)
	req.Header.Add("Authorization", signature)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("x-ms-date", dateString)

	resp, err := client.Do(req)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}
	defer resp.Body.Close()

	bv, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		direktivapps.RespondWithError(w, code, fmt.Sprintf("Response Message: %s, Response Code: %v \nResponseBody: %s", resp.Status, resp.StatusCode, bv))
		return
	}
	direktivapps.Respond(w, bv)
}

func BuildSignature(message, secret string) (string, error) {

	keyBytes, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		return "", err
	}

	mac := hmac.New(sha256.New, keyBytes)
	mac.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil)), nil
}
