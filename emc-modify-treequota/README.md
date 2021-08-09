---
{
    "image": "vorteil/emc-modify-treequota",
    "desc": "Modifies the hardlimit, softlimit or description on a treeQuota."
}
---

# EMC Modify TreeQuota

Modifies the hardlimit, softlimit or description on a treeQuota.

## Direktiv

An example workflow that adjusts both the hardlimit, softlimit by adding an extra 10 gigabytes to it and changing the description to 'This is a new description.'.

```yaml
id: modify-tree-quota
functions:
- id: modify
  type: reusable
  image: trentis/emc-test-quota
description: "Modify the tree quota."
states:
  - id: get
    type: action
    action: 
      secrets: ["PASSWORD"]
      function: modify
      input: 
        url: "https://unity.direktiv.io"
        username: "admin"
        password: jq(.secrets.PASSWORD)
        path: "/direktivquota"
        description: "This is a new description."
        hardlimit: "10G"
        softlimit: "10G"
```

## Input

```yaml
input:
  url: "https://unity.direktiv.io"
  username: "admin"
  password: jq(.secrets.PASSWORD)
  path: "/direktivquota"
  description: "This is a new description."
  hardlimit: "10G"
  softlimit: "10G"
```

**NOTE:** 'description' is an optional field.

## Output

Upon success the following output is the size of the new quota adjustments.

```json
{
	"return": {
		"description": "This is a new description.",
		"hardlimit": "20G",
		"softlimit": "18G"
	}
}
```

## Error

In the case that an error is encountered, it will present in the following format: 

```json
{
    "errorCode": "com.emc-modify-tree-quota.%s.error",
    "errorMsg": "Something went wrong"
}
```

**NOTE:** '%s' will be replaced contextually with what error has been received.