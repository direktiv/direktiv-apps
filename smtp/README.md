---
{
  "image": "direktiv/smtp",
  "desc": "A simple smtp client to send an email"
}
---

# SMTP

## Direktiv
An example workflow of using the container in a workflow on Direktiv.

```yaml
id: send-email
functions:
- id: send
  image: direktiv/smtp:v1
description: "Initializes a smtp server to send an email"
states:
- id: SendEmail
  type: action
  action:
    function: send
    input:
      to: ["recipient@example.com"]
      subject: "SMTP Email"
      template: false
      message: "Hello, world!" ## This can also be used as a base64 string for templating need to pass the template var aswell
      # Args is used to parse variables to the template args:
      from: "sender@example.com"
      user: "smtpuser"
      password: "NOT4REALPW"
      server: "smtp.example.com"
      images: ["logo.png"] ## Any images you want to embed into the body (will only work with a template)
      port: 587
```

**NOTE:** if you are using images in the template you need to provide a string array for the image names. Wherever in the HTML the image is being declared you need the src to be `cid:logo.png`

## Input

Sends an email via SMTP according to the input object:

```yaml
input:
  to: "recipient@example.com"
  subject: "SMTP Email"
  message: "Hello, world!"
  from: "sender@example.com"
  password: "NOT4REALPW"
  server: "smtp.example.com"
  port: 587
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
