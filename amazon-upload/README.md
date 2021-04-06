+++
{
  "image": "vorteil/amazon-upload",
  "desc": "Uploads a base64 encoded string to a blob on Amazon."
}
+++

# Amazon Upload

Uploads a base64 encoded string to a blob on Amazon.

## Direktiv

An example workflow of uploading a base64 encoded string to Amazon.

```yaml
id: upload-object
functions:
- id: upload
  image: vorteil/amazon-upload:v1
description: "Uploads a base64 encoded string to Amazon"
states:
- id: upload-file
  type: action
  action:
    function: upload
    input: .
```

## Input

The input required to run the above workflow properly is the following:

```json
{
    "bucket": "test",
    "region": "ap-southeast-2",
    "upload-name": "upload-name",
    "key": .secrets.AMAZON_KEY,
    "secret": .secrets.AMAZON_SECRET,
    "data": "base64 encoded string"
}
```

## Output

If the action is successful the no output will be printed.

## Error

In the case that an error is encountered, it will present in the following format:

```json
{
    "errorCode": "com.amazon-upload.error",
    "errorMsg": "Something went wrong"
}
```