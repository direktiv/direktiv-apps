+++
{
  "image": "vorteil/influxdb",
  "desc": "A container that provides the opportunities to write or query data on an InfluxDB instance."
}
+++
# Influxdb

A container that provides the opportunities to write or query data on an InfluxDB instance.

## Direktiv

An example workflow to call the container to perform an operation on InfluxDB.

```yaml
id: get-influx-db
functions:
- id: influxdb
  image: vorteil/influxdb:v1
description: "Writes or queries data inside a Influxdb instance"
states:
- id: listInfluxDBInstance
  type: action
  action:
    function: influxdb
    input: .
```

## Query Data on InfluxDB

### Input

The input needed to run a query on an InfluxDB instance is the following: 

```json
{
    "url": "http://localhost:8089",
    "token": .secrets.AUTH_TOKEN,
    "type": "query",
    "organisation": "vorteil.io",
    "bucket": "vorteil",
    "query": "from(bucket:\"my-bucket\")|> range(start: -1h) |> filter(fn: (r) => r._measurement == \"stat\""
}
```

### Output 

The output is a custom JSON marshal that returns the field, value and time of the record under a records key array. Like the following json:

```json
{
    "records": [
        {
            "avg":24.5
        }, 
        {
            "max":45
        }
    ]
}
```

## Write Data to InfluxDB

### Input

The input needed to write data to an InfluxDB instance is the following:

```json
{
    "url": "http://localhost:8089",
    "token": .secrets.AUTH_TOKEN,
    "type": "query",
    "organisation": "vorteil.io",
    "bucket": "vorteil",
    "measurement": "stat",
    "tags": {
        "unit": "temperature",
    },
    "fields": {
        "avg": 24.0,
        "max": 45.0
    }
}
```

### Output

When writing data to an InfluxDB instance nothing gets returned so the golang application will write json like the following:

```json
{
    "success": true
}
```

## Error

In the case that an error is encountered, it will present in the following format:

```json
{
    "errorCode": "com.influxdb.error",
    "errorMsg": "Something went wrong"
}
```