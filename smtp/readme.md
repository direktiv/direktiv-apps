# SMTP

## Input

A simple go application that sends an email via SMTP client. The following json input can be used to send an email.

```json
{
    "from": "",
    "to": "",
    "subject": "SMTP Email",
    "message": "pog champ",
    "server": "smtp.gmail.com",
    "port": 587,
    "password": ""
}
```

## Output

No output is provided upon success.

### Error

If an error was to happen during the container running the output is written back to direktiv as the following

```json
{
    "errorCode": "com.smtp.error",
    "errorMsg": "Something went wrong"
}
```