# Post Message to slack channel

- Create an app assign it to your team for slack.
- Use the requester program to send the following input to hit your application to push a message
```json
{
    "method": "POST",
    "host": "URL TO YOUR SLACK WEBHOOK",
    "body" : {
        "text":"helloworld test"
    },
    "headers" : {
        "Content-type": "application/json"
    }
}
```


- The response should be something along the lines of the follow json struct and you should see the message appear in your slack channel.

```json
{
    "error":"",
    "response":"ok",
    "statusCode":200,
    "status":"200 OK"
}
```