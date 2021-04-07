---
{
  "image": "vorteil/awsgo",
  "desc": "Executes any cli command with aws using a golang wrapper to provide the authentication before executing."
}
---

# Awscli

Executes any cli command with aws using a golang wrapper to provide the authentication before executing.

## Direktiv

An example workflow of listing the iam account summary tied to the access key and secret that was provided.

```yaml
id: get-account-data
functions:
- id: post
  image: vorteil/awsgo
description: "Displays account information about aws"
states:
- id: listAccountSummary
  type: action
  action:
    function: post
    input: .
```

## Input

The input needed to run the above workflow properly is the following:

```json
{
    "access-key": .secrets.ACCESS_KEY,
    "access-secret": .secrets.ACCESS_SECRET,
    "command": ["iam", "get-account-summary"]
}
```

## Output

If the action is successful, the awsgo container shall respond with the command execution.

## Error

In the case that an error is encountered, it will present in the following format:

```json
{
    "errorCode": "com.awscli.error",
    "errorMsg": "Something went wrong"
}
```