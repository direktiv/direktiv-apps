# Examples to Write a simple 'Greeting' container 

The examples will provide you with how to write it in a specific language and the Dockerfile required to build it into a container.

The following action takes a name and responds with a greeting message.

## Input

```json
{
    "name": "Trent"
}
```

## Output

```json
{
    "greeting": "Welcome to Direktiv, Trent!"
}
```

## Error

```json
{
    "errorCode": "com.greeting-%s.error",
    "errorMsg": "Something went wrong"
}
```

**NOTE** '%s' in the errorCode is based on what is currently happening in the container.

## Languages

| Language | README |
| -------- | ------ |
| golang | [readme](https://github.com/vorteil/direktiv-apps/tree/master/examples/golang) |