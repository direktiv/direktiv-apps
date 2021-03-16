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

type AzureLogDetails struct {
	Message     string `json:"message"`
	WorkspaceID string `json:"workspace-id"` // log analytics workspace ID
	Key         string `json:"key"`
	Type        string `json:"type"` // name of the log type you wish to record on azure
}

func main() {

	g := direktivapps.ActionError{
		ErrorCode:    "com.azlog.error",
		ErrorMessage: "",
	}

	obj := new(AzureLogDetails)
	direktivapps.ReadIn(obj, g)

	var msgData struct {
		Message string `json:"msg"`
	}
	msgData.Message = obj.Message

	// marshal message object
	data, err := json.Marshal(msgData)
	if err != nil {
		g.ErrorMessage = err.Error()
		direktivapps.WriteError(g)
	}

	dateString := time.Now().UTC().Format(time.RFC1123)
	dateString = strings.Replace(dateString, "UTC", "GMT", -1)

	hashString := fmt.Sprintf("POST\n%v\napplication/json\nx-ms-date:%s\n/api/logs", strconv.Itoa(len(data)), dateString)
	hashedString, err := BuildSignature(hashString, obj.Key)
	if err != nil {
		g.ErrorMessage = err.Error()
		direktivapps.WriteError(g)
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
		g.ErrorMessage = err.Error()
		direktivapps.WriteError(g)
	}

	req.Header.Add("Log-Type", obj.Type)
	req.Header.Add("Authorization", signature)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("x-ms-date", dateString)

	resp, err := client.Do(req)
	if err != nil {
		g.ErrorMessage = err.Error()
		direktivapps.WriteError(g)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		g.ErrorMessage = fmt.Sprintf("Response Message: %s, Response Code: %v \nResponseBody: %s", resp.Status, resp.StatusCode, body)
		direktivapps.WriteError(g)
	}

	bv, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		g.ErrorMessage = err.Error()
		direktivapps.WriteError(g)
	}

	direktivapps.WriteOut(bv, g)
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
