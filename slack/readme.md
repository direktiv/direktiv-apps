# Post a message to a Slack channel

- In Slack, create an app and assign it to your team.
- Use `vorteil/request` to send the following request, substituting values where necessary.

```json
{
    "method"    : "POST",
    "host"      : "SLACK_WEBHOOK_URL",
    "body"      : {
        "text"  : "Hello, world!"
    },
    "headers"   : {
        "Content-Type" : "application/json"
    }
}
```

The action will return a Status Code 200 if successful.