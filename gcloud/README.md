---
{
  "image": "direktiv/gcloud",
  "desc": "Executes a gcloud cli command using a golang wrapper to provide authentication via service account key file."
}
---


# Googlecli

Executes a gcloud cli command using a golang wrapper to provide authentication via service account key file.

## Direktiv

An example workflow that lists all instances in a Google cloud project

```yaml
id: get-instances-list
functions:
- id: post
  image: direktiv/gcloud:v2
  size: medium
description: "Displays a gcp project instances"
states:
- id: listProjectInstances
  type: action
  action:
    function: post
    input: jq(.)
```

## Input

The input needed to run the above workflow properly is the following:

```yaml
input:
  project: "GCP_PROJECT"
  serviceAccountKey: "SERVICE_ACCOUNT_KEY"
  command: ["compute", "instances", "list"]
```

When authenticating with the CLI nothing will be logged to the stdout of the container.

The `SERVICE_ACCOUNT_KEY` is a string of a Google Cloud service account json key.

To generate a service account key please refer to [Google Clouds documentation](https://cloud.google.com/iam/docs/creating-managing-service-account-keys).

## Output

The gcloud container will append "--format=json" to the command array when executing, so the container will respond in `json` format.


## Error 

In the case that an error is encountered, it will present in the following format:

```json
{
    "errorCode": "com.googlecli.error",
    "errorMsg": "Something went wrong"
}
```