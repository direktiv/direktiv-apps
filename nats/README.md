---
{
  "image": "vorteil/nats",
  "desc": "Writes a message to a NATS Service"
}
---

# NATS

Writes a message to a NATS Service

## Direktiv

An example workflow of writing to a NATS service.

```yaml
id: nats
functions:
- id: nats-write
  type: reusable
  image: vorteil/nats:v1
description: "Writes a message to a NATS service"
states:
- id: natswrite
  type: action
  action:
    function: nats-write
    input: jq(.)
```

## Input

Writes a message to the NATS service hosted on the address. The fields `username` and `password` are optional fields used for user authentication.

```yaml
input:
  address: "localhost:4222"
  username: "test"
  password: "test"
  subject: "direktiv"
  message: "Hello World :)!"
```


## Output

If the output is successful nothing will be written but the message should appear in queue on NATS.

## Error

In the case that an error is encountered, it will present in the following format:

```json
{
    "errorCode": "com.nats.error",
    "errorMsg": "Something went wrong"
}
```
