# Github CLI

This container executes github cli commands

## Direktiv

*Script Example*
```yaml
functions:
- id: gh
  image: direktiv/github
  type: reusable
states:
- id: gh 
  type: action
  action:
    function: scp
    secrets: ["token"]
    input: 
      token: jq(.secrets.token)
      cmds:
      - issue list --json number,title --repo direktiv/direktiv
```

If the command returns JSON, e.g. via the 


## Configuration Values

| Attribute | Function |
| -- | -- |
|token|Github access token|
|continue|If set to true the container continues on errors in one command|
|cmds|Array of gh commands to run without `gh`, e.g. `issue list`|


## Error on Execution

On failure the action returns the error `io.direktiv.github.cmd.error`.
