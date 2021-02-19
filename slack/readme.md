# Use request container to post to slack

- Create an app assign it to your team for slack.
- Use the request program to send the following input to hit your application to push a message

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


- The response will be a 200 if it is successful.