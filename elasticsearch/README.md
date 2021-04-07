---
{
  "image": "vorteil/elasticsearch",
  "desc": "Writes a new JSON document to an index or queries an index in its entirety."
}
---

# Elasticsearch

Writes a new document to an index or queries an index in its entirety.

## Direktiv

An example workflow of using the Elasticsearch container for either options

```yaml
id: test2
description: "" 
functions:
- id: es
  image: trentis/elasticsearch
states:
- id: hello
  type: action
  action:
    function: es
    input: .
```

## Input

### Query

Searches the index and returns all the results.

```json
{
    "method": "query",
    "index": "direktiv",
    "addresses": ["http://192.168.1.30:9200"]
}
```

### Store

Stores the record under the index provided.

```json
{
    "method": "store",
    "index": "direktiv",
    "addresses": ["http://192.168.1.30:9200"],
    "record": {
        "test": "test2",
        "thisiscool": "yeahcool"
    }
}
```

## Output
The output of both these operations returns what elasticsearch sends back.

### Query
```json
{
  "return": {
    "_shards": {
      "failed": 0,
      "skipped": 0,
      "successful": 1,
      "total": 1
    },
    "hits": {
      "hits": [
        {
          "_id": "Ya0eqngBJ4CrZdCbJqyO",
          "_index": "direktiv",
          "_score": 1,
          "_source": {
            "test": "test2"
          },
          "_type": "_doc"
        },
        {
          "_id": "Yq0fqngBJ4CrZdCbZ6xR",
          "_index": "direktiv",
          "_score": 1,
          "_source": {
            "test": "test2",
            "thisiscool": "yeahcool"
          },
          "_type": "_doc"
        }
      ],
      "max_score": 1,
      "total": {
        "relation": "eq",
        "value": 2
      }
    },
    "timed_out": false,
    "took": 1
  }
}{
  "return": {
    "_shards": {
      "failed": 0,
      "skipped": 0,
      "successful": 1,
      "total": 1
    },
    "hits": {
      "hits": [
        {
          "_id": "Ya0eqngBJ4CrZdCbJqyO",
          "_index": "direktiv",
          "_score": 1,
          "_source": {
            "test": "test2"
          },
          "_type": "_doc"
        },
        {
          "_id": "Yq0fqngBJ4CrZdCbZ6xR",
          "_index": "direktiv",
          "_score": 1,
          "_source": {
            "test": "test2",
            "thisiscool": "yeahcool"
          },
          "_type": "_doc"
        }
      ],
      "max_score": 1,
      "total": {
        "relation": "eq",
        "value": 2
      }
    },
    "timed_out": false,
    "took": 1
  }
}
```


### Store
```json
{
  "return": {
    "_id": "Yq0fqngBJ4CrZdCbZ6xR",
    "_index": "direktiv",
    "_primary_term": 1,
    "_seq_no": 1,
    "_shards": {
      "failed": 0,
      "successful": 1,
      "total": 2
    },
    "_type": "_doc",
    "_version": 1,
    "result": "created"
  }
}
```

## Error

In the case that an error is encountered, it will present in the following format:

```json
{
    "errorCode": "com.elasticsearch.error",
    "errorMsg": "Something went wrong"
}
```