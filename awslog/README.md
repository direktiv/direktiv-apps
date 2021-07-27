---
{
  "image": "vorteil/awslog",
  "desc": "Writes a log line to aws cloudwatch logs with provided log stream and group."
}
---

# AWSLog

Writes a log line to aws cloudwatch logs with provided log stream and group.

## Direktiv

An example workflow of using the container on Direktiv.

```yaml
id: log-to-aws
functions:
- id: log
  image: vorteil/awslog:v2
description: "Log a message to AWS cloudwatch"
states:
- id: log-hello
  type: action
  action:
    function: log
    input: jq(.)
```

## Input

The input needed to run the above workflow properly is the following:

```yaml
input:
  key: jq(.secrets.AWS_KEY)
  secret: jq(.secrets.AWS_SECRET)
  region: "us-east-2"
  "log-group": "vorteil"
  "log-stream": "direktiv"
  message: "Hello"
```

## Output

If the action is successful, no output will be logged to the direktiv container. But the log should appear under the log stream on AWS.

## Error

In the case that an error is encountered, it will present in the following format:

```json
{
    "errorCode": "com.awslog.error",
    "errorMsg": "Something went wrong"
}
```