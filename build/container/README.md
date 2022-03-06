# Container

Building and pushing containers in Direktiv to container registries

## Direktiv
```yaml
functions:
- id: container
  image: direktiv/container
  type: reusable
states:
- id: builder 
  type: action
  action:
    function: container
    secrets: ["docker"]
    input: 
      buildkit: true
      args:
      - "--build-arg=myarg=value"
      registries:
      - user: gerke74
        password: jq(.secrets.docker)
        registry: https://index.docker.io/v1/
      tag: gerke74/direktiv-build
      context: "https://github.com/jensg-st/hello-world.git#refs/heads/main"
```

## Input

### args

A list of additional args for the build process. Args have to be provided in one string without spaces. For example `--build-arg hello=value` is not valid.

```yaml
args:
- "--build-arg=hello=value"
- "--pull"
```

### dockerfile-arg

The build process can be executed with a completely different Dockerfile provided by this arguments. This is only possible if a tar is used. It is NOT possible for git or http repositories.

```yaml
dockerfile-arg:
  type: plain
  data: |-
    FROM alpine
     CMD ["sleep"]
```

### tar

Instead of referencing a file, external URL or git a tar file can be posted to the build process. This tar file has to be base64 encoded. The project needs to be in a a folder and the name provided has to be folder name. For example the tar folder structure needs to be:

- myproject
  - Dockerfile 
  - file.txt
  - anothefile.sh

The YAML for this would look like the following:

```yaml
tar:
  type: base64
  name: myproject
  data: 123BASEDATA==
```

### tag

Tag used for build and push. It is equal to the `-t` argument in docker. 

```yaml
tag: mycompany/myapp
```

### registries

Registries are being used for pushing and pulling the images. The are getting created per instance so multiple build jobs can run in parrallel. 

```yaml
registries:
- user: dockerhub-user
  password: mytoken
  registry: https://index.docker.io/v1/
registries:
- user: _json_key
  password: myjsontoken
  registry: us-central1-docker.pkg.dev

```

### context

Context can be an external URL or git repository to build. This is used if there is no tar file provided by the request. 

```yaml
context: https://github.com/jensg-st/hello-world.git#refs/heads/main
```

### buildkit

Enables building with Buildkit. 

```yaml
buildkit: true
```

## Output

This function returns nil/null if successful.

## Error

The following errors can be returned:

- com.docker.queue.error: work queue was full
- com.docker.timeout.error: docker did not start in time
- com.docker.build.error: error during build
- com.docker.nobuild.error: build did not return in time

