---
{
  "image": "direktiv/go",
  "desc": "Reads a variable in and go builds"
}
---

# Go

Able to compile go binaries

## Direktiv

An example workflow of using the container on Direktiv.

```yaml
id: go-build
functions:
- id: go
  image: direktiv/go:v1
  files:
  - key: helloworld
    scope: instance
    type: tar.gz
description: Builds a go binary from the files
states:
- id: go
  type: action
  action:
    function: go
    input:
      args: ["build", "-o", "helloworldserver"]
      execution-folder: helloworld
      variable: helloworldserver
      variable-type: instance
```

## Input

The input required to run the container needs a json struct like the following:

```json
{
    "args": ["build", "-o", "helloworldserver"],
    "execution-folder": "helloworld",
    "variable": "helloworldserver",
    "variable-type": "instance"
}
```

## Output

The output of the binary gets saved as the `variable` and `variable-type`.

## Error

In the case that an error is encountered it will present in the following format:

```json
{
    "errorCode": "com.go.error",
    "errorMsg": "Something went wrong"
}
```