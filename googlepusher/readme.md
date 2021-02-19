# Store

Store appends to a google spreadsheet.

## Input

- Grab a service account key.
- Look at the client_email in the credentials.json and provide access for that email to the spreadsheet ID you wish to edit.
- Enable the google sheets api for the project listed in the service account key.
- Send the following input file so the program knows what to access and write to the spreadsheet

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
    "values": ["Hehe", "hehe2", "hehe3", "hehe4", "hehe5"]
}
```

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