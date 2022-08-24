# Python

Run python scripts in Direktiv.

## Direktiv

This container runs python scripts provided by the workflow.

*Script Example*
```yaml
functions:
- id: run
  image: direktiv/python:v1
  type: knative-workflow
states:
- id: getter 
  type: action
  action:
    function: run
    input: 
      script: 
         name: simple
         data: |
          print('Hello')
          print('World')
         type: plain
```

The file can be provided in different formats anc can be downloaded from external sources.

*Command Example*
```yaml
functions:
- id: getfile
  type: reusable
  image: direktiv/request:v4
- id: run
  image: direktiv/python:v1
  type: reusable
states:
- id: getfile
  type: action
  action:
    function: getfile
    input:
      method: "GET"
      url: "https://gist.githubusercontent.com/jensg-st/4a14fa88cc96855358d25a90572f81b2/raw/16be9710c0f941b3482bdc9a2245ff2ea5a3ac9f/simplepython"
    retries:
      max_attempts: 3
      codes: 
        - "*"
      delay: PT5S
  transition: log
- id: log
  type: noop
  transform: 
    script: jq(.return.body)
  transition: run
- id: run
  type: action
  action:
    function: run
    input: 
      script: 
         name: runit
         data: jq(.script)
         type: base64
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