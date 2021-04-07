---
{
  "image": "vorteil/awsgo",
  "desc": "Run an amazon ec2 instance on aws."
}
---

# AWS EC2 Run Instance

Run a new amazon ec2 instance using the aws cli and a simple golang app to provide the authentication before executing.

## Direktiv

An example workflow that run an example instance named `example-instance-name`.

```yaml
id: aws-run-ec2-instance
description: "Run a new ec2 instance"
functions:
- id: awsRunInstance
  image: vorteil/aws-ec2-run:v1
  size: medium
states:
- id: runInstance
  type: action
  action:
    function: awsRunInstance
    input: .
```

## Input

The input needed to run the above workflow properly is the following:

```json
{
    "access-key": .secrets.ACCESS_KEY,
    "access-secret": .secrets.ACCESS_SECRET,
    "image-id": "ami-0abcdef1234567890",
    "instance-type": "t2.micro",
    "key-name": "MyKeyPair",
}
```

### Input - Optional Fields

There are additional fields that can be also used in the input for more control over the created instance.

| Field              | Description                                       | Type              |
|--------------------|---------------------------------------------------|-------------------|
| subnet-id          | The ID of the subnet to launch the instance into. | String            |
| security-group-ids | The IDs of the security groups.                   | String            |
| tag-specifications | Tags for the created instance.                    | Map[string]string |

## Output

If the action is successful, the container will respond with details on the operation in JSON format.

## Error

In the case that an error is encountered, it will present in the following format:

```json
{
    "errorCode": "com.aws-ec2-run.error",
    "errorMsg": "Something went wrong"
}
```