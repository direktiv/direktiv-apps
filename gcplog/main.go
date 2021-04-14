package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"cloud.google.com/go/logging"
	"github.com/vorteil/direktiv-apps/pkg/direktivapps"
	"github.com/vorteil/direktiv-apps/pkg/validator"
	"google.golang.org/api/option"
)

type GCPLogging struct {
	ServiceAccountKey string `json:"serviceAccountKey" validate:"required"`
	Message           string `json:"message" validate:"required"`
	ProjectID         string `json:"project-id" validate:"required"`
	LogName           string `json:"log-name" validate:"required"`
	LogLevel          string `json:"log-level"`
}

const code = "com.gcplog.error"

var severityList = []string{"Default", "Debug", "Info", "Notice", "Warning", "Error", "Critical", "Alert", "Emergency"}

func main() {
	direktivapps.StartServer(GCPLog)
}

func GCPLog(w http.ResponseWriter, r *http.Request) {
	obj := new(GCPLogging)
	aid, err := direktivapps.Unmarshal(obj, r)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	// Validate Input
	v := validator.CreateValidator()

	if missingFields := v.ValidateRequired(obj); len(missingFields) > 0 {
		for _, mf := range missingFields {
			direktivapps.Log(aid, fmt.Sprintf("Input Error: %s is required\n", mf))
		}

		direktivapps.RespondWithError(w, code, fmt.Sprintf("Invalid input: Fields [%s] are required", strings.Join(missingFields, ",")))
		return
	}

	var validLevel bool

	// Set default severity level
	if obj.LogLevel == "" {
		obj.LogLevel = severityList[0]
		validLevel = true
	}

	// Check Severity is valid
	for i := range severityList {
		if obj.LogLevel == severityList[i] || validLevel {
			validLevel = true
			break
		}
	}

	if !validLevel {
		direktivapps.RespondWithError(w, code, fmt.Sprintf("log-level: '%s' is not valid. VAlID Levels: %s", obj.LogLevel, strings.Join(severityList, ", ")))
		return
	}

	// Create GCP Client
	ctx := context.Background()
	client, err := logging.NewClient(ctx, obj.ProjectID, option.WithCredentialsJSON([]byte(obj.ServiceAccountKey)))
	if err != nil {
		direktivapps.RespondWithError(w, code, fmt.Sprintf("GCP Client Error: %v", err))
		return
	}
	defer client.Close()

	// Create GCP Logger
	logger := client.Logger(obj.LogName)

	// Set Logger Severity
	targetLogger := logger.StandardLogger(logging.ParseSeverity(obj.LogLevel))
	targetLogger.Println(obj.Message)

	direktivapps.Log(aid, "Writing Log Message to gcp")
	// Flush Logger
	if err := logger.Flush(); err != nil {
		direktivapps.RespondWithError(w, code, fmt.Sprintf("GCP Logger Error: %v", err))
		return
	}

	// Write empty to notified its finished
	direktivapps.Respond(w, []byte{})
}
