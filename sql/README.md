---
{
  "image": "direktiv/sql",
  "desc": "Send sql commands to MySQL, MSSQL and Postgresql"
}
---

# SQL

This function can send SQL queries to supported databases (MySQL, Postgres, MSSQL). If the query returns rows they will be in the response of the function. The function can run multiple queries and appends them in the response in order they were provided.

The main configuration is:

```yaml
# database sepcific connection string
connection: host=myhost port=5432 user=myuser password=mypassword dbname=mydb sslmode=require
# allowed values "postgres", "mysql", "mssql"
db: postgres
 # if a query fails the function should return with an error or keep executing the queries
fail: false
```

## Postgres

```yaml
connection: host=localhost port=5432 user=myuser password=mypassword dbname=mydb sslmode=require
db: postgres
fail: false
queries:
- tx: true
  query: CREATE TABLE if not exists testme (name varchar(32));
- tx: true
  query: INSERT INTO testme VALUES ('myname') RETURNING *;
- tx: true
  query: select * from testme;
```

## MySQL example

```yaml
connection: root:password@/testme
db: mysql
fail: true
queries:
- tx: false
  query: select * from testme;
```

## MSSQL example

```yaml
connection: sqlserver://sa:Password12345@localhost:1433
db: mssql
fail: false
queries:
- tx: false
  query: CREATE DATABASE testme
- tx: true
  query: CREATE TABLE testme (name varchar(32));
- tx: true
  query: INSERT INTO testme VALUES ('myname');
- tx: true
  query: select * from testme;
```

## Full Example

Some of the return value can be in Base64 format if the type can not be detected, e.g. UUID. In this case the value can be decoded with JQ. This full example converts the ID which is in Base64 into a string.

```yaml
description: A simple 'action' state that sends a get request
functions:
- id: get
  image: direktiv/sql
  type: reusable
states:
- id: getter
  type: action
  log: INSERT INTO testme VALUES ('Hello World', 'jq(now | strftime("%Y-%m-%d %H:%M:%S"))') RETURNING *
  transition: print-uuid
  action:
    function: get
    input:
      connection: "host=mydatabasehost port=5432 user=myuser password=mypasswd dbname=mydb sslmode=require"
      db: postgres
      queries:
      - tx: true
        query: |
          CREATE TABLE IF NOT EXISTS testme (
              name varchar(32),
              time timestamp,
              id uuid
          );
      - tx: true
        query: |
          INSERT INTO testme VALUES ('Hello World', 'jq(now | strftime("%Y-%m-%d %H:%M:%S"))', '8853e918-68bf-4410-990b-cd1c08f6fa82') RETURNING *
- id: print-uuid
  type: noop
  log: jq(.return.results[1].result[0].id | @base64d)
```

## Error on Execution

If the container was to error out during the action an error will be thrown matching:

```json
{
    "errorCode": "io.direktiv.%s",
    "errorMsg": "Something went wrong"
}
```

The following errors can be thrown:

- io.direktiv.db
- io.direktiv.query
- io.direktiv.unmarshal
- io.direktiv.marshal
