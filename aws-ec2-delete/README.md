---
{
  "image": "vorteil/aws-ec2-delete",
  "desc": "Delete a amazon ec2 instance."
}
---

# AWS EC2 Delete Instance

Delete a amazon ec2 instance using the aws cli and a simple golang app to provide the authentication before executing.

## Direktiv

An example workflow that deletes a ec2 instance with the id `i-090dd9255a7aaa582`.

```yaml
id: aws-delete-ec2-instance
description: "Delete a ec2 instance"
functions:
- id: aws-delete-instance
  image: vorteil/aws-ec2-delete:v2
  size: medium
states:
- id: deleteInstance
  type: action
  action:
    secrets: ["ACCESS_KEY", "ACCESS_SECRET"]
    function: aws-delete-instance
    input: jq(.)
```

## Input

The input needed to run the above workflow properly is the following:

```yaml
input:
  "access-key": jq(.secrets.ACCESS_KEY)
  "access-secret": jq(.secrets.ACCESS_SECRET)
  region: "us-east-1"
  "instance-id": "i-090dd9255a7aaa582"
```

## Output
The output will be in JSON format if successful, otherwise an error will be outputed to the logs.

Below is an example output of a successful deletion:
```json
{
  "return": {
    "TerminatingInstances": [
      {
        "CurrentState": {
          "Code": 32,
          "Name": "shutting-down"
        },
        "InstanceId": "i-090dd9255a7aaa582",
        "PreviousState": {
          "Code": 16,
          "Name": "running"
        }
      }
    ]
  }
}
```

## Error

In the case that an error is encountered, it will present in the following format:

```json
{
    "errorCode": "com.aws-ec2-delete.error",
    "errorMsg": "Something went wrong"
}
```