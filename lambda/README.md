---
{
  "image": "direktiv/lambda",
  "desc": "Executes a cloud function on aws using their golang SDK."
}
---
# LambdaInvoke

Executes a cloud function on aws using their golang SDK.

## Direktiv

```yaml
id: invoke-function
functions:
- id: post
  image: direktiv/lambda:v2
description: "Invokes a cloud function based on the given details"
states:
- id: invoke-cloud
  type: action
  action:
    function: post
    input: jq(.)
```

## Input

The following input is needed for the cloud function to be invoked successfully.

```yaml
input:
  key: jq(.secrets.AWS-KEY)
  secret: jq(.secrets.AWS-SECRET)
  region: "us-east-2"
  function: "helloworld"
  body: 
    any: "data"
```

**NOTE:** The `body`  field is optional.


## Output

The output will be the response of the cloud function json marshalled for the ability to jq the response.

## Error

In the case that an error is encountered, it will present in the following format:

```json
{
    "errorCode": "com.lambdainvoke.error",
    "errorMsg": "Something went wrong"
}
```