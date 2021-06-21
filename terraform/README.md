---
{
    "image": "vorteil/terraform",
    "desc": "A simple container to use terraform on Direktiv."
}
---

# Terraform

## Direktiv
An example workflow of using the container on Direktiv. Notice we're providing the terraform files from the functions array.

```yaml
id: spawn-gcp-ubuntu-instance
functions:
- id: tfrun
  image: vorteil/terraform:v1
  files:
  - key: main.tf
    scope: workflow
    type: plain
  - key: variables.tf
    scope: workflow
    type: plain
  - key: outputs.tf
    scope: workflow
    type: plain
description: spawns a new gcp instance that returns the ip and user to ssh
states:
  - id: spawn-gcp-instance
    type: action
    action:
      secrets: ["GCP_SSH_PUB_KEY", "GCP_SSH_USER", "GCP_SERVICE_ACCOUNT_KEY", "GCP_PROJECT_ID", "GCP_INSTANCE_PW"]
      function: tfrun
      input: |
        {
          "action": "apply",
          "tfstate": "gcp-instance-deployed",
          "google-auth": .secrets.GCP_SERVICE_ACCOUNT_KEY,
          "variables": {
            "project_id": .secrets.GCP_PROJECT_ID,
            "gce_ssh_user": .secrets.GCP_SSH_USER,
            "gce_ssh_pub_key": .secrets.GCP_SSH_PUB_KEY,
            "password": .secrets.GCP_INSTANCE_PW
          }
        }
```

## Authentication
Three different ways of authenticating for three different cloud providers. Add these to the input section. You may use all three if you are provisioning to every cloud platform.

Azure
```json
{
    "azure-auth": {
        "client-id": "",
        "client-secret": "",
        "subscription-id": "",
        "tenant-id": ""
    }
}
```

Google
```json
{
    "google-auth": "ENTIRE_SERVICE_ACCCOUNT_KEY"
}
```

Amazon
```json
{
    "amazon-auth": {
        "access-key": "",
        "secret-key": ""
    }
}
```

## Input

```json
{
    "action": "apply",
    "tfstate": "gcp-instance-deployed",
    "google-auth": .secrets.GCP_SERVICE_ACCOUNT_KEY,
    "variables": {
        "project_id": .secrets.GCP_PROJECT_ID,
        "gce_ssh_user": .secrets.GCP_SSH_USER,
        "gce_ssh_pub_key": .secrets.GCP_SSH_PUB_KEY,
        "password": .secrets.GCP_INSTANCE_PW
    }
}
```

**NOTE: 'action' can be apply, destroy, validate or plan. 'tfstate' is the name of the state variable if it doesn't exist it will create a new workflow variable upon completion.** 

## Output

The output of the container will be matching the 'outputs.tf' file provided. For an example with the provided terraform file.

```hcl
output "ip-address" {
    value = google_compute_instance.default.network_interface[0].access_config[0].nat_ip
}
```

Will output as the following json..
```json
{
    "ip-address": {
        "value": "192.168.1.30"
    }
}
```

## Error

In the case that an error is encountered, it will present in the following format:

```json
{
    "errorCode": "com.terraform.%s.error",
    "errorMsg": "Something went wrong"
}
```

**Note: '%s' will be replaced with a contextual error.** 