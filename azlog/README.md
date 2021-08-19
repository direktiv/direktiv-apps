---
{
  "image": "vorteil/azlog",
  "desc": "Writes a log line to Azure Log Analytics Workspace."
}
---

# AzureLog

Writes a log line to Azure Log Analytics Workspace.

## Direktiv

An example workflow of using the container on Direktiv.

```yaml
id: log-to-azure
functions:
- id: log
  image: vorteil/azlog:v2
description: "Log a message to Azure Analytics Workspace"
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
  message: "Hello"
  "workspace-id": jq(.secrets.WORKSPACE_ID)
  key: jq(.secrets.WORKSPACE_KEY)
  type: "direktiv"
}
```

**NOTE:** The `type` field is used so you can query based from it in the Analytics Workspace.

## Output

If the request is successful. Not output will be returned and the logs should appeared if queried on Azure.

## Error

In the case that an error is encountered, it will present in the following format:

```json
{
    "errorCode": "com.azlog.error",
    "errorMsg": "Something went wrong"
}
```