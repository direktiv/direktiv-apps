---
{
  "image": "vorteil/googlemsg",
  "desc": "Writes a google message to a webhook URL."
}
---

# GoogleMessage

Writes a google message to a webhook URL.

## Direktiv

An example workflow of using the container on Direktiv.

```yaml
id: post-to-google
functions:
- id: post
  image: vorteil/googlemsg
description: "Send a message to a google chat"
states:
- id: postToGoogle
  type: action
  action:
    function: post
    input: .
```

## Input

The input needed to run the above workflow properly is the following:

```json
{
    "message": "This message will be posted to a google chat",
    "url": "GOOGLE_CHAT_WEBHOOK_URL"
}
```

## Output

If the action is successful, the google chat message should appear in the appropriate channel it was sent to.

## Error

In the case that an error is encountered, it will present in the following format:

```json
{
    "errorCode": "com.googlechat.error",
    "errorMsg": "Something went wrong"
}
```