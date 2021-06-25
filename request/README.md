---
{
  "image": "vorteil/request",
  "desc": "Perform a basic HTTP/S request."
}
---
# Request

Perform a basic HTTP/S request.

## Direktiv

An example workflow of using the container in a workflow on Direktiv.


```yaml
id: request
functions:
- id: myrequest
  image: vorteil/request:v6
description: "send a get request" 
states:
- id: hello
  type: action
  action: 
    function: myrequest
    input: | 
        {
            "method": "GET",
            "host"  : "https://jsonplaceholder.typicode.com/posts"
        }
```

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
    },
    "params": {
        "Hello": "world"
    },
    "username": "test",
    "password": "test2",
    "insecureSkipVerify": false
}
```

**NOTE:** The `method` and `url` fields are mandatory.

## Output

Output for a request that returns JSON.
```json
{
    "return": {
        "body": {
            "hello": "world"
        },
        "headers":{
            "Age": 0
        },
        "status-code": 200,
        "status": "200 OK"
    }
}
```

Output for a request that isn't JSON.
```json
{
    "return": {
        "data": "PGh0bWw+PGJvZHk+SGVsbG8gV29ybGQhPC9ib2R5PjwvaHRtbA==",
        "headers":{
            "Age": 0
        },
        "status-code": 200,
        "status": "200 OK"
    }
}
```

## Error

In the case that an error is encountered, it will present in the following format:

```json
{
    "errorCode": "com.%s.error",
    "errorMsg": "Something went wrong"
}
```

**NOTE:** `%s` indicating where in the container logic the request failed.