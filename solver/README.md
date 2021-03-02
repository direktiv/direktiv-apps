## Solve

Solves math expressions using the [Knetiv/govaluate](https://github.com/Knetic/govaluate) Golang package. 

Expects input in the following format:

```json
{
	"expressions": ["2+2", "4-1", "10*3", "20/2"]
}
```

Returns output in the following format:

```json
["4", "3", "30", "10"]
```
