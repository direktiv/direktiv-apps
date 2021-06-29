---
{
  "image": "vorteil/aws-ec2-stop",
  "desc": "Stops a amazon ec2 instance."
}
---

# AWS EC2 Stop Instance

Stops a amazon ec2 instance using the aws cli and a simple golang app to provide the authentication before executing.

## Direktiv

An example workflow that stops a ec2 instance with the id `i-090dd9255a7aaa582`.

```yaml
id: aws-stop-ec2-instance
description: "Stop a ec2 instance"
functions:
- id: aws-stop-instance
  image: vorteil/aws-ec2-stop:v2
  size: medium
states:
- id: stopInstance
  type: action
  action:
    secrets: ["ACCESS_KEY", "ACCESS_SECRET"]
    function: aws-stop-instance
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

Below is an example output of a successful stop:
```json
{
  "return": {
    "StoppingInstances": [
      {
        "CurrentState": {
          "Code": 64,
          "Name": "stopping"
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
    "errorCode": "com.aws-ec2-stop.error",
    "errorMsg": "Something went wrong"
}
```