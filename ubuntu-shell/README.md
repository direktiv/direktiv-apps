---
{
  "image": "vorteil/ubuntu-shell",
  "desc": "Executing shell scripts"
}
---
# Request

Executing shell scripts on Ubuntu

## Direktiv

An example workflow of using the container in a workflow on Direktiv.


```yaml
id: shell
functions:
- id: myshell
  image: vorteil/ubuntu-shell:v1
states:
- id: hello
  type: action
  action:
    function: shell
    input: |
        {
          "script": (.var.script),
          "args": ["arg1", "arg2"]
        }'
```

## Input

The action takes two parameters. The first one is 'script' which is a base64 encoded shell script. Usually it comes from a workflow-scoped variable. The seoncd is an array of arguments for the script.

## Output

The action returns either JSON with an output of the shell script or a valid JSON if the script return valid json.

#### Output return

```sh
#!/bin/bash

ls -la
pwd
```

```json
{
  "output": "output of the shell script\n"
}
```

#### JSON return

```sh
#!/bin/bash

echo direktiv > /dev/null
echo '{ "direktiv":"'"$1"'" }'
```

```json
{
  "direktiv": "arg"
}
```

## Error

In the case that an error is encountered, it will present in the following format:

```json
{
    "errorCode": "com.ubuntu.error",
    "errorMsg": "Something went wrong"
}
```
