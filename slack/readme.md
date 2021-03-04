## Direktiv
An example workflow of using the container in a workflow on Direktiv.

```yaml
id: send-message
functions:
- id: send
  image: vorteil/slack
description: "Sends the input" 
states:
- id: send-msg
  type: action
  action:
    function: send
    input: '{
        "url": .secrets.SLACK_HELLO,
        "message" : "this message gets posted to slack url"
    }'
```

## Post a message to a Slack channel

- In Slack, create an app and assign it to your team.
- Use `vorteil/slack` to send the following request, substituting values where necessary.

```json
{
   "url": "SLACK_WEBHOOK_URL",
   "message": "This message gets posted in the slack channel the webhook is attached to"
}
```

The action will return a Status Code 200 if successful.