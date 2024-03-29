---
{
  "image": "direktiv/twilio",
  "desc": "Sends an email or SMS message using Twilio."
}
---

# Twilio
Sends an email or SMS message using Twilio.

## Direktiv
An example workflow of using the container in a workflow on Direktiv.


```yaml
id: send-approval
functions:
- id: send
  image: direktiv/twilio:v1
description: "Sends a sms to provided number" 
states:
- id: fire
  type: action
  action:
    secrets: ["TWILIO_TOKEN", "TWILIO_SID", "TWILIO_PROVIDED_NUMBER"]
    function: send
    input:
      typeof: "sms"
      sid: jq(.secrets.TWILIO_SID)
      token: jq(.secrets.TWILIO_TOKEN)
      message: "A fun SMS message!"
      from: jq(.secrets.TWILIO_PROVIDED_NUMBER)
      to: jq("+" + (.number|tostring))
```

## Input

### Email

- Generate a token [here](https://app.sendgrid.com/settings/api_keys).
- Use the following input object, substituting the values of the `token`, `to`, and `from` fields appropriately.
- Adding debug to the json struct will output more of the application.
```yaml
input:
  typeof: "email"
  token: "TOKEN_STRING"
  subject: "Subject Header of Email"
  message: "Normal string message"
  htmlMessage: "<html><body><h5>HI</h5></body></html"
  from: "FROM@EXAMPLE.COM"
  to: "TO@EXAMPLE.COM"   
```

- To send a basic (plaintext) message, populate the `message` field and ensure that `htmlMessage` is unused. In the event that both `message` and `htmlMessage` are populated, the `htmlMessage` takes precedence. 

### SMS

- Acquire your Twilio SID, token, and 'from' number.

```yaml
input:
  typeof: "sms"
  sid: "TWILIO_SID"
  token: "TWILIO_TOKEN"
  message: "RECIPIENT_NAME"
  from: "TWILIO_FROM_NUM"
  to: "RECIPIENT_NUM"
```

## Output

No output is returned upon success.
In the case that an error is encountered, it will present in the following format:

```json
{
    "errorCode" : "com.request.error",
    "errorMsg"  : "Something went wrong"
}
```