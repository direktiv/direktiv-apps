---
{
  "image": "vorteil/azgo",
  "desc": "Executes any azure cli command using a golang wrapper to provide the authentication before executing."
}
---

# Azcli

Executes any azure cli command using a golang wrapper to provide the authentication before executing.

## Direktiv

An example workflow of listing the account data tied to my service account principal.

```yaml
id: get-account-data
functions:
- id: post
  type: reusable
  image: vorteil/azgo:v2
description: "Displays account information about your azure account"
states:
- id: listAccountInfo
  type: action
  action:
    function: post
    input: jq(.)
```

## Input

The input needed to run the above workflow properly is the following:

```yaml
input:
  name: "SERVICE_PRINCIPAL_APPNAME"
  password: "SERVICE_PRINCIPAL_APPPW"
  tenant: "SERVICE_PRINCIPAL_APPTENANT"
  command: ["account", "show"]
```
When authenticating with the CLI nothing will be logged to the stdout of the container.

To generate a service principal account the following command in the azure cli should be able to.

```sh
az ad sp create-for-rbac --name ServicePrincipalName
```

Should output an object like the following,

```json
{
  "appId": "",
  "displayName": "",
  "name": "",
  "password": "",
  "tenant": ""
}
```

Name, password and tenant are the fields you need to authenticate.

Command are the arguments you wish to pass to 'az'. So when using this command to execute 'account show' you only need to provide these two strings rather than 'az' aswell.

## Output

If the action is successful, the azgo container shall respond with the command execution.

## Error

In the case that an error is encountered, it will present in the following format:

```json
{
    "errorCode": "com.azcli.error",
    "errorMsg": "Something went wrong"
}
```
