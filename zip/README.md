---
{
  "image": "direktiv/zip",
  "desc": "Prints all the information that direktiv sends to the container"
}
---

# Zip

Create a zip file with the files provided.

## Direktiv

An example workflow of using the container on Direktiv.

```yaml
id: zip-file
functions:
- id: zip
  image: direktiv/zip:v3
states:
- id: zip-data
  type: action
  action:
    function: zip
    input:
      password: jq(.password)
      files: jq(.smtp.attachments)
      scope: workflow
      name: myzip
```

### Input

- password: the string that will encrypt the zip file
- files
  - name: the name of the file
  - data: the data for the file
  - type: base64 to decode or blank

### Output

Empty but the action stores the zip file under the provided name and scope. 

## Error

In the case that an error is encountered it will present in the following format:

```json
{
    "errorCode": "com.zip.%s.error",
    "errorMsg": "Something went wrong"
}
```

Replacing '%s' contextually with where the error occurred.