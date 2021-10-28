---
{
  "image": "direktiv/aws-ec2-start",
  "desc": "Start a amazon ec2 instance."
}
---

# AWS EC2 Start Instance

Start a amazon ec2 instance using the aws cli and a simple golang app to provide the authentication before executing.

## Direktiv

An example workflow that starts a ec2 instance with the id `i-090dd9255a7aaa582`.

```yaml
id: aws-start-ec2-instance
description: "Start a ec2 instance"
functions:
- id: aws-start-instance
  image: direktiv/aws-ec2-start:v1
  size: medium
states:
- id: startInstance
  type: action
  action:
    secrets: ["ACCESS_KEY", "ACCESS_SECRET"]
    function: aws-start-instance
    input: jq(.)
```

## Input

The input needed to run the above workflow properly is the following:

```yaml
input:
  "access-key": jq(.secrets.ACCESS_KEY)
  "access-secret": jq(.secrets.ACCESS_SECRET)
  "region": "us-east-1"
  "instance-id": "i-090dd9255a7aaa582"
```

## Output

The output will be in JSON format if successful, otherwise an error will be outputed to the logs.

Below is an example output of a successful start:
```json
{
  "return": {
    "StartingInstances": [
      {
        "CurrentState": {
          "Code": 0,
          "Name": "pending"
        },
        "InstanceId": "i-090dd9255a7aaa582",
        "PreviousState": {
          "Code": 80,
          "Name": "stopped"
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
    "errorCode": "com.aws-ec2-start.error",
    "errorMsg": "Something went wrong"
}
```