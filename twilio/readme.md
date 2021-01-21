# Twilio

Takes multiple arguments to be able to send emails and sms. The following example below is what an email input file would look like.

```json
    {
        "typeof": "email",
        "token": "GENERATED SENDGRID API KEY",
        "subject": "Subject Header of Email",
        "message": "Normal string message",
        "htmlMessage": "<html><body><h5>HI</h5></body></html",
        "from": "trent.hilliam@vorteil.io",
        "to": "trent.hilliam@vorteil.io"   
    }
```

One of the message fields need to be used. htmlMessage takes priority over the other.

Another type of input field we could use is for messaging mobiles via Twilio. Notice the extra field "sid" we need to use this as this functionality requires a different authentication.

```json
    {
        "typeof": "sms",
        "sid": "TWILIO SID",
        "token": "TWILIO TOKEN",
        "message": "George Smith",
        "from": "+15097742730",
        "to": "+61435545810"
    }
```

The output file will look something like

```json
{
    "Error": "",
    "Response": "",
    "Status": "",
    "StatusCode": ""
}
```

Error is only used if the request encountered an error the process
Response, Status and Status code are all returned from the API we are requesting