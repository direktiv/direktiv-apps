# LambdaInvoke

Executes a cloud function on aws using their golang SDK.

## Direktiv

```yaml
id: invoke-function
functions:
- id: invoke
  image: vorteil/lambda
description: "Invokes a cloud function based on the given details"
states:
- id: invoke-cloud
  type: action
  action:
    function: post
    input: .
```

## Input

The following input is needed for the cloud function to be invoked successfully.

```json
{
    "key": .secrets.AWS-KEY,
    "secret": .secrets.AWS-SECRET,
    "region": "us-east-2",
    "function": "helloworld",
    "body": {
        "any": "data"
    }
}
```

**NOTE:** The `body`  field is optional.


## Output

The output will be the response of the cloud function json marshalled for the ability to jq the response.

## Error

In the case that an error is encoutnered, it will present in the following format:

```json
{
    "errorCode": "com.lambdainvoke.error",
    "errorMsg": "Something went wrong"
}
```