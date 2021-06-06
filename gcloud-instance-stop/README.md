---
{
  "image": "vorteil/gcloud-instance-stop",
  "desc": "Stop a compute engine instance on google cloud."
}
---
# GCloud Instance Stop

Stop a google cloud compute instance using googe apis and a simple golang app to provide authentication via service account key file.

## Direktiv

An example workflow that stops an example instance named `example-instance-name` in the `us-central1-a` zone.

```yaml
id: gcloud-stop-compute-instance
description: "Stop a gcp project instance"
functions:
- id: gcp-stop-instance
  image: vorteil/gcloud-instance-stop:v2
  size: medium
states:
- id: stopInstance
  type: action
  action:
    function: gcp-stop-instance
    secrets: ["GCP_SERVICE_ACCOUNT_KEY"]
    input: .
```

## Input

The input needed to run the above workflow properly is the following:

```json
{
       "project": "GCP_PROJECT",
       "serviceAccountKey": "SERVICE_ACCOUNT_KEY",
       "instanceID": "example-instance-name",
       "zone": "us-central1-a"
}
```

When authenticating with the gcloud nothing will be logged to the stdout of the container.

The `SERVICE_ACCOUNT_KEY` is a string of a Google Cloud service account json key.

To generate a service account key please refer to [Google Clouds documentation](https://cloud.google.com/iam/docs/creating-managing-service-account-keys).

## Output
The output will be in JSON format if successful, otherwise an error will be outputed to the logs.

## Error 

In the case that an error is encountered, it will present in the following format:

```json
{
    "errorCode": "com.gcloud-instance-stop.error",
    "errorMsg": "Something went wrong"
}
```