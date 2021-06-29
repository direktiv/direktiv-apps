---
{
  "image": "vorteil/gcplog",
  "desc": "Writes a log line to a stackdriver logging implementation"
}
---

# GCPLog

Writes a log line to a stackdriver logging implementation

## Direktiv

An example workflow of using the container on Direktiv.

```yaml
id: log-to-google
functions:
- id: log
  image: vorteil/gcplog:v2
description: "Log a message to google cloud"
states:
- id: log-hello
  type: action
  action:
    function: log
    input: jq(.)
```

## Input

The input needed to run the above workflow properly is the following:

```yaml
input:
  serviceAccountKey: jq(.secrets.SERVICEACCOUNTKEY)
  message: "hello"
  "project-id": "direktiv"
  "log-name": "NameOfLog"
  "log-level": "Info"
```

## Input Optional - Log Level
The field `log-level` is an optional field to set the serverity level of the log being sent.

If the field `log-level` is omitted it will be set to the GCP Log Level Default.

The supported GCP log levels are: "Default", "Debug", "Info", "Notice", "Warning", "Error", "Critical", "Alert", "Emergency"

## Output

If the action is successful, no output will be logged to the direktiv container. But the log should appear under the metrics on Google Cloud.

## Error

In the case that an error is encountered, it will present in the following format:

```json
{
    "errorCode": "com.gcplog.error",
    "errorMsg": "Something went wrong"
}
```