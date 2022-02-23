# Git

Perform git commands in Direktiv. This container supports asubset of the original git commands. This container returns json for all commands unlike the original git command. If advanced git commands are required then the ubuntu-shell conatiner can be used which has git installed. 

## Direktiv
```yaml

functions:
- id: git
  image: direktiv/git:v5
  type: reusable
states:
- id: git 
  type: action
  action:
    function: git
    input: 
      clone:
        repo: https://github.com/direktiv/direktiv.git
      cmds:
      - cmd: list-tags
```

## Input

The *clone* part of the input requires the URL of the repository and additional optional attributes:

- depth: Create a shallow clone with a history truncated to the specified number of commits
- ref: Clones only a single branch
- username: Username of the user for private repositories
- pwd: Password for private repositories or Github peronal access token.

```yaml
input: 
  clone:
    repo: https://github.com/direktiv/direktiv.git
    depth: 1
    ref: main
```


Commands are provided via the *cmd* field in input. It consists of the main command, listed in the table below, and arguments. They can have an additional attribute *continue*. If this value is set to true the container proceeds to the next command if an errror occurs and does not return an error.

```yaml
input: 
  cmds:
  - cmd: checkout
    args:
    - tag
    - v0.1
```

|Command|Description|
|---|---|
|list-tags|Lists all tags of the project|
|delete-tag **tag-name**|Deletes a tag|
|create-tag **tag-name**|Creates a tag|
|list-branches|Lists all branches of the project|
|delete-branch **branch-name**|Deletes a branch|
|create-branch **branch-name**|Creates a branch|
|checkout **type** **ref**|Checks-out the type (tag, branch, ref), e.g. *checkout branch main*|
|get-file **filename**|Returns the file content as base64|
|get-files|Lists all files in the project|
|list-commits|Lists all commits in the project|
|add **filename**|Adds a file to the commit|
|commit **name** **email** **message**|Commits the changes|
|status|Returns the status of the tree|
|diff **ref** **ref**|Returns the diff between two trees. The refernce can be a commit hash or tag, e.g. *diff v0.1.0 v0.1.1*|
|logs **type** **ref**|Returns logs. The args *type* and *ref* are optional. If provided it returns the logs from that reference| 

A special command is *script*. This can run shell scripts if files need to be changed or added. 

```yaml
input: 
  clone:
    repo: https://github.com/direktiv/direktiv.git
  cmds:
  - cmd: script
    script:
      data: |
        #!/bin/bash
        echo "new file" > nf.txt
        ls -la
      type: plain
  - cmd: status
```

## Output

The output is JSON if the command is successful. 

```json
{
  "return": {
    "0": {
      "v0.1": "ba33ba8e83fb948afb943600687965d7bdef6d11",
      "v0.1.1": "759e88542a68e309d03a18e5d87a06376844a478",
      "v0.1.2": "c22df2b2aceda8daa5ae012c425dd0ced44c2ac3",
    }
  }
}
```

## Full Example

The following exmaple shows a workflow whic requires an input like the following: 

```json
{
  "name": "myname"
}
```

The container clones the repository and creates a file *$name*.txt, adds it and commits it.

```yaml
functions:
- id: git
  image: direktiv/git:v5
  type: reusable
states:
- id: git 
  type: action
  action:
    function: git
    input: 
      clone:
        repo: https://github.com/mysecret/repo.git
        user: github-user
        pwd: ghp_123
      cmds:
      - cmd: script
        script:
          data: |
            #!/bin/bash
            echo jq(.name) > jq(.name).txt
            ls -la
          type: plain
      - cmd: status
      - cmd: add
        args:
        - jq(.name).txt
      - cmd: commit
        args:
        - jq(.name)
        - my@email.com
        - committed by jq(.name)
      - cmd: status
      - continue: true
        cmd: push
        args:
        - main
```


## Error

In the case that an error is encountered, it will present in the following format and stops the execution of commands:

```json
{
    "errorCode": "com.git.error",
    "errorMsg": "Something went wrong"
}
```
