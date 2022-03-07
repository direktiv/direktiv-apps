---
{
  "image": "direktiv/vmware-powercli",
  "desc": "Ability to send multiple vmware power-cli commands."
}
---

# VMware Power CLI

The ability to send multiple `vmware-power-cli` commands and execute them. If 'host' is empty it is a regular powershell.

## Direktiv

*Script Example*
```yaml
functions:
- id: power
  image: direktiv/vmware-powercli:v3
  type: reusable
states:
- id: power 
  type: action
  action:
    function: power
    input: 
      host: 192.168.0.39
      user: root
      password: mypassword
      scripts: 
        - script: 
            name: script1.ps1
            data: "Get-VM -Name ubuntu1 | ConvertTo-Json  -Depth 1 -AsArray"
            type: plain
        - script: 
            name: script2.ps1
            data: V3JpdGUtT3V0cHV0ICRhcmdz
            type: base64
          args: ["ARGUMENT", jq(.value)]
        - script: 
            name: script3.ps1
            data: |
              '{ "true": "false" }' | Out-File -FilePath out.json
            type: plain
          output: out.json
```

The script array can take a script in the following format: 

```yaml
name: script1.ps1
data: XYZ
type: plain
```

The type can be plain, base64 a file or a variable in SCOPE/NAME format, e.g. *workflow/myscript*. The *output* value defines a file which is read after the script has finished. It needs to be in JSON format and will be returned in *output* again. 

Script results will be converted to JSON as well if the result is JSON. If the result is text it will look like the following snippet:

```json
"return": {
		"script.ps1": {
			"stdout": "connecting to server\n\nName  Port  User\n----  ----  ----\nserver   443   root\nWARNING: The 'Version' property of VirtualMachine type is deprecated. Use the 'HardwareVersion' property instead.\n\nName   : ubuntu",
      "stderr": "",
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
- id: power
  image: direktiv/vmware-powercli:v1
  type: reusable
states:
- id: power
  type: action
  action:
    secrets: ["ESXI_PWD"]
    function: power
    input:
      host: my.esxi.server
      user: root
      password: jq(.secrets.ESXI_PWD)
      scripts:
        - script: 
            name: script1.ps1
            data: "Get-VM -Name ubuntu1 | ConvertTo-Json  -Depth 1 -AsArray"
            type: plain
```


## Additional Configuration Values

| Attribute | Function |
| -- | -- |
|Continue|if set to true a script failure does not throw an error|


## Error on Execution

If the container was to error out during the action an error will be thrown matching:

```json
{
    "errorCode": "com.powershell.%s.error",
    "errorMsg": "Something went wrong"
}
```

The `%s` will be replaced with what is currently happening in the container that created the error e.g. `unmarshalling-input`.
