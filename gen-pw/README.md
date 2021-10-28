---
{
  "image": "direktiv/genpw",
  "desc": "Generates a random password"
}
---

# Genpw

Generates a random password.

## Direktiv

An example using this container with the following workflow on Direktiv.

```yaml
id: genpw
description: Generates a random password
functions:
- id: genpw
  image: direktiv/gen-pw:v1
states:
- id: gen-password
  type: action
  action:
    function: genpw
    input:
      length: 10
      symbols: 0
      digits: 3
```

### Input
- length: the length you want the string to be
- symbols: how many symbols do you want in the string?
- digits: how many digits do you want in the string?

### Ouput

```json
{
    "password": "8cpPy43KmR"
}
```

## Error

In the case that an error is encountered it will present in the following format:

```json
{
    "errorCode": "com.genpw.error",
    "errorMsg": "Something went wrong"
}
```