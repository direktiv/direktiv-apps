
# GCloud Instance Create

Creates a new google cloud compute instance a golang wrapper to provide authentication via service account key file.

## Direktiv

An example workflow that creates an example instance named `example-instance-name` in the `us-central1-a` zone.

```yaml
id: gcloud-new-compute-instance
description: "Creates new gcp project instance"
functions:
- id: gcpNewInstance
  image: vorteil/gcloud-instance-create:v1
  size: medium
states:
- id: createInstance
  type: action
  action:
    function: gcpNewInstance
    input: .
```

## Input

The input needed to run the above workflow properly is the following:

```json
{
       "project": "GCP_PROJECT",
       "serviceAccountKey": "SERVICE_ACCOUNT_KEY",
       "name": "example-instance-name",
       "zone": "us-central1-a",
       "machineType": "e2-medium",
       "image": "example_image",
}
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


## Error 

In the case that an error is encountered, it will present in the following format:

```json
{
    "errorCode": "com.gcloud-instance-create.error",
    "errorMsg": "Something went wrong"
}
```