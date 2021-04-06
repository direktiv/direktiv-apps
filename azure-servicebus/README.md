+++
{
  "image": "vorteil/azure-servicebus",
  "desc": "Writes a message to Azure's Servicebus"
}
+++

# Azure Servicebus

Writes a message to Azure's Servicebus

## Direktiv

An example workflow of writing 'Hello World!' to the servicebus.

```yaml
id: write-helloworld
functions:
- id: write
  image: vorteil/azure-servicebus
description: "Writes 'Hello World!' to an Azure Servicebus."
states:
- id: write-message
  type: action
  action:
    function: write
    input: .
```

## Input

The input required to run the above workflow properly is the following:

```json
{
    "policy": .secrets.AZURE_POLICY,
    "key": .secrets.AZURE_POLICY_KEY,
    "namespace": "direktiv",
    "message": "Hello World!",
    "queue": "listener"
}
```

## Output

If the action is successful, no output will be returned.

## Error

In the case that an error is encountered, it will present in the following format:

```json
{
    "errorCode": "com.azure-servicebus.error",
    "errorMsg": "Something went wrong"
}
```