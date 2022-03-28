# SNMP

Sends a snmp trap to an address

## Direktiv

An example workflow of sending `Hello World!` as a SNMP Trap:

```yaml
id: send-snmp-message
functions:
- id: snmp
  image: direktiv/snmp:v3
description: Sends a SNMP trap to an address
states:
- id: write-message
  type: action
  action:
    function: snmp
    input:
      server: 192.168.0.177
      port: 161
      version: 1
      enterprise: 1.3.6.1.2.1.1.6
      specific: 1
      generic: 1
      variables:
      - name: 1.3.6.1.2.1.1.6.1
        type: 4
        value: Hello
```


## Configuration Values

| Attribute | Function |
| -- | -- |
|server|Server to send the trap to|
|port|Port to use, defaults to 161|
|version|SNMP version. Can be 1, 2 or 3. Defaults to 2 which means 2c|
|transport|TCP or UDP, defaults to UDP|
|community|Community string, defaults to public|
|inform|Boolean to use it as a inform request|
|generic|Number of generic type, e.g. 6|
|specific|Number of specific type, e.g. 100|
|enterprise|Type of managed object, e.g. 1.3.6.1.4.1.3224|


## Error on Execution

This container throws a `direktiv.snmp.error` if there are connectivity issues. 
