---
{
  "image": "vorteil/gcloud-instance-delete",
  "desc": "Delete a compute engine instance on google cloud."
}
---
# GCloud Instance Delete

Delete a new google cloud compute instance using googe apis and a simple golang app to provide authentication via service account key file.

## Direktiv

An example workflow that deletes an example instance named `example-instance-name` in the `us-central1-a` zone.

```yaml
id: gcloud-new-compute-instance
description: "Deletes a gcp project instance"
functions:
- id: gcpDeleteInstance
  image: vorteil/gcloud-instance-delete:v1
  size: medium
states:
- id: deleteInstance
  type: action
  action:
    function: gcpDeleteInstance
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
    "errorCode": "com.gcloud-instance-delete.error",
    "errorMsg": "Something went wrong"
}
```