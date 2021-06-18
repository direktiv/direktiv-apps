---
{
  "image": "vorteil/aws-ec2-create",
  "desc": "Creates an amazon ec2 instance on aws."
}
---

# AWS EC2 Create Instance

Create a new amazon ec2 instance using the aws cli and a simple golang app to provide the authentication before executing.

## Direktiv

An example workflow that creates an example instance from the ami `ami-0abcdef1234567890`.

```yaml
id: aws-create-ec2-instance
description: "Create a new ec2 instance"
functions:
- id: aws-create-instance
  image: vorteil/aws-ec2-create:v3
  size: medium
states:
- id: createInstance
  type: action
  action:
    secrets: ["ACCESS_KEY", "ACCESS_SECRET"]
    function: aws-create-instance
    input: .
```

## Input

The input needed to run the above workflow properly is the following:

```json
{
    "access-key": .secrets.ACCESS_KEY,
    "access-secret": .secrets.ACCESS_SECRET,
    "image-id": "ami-0abcdef1234567890",
    "region": "us-east-1",
    "instance-type": "t2.micro",
}
```

### Input - Optional Fields

There are additional fields that can be also used in the input for more control over the created instance.

| Field              | Description                                       | Type              |
|--------------------|---------------------------------------------------|-------------------|
| name               | The name of the instance being created            | String            |
| subnet-id          | The ID of the subnet to launch the instance into. | String            |
| key-name           | The name of the key pair. Warning: If not set you can't connect to the instance unless you choose an AMI that is configured to allow users another way to log in. | String            |
| security-group-ids | The IDs of the security groups.                   | String            |
| tag-specifications | Tags for the created instance.                    | Map[string]string |

## Output
The output will be in JSON format if successful, otherwise an error will be outputed to the logs.

Below is an example trimmed down output of a successful instance creation with sensitive information removed:
```json
{
  "return": {
    ...
    "Instances": [
      {
        ...
        "Hypervisor": "xen",
        "ImageId": "ami-0abcdef1234567890",
        "InstanceId": "i-04303562db5f8335f",
        "InstanceType": "t2.micro",
        "LaunchTime": "2021-04-04T00:19:23.000Z",
        "MetadataOptions": {
          "HttpEndpoint": "enabled",
          "HttpPutResponseHopLimit": 1,
          "HttpTokens": "optional",
          "State": "pending"
        },
        ...
      }],
  }
}
```

## Error

In the case that an error is encountered, it will present in the following format:

```json
{
    "errorCode": "com.aws-ec2-create.error",
    "errorMsg": "Something went wrong"
}
```