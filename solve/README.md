---
{
  "image": "vorteil/solve",
  "desc": "Solves math expressions"
}
---

## Solve

Solves math expressions using the [Knetiv/govaluate](https://github.com/Knetic/govaluate) Golang package. 

## Direktiv
An example workflow of using the container solve container to solve basic math.

```yaml
id: solver
functions: 
- id: solve-math-expression
  image: vorteil/solve:v2
states:
- id: solve
  type: action
  action:
    function: solve-math-expression
    input:
	    x: "10*3"
```

# Input
Expects input in the following format where x contains the math expression to solve:

```yaml
input:
	x: "10*3"
```

## Output

Returns output in the following format:

```json
{
  "return": 30
}
```
