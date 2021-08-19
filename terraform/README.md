---
{
    "image": "vorteil/terraform",
    "desc": "A simple container to use terraform on Direktiv."
}
---

# Terraform

## Direktiv
An example workflow of using the container on Direktiv. Notice we're providing the terraform file from the functions array. Which puts it into a temp directory where we call Terraform from.

```yaml
id: spawn-gcp-ubuntu-instance
functions:
- id: tfrun
  image: vorteil/terraform:v1
  files:
  - key: main.tf
    scope: workflow
    type: plain
description: spawns a new gcp instance that returns the ip
states:
  - id: spawn-gcp-instance
    type: action
    action:
      secrets: ["GCP_PROJECT_ID"]
      function: tfrun
      input: 
        action: "apply"
        "args-on-init":
          - "-backend-config=address=http://localhost:8001/terraform-gcp-instance"
        variables:
          "state-name": "terraform-gcp-instance"
          "project_id": jq(.secrets.GCP_PROJECT_ID)
```

Passing the above 'args-on-init' and variable 'state-name' will allow the container to use a http backend for Terraform to store the state. If not provided the state will be returned in the JSON output.

**Note: the variable 'state-name' needs to match the path on the 'args-on-init' path.**

## Output

If the action is successful, the output should contain any outputs written in Terraform and the Terraform state.

```json
{
    "output": {},
    "tfstate": {}
}
```

## Error

In the case that an error is encountered it will present in the following format:

```json
{
    "errorCode": "com.terraform.%s.error",
    "errorMsg": "Something went wrong"
}
```

**Note: '%s' will be replaced with a contextual error.**
