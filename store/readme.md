# Store

Add a row to a 'Google Sheets' spreadsheet.

## Direktiv
An example workflow of using the container in a workflow on Direktiv.
```yaml
id: store-excel
functions: 
- id: store
  image: vorteil/store
description: "Stores who accessed into excel"
states:
- id: storeExcel
  type: action
  action:
    function: store
    input: '{ 
    "authentication" : {
        "type": "service_account",
        "private_key": "PRIVATE KEY",
        "client_email": "CLIENT_EMAIL",
        "token_uri": "https://oauth2.googleapis.com/token",
    },
    "spreadSheetID": "SPREADSHEET_ID",
    "range": "A1",
    "values": ["These", "are", "individual", "columns", "!"]
}'
```

## Input

- On a Google Cloud project, create a Service Account.
- Download a JSON-format Service Account key file.
- Enable the 'Google Sheets API' on the Google Cloud project.
- Populate the following input object, using the `client_email` and `private_key` fields for their respective matching fields.

```json
{
    "authentication" : {
        "type": "service_account",
        "private_key": "PRIVATE KEY",
        "client_email": "CLIENT_EMAIL",
        "token_uri": "https://oauth2.googleapis.com/token",
    },
    "spreadSheetID": "SPREADSHEET_ID",
    "range": "A1",
    "values": ["These", "are", "individual", "columns", "!"]
}
```

*Note: Replace `SPREADSHEET_ID` with the ID of whichever Google Sheets document should be targeted.*
*Note: Providing 'debug' with the json struct will print everything about what the application is doing.

## Output

No output will be provided if successful

### Error

If an error was to happen during the container running the output is written back to direktiv as the following

```json
{
    "errorCode": "com.store.error",
    "errorMsg": "Something went wrong"
}
```