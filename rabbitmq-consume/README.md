---
{
  "image": "vorteil/rabbitmq-consume",
  "desc": "Consumes messages from a RabbitMQ Service"
}
---

# RabbitMQ Consume

Consumes messages from a RabbitMQ Service

## Direktiv

An example workflow of consuming from a RabbitMQ service.

```yaml
id: rabbitmq
functions:
- id: rabbit-consume
  type: reusable
  image: vorteil/rabbitmq-consume:v1
description: "Read messages from a RabbitMQ service"
states:
- id: readmsg
  type: action
  action:
    function: rabbit-consume
    input: jq(.)
```

## Input

Consume messages from a Rabbit MQ service hosted on the address. The property `timeout` is optional(default=1000) and is the amount of time in milliseconds to consume messages.

```yaml
input:
  address: "192.168.1.30:5672"
  username: "test"
  password: "test"
  queue: "direktiv"
  timeout: 1000
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
    "errorCode": "com.rabbitmq-consume.error",
    "errorMsg": "Something went wrong"
}
```
