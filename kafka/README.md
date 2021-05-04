---
{
  "image": "vorteil/kafka",
  "desc": "Writes a message to a Kafka service"
}
---

# Kafka

Writes a message to a Kafka service


## Direktiv

An example workflow of using the Kafka container.

```yaml
id: kafka
functions:
- id: kafkaWrite
  image: vorteil/kafka:v2
description: "Writes a message to a Kafka service"
states:
- id: writemsg
  type: action
  action:
    function: req
    input: .
```

## Input

Writes a message to a Kafka service hosted on the address.

```json
{
        "topic": "quickstart-events",
        "partition": 0,
        "message": "Hello World :)!",
        "address": "localhost:9092"
}
```

## Output

If the output is successful nothing will be written but the message should appear when listening to that topic.

## Error

In the case that an error is encountered, it will present in the following format:

```json
{
    "errorCode": "com.kafka.error",
    "errorMsg": "Something went wrong"
}
```
