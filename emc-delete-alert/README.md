---
{
  "image": "vorteil/emc-delete-alert",
  "desc": "Delete an alert"
}
---

# EMC Delete Alert

Deletes an alert on EMC.

## Direktiv

An example workflow that deletes an alert on EMC.

```yaml
id: delete-alert
functions:
- id: delete
  image: vorteil/emc-delete-alert:v1
description: "Delete the alert"
states:
  - id: delete-alert-test
    type: action
    action: 
      secrets: ["PASSWORD"]
      function: delete
      input: 
        username": "admin"
        password": jq(.secrets.PASSWORD)
        id": "alert_48"
```

## Input

```yaml
input:
  username: "admin"
  password: jq(.secrets.PASSWORD)
  id: "alert_48"
```

## Output

No output is returned upon success but the alert should now be cleared from the list.

## Error

```json
{
    "errorCode": "com.emc.delete-alert.%s.error",
    "errorMsg": "Something went wrong"
}
```

**NOTE:** '%s' will be replaced contextually with what error has been received.