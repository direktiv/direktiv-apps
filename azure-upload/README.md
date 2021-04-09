---
{
  "image": "vorteil/azure-upload",
  "desc": "Uploads a base64 encoded string to a blob on Azure."
}
---

# Azure Upload

Uploads a base64 encoded string to a blob on Azure.

## Direktiv

An example workflow of uploading a base64 encoded string to Azure.

```yaml
id: upload-object
functions:
- id: upload
  image: vorteil/azure-upload:v1
description: "Uploads a base64 encoded string to blob storage"
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
    "container": "azure-container",
    "storage-account": "storage-account",
    "storage-account-key": .secrets.AZURE_STORAGE_KEY,
    "data": "base64 decoded string",
    "upload-name": "test"
}
```

## Output

If the action is successful the output will be the response from the Azure upload.

## Error

In the case that an error is encountered, it will present in the following format:

```json
{
    "errorCode": "com.azure-upload.error",
    "errorMsg": "Something went wrong"
}
```