+++
{
  "image": "vorteil/simplepostgresclient",
  "desc": "Perform simple interactions on a PostgreSQL database. "
}
+++
# Simple Postgres Client

Perform simple interactions on a PostgreSQL database. 

## Direktiv 

Here's an example of how this could be used in Direktiv to store event data in a database.

```yaml
id: eventOnboarder
start:
  type: event
  event: 
    type: com.github.pull.create
    filters:
      source: "https://github.com/cloudevents/spec/pull"
functions:
- id: postgres
  image: vorteil/simple-postgres-client:v1
states:
- id: store
  type: action
  action:
    function: postgres
    secrets: ["DB_PASSWORD"]
    input: '{
      conn: ("postgres://vorteil:" + .secrets.DB_PASSWORD + "@203.0.113.5:5432/github_events_db"),
      table: "github_events",
      transaction: [{
        type: "insert",
	data: ."com.github.pull.create"
      }]
    }'
```

## Input 

The Simple Postgres Client can support a handful of different operations, and can chain them together into simple transactions. Here's an example input, which will be explained in further detail below.

```json
{
    "conn": "postgres://vorteil:password@localhost:5432/postgres?sslmode=disable",
    "table": "test",
    "transaction": [
        {
            "type": "update",
            "where": {"name": null},
            "set": {"name": "jon", "age": 24}
        }
    ]
}
```

### Top-Level Parameters

The following fields can appear on the first level of the input JSON:

The `conn` field must be a string that can be used to connect to a Postgres database. 

The `table` field must identify a table within the database. All functions supported by the Simple Postgres Client work on the same table. 

The `debug` field may be either `true` or `false`, or omitted entirely. If set to true, the queries logged will not obscure values.

The real input goes into the `transaction` array. Each element in the transaction array must be an object that conforms to one of the four operation inputs outlined below (`delete`, `insert`, `select`, or `update`). All "steps" in this transaction array will be performed within a single database transaction, meaning they'll all be committed or rolled back together.

### Delete Input 

An `delete` step is identified by the required `"type": "delete"` parameter. It also requires a `where` parameter, which must be an object that contains only key-value pairs of primitive types. 

```json 
{
  "type": "delete",
  "where": {
    "a": 1
  }
}
```

### Insert Input 

An `insert` step is identified by the required `"type": "insert"` parameter. It also requires `data`, which can be either a single object or an array of objects. All data objects must contain only key-value pairs of primitive types. 

```json 
{
  "type": "insert",
  "data": [{
    "b": "B",
    "c": 2,
    "d": null
  }]
}
```

Each data object will be inserted into the database with the keys used to identify columns.

### Select Input 

A `select` step is identified by the required `"type": "select"` parameter. It also requires `where` and `fields` parameters. The `where` parameter must be an object that contains only key-value pairs of primitive types. The `fields` parameter can be either `"fields": "*"` or an array of strings that select specific columns to appear in the output of the step.

```json 
{
  "type": "select",
  "where": {
    "a": 1
  },
  "fields": "*"
}
```

### Update Input

An `update` step is identified by the required `"type": "update"` parameter. It also requires `where` and `set` parameters, which are both objects that must contain only key-value pairs of primitive types. 

```json 
{
  "type": "update",
  "where": {
    "a": 1
  },
  "set": {
    "b": "B",
    "c": 2,
    "d": null
  }
}
```

## Output

Output data will be stored as a JSON array with each element representing the results of one step in the transaction. 

### Example

```json
[{ "rowsAffected": 3 }, { "rows": [{"id": 1}, {"id": 2}] }]
```

### Delete Output

The output from a `delete` step will only contain the number of `rowsAffected` by the delete operation.

### Insert Output

The output from an `insert` step will be an empty object `{}`.

### Select Output

The output from a `select` step will be an object containing an array of `rows`, each of which is itself an object with keys and values extracted from the database.

### Update Output

The output from an `update` step will only contain the number of `rowsAffected` by the update operation.

## Error 

If an error occurs, the transaction will be rolled-back and the error will be logged. It will also be found in `/direktiv-data/error.json` and have the following structure:

```json
{
  "errorCode": "error.conn",
  "errorMsg": "Failed to connect to database."
}
```