# Requester

Takes two arguments on being the input file the other being the output file.

It reads the input file to workout the request it needs to send via HTTP. The below struct is what is required for the input file.

```json
{
    "method": "GET",
    "host": "https://vorteil.io",
    "body" : {
        "hello":"world"
    },
    "headers" : {
        "TEST" "test2"
    }
}
```
The body is the body we send with the request that will be a json object
The headers is a key value pair that we will add to the request

The output file will look something like

```json
{
    "Error": "",
    "Response": "",
    "Status": "",
    "StatusCode": ""
}
```

Error is only used if the request encountered an error during the http request.
Response, Status and Status code are all returned from the API we are requesting