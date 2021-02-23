# SMTP


## Direktiv
An example workflow of using the container in a workflow on Direktiv.

```yaml
id: send-email
functions:
- id: send
  image: vorteil/smtp
description: "Initializes a smtp server to send an email" 
states:
- id: SendEmail
  type: action
  action:
    function: send
    input: '{
    "to"      : "recipient@example.com",
    "subject" : "SMTP Email",
    "message" : "Hello, world!",
    "from"    : "sender@example.com",
    "password": "NOT4REALPW",
    "server"  : "smtp.example.com",
    "port"    : 587
}'
```

## Input

Sends an email via SMTP according to the input object:

```json
{
    "to"      : "recipient@example.com",
    "subject" : "SMTP Email",
    "message" : "Hello, world!",
    "from"    : "sender@example.com",
    "password": "NOT4REALPW",
    "server"  : "smtp.example.com",
    "port"    : 587
}
```

## Output

No output is provided on success.

## Error

In the case that an error is encountered, it will present in the following format:

```json
{
    "errorCode": "com.smtp.error",
    "errorMsg": "Something went wrong"
}
```