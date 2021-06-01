---
{
  "image": "vorteil/google-terraform",
  "desc": "Executes a terraform script."
}
---

# google-Terraform

Executes a terraform script provided by a URL.

## Direktiv

```yaml
id: create-vm
functions:
- id: terraform
  image: vorteil/google-terraform:v1
description: "Executes a terraform script and performs the actions on google"
states:
- id: run
  type: action
  secrets: ["key"]
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
    "service-account-key": .secrets.key,
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