# Write data to a google spreadsheet

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

- Output.json file should look like something below or contain an error depending on your configuration

```json
{
    "error":"",
    "response":"",
    "statusCode":0,
    "status":""
}
```