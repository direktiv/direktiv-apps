---
{
  "image": "vorteil/debug",
  "desc": "Prints all the information that direktiv sends to the container"
}
---

# Debug

Prints all the information that direktiv sends to the container

## Direktiv

An example workflow of using the container on direktiv

```yaml
id: debugging
functions:
- id: debug
  image: vorteil/debug:v2
description: "Prints all input the container receives"
states:
- id: print
  type: action
  action:
    function: print
    input: .
```

## Input & Output

The input will get printed as output for debugging purposes.

```json
{
    "hello":"world"
}
```

## Error

In the case that an error is encountered it will present in the following format:

```json
{
    "errorCode": "com.debug.error",
    "errorMsg": "Something went wrong"
}
```