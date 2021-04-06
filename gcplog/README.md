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
  image: vorteil/gcplog:v1
description: "Log a message to google cloud"
states:
- id: log-hello
  type: action
  action:
    function: log
    input: .
```

## Input

The input needed to run the above workflow properly is the following:

```json
{
    "serviceAccountKey": .secrets.SERVICEACCOUNTKEY,
    "message": "hello",
    "project-id": "direktiv",
    "log-name": "NameOfLog"
}
```

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