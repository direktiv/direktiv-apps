# govc

## Direktiv

An example below of changing my sso-admin password on vcenter using Direktiv.

```yaml
functions:
- id: govc
  type: reusable
  image: direktiv/govc:v1
description: change password on vcenter using govc
states:
- id: change
  type: action
  action:
    function: govc
    secrets: ["USERNAME", "NEWPASSWORD", "PASSWORD"]
    input:
      username: jq(.secrets.USERNAME)
      password: jq(.secrets.PASSWORD)
      commands: 
      - jq("sso.user.update -p="+ .secrets.NEWPASSWORD + "-k=true trent")
      url: "https://vcsa.vorteil.local"
```

- **-p=** is where we provided the password update
- **-k=true** is a flag to allow insecure connections to vcenter

## Input

```yaml
input:
  username: test # Used to login
  password: x # Used to login
  commands: # commands that fork the govc binary
  - test
  - test2
  url: https://vcsa.vorteil.local #VCenter address to connect to
```

## Output
The output for this action is an array of strings which contains the output of each command executed on the application.

```json
{
    "return": [
        "execution output of the first command in the commands array",
        "execution output of the second command"
    ]
}
```

## Error

In the case that an error is encountered, it will present in the following format:

```json
{
  "errorCode": "com.%s-govc.error",
  "errorMsg": "Something went wrong"
}
```

**NOTE:** `%s` indicating where in the container logic the application failed.