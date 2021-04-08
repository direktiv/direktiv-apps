---
{
  "image": "vorteil/aws-ec2-start",
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
- id: awsStartInstance
  image: vorteil/aws-ec2-start:v1
  size: medium
states:
- id: startInstance
  type: action
  action:
    function: awsStartInstance
    input: .
```

## Input

The input needed to run the above workflow properly is the following:

```json
{
    "access-key": .secrets.ACCESS_KEY,
    "access-secret": .secrets.ACCESS_SECRET,
    "region": "us-east-1",
    "instance-id": "i-090dd9255a7aaa582",
}
```

## Output

The output will be in JSON format if successful, otherwise an error will be outputed to the logs.

## Error

In the case that an error is encountered, it will present in the following format:

```json
{
    "errorCode": "com.aws-ec2-start.error",
    "errorMsg": "Something went wrong"
}
```