# Examples to Write a simple 'Greeting' container 

Each example demonstrates how to write a basic app in a specific language, and provides the Dockerfile required to build it into a container. The following Direktiv action takes a name and responds with a greeting:

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

***NOTE:*** The '%s' will be replaced by a more appropriate message depending on the nature of the error encountered.

| Languages |
| -------- |
| [Golang](https://github.com/vorteil/direktiv-apps/tree/master/examples/golang), [NodeJS](https://github.com/vorteil/direktiv-apps/tree/master/examples/nodejs), [.NET](https://github.com/vorteil/direktiv-apps/tree/master/examples/dotnet), [Rust](https://github.com/vorteil/direktiv-apps/tree/master/examples/rust), [Python](https://github.com/vorteil/direktiv-apps/tree/master/examples/python), [Java](https://github.com/vorteil/direktiv-apps/tree/master/examples/java) |