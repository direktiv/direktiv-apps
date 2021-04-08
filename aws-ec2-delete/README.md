---
{
  "image": "vorteil/aws-ec2-delete",
  "desc": "Delete a amazon ec2 instance."
}
---

# AWS EC2 Delete Instance

Delete a amazon ec2 instance using the aws cli and a simple golang app to provide the authentication before executing.

## Direktiv

An example workflow that deletes a ec2 instance named `example-instance-name`.

```yaml
id: aws-delete-ec2-instance
description: "Delete a ec2 instance"
functions:
- id: awsDeleteInstance
  image: vorteil/aws-ec2-delete:v1
  size: medium
states:
- id: deleteInstance
  type: action
  action:
    function: awsDeleteInstance
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

If the action is successful, the container will finish.

## Error

In the case that an error is encountered, it will present in the following format:

```json
{
    "errorCode": "com.aws-ec2-delete.error",
    "errorMsg": "Something went wrong"
}
```