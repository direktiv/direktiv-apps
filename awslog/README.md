# AWSLog

Writes a log line to aws cloudwatch logs with provided log stream and group.

## Direktiv

An example workflow of using the container on Direktiv.

```yaml
id: log-to-aws
functions:
- id: log
  image: vorteil/awslog:v1
description: "Log a message to AWS cloudwatch"
states:
- id: log-hello
  type: action
  action:
    function: log
    input: .
```

## Input

The input needed to run the above workflow properly is the following:

```json
{
    "key": .secrets.AWS_KEY,
    "secret": .secrets.AWS_SECRET,
    "region": "us-east-2",
    "log-group": "vorteil",
    "log-stream": "direktiv",
    "message": "Hello"
}
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