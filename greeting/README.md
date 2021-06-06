---
{
  "image": "vorteil/greeting",
  "desc": "Outputs a greeting"
}
---

# Greeting

Outputs a greeting

## Direktiv

An example workflow of using the container in a workflow on Direktiv.

```yaml
id: greeting
functions: 
- id: greeting-function
  image: vorteil/greeting:v2
states:
- id: Greet
  type: action
  action:
    function: greeting-function
    input: '{ name: .person.name }'
  transform: '{ greeting: .return.greeting }'
```

## input

```json
{
    "name": "Trent"
}
```

## Output

```json
{
    "greeting": "Welcome to Serverless Workflow, Trent!"
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