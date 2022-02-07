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
      password: jq(.password) # password encrypts the zip
      files:   # files to zip
      - name: mypdf.pdf 
        data: workflow/pdffile 
        type: var
      - name: hellworld.txt
        type: plain
        data: this is plain text
      - name: base64.txt
        type: base64
        data: aHR0cHM6Ly93d3cueW91dHViZS5jb20vd2F0Y2g/dj1kUXc0dzlXZ1hjUQ==
      - name: secondpdf.pdf
        type: file
        data: /test.pdf
      scope: workflow # stores the zip as variable
      name: myzip 
      return: true # returns zip as base64 if true
```

### Input

- password: the string that will encrypt the zip file
- files
  - name: the name of the file
  - data: the data for the file
  - type: plain, var, basse64, file

### Output

Zip file as base64 in return and zip can be stored as variable.

## Error

In the case that an error is encountered it will present in the following format:

```json
{
    "errorCode": "com.zip.%s.error",
    "errorMsg": "Something went wrong"
}
```

Replacing '%s' contextually with where the error occurred.