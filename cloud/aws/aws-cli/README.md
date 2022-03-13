# Python

Run aws cli commands in Direktiv

## Direktiv

*Script Example*
```yaml
description: creates a test vm in aws
functions:
- id: aws
  image: main-server:5443/aws-cli
  type: reusable
- id: delete
  workflow: delete-vm
  type: subflow
states:
- id: validate
  type: validate
  schema:
    title: Files
    type: object
    required:
      - name
      - size
    properties:
      name:
        type: string
        default: deleteme
      size:
        type: string
        default: t2.xlarge
  transition: create-vm
- id: create-vm 
  type: action
  action:
    function: aws
    input: 
      access-key: MYACCESSKEY
      secret-key: sUpERSecrETKEY
      region: eu-central-1
      commands: 
        - ec2 create-security-group --group-name jq(.name) --description jq(.name)
          --tag-specifications ResourceType=security-group,Tags=[{Key=direktiv,Value=build},{Key=name,Value=jq(.name)}]
        - ec2 authorize-security-group-ingress --group-name jq(.name) --cidr 0.0.0.0/0 --protocol tcp --port 443
      script:
        data: |
        aws ec2 delete-key-pair --key-name jq(.name)
        type: plain
```


## Additional Configuration Values

| Attribute | Function |
| -- | -- |
|access-key|AWS access key|
|secret-key|AWS secret key|
|region| Region, default `us-east-1` |
|commands| List of commands to excecute without `aws` |
|continue| if set to true command execution continues if one command fails|
|script| authenticated script to run, runs as bash but #! does not have to be part of the script data |
|print| print result to stdout, bool|

## Error on Execution

Only one error `com.aws.execute` will be thrown if something if a command fails. Commands will be executed and CLI json or errors returned.

