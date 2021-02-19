# Request

## Input

Request takes a simple json input like the one below.

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

Providing body and headers are optional.

## Output

The output of the following container will be the response of the http request.

### Error

If an error was to happen during the container running the output is written back to direktiv as the following

```json
{
    "errorCode": "com.request.error",
    "errorMsg": "Something went wrong"
}
```