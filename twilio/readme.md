# Twilio

## Input
A simple go application that either sends an email or sms using Twilio infrastructure.

### Email
To send an email using twilio it requires generating a token from [here](https://app.sendgrid.com/settings/api_keys) and sending the following input. (The email used to create this token is the email you will be sending the "from" tag from)

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

One of the message fields needs to be used. htmlMessage takes priority over the other.


### SMS
To send an sms you will need your twilio sid and token. Twilio should also provide you with a number that they use to send the sms messages.

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

## Output

Both requests here return no output from the response if successful.

If an error was to happen during the container running the output is written back to direktiv as the following

```json
{
    "errorCode": "com.request.error",
    "errorMsg": "Something went wrong"
}
```