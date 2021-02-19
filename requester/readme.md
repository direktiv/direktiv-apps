# Request

Perform a basic HTTP/S request.

## Input

```json
{
    "method": "GET",
    "url"   : "https://vorteil.io",
    "body"  : {
        "hello": "world"
    },
    "headers": {
        "Content-Type": "application/json" 
    }
}
```

**NOTE:** The `body` and `header` fields are optional.

## Output

If the request is successful, the response will be wrapped inside of a JSON object, within the `return` field.
In the event that the response body is not in JSON format, the data will be base64 encoded and stored as the value of the `return` field.


```json
{
    "return": {
        ...
    }
}
```

## Error

In the case that an error is encountered, it will present in the following format:

```json
{
    "errorCode": "com.request.error",
    "errorMsg": "Something went wrong"
}
```