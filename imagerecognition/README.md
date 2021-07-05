---
{
  "image": "vorteil/imagerecognition",
  "desc": "Checks an image to see if it is safe for work and responds with content is likely to be racy, adult or violence themed."
}
---
# ImageRecognition

Checks an image to see if it is safe for work and responds with content is likely to be racy, adult or violence themed.

## Direktiv

An example workflow of using the container on Direktiv.

```yaml
id: check-image
functions:
- id: check
  image: vorteil/imagerecognition:v2
description: "Evaluates an image using Google Vision API"
states:
- id: checkImage
  type: action
  action:
    function: check
    input: jq(.)
```

## Input

The input needed to run the above workflow properly is the following:

```yaml
input:
  url: "IMAGE_TO_URL"
  serviceAccountKey: jq(.secrets.SERVICE_ACCOUNT_KEY)
```

## Output

If the action is successful, the image recognition container will respond with the following JSON.

```json
{
    "safeForWork": false,
    "violenceLikelihood": "LIKELY",
    "racyLikelihood": "VERY_LIKELY",
    "adultLikelihood": "UNLIKELY"
}
```

## Error

In the case that an error is encountered, it will present in the following format:

```json
{
    "errorCode": "com.imagerecognition.error",
    "errorMsg": "Something went wrong"
}
```