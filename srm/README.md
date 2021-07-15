---
{
  "image": "vorteil/srm",
  "desc": "Sends metrics to a Dell SRM service"
}
---

# SRM

Sends metrics to a Dell SRM service

## Direktiv

An example workflow of sending metrics to the SRM service:

```yaml
id: dump-metrics
functions:
- id: store
  image: vorteil/srm:v2
description: "Sends metrics to a Dell SRM service"
states:
- id: write-metrics
  type: action
  action:
    function: store
    input:
      srm-lb-address: "192.168.1.30"
      srm-lb-port: "2000"
      group: "group"
      timestamp: 1625900409
      variable: "rcc-latency-collector.au2004lpablah3.nsx.elapsed"
      value: 100.01
      name: "Availability"
      source: "rcc-latency-collector"
      device: "au2004lpablah3"
      devtype: "host"
	  function: "nsx"
	  saas: "CIaaS"
      part: "nsx"
	  unit: "ms"
      parttype: "RESTAPI"
```

- 'srm-lb-address' is a required field
- 'srm-lb-port' is a required field
- 'value' is a required field
- 'variable' is a required field
- 'timestamp' is an optional field, if not provided will be set to the current time in UNIX.
- 'group' is an optional field, if not provided will be set to 'group'
- any other fields are optional and will be added to the end of the TCP message as KEY=VALUE

## Output

If the action is successful, no output will be returned.

## Error

In the case that an error is encountered, it will present in the following format:

```json
{
    "errorCode": "com.srm.%s.error",
    "errorMsg": "Something went wrong"
}
```