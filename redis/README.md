---
{
  "image": "vorteil/redis",
  "desc": "Sets and gets on a redis memory store."
}
---
## Redis

The following container can do the below

- Gets a specific key's value
- Sets an array of KEY=VALUE strings

## Direktiv

An example workflow of using the container on Direktiv.

```yaml
id: redis
function:
- id: redis-op
  image: vorteil/redis:v2
description: "Gets a key value or set an array of KEY=VALUE strings"
states:
- id: redis-obj
  type: action
  action:
    function: msg
    input: .
```

## Input

### Get

```json
{
    "method": "get",
    "addr": "localhost:6379",
    "password": "",
    "db": 0,
    "values": "TEST"
}
```

### Set

```json
{
    "method": "set",
    "addr": "localhost:6379",
    "password": "",
    "db": 0,
    "values": ["TEST=TEST2", "TEST2=TEST3"]
}
```

## Output

### Set

If the action is successful, the redis should have stored each key and value. No output will be returned to the user.

### Get

If the action is successful, the output will be a json object like the following.

```json
{
    "value": "TEST2"
}
```

## Error

In the case that an error is encountered, it will present in the following format:

```json
{
    "errorCode": "com.redis.error",
    "errorMsg": "Something went wrong"
}
```