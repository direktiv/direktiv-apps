```id: pod
description: simple workflow with isolated pod"
functions:
- id: get
  type: isolated
  image: localhost:5000/podtest:latest
states:
- id: getter
  type: action
  action:
    function: get
    input:
      data: mydata
```
