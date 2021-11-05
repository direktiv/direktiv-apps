---
{
  "image": "direktiv/vmware-powercli",
  "desc": "Ability to send multiple vmware power-cli commands."
}
---

# VMware Power CLI

The ability to send multiple `vmware-power-cli` commands and execute them.

## Direktiv

This container can run scripts and command with PowerShell. The VMWare plugin has been already imported.

Using a script needs a reference to a namespace or global variable. The shell runs unauthenticated and does not create a connection to a vCenter cluster or ESXi host.

*Script Example*
```yaml
description: A simple 'action' state that sends a get request
functions:
- id: get
  image: direktiv/vmware-powercli:v1
  type: reusable
  files:
    - key: myscript"
      scope: namespace
      type: plain
      as: "ps-script.ps1"
states:
- id: getter
  type: action
  action:
    function: get
    input:
      print: true
      scripts:
        - name: "ps-script.ps1"
          args:
            - "arg1"
```

Script results will be converted to JSON as well if the result is JSON. If the result is text it will look like the following snippet:

```json
"return": {
		"script-0": {
			"output": "connecting to server\n\nName  Port  User\n----  ----  ----\nserver   443   root\nWARNING: The 'Version' property of VirtualMachine type is deprecated. Use the 'HardwareVersion' property instead.\n\nName   : ubuntu",
			"result": "success"
		}
	}
```

The command alternative runs commands in order they are listed. If they return JSON the response will be added as an addressable JSON object in the response to Direktiv. The command `Get-VM -Name ubuntu | ConvertTo-Json  -Depth 1 -AsArray` will yield the following output:

```json
{
	"return": {
		"0": {
			"output": [
				{
					"BootDelayMillisecond": 0,
					"CoresPerSocket": 1,
					"CreateDate": "2021-11-02T21:58:09.800272Z",
					"CustomFields": "",
					"DatastoreIdList": "Datastore-6181b457-fb4061f6-f9b8-000c291013a9",
					"DrsAutomationLevel": null,
				}
			],
			"result": "success"
		}
	}
}
```
It is important to add *-Confirm:$false* if the commands being called needs confirmation.

*Command Example*
```yaml
description: A simple 'action' state that sends a get request
functions:
- id: get
  image: direktiv/vmware-powercli:v1
  type: reusable
states:
- id: getter
  type: action
  action:
    secrets: ["ESXI_PWD"]
    function: get
    input:
      host: my.esxi.server
      user: root
      password: jq(.secrets.ESXI_PWD)
      on-error: stop
      full-command: true
      run:
        - Get-VM -Name ubuntu
  transition: check
```


## Additional Configuration Values

| Attribute | Function |
| -- | -- |
|on-error|if set to 'stop' the function fails on the fist error|
|full-command|Prints the full command. If set to false only the first 10 characters will be printed|
|print| Prints stdout, stderr of the command / script|

## Error on Execution

If the container was to error out during the action an error will be thrown matching:

```json
{
    "errorCode": "com.vmware-power-cli.%s.error",
    "errorMsg": "Something went wrong"
}
```

The `%s` will be replaced with what is currently happening in the container that created the error e.g. `unmarshalling-input`.
