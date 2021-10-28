---
{
  "image": "direktiv/smtp-bare",
  "desc": "Similar to direktiv/smtp but right now does not allow templating via HTML."
}
---

# SMTP-bare

Send a message via SMTP.

## Direktiv

An example workflow of using this container on Direktiv.

```yaml
id: new-smtp-event
description: Send an email
functions:
- id: send
  image: direktiv/smtp-bare:v1
states:
- id: test-input
  type: action
  action:
    secrets: ["EMAIL_PW"]
    function: send
    input:
      tls: true
      to: 'jq(.smtp.to)'
      from: EMAIL_ADDDRESS
      subject: hello world
      password: jq(.secrets.EMAIL_PW)
      address: smtp.gmail.com
      port: 587
      attachments:
      - name: attach.zip
        data: jq(.data)
        type: base64
```

### Input

- tls: enable tls encryption
- to: an array of email addresses to send a message
- from: an email address that you want the message sent from
- password: password required for the email address you are sending from.
- subject: the subject header in the email
- address: the address of the smtp server you want to use
- port: the port of the smtp server
- attachments
  - name: name of the attachment
  - data: the data of the attachment
  - type: 'base64' attachment or leave blank 

### Output

If successful no output will be returned when the message has been sent.

## Error

In the case that an error is encountered it will present in the following format:

```json
{
    "errorCode": "com.smtp-bare.%s.error",
    "errorMsg": "Something went wrong"
}
```

Replaces '%s' contextually with an error that happens during the runtime.