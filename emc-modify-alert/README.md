---
{
    "image": "direktiv/emc-modify-alert",
    "desc": "Acknowledges an alert."
}
---

# EMC Modify Alert

Acknowledges an alert.

## Direktiv

An example workflow that acknowledges an alert on EMC.

```yaml
id: acknowledge-alert
functions:
- id: acknowledge
  image: direktiv/emc-modify-alert:v1
description: "Acknowledge the alert."
states:
  - id: acknowledge-alert-test
    type: action
    action: 
      secrets: ["PASSWORD"]
      function: acknowledge
      input: 
        url: "https://unity.direktiv.io"
        username: "admin"
        password: jq(.secrets.PASSWORD)
        message: "This is a test message to be sent in an SNMP trap."
```

## Input

```yaml
input:
  url: "https://unity.direktiv.io"
  username: "admin"
  password: jq(.secrets.PASSWORD)
  message: "This is a test message to be sent in an SNMP trap."
```

## Output

Upon success the following output is the ID of the alert so we can delete it later on in the workflow.

```json
{
    "return": {
        "id": "alert_50"
    }
}
```

## Error

In the case that an error is encountered, it will present in the following format:

```json
{
    "errorCode": "com.emc.modify-alert.%s.error",
    "errorMsg": "Something went wrong"
}
```

**NOTE:** '%s' will be replaced contextually with what error has been received.