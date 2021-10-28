---
{
  "image": "direktiv/azinvoke",
  "desc": "Executes a cloud function on azure using the function name, function app and function key as authentication to do so."
}
---

# AzureInvoke

Executes a cloud function on azure using the function name, function app and function key as authentication to do so. 

## Direktiv

```yaml
id: invoke-function
functions:
- id: invoke
  image: direktiv/azinvoke:v1
description: "Invokes a cloud function based on the given details"
states:
- id: invoke-cloud
  type: action
  action:
    function: invoke
    input: jq(.)
```

## Input 

The following input is needed for the cloud function to be invoked successfully.

```yaml
input:
  "function-app": "hello-direktiv"
  "function-name": "direktivTrigger"
  "function-key": jq(.secrets.HELLO_DIREKTIV_KEY)
  body:
    data: "sent to function"
    x: "y"
```

**NOTE:** The `body` field is optional.

## Output

The output will be the response of the cloud function json marshalled for the ability to jq the response.

## Error

In the case that an error is encountered, it will present in the following format:

```json
{
    "errorCode": "com.azureinvoke.error",
    "errorMsg": "Something went wrong"
}
```