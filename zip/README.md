---
{
  "image": "vorteil/zip",
  "desc": "Prints all the information that direktiv sends to the container"
}
---

# Zip

Create a zip file with the files provided.

## Direktiv

An example workflow of using the container on Direktiv.

```yaml
id: zip-file
description: Listens for a smtp-event to send the email appropriately.
functions:
- id: zip
  type: reusable
  image: vorteil/zip:v1
states:
- id: zip-data
  type: action
  transform: 
    data: jq(.return)
    smtp: jq(.smtp)
    password: jq(.password)
  action:
    function: zip
    input:
      password: jq(.password)
      files: jq(.smtp.attachments)
```

### Input

- password: the string that will encrypt the zip file
- files
  - name: the name of the file
  - data: the data for the file
  - type: base64 to decode or blank

### Output

The output returns the entire file which comes back as base64 string.

## Error

In the case that an error is encountered it will present in the following format:

```json
{
    "errorCode": "com.zip.%s.error",
    "errorMsg": "Something went wrong"
}
```

Replacing '%s' contextually with where the error occurred.