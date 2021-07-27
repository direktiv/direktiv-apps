---
{
  "image": "vorteil/teams",
  "desc": "Sends a JSON body to a webhook provided by Microsoft Teams"
}
---

# Teams

Sends a JSON body to a webhook provided by Microsoft Teams.

## Direktiv

An example workflow of sending `Hello World!` as plain text:

```yaml
id: post-to-teams
functions:
- id: post
  image: vorteil/teams:v1
description: Sends a JSON body to a webhook provided by Microsoft Teams.
states:
- id: write-message
  type: action
  action:
    function: post
    input:
      body:
        text: Hello World!
      url: https://URL_TO_WEBHOOK
```

Another example workflow that sends `Hello World!` in an AdaptiveCard.

```yaml
id: post-to-teams
functions:
- id: post
  image: vorteil/teams:v1
description: Sends a JSON body to a webhook provided by Microsoft Teams.
states:
- id: write-message
  type: action
  action:
    function: post
    input:
      body:
        type: message
        attachments:
        - contentType: application/vnd.microsoft.card.adaptive
          content:
            $schema: http://adaptivecards.io/schemas/adaptive-card.json
            type: AdaptiveCard
            version: "1.2"
            body:
            - type: TextBlock
              text: Hello World!
      url: https://URL_TO_WEBHOOK
```


**NOTE:** The `body` field can contain alot more to align with adaptive cards or a simple text object that contains a string.

## Output of Execution

The output of the execution is the response of the `POST` request to the webhook.

## Error on Execution

If the container was to error out during the action an error will be thrown matching:

```json
{
    "errorCode": "com.teams.%s.error",
    "errorMsg": "Something went wrong..."
}
```

**NOTE:** `%s` will be replaced with what is currently happening in the container that errored out e.g. `unmarshalling-input`.