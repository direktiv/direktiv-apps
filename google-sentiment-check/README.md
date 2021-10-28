---
{
  "image": "direktiv/google-sentiment-check",
  "desc": "Reads a string and tells you the sentiment of the written text."
}
---

# Google Sentiment Check

Reads a string and tells you the sentiment of the written text.

## Direktiv

An example workflow of reading the sentiment of a text message.

```yaml
id: read-message
functions:
- id: read
  image: direktiv/google-sentiment-check:v1
states:
- id: readmsg
  type: action
  action:
    function: read
    input: jq(.)
```

## Input

The input required to run the above workflow properly is the following:

```yaml
input:
  message: "Hello World!"
  serviceAccountKey: jq(.secrets.SERVICE_ACCOUNT_KEY)
```

## Output

If the action is successful the following output will be returned.

```json
{
    "feeling": "Positive",
    "score": 0.6,
    "magnitude": 3.0,   
}
```

## Error

In the case that an error is encountered, it will present in the following format:

```json
{
    "errorCode": "com.google-sentiment-check.error",
    "errorMsg": "Something went wrong"
}
```