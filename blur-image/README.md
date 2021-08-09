---
{
  "image": "vorteil/blur",
  "desc": "Takes a URL of an image and blurs it"
}
---

# Blur

Takes a URL of an image and blurs it.

## Direktiv

An example workflow of using the container on Direktiv

```yaml
id: blur-image
functions:
- id: blur
  type: reusable
  image: vorteil/blur:v1
description: Takes a URL of an image and blurs it
states:
- id: blur
  type: action
  action: 
    function: blur
    input:  
      image: "http://urlimage"
```

## Input

The input required to run the container needs a json struct that requires an image URL.

```json
{
    "image": "http://urlimage" 
}
```

## Output

Returns the entire image blurred

## Error

In the case than error is encountered it will present in the following format:

```json
{
    "errorCode": "com.blurimage.error",
    "errorMsg": "Something went wrong"
}
```