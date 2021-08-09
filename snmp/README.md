---
{
  "image": "vorteil/snmp",
  "desc": "Sends a snmp trap to an address"
}
---

# SNMP

Sends a snmp trap to an address


## Direktiv

An example workflow of sending `Hello World!` as a SNMP Trap:

```yaml
id: send-snmp-message
functions:
- id: snmp
  type: reusable
  image: vorteil/snmp:v1
description: Sends a SNMP trap to an address
states:
- id: write-message
  type: action
  action:
    function: snmp
    input:
      url: "0.0.0.0"
      port: 9442
      variables:
      - name: ".1.3.6.1.4.9.27"
        type: 4
        value: Hello World!
```
- **NOTE:** if you want to define the snmp to be an `inform-request` simply provide that boolean variable.

- **NOTE:** the variable `type` is an integer following the BER encoded types the table below defines the ints matching fields.

    | Type | Int |
    | ---- | --- |
    | End of Contents, Unknown Type | 0 |
    | Boolean | 1 |
    | Integer | 2 |
    | BitString | 3 |
    | OctetString | 4 |
    | Null | 5 |
    | ObjectIdentifier | 6 |
    | ObjectDescription | 7 |
    | IPAddress | 40 |
    | Counter32 | 41 | 
    | Gauge32 | 42 |
    | TimeTicks | 43 |
    | Opaque | 44 |
    | NsapAddress | 45 |
    | Counter64 | 46 |
    | Uinteger32 | 47 |
    | OpaqueFloat | 78 |
    | OpaqueDouble | 79 |
    | NoSuchObject | 80 |
    | NoSuchInstance | 81 |
    | EndOfMibView | 82 |

- **NOTE:** `variables` can also take a field called `snmpv1-headers` that has the following fields.

    ```yaml
  snmpv1-headers:
  enterprise: string
  agent-address: string
  generic-trap: int
  specific-trap: int
  timestamp: uint
  ```

## Output of Execution

The output of the execution if successfuly is the packet itself from sending the snmp trap.

## Error on Execution

If the container was to error out during the action an error will be thrown matching:

```json
{
    "errorCode": "com.snmp.%s.error",
    "errorMsg": "Something went wrong..."
}
```

**NOTE:** `%s` will be replaced with what is currently happening in the container that errored out e.g. `unmarshalling-input`.