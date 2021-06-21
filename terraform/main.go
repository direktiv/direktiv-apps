package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"sync"
	"syscall"
	"time"

	"github.com/vorteil/direktiv-apps/pkg/direktivapps"
)

// TerraformInput takes different arguments for each authentication service.
type TerraformInput struct {
	Action                    string                 `json:"action"`    // plan, validate, apply, destroy
	TFVars                    map[string]interface{} `json:"variables"` // the tf variables from the input.
	TFStateName               string                 `json:"tfstate"`   // the name of the variable uses for tfstate
	GoogleCloudAuthentication string                 `json:"google-auth,omitempty"`
	AzureAuthentication       AzureAuthentication    `json:"azure-auth,omitempty"`
	AmazonAuthentication      AmazonAuthentication   `json:"amazon-auth,omitempty"`
}

// AzureAuthentication provides clientID, client secret, subscriptionid and tenantid
type AzureAuthentication struct {
	ClientID       string `json:"client-id"`
	ClientSecret   string `json:"client-secret"`
	SubscriptionID string `json:"subscription-id"`
	TenantID       string `json:"tenant-id"`
}

// AmazonAuthentication provides the accesskey and secret key
type AmazonAuthentication struct {
	AccessKey string `json:"access-key"`
	SecretKey string `json:"secret-key"`
}

var code = "com.terraform.%s.error"
var terraformBin = "/terraform"

var aidGlobal string
var tfnameGlobal string
var locked = false

// CMDWriter allows us to log to the action id provided.
type CMDWriter struct {
	Aid string `json:"aid"`
}

// Write writes to the new api for logging direktiv apps
func (c *CMDWriter) Write(p []byte) (n int, err error) {
	direktivapps.Log(c.Aid, string(p))
	return len(p), nil
}

func main() {
	direktivapps.StartServer(TerraformHandler)
}

func setupGCPAuth(auth string) error {
	err := ioutil.WriteFile("/credentials", []byte(auth), os.ModePerm)
	if err != nil {
		return err
	}
	err = os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "/credentials")
	if err != nil {
		return err
	}
	return nil
}

func setupAzureAuth(auth *AzureAuthentication) error {
	err := os.Setenv("ARM_CLIENT_ID", auth.ClientID)
	if err != nil {
		return err
	}
	err = os.Setenv("ARM_CLIENT_SECRET", auth.ClientSecret)
	if err != nil {
		return err
	}
	err = os.Setenv("ARM_SUBSCRIPTION_ID", auth.SubscriptionID)
	if err != nil {
		return err
	}
	err = os.Setenv("ARM_TENANT_ID", auth.TenantID)
	if err != nil {
		return err
	}
	return nil
}

func setupAmazonAuth(auth *AmazonAuthentication) error {
	err := os.Setenv("AWS_ACCESS_KEY_ID", auth.AccessKey)
	if err != nil {
		return err
	}

	err = os.Setenv("AWS_SECRET_ACCESS_KEY", auth.SecretKey)
	if err != nil {
		return err
	}

	return nil
}

func httpBackend() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", TFStateHandler)

	srv := &http.Server{
		Addr:    ":8001",
		Handler: mux,
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM)

	go func() {
		<-sigs
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		srv.Shutdown(ctx)
	}()

	srv.ListenAndServe()
}

var stateAid sync.Mutex

func TFStateHandler(w http.ResponseWriter, r *http.Request) {
	url := fmt.Sprintf("http://localhost:8889/var?aid=%s&scope=workflow&key=%s", aidGlobal, tfnameGlobal)

	switch r.Method {
	case http.MethodGet:
		direktivapps.Log(aidGlobal, "Fetching tfstate variable...")
		resp, err := http.Get(url)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		defer resp.Body.Close()
		data, _ := ioutil.ReadAll(resp.Body)
		w.Write(data)
	case http.MethodPost:
		direktivapps.Log(aidGlobal, "Saving new tfstate variable...")
		req, err := http.NewRequest("POST", url, r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			direktivapps.RespondWithError(w, fmt.Sprintf(code, "resp-do-tfstate"), err.Error())
			return
		}
		defer resp.Body.Close()
		return
	case http.MethodDelete:
		req, err := http.NewRequest("POST", url, bytes.NewReader([]byte("")))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			direktivapps.RespondWithError(w, fmt.Sprintf(code, "resp-do-tfstate"), err.Error())
			return
		}
		defer resp.Body.Close()
		return
	}
}

func checkBackendIsAlive() bool {
	timeout := 1 * time.Second
	_, err := net.DialTimeout("tcp", "localhost:8001", timeout)
	if err != nil {
		return false
	}
	return true
}

func TerraformHandler(w http.ResponseWriter, r *http.Request) {

	obj := new(TerraformInput)

	aid, err := direktivapps.Unmarshal(obj, r)
	if err != nil {
		direktivapps.RespondWithError(w, fmt.Sprintf(code, "unmarshal-input"), err.Error())
		return
	}

	if locked {
		direktivapps.RespondWithError(w, fmt.Sprintf(code, "inuse"), "terraform container in use")
		return
	}
	// wait till the server is alive
	stateAid.Lock()
	aidGlobal = aid
	tfnameGlobal = obj.TFStateName
	locked = true

	direktivapps.Log(aid, "Finding path to call terraform from...")
	terraformPath := r.Header.Get("Direktiv-TempDir")
	direktivapps.Log(aid, fmt.Sprintf("Found '%s'", terraformPath))

	// url := fmt.Sprintf("http://localhost:8889/var?aid=%s&scope=workflow&key=terraform.state", aid)
	direktivapps.Log(aid, "Checking if tfstate service http backend is alive...")
	alive := checkBackendIsAlive()
	// if backend not alive spawn backend
	if !alive {
		go httpBackend()
	}

	direktivapps.Log(aid, "Wait till backend service is functional")
	for !alive {
		alive = checkBackendIsAlive()
	}

	direktivapps.Log(aid, "Handling Authentication...")
	// Handle authentication with different cloud providers
	// check if google auth exists and if so add it
	if obj.GoogleCloudAuthentication != "" {
		direktivapps.Log(aid, "Adding Google Cloud Authentication")
		err = setupGCPAuth(obj.GoogleCloudAuthentication)
		if err != nil {
			direktivapps.RespondWithError(w, fmt.Sprintf(code, "gcp-auth"), err.Error())
			return
		}
	}

	// check if azure auth exists and if so add it
	if obj.AzureAuthentication != (AzureAuthentication{}) {
		direktivapps.Log(aid, "Adding Azure Authentication")
		err = setupAzureAuth(&obj.AzureAuthentication)
		if err != nil {
			direktivapps.RespondWithError(w, fmt.Sprintf(code, "azure-auth"), err.Error())
			return
		}
	}

	// check if amazon auth exists and if so add it
	if obj.AmazonAuthentication != (AmazonAuthentication{}) {
		direktivapps.Log(aid, "Adding Amazon Authentication")
		err = setupAmazonAuth(&obj.AmazonAuthentication)
		if err != nil {
			direktivapps.RespondWithError(w, fmt.Sprintf(code, "amazon-auth"), err.Error())
			return
		}
	}

	direktivapps.Log(aid, "Initializing terraform....")

	cmdW := &CMDWriter{
		Aid: aid,
	}
	dirArg := fmt.Sprintf("-chdir=%s", terraformPath)

	direktivapps.Log(aid, fmt.Sprintf("Reading in TFVars.json..."))

	data, err := json.Marshal(obj.TFVars)
	if err != nil {
		direktivapps.RespondWithError(w, fmt.Sprintf(code, "tfvars-marshal"), err.Error())
		return
	}

	err = ioutil.WriteFile(path.Join(terraformPath, "terraform.tfvars.json"), data, os.ModePerm)
	if err != nil {
		direktivapps.RespondWithError(w, fmt.Sprintf(code, "tfvars"), err.Error())
		return
	}

	init := exec.Command(terraformBin, dirArg, "init")
	init.Stderr = cmdW
	init.Stdout = cmdW
	err = init.Run()
	if err != nil {
		direktivapps.RespondWithError(w, fmt.Sprintf(code, "init"), err.Error())
		return
	}

	direktivapps.Log(aid, fmt.Sprintf("Executing '%s' for terraform", obj.Action))
	switch obj.Action {
	case "apply":
		fallthrough
	case "destroy":
		cmd := exec.Command(terraformBin, dirArg, obj.Action, "-auto-approve")
		cmd.Stdout = cmdW
		cmd.Stderr = cmdW
		err = cmd.Run()
		if err != nil {
			direktivapps.RespondWithError(w, fmt.Sprintf(code, "apply-destroy"), err.Error())
			return
		}
	case "validate":
		fallthrough
	case "plan":
		cmd := exec.Command(terraformBin, dirArg, obj.Action)
		cmd.Stdout = cmdW
		cmd.Stderr = cmdW
		err = cmd.Run()
		if err != nil {
			direktivapps.RespondWithError(w, fmt.Sprintf(code, "validate-plan"), err.Error())
			return
		}
	}

	direktivapps.Log(aid, "Sending output back to direktiv...")
	cmd := exec.Command(terraformBin, dirArg, "output", "-json")
	data, err = cmd.Output()
	if err != nil {
		direktivapps.RespondWithError(w, fmt.Sprintf(code, "output"), err.Error())
		return
	}

	stateAid.Unlock()
	locked = false

	direktivapps.Respond(w, data)
}
