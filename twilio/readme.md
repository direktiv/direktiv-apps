# Twilio
Sends an email or SMS message using Twilio.

## Input

### Email

- Generate a token [here](https://app.sendgrid.com/settings/api_keys).
- Use the following input object, substituting the values of the `token`, `to`, and `from` fields appropriately.

```json
    {
        "typeof"        : "email",
        "token"         : "TOKEN_STIRNG",
        "subject"       : "Subject Header of Email",
        "message"       : "Normal string message",
        "htmlMessage"   : "<html><body><h5>HI</h5></body></html",
        "from"          : "FROM@EXAMPLE.COM",
        "to"            : "TO@EXAMPLE.COM"   
    }
```

- To send a basic (plaintext) message, populate the `message` field and ensure that `htmlMessage` is unused. In the event that both `message` and `htmlMessage` are populated, the `htmlMessage` takes precedence. 

### SMS

- Acquire your Twilio SID, token, and 'from' number.

```json
    {
        "typeof"    : "sms",
        "sid"       : "TWILIO_SID",
        "token"     : "TWILIO_TOKEN",
        "message"   : "RECIPIENT_NAME",
        "from"      : "TWILIO_FROM_NUM",
        "to"        : "RECIPIENT_NUM"
    }
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