---
{
  "image": "direktiv/vmware-powercli",
  "desc": "Ability to send multiple vmware power-cli commands."
}
---

# VMware Power CLI

The ability to send multiple `vmware-power-cli` commands and execute them.

## Direktiv

An example workflow of getting the version of `vmware-powercli` we run.

```yaml
id: get-version
functions:
- id: power-cli
  image: direktiv/vmware-powercli:v1
description: The ability to send multiple `vmware-powercli` commands.
states:
- id: execute
  type: action
  action:
    function: power-cli
    input:
      run:
      - Get-PowerCLIVersion
```

**NOTE:** the `run` variable is an array of powershell strings which is wrapped in a `-Command` on the container itself.

## Output of Execution

The output of the execution is a JSON object that takes your commands you ran as a key to build the structure.

An example of the above workflow that we ran.

```json
{
	"return": {
		"Get-PowerCLIVersion": "\nPowerCLI Version\n----------------\n   VMware PowerCLI 12.2.0 build 17538434\n---------------\nComponent Versions\n---------------\n   VMware Common PowerCLI Component 12.3 build 17838947\n   VMware Cis Core PowerCLI Component PowerCLI Component 12.3 build 17839331\n   VMware VimAutomation VICore Commands PowerCLI Component PowerCLI Component 12.3 build 17839688\n\n\n"
	}
}
```

## Error on Execution

If the container was to error out during the action an error will be thrown matching:

```json
{
    "errorCode": "com.vmware-power-cli.%s.error",
    "errorMsg": "Something went wrong"
}
```

**NOTE:** `%s` will be replaced with what is currently happening in the container that errored out e.g. `unmarshalling-input`.