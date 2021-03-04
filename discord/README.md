# DiscordMessage

Writes a discord message to a webhook URL.

## Direktiv

An example workflow of using the container on Direktiv.

```yaml
id: post-to-discord
functions:
- id: post
  image: vorteil/discordmsg
description: "Sends a message to a discord channel"
states:
- id: postToDiscord
  type: action
  action:
    function: post
    input: .
```

## Input 

The input needed to run the above workflow properly is the following:

```json
{
    "message": "",
    "tts": false,
    "url": "DISCORD_CHAT_WEBHOOK_URL"
}
```

**NOTE:** The `tts` field is optional. Its whether you want discord to announce the message e.g. talk to speech.

## Output

If the action is successful, the discord message should appear in the appropriate channel it was sent to. 

## Error

In the case that an error is encountered, it will present in the following format:

```json
{
    "errorCode": "com.discord.error",
    "errorMsg": "Something went wrong"
}
```