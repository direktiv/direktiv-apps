---
{
  "image": "vorteil/rabbitmq",
  "desc": "Writes a message to the RabbitMQ Service"
}
---

# RabbitMQ

Writes a message to the RabbitMQ Service

## Direktiv

An example workflow of using the RabbitMQ service.

```yaml
id: rabbitmq
functions:
- id: rabbitMQWrite
  image: vorteil/rabbitmq:v2
description: "Writes a message to a RabbitMQ service"
states:
- id: writemsg
  type: action
  action:
    function: rabbitMQWrite
    input: .
```

## Input

Writes a message to the Rabbit MQ service hosted on the address.

```json
{
    "address": "192.168.1.30:5672",
    "username": "test",
    "password": "test",
    "queue": "direktiv",
    "message": "Hello World :)!"
}
```

## Output

If the output is successful nothing will be written but the message should appear in queue on RabbitMQ.

## Error

In the case that an error is encountered, it will present in the following format:

```json
{
    "errorCode": "com.rabbitmq.error",
    "errorMsg": "Something went wrong"
}
```
