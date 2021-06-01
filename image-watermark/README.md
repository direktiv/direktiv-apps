---
{
  "image": "vorteil/watermark",
  "desc": "Generates a basic watermark."
}
---
# Watermark



## Direktiv

An example workflow of using the container in a workflow on Direktiv.


```yaml
id: watermark
functions:
- id: mywatermark
  image: vorteil/watermark:v2
states:
- id: hello
  type: action
  action:
    function: mywatermark
    input: '{
      "text": "Hello",
      "img": "bas64image==",
      "color": "#ffffffC0"
    }'
```

The 'img' value has to be a base64 encoded image (png/jpg).

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
    "errorCode": "com.img.watermark.error",
    "errorMsg": "Something went wrong"
}
```
