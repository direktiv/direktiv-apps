---
{
  "image": "vorteil/rabbitmq-consumer",
  "desc": "Consumes messages from a RabbitMQ Service"
}
---

# RabbitMQ Consumer

Consumes messages from a RabbitMQ Service

## Direktiv

An example workflow of consuming from a RabbitMQ service.

```yaml
id: rabbitmq
functions:
- id: rabbit-consume
  image: vorteil/rabbitmq-consumer:v1
description: "Read messages from a RabbitMQ service"
states:
- id: readmsg
  type: action
  action:
    function: rabbit-consume
    input: .
```

## Input

Consume messages to the Rabbit MQ service hosted on the address. The property `waitTime` is optional(default=1000) and is the amount of time in milliseconds to consume messages.

```json
{
    "address": "192.168.1.30:5672",
    "username": "test",
    "password": "test",
    "queue": "direktiv",
    "waitTime": 1000
}
```

## Output

If the output is successful, the messages consumed will be returned as a array with a timestamp.

```json
{
	"return": [
		{
			"message": "Hello World :)!",
			"timestamp": "0001-01-01 00:00:00 +0000 UTC"
		}
	]
}
```


## Error

In the case that an error is encountered, it will present in the following format:

```json
{
    "errorCode": "com.rabbitmq-consumer.error",
    "errorMsg": "Something went wrong"
}
```
