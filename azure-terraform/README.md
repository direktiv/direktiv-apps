---
{
  "image": "vorteil/azure-terraform",
  "desc": "Executes a terraform script."
}
---

# Azure-Terraform

Executes a terraform script provided by a URL.

## Direktiv

```yaml
id: create-vm
functions:
- id: terraform
  image: vorteil/azure-terraform:v1
description: "Executes a terraform script and performs the actions on azure"
states:
- id: run
  type: action
  secrets: ["client-id", "client-secret", "sub-id", "tenant-id"]
  action:
    function: terraform
    input: .
```

## Input

The input required to run the above workflow properly is the following:

```json
{
    "tf": "https://github.com/trentis/test.tf",
    "action": "apply",
    "client_id": .secrets.client-id,
    "client_secret": .secrets.client-secret,
    "subscription_id": .secrets.sub-id,
    "tenant_id": .secrets.tenant-id,
    "tfstate": "",
}
```

**Leave tfstate as an empty string when running the first terraform execution**

## Output

If the action is successful the output will be the tf state.

## Error

In the case that an error is encountered, it will present in the following format:

```json
{
    "errorCode": "com.terraform.error",
    "errorMsg": "Something went wrong"
}
```