---
{
  "image": "direktiv/sendgrid",
  "desc": "A simple Twilio sendgrid client"
}
---

# Sendgrid

## Direktiv
An example workflow of using the container in a workflow on Direktiv. The message content can be defined by three different values:

- message: plain text message
- html-message: html content as message
- template & template-data: *template* defines a template in [golang templating format](https://pkg.go.dev/text/template) and template-data defines the input data for the template.

*Templating example*
```yaml
template: "hello {{.name}}! {{range .items}}{{.}} {{end}}"
template-data:
  name: Testname
  items:
  - one
  - two
```

*Workflow Example*
```yaml
description: send emails
functions:
- id: mail
  image: direktiv/sendgrid
  type: reusable
states:
- id: sendemail
  type: action
  action:
    function: mail
    secrets: ["SENDGRID"]
    input:
      sender-name: John Sender
      sender-email: sender@sendit.xyz
      subject: "test"
      message: "this is a message"
      recv-name: Jim Receiver
      recv-email: receiver@abcdefgh.abc
      apikey: jq(.secrets.SENDGRID)

```

## Input

Sends an email via Sendgrid according to the input object:

```yaml
input:
  sender-name: John Sender
  sender-email: sender@sendit.xyz
  subject: "test"
  message: "this is a message"
  recv-name: Jim Receiver
  recv-email: receiver@abcdefgh.abc
  apikey: jq(.secrets.SENDGRID)
```

## Output

No output is provided on success.

## Error

In the case that an error is encountered, it will present in the following format:

```json
{
    "errorCode": "com.sendgrid.error",
    "errorMsg": "Something went wrong"
}
```
