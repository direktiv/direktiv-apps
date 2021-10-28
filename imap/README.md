---
{
  "image": "direktiv/imap",
  "desc": "reads the first message of an email inbox."
}
---

# IMAP

Reads the first message of an email inbox.

## Direktiv

An example workflow of using the container on Direktiv.

```yaml
id: imap-listener
functions:
- id: imap
  image: direktiv/imap:v1
description: Reads the first message of an email inbox
states:
- id: read-msg
  type: action
  action: 
    function: imap
    secrets: ["EMAIL", "PASSWORD"]
    input:
      email: jq(.secrets.EMAIL)
      password: jq(.secrets.PASSWORD)
      imap-address: imap.gmail.com:993
```

## Input

The input required to run the cotnainer needs a json struct like the following:

```json
{
    "email": "jq(.secrets.EMAIL)",
    "password": "jq(.secrets.PASSWORD)",
    "imap-address": "imap.gmail.com:993"
}
```

## Output

Returns the contents of an email body

## Error

In the case that an error is encountered it will present in the following format:

```json
{
    "errorCode": "com.imap.error",
    "errorMsg": "Something went wrong"
}
```