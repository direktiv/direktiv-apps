---
{
    "image": "direktiv/regex",
    "desc": "Takes a regex string and returns the results"
}
---

# Regex

Takes a string and regex and returns the results that matched.

## Direktiv

An example workflow of using the container on Direktiv.

```yaml
id: regex-msg
functions:
- id: regex
  image: direktiv/regex:v1
description: Checks a string against regex and returns the matches
states:
- id: reg
  type: action
  action:
    function: regex
    input:
      msg: jq(.return.msg)
      regex: "\\+[0-9]{1,2}[0-9]{9}"
```

## Input

The input required to run the container needs a json struct like the following:

```json
{
    "msg": "hello test 053443533",
    "regex": "\\+[0-9]{1,2}[0-9]{9}"
}
```

## Output

The output of the container is an array called `results`.

```json
{
    "return": {
        "results": ["f", "x"]
    }
}
```

## Error

In the case that an error is encountered it will present in the following format:

```json
{
    "errorCode": "com.regex.error",
    "errorMsg": "Something went Wrong"
}
```