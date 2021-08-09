---
{
  "image": "vorteil/kafka-consume",
  "desc": "Consume messages from a Kafka service"
}
---

# Kafka

Consumes messages from a Kafka service


## Direktiv

An example workflow of consuming from a Kafka service.

```yaml
id: kafka
functions:
- id: kafka-consume
  type: reusable
  image: vorteil/kafka-consume:v2
description: "Read messages from a Kafka service"
states:
- id: readmsg
  type: action
  action:
    function: kafka-consume
    input: jq(.)
```

## Input

Consume messages from a Kafka service hosted on the address. The property `timeout` is optional(default=1000) and is the amount of time in milliseconds to consume messages.

```yaml
input:
  topic: "quickstart-events"
  partition: 0
  address: "localhost:9092"
  timeout: 1000
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
    "errorCode": "com.kafka-consume.error",
    "errorMsg": "Something went wrong"
}
```
