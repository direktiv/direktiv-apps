---
{
  "image": "vorteil/direktiv-cloud-event",
  "desc": "Send a cloud event to a different namespace mid workflow."
}
---
# Direktiv Cloud Event

Send a cloud event to a different namespace mid workflow.

## Direktiv

An example workflow sending a cloud event via the container

```yaml
id: send-event
functions:
- id: sendcloud
  image: vorteil/direktiv-cloud-event:v1
description: Send a cloud event to a different namespace 
states:
- id: send
  type: action
  action: 
    function: sendcloud
    input: 
      type: "direktiv-cloud"
      source: "another workflow"
      namespace: "trent"
      "access_token": jq(.secrets.at)
      data:
        hello: "test"
```


An example workflow listening to an event the container will send based off the previous workflow.


```yaml
id: listen-for-event
description: Listen to a custom cloud event
start:
  type: event
  state: helloworld
  event:
    type: direktiv-cloud
states:
  - id: helloworld
    type: noop
    transform: 'jq({result: "Hello world!"})'
```

## Input

Providing debug to the json struct will print everything during the process it runs.

```yaml
input:
  type: "direktiv-cloud"
  source: "another workflow"
  namespace: "trent"
  data: 
    hello": "test"
  "access_token": jq(.secrets.at)
```


## Output

If the container is successful, the response will be wrapped inside of a JSON object, within the `return` field.


```json
{
    "return": {
        ...
    }
}
```

## Error

In the case that an error is encountered, it will present in the following format:

```json
{
    "errorCode": "com.cloud-event.%s.error",
    "errorMsg": "Something went wrong"
}
```

**%s denoting something wrong with the container**