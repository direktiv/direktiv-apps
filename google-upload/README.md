---
{
  "image": "direktiv/google-upload",
  "desc": "Uploads a base64 encoded string to a bucket on Google."
}
---

# Google Upload

Uploads a base64 encoded string to a bucket on Google.

## Direktiv

An example workflow of uploading a base64 encoded string to Google.

```yaml
id: upload-object
functions:
- id: upload
  image: direktiv/google-upload:v2
description: "Uploads a base64 encoded string to bucket storage"
states:
- id: upload-file
  type: action
  action:
    function: upload
    input: jq(.)
```

## Input

The input required to run the above workflow properly is the following:

```yaml
input:
  bucket: "direktiv"
  data: "base64 decoded string"
  "upload-name": "test"
  serviceAccountKey: jq(.secrets.GOOGLE_SERVICE_ACCOUNT)
```

## Output

If the action is successful not output is returned.

## Error

In the case that an error is encountered, it will present in the following format:

```json
{
    "errorCode": "com.google-upload.error",
    "errorMsg": "Something went wrong"
}
```