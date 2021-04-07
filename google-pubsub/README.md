---
{
  "image": "vorteil/google-pubsub",
  "desc": "Writes a message to Google's Pubsub service."
}
---


# Google Pubsub

Writes a message to Google's Pubsub service.

## Direktiv

An example workflow of writing 'Hello World!' to the service.

```yaml
id: write-helloworld
functions:
- id: write
  image: vorteil/google-pubsub
description: "Writes 'Hello World!' to Google Pubsub service."
states:
- id: write-message
  type: action
  action:
    function: write
    input: .
```

## Input

The input required to run the above workflow properly is the following:

```json
{
    "serviceAccountKey": .secrets.SERVICE_ACCOUNT_KEY,
    "message": "Hello World!",
    "project-id": "direktiv",
    "topic-id": "listener",
    "attributes": {
        "weight": "1"
    }
}
```

**NOTE:** The `attributes` field is optional its used for filtering. It is also a map of string to string.

## Output

If the action is successful, no output will be returned.

## Error

In the case that an error is encountered, it will present in the following format:

```json
{
    "errorCode": "com.google-pubsub.error",
    "errorMsg": "Something went wrong"
}
```