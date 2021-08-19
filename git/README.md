---
{
  "image": "vorteil/git",
  "desc": "Perform git commands"
}
---
# Request

Perform git commands.

## Direktiv

An example workflow of using the container in a workflow on Direktiv.


```yaml
id: git
functions:
- id: mygit
  image: vorteil/git:v1
description: "send git commands"
states:
- id: hello
  type: action
  action:
    function: mygit
    input: 
      cmds:
        - "clone https://github.com/vorteil/direktiv.git"
        - "-C direktiv tag"
```

## Input

The action takes a list of git commands to run. The working directory is not being changed during the process. If subsequent git commands are being executed after a clone the need the "-C" argument to point to the cloned directory. The variable '$out' can be used to clone into workflow variables, e.g.:

```yaml
input:
  cmds: 
    - "clone https://github.com/vorteil/direktiv.git $out/instance/direktiv"
    - "-C $out/instance/direktiv tag"
```


#### Using secrets for token access

```yaml
action:
      secrets: ["gittoken"]
      function: get
      input: 
        cmds:
          - jq("clone https://" + .secrets.gittoken + "@github.com/vorteil/direktiv.git")
```



## Output

If the request is successful, the response lists the commands with 'cmdX' as key. It contains the result of the command as well as the command itself.


```json
{
  "return": {
    "cmd0": {
      "cmd": "clone https://github.com/vorteil/direktiv.git",
      "output": "direktiv"
    },
    "cmd1": {
      "cmd": "-C direktiv tag",
      "output": [
        "v0.1"
      ]
    }
  }
}
```

## Error

In the case that an error is encountered, it will present in the following format and stops the execution of commands:

```json
{
    "errorCode": "com.git.error",
    "errorMsg": "Something went wrong"
}
```
