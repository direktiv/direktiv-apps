# SSH & SCP

This container can run scripts and commands via SSH and can copy files via SCP. Users can autheticate via password or certificate.

## Direktiv

*Script Example*
```yaml
functions:
- id: scp
  image: direktiv/sshscp
  type: reusable
states:
- id: scp 
  type: action
  action:
    function: scp
    secrets: ["password"]
    input: 
      actions:
      - type: ssh
        host: direktiv@192.168.0.35
        auth: jq(.secrets.password)
        output: "/tmp/out"  # reading /tmp/out at the end
        files:
        - name: myscript.sh
          type: plain
          data: |- # this script pipes to /tmp/out and will be read into *output* in the response
            ls | jq -R -s -c 'split("
            ")[:-1]' > /tmp/out 
```

The files array can take files in the following format: 

```yaml
name: script1.sh
data: XYZ
type: plain
```

The type can be plain, base64 a file or a variable in SCOPE/NAME format, e.g. *workflow/myscript*. The *output* value defines a file which is read after the script has finished. If it is in JSON it will be added as JSON within the return block. Otherwise it will be handled as string.

```json
"return": {
    "scp": {},
    "ssh": {
      "192.168.0.35": [
        {
          "script": "myscript.sh",
          "stdout": "[\"Desktop\",\"Documents\",\"Downloads\",\"Music\",\"Pictures\",\"Public\",\"Templates\",\"Videos\"]\n",
          "success": true
        }
      ]
    }
```

## SCP

The difference for SCP is the pattern of the hostname. The hostname needs a file path at the end like it is used in a command line SCP command. 

```yaml
functions:
- id: scp
  image: direktiv/sshscp
  type: reusable
states:
- id: scp 
  type: action
  action:
    function: scp
    secrets: ["password", "aws"]
    input: 
      actions:
      - continue: true
        type: scp
        host: ec2-user@ec2-1-2-3-4.mars-orbit-1.compute.amazonaws.com:/tmp
        auth: direktiv-cert.pem
        files:
        - name: file.txt
          type: plain
          data: mydata
        - name: direktiv-cert.pem
          type: plain
          data: jq(.secrets.aws)
      - continue: true
        type: scp
        host: direktiv@192.168.0.35:/tmp
        auth: jq(.secrets.password)
        files:
        - name: file2.txt
          type: plain
          data: |
            this is data
```

## Authentication

The client can authenticate with password or certificate. If there is a file in the file list with an identical name specified in *auth* it will be used as a certificate to authenticat, e.g.

```yaml
auth: direktiv-cert.pem
files:
- name: direktiv-cert.pem
    type: plain
    data: jq(.secrets.aws)
```

If there is no file with that name *auth* is handled as password. 



## Additional Configuration Values

| Attribute | Function |
| -- | -- |
|Continue|if on host fails the action tries the next one if there are more than one|
|actions.Type|*scp* or *ssh*|
|actions.host|hostname to connect to|
|actions.Auth|password or certificate to use|
|actions.Port|SSH/SCP port, default 22 if not set|
|actions.Output|SSH only, if the script creates a file with that name, the content will be used as return value|
|actions.Continue|If set to true the container continues on errors in one command or host|

## Error on Execution

If the container was to error out during the action an error will be thrown matching:

```json
{
    "errorCode": "com.sshscp.%s.error",
    "errorMsg": "Something went wrong"
}
```

