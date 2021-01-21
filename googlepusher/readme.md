# Push Data to Excel Sheet

- Grab a service account key.
- Look at the client_email in the credentials.json and provide access for that email to the spreadsheet ID you wish to edit.
- Enable the google sheets api for the project listed in the service account key.
- Send the following input file so the program knows what to access and write to the spreadsheet

```json
{
    "authentication" : {
        "type": "service_account",
        "project_id": "PROJECT_ID",
        "private_key_id": "PRIVATE_KEY_ID",
        "private_key": "PRIVATE KEY",
        "client_email": "CLIENT_EMAIL",
        "client_id": "CLIENT_ID",
        "auth_uri": "https://accounts.google.com/o/oauth2/auth",
        "token_uri": "https://oauth2.googleapis.com/token",
        "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
        "client_x509_cert_url": "CLIENT_X509_CERT_URL"
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