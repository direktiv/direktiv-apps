# Request

Takes a simple json file like the one below to send a HTTP request using Direktiv.

```json
{
    "method": "GET",
    "host": "https://vorteil.io",
    "body" : {
        "hello":"world"
    },
    "headers" : {
        "Authorization": "Bearer Test"
    }
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
