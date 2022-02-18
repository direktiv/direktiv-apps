# Python

Run python scripts in Direktiv.

## Direktiv

This container runs python scripts provided by the workflow.


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

<!-- Script results will be converted to JSON as well if the result is JSON. If the result is text it will look like the following snippet:

```json
"return": {

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
It is important to add *-Confirm:$false* if the commands being called needs confirmation. -->

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
|pre-run|Script running before the python script, see *Providing the script*, e.g. to pip install|
|post-run|Script running after the python script, see *Providing the script*, to move a file into the  variable folder|
|script| The python script to use, see *Providing the script* |
|args| Arguments for the script|
|envs| Environment variables for the script, applies to pre and post run scripts too|
|output| This JSON file will be used as JSON response from the script, base path is the working directory|

## Error on Execution

| Error | Reason |
| -- | -- |
|com.python.pre-run.error| The script provided in pre-run returned an error |
|com.python.post-run.error| The script provided in post-run returned an error |
|com.python.execute.error| Executing the python script resulted in an error |

## Providing the script

A script can be provided in different formats:

### Plain text

Ideal for smaller scripts.

```python
script: 
    name: jens.py
    data: |
        print('we are doing a request')
    type: plain
```

### Variable

Reads the variable from the scope and uses this as script. The format is SCOPE/VARIABLE-NAME.

```python
script: 
    name: jens.py
    data: workflow/myscript
    type: var
```

### File

Executes an existing file on the file system.

```python
script: 
    data: /myfile.py
    type: var
```

### Base64

Decodes the provided the base64 string and uses this as the script. Ideal if the script is part of the workflow state. 

```python
script: 
    data: cHJpbnQoJ3dlIGFyZSBkb2luZyBhIHJlcXVlc3QnKQ==
    type: var
```