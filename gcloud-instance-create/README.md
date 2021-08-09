---
{
  "image": "vorteil/gcloud-instance-create",
  "desc": "Create a compute engine instance on google cloud."
}
---
# GCloud Instance Create

Creates a new google cloud compute instance using google apis and a simple golang app to provide authentication via service account key file.

## Direktiv

An example workflow that creates an example instance named `example-instance-name` in the `us-central1-a` zone.

```yaml
id: gcloud-new-compute-instance
description: "Creates new gcp project instance"
functions:
- id: gcp-new-intance
  type: reusable
  image: vorteil/gcloud-instance-create:v2
  size: medium
states:
- id: createInstance
  type: action
  action:
    function: gcp-new-intance
    input: jq(.)
```

## Input

The input needed to run the above workflow properly is the following:

```yaml
input:
  project: "GCP_PROJECT"
  serviceAccountKey: "SERVICE_ACCOUNT_KEY"
  name: "example-instance-name"
  zone: "us-central1-a"
  machineType: "e2-medium"
  image: "example_image"
```

When authenticating with the gcloud nothing will be logged to the stdout of the container.

The `SERVICE_ACCOUNT_KEY` is a string of a Google Cloud service account json key.

To generate a service account key please refer to [Google Clouds documentation](https://cloud.google.com/iam/docs/creating-managing-service-account-keys).

### Input - Optional Fields

There are additional fields that can be also used in the input for more control over the created instance.

| Field        | Description                                               | Type         | Default Value           |
|--------------|-----------------------------------------------------------|--------------|-------------------------|
| diskType     | Specify what Google Compute Engine disk type to use       | String       | pd-balanced             |
| externalIP   | Whether if instance should be created with an external IP | Boolean      | false                   |
| tags         | Tags for the created instance                             | String Array | []                      |
| imageProject | GCP project where image is stored.                        | String       | same as "project" field |


## Output
The output will be in JSON format if successful, otherwise an error will be outputed to the logs.

Below is an example output of a successful instance creations with sensitive information removed:
```json
{
  "return": {
    "id": "**********************",
    "insertTime": "2021-04-06T17:36:11.511-07:00",
    "kind": "compute#operation",
    "name": "**************************************",
    "operationType": "insert",
    "progress": 0,
    "selfLink": "https://www.googleapis.com/compute/v1/projects/****/zones/us-central1-a/operations/******************************",
    "startTime": "2021-04-06T17:36:11.513-07:00",
    "status": "RUNNING",
    "targetId": "**************************",
    "targetLink": "https://www.googleapis.com/compute/v1/projects/****/zones/us-central1-a/instances/********************",
    "user": "********************************************",
    "zone": "https://www.googleapis.com/compute/v1/projects/****/zones/us-central1-a"
  }
}
```

## Error 

In the case that an error is encountered, it will present in the following format:

```json
{
    "errorCode": "com.gcloud-instance-create.error",
    "errorMsg": "Something went wrong"
}
```