# Googlecli

Executes a gcloud cli command using a golang wrapper to provide authentication via service account key file.

## Direktiv

An example workflow of listing the instances on a gcloud project.

```yaml
id: get-instances-list
functions:
- id: post
  image: vorteil/gcloud:latest
  size: medium
description: "Displays a gcp project instances"
states:
- id: listProjectInstances
  type: action
  action:
    function: post
    input: .
```

## Input

The input needed to run the above workflow properly is the following:

```json
{
    "project": .secrets.PROJECT,
    "serviceAccountKey": .secrets.SERVICE_ACCOUNT_KEY,
    "command": ["compute", "instances", "list"]
}
```

## Output

If the action is successfuly, the gcloud container shall respond with the command execution as `json`.

## Error 

In the case that an error is encountered, it will present in the following format:

```json
{
    "errorCode": "com.googlecli.error",
    "errorMsg": "Something went wrong"
}
```