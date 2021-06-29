---
{
  "image": "vorteil/googleinvoke",
  "desc": "Executes a cloud function on google using a client authenticated via a service account key."
}
---

# GoogleInvoke

Executes a cloud function on google using a client authenticated via a service account key.

## Direktiv

```yaml
id: invoke-function
functions:
- id: invoke
  image: vorteil/googleinvoke:v2
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
  region: "us-east1"
  function: "helloworld"
  serviceAccountKey: jq(.secrets.SERVICE_ACCOUNT_KEY)
  method: "POST"
  body: 
    message: "hello"
```

**NOTE:** The `body` field is optional.

## Output

The output will be the response of the cloud function json marshalled for the ability to jq the response.

## Error

In the case that an error is encountered, it will present in the following format:

```json
{
    "errorCode": "com.googleinvoke.error",
    "errorMsg": "Something went wrong"
}
```