---
{
  "image": "vorteil/google-translator",
  "desc": "Reads a string and converts it to the targeted language."
}
---

# Google Translator

Reads a string and converts it to the targeted language.

## Direktiv

An example workflow of reading the string and converting to a target language.

```yaml
id: read-message
functions:
- id: translate
  image: vorteil/google-translator:v1
states:
- id: read
  type: action
  action:
    function: translate
    input: .
```

## Input

The input required to run the above workflow properly is the following:

```json
{
    "message": "Hello World!",
    "target-language": "ja",
    "serviceAccountKey": .secrets.SERVICE_ACCOUNT_KEY
}
```

## Output

Upon successfully translating the string using Google's API. THe following output will appear

```json
{
    "message": "こんにちは世界"
}
```

## Error

In the case that an error is encountered, it will present in the following format:

```json
{
    "errorCode": "com.google-translator.error",
    "errorMsg": "Something Went Wrong"
}
```
