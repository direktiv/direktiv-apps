---
{
  "image": "direktiv/kafka",
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
- id: kafka-write
  image: direktiv/kafka:v1
description: "Writes a message to a Kafka service"
states:
- id: writemsg
  type: action
  action:
    function: kafka-write
    input: jq(.)
```

## Input

Writes a message to a Kafka service hosted on the address.

```yaml
input:
  topic: "quickstart-events"
  partition: 0
  message: "Hello World :)!"
  address: "localhost:9092"
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
