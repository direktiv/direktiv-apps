---
{
	"image": "vorteil/servicenow",
	"desc": "perform a Service Now API request"
}
---

# Service Now

The `vorteil/service-now` image is a work-in-progress, but is fairly straightforward to use. This example of the input provided to the isolate contains fields to set the method of the request, and payload (body), basic authentication, and more. 

## Example Input

```json
{
  "method": "PATCH",
  "api": "change-request-normal",
  "insecureSkipVerify": true,
  "user": "exampleUser",
  "password": "examplePassword",
  "instance": "https://example.servicenow.com",
  "sys_id": "example_sys_id",
  "payload": {
    "short_description": "This sets the 'short_description' field of the targeted normal change request!"
  }
}
```

The above input instructs the container to send a PATCH request to the ServiceNow server hosted at `example.servicenow.com`, and to skip certificate verification. The `api` field is used to append a particular API prefix to the URL being used for the API request, followed by the value of the `sys_id` field. The API request body will contain a JSON payload defined by the `payload` field.

Knowing that the API request URL is created by combining `instance`, `api`, and `sys_id` allows users to send more advanced requests, such as by omitting a real `sys_id` value but instead including a number of query parameters instead. For example:

```json
{
  ...
  "api": "change-request-normal",
  "instance": "https://example.servicenow.com",
  "sys_id": "?short_description=HelloWorld"
}
```

The resulting URL used for sending the API request would look like:

```
https://example.servicenow.com/api/sn_chg_rest/change/normal?short_description=HelloWorld
```

Currently, the available values for the `api` field are limited to the following:

```yaml
cmdb-instance: api/now/cmdb/instance/cmdb_ci_storage_server
change-request-normal: api/sn_chg_rest/change/normal
change-request: api/sn_chg_rest/change
incident: api/now/table/incident
user: api/now/table/sys_user
```

## Example Usage

```yaml
id: example
functions:
- id: servicenow
  image: vorteil/servicenow:v8
description: "send a get request" 
states:
- id: example-state
  type: action
  action: 
    function: servicenow
    secrets: ["SERVICENOW_USER","SERVICENOW_PASSWORD","SERVICENOW_INSTANCE"]
    input: | 
      {
        "method": "PATCH",
        "api": "change-request-normal",
        "insecureSkipVerify": true,
        "user": .secrets.SERVICENOW_USER,
        "password": .secrets.SERVICENOW_PASSWORD,
        "instance": .secrets.SERVICENOW_INSTANCE,
        "sys_id": "example_sys_id",
        "payload": {
          "short_description": "This sets the 'short_description' field of the targeted normal change request!"
      }
}
```