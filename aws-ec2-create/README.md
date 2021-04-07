---
{
  "image": "vorteil/aws-ec2-create",
  "desc": "Creates an amazon ec2 instance on aws."
}
---

# AWS EC2 Create Instance

Create a new amazon ec2 instance using the aws cli and a simple golang app to provide the authentication before executing.

## Direktiv

An example workflow that creates an example instance named `example-instance-name`.

```yaml
id: aws-create-ec2-instance
description: "Create a new ec2 instance"
functions:
- id: awsCreateInstance
  image: vorteil/aws-ec2-create:v1
  size: medium
states:
- id: createInstance
  type: action
  action:
    function: awsCreateInstance
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
    "errorCode": "com.aws-ec2-create.error",
    "errorMsg": "Something went wrong"
}
```