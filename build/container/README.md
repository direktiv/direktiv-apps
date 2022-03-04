# Container

Building and pushing containers in Direktiv

addiotonal args need to be like --build-arg=jens=gerke

## Direktiv
```yaml

```

## Input

git://[repository url][#reference][#commit-id]


|Command|Description|
|---|---|
|list-tags|Lists all tags of the project|



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


## Error

In the case that an error is encountered, it will present in the following format and stops the execution of commands:

```json
{
    "errorCode": "com.git.error",
    "errorMsg": "Something went wrong"
}
```
