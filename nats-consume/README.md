---
{
  "image": "vorteil/nats-consume",
  "desc": "Consumes messages from a NATS Service"
}
---

# NATS

Consumes messages from a NATS Service

## Direktiv

An example workflow of consuming from a NATS service.

```yaml
id: nats
functions:
- id: nats-consume
  image: vorteil/nats-consume:v1
description: "Read messages from a NATS service"
states:
- id: natsread
  type: action
  action:
    function: nats-consume
    input: .
```

## Input

Writes a message to the NATS service hosted on the address. The fields `username` and `password` are optional fields used for user authentication. 

The property `timeout` is optional(default=1000) and is the amount of time in milliseconds to consume messages.

The property `max-messages` is optional(default=10) and is the amount of messages to wait to be consumed.

This app will consume messages until either the `max-messages` are consumed or the timeout is reached.


```json
{
    "address": "localhost:4222",
    "username": "test",
    "password": "test",
    "subject": "direktiv",
    "timeout": 1000,
    "max-messages": 10
}
```

## Output

If the output is successful, the messages consumed will be returned as a string array.

```json
{
	"return": {
		"messages": [
			"Hello World :)!",
			"Goodbye"
		]
	}
}
```

## Error

In the case that an error is encountered, it will present in the following format:

```json
{
    "errorCode": "com.nats-consume.error",
    "errorMsg": "Something went wrong"
}
```
