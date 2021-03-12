# Direktiv Apps

Simple containers that run on Direktiv.

## Cloud Functions

| Image | Description | How To Use |
| ------------- | ------------- | ---------- |
| [vorteil/azinvoke](https://hub.docker.com/r/vorteil/azinvoke) | Executes an azure cloud function via authentication using function keys. | [README](https://github.com/vorteil/direktiv-apps/tree/master/azinvoke) |
| [vorteil/lambda](https://hub.docker.com/r/vorteil/lambda) | Executes an aws cloud function via authentication using access tokens. | [README](https://github.com/vorteil/direktiv-apps/tree/master/lambdainvoke) |
| [vorteil/googleinvoke](https://hub.docker.com/r/vorteil/googleinvoke) | Executes a google cloud function via authentication using service account key | [README](https://github.com/vorteil/direktiv-apps/tree/master/googleinvoke) |

## CLIs

| Image | Description | How To Use |
| ------------- | ------------- | ---------- |
| [vorteil/azgo](https://hub.docker.com/r/vorteil/azgo)  | Executes azure CLI commands using authentication provided by a service principal account.  | [README](https://github.com/vorteil/direktiv-apps/tree/master/azcli) |
| [vorteil/gcloud](https://hub.docker.com/r/vorteil/gcloud) | Executes a gcloud command using authentication provided by the service account key file.| [README](https://github.com/vorteil/direktiv-apps/tree/master/googlecli)|
| [vorteil/awsgo](https://hub.docker.com/r/vorteil/awsgo) | Executes an aws cli command using authentication provided by an access token. | [README](https://github.com/vorteil/direktiv-apps/tree/master/awscli) |


## Utility

| Image | Description | How To Use |
| ------------- | ------------- | ---------- |
| [vorteil/request](https://hub.docker.com/r/vorteil/request) | Sends a http request and returns the entire output. | [README](https://github.com/vorteil/direktiv-apps/tree/master/requester)|
| [vorteil/twilio](https://hub.docker.com/r/vorteil/twilio) | Sends an email or a sms message. | [README](https://github.com/vorteil/direktiv-apps/tree/master/twilio) | 
| [vorteil/tweet](https://hub.docker.com/r/vorteil/tweet) | Tweets a message out. | [README](https://github.com/vorteil/direktiv-apps/tree/master/twitter) |
| [vorteil/store](https://hub.docker.com/r/vorteil/store) | Appends to a google spreadsheet. | [README](https://github.com/vorteil/direktiv-apps/tree/master/googlepusher) |
| [vorteil/smtp](https://hub.docker.com/r/vorteil/smtp) | Creates a smtp client to send an email. | [README](https://github.com/vorteil/direktiv-apps/tree/master/smtp) |

## Chat Applications

| Image | Description | How To Use |
| ------------- | ------------- | ---------- |
| [vorteil/discordmsg](https://hub.docker.com/r/vorteil/discordmsg) | Sends a message to a discord channel from the provided webhook url. | [README](https://github.com/vorteil/direktiv-apps/tree/master/discord)|
| [vorteil/googlemsg](https://hub.docker.com/r/vorteil/googlemsg) | Sends a message to a google chat channel from the provided webhook url. | [README](https://github.com/vorteil/direktiv-apps/tree/master/googlechat) |
| [vorteil/slackmsg](https://hub.docker.com/r/vorteil/slackmsg) | Sends a message to a slack chat channel from the provided webhook url. | [README](https://github.com/vorteil/direktiv-apps/tree/master/slack) |


## AI

| Image | Description | How To Use |
| ------------- | ------------- | ---------- |
| [vorteil/imagerecognition](https://hub.docker.com/r/vorteil/imagerecognition) | Using the google vision api it will read an image URL and return whether it is safe for work. | [README](https://github.com/vorteil/direktiv-apps/tree/master/image-recognition)


## Databases

| Image | Description | How To Use |
| ------------- | ------------- | ---------- |
| [vorteil/simplepostgresclient](https://hub.docker.com/r/vorteil/simplepostgresclient) | Connects to a postgres database and allows simple sql commands. | [README](https://github.com/vorteil/direktiv-apps/tree/master/simple-postgres-client)


## General Demos

| Image | Description | How To Use |
| ------------- | ------------- | ---------- |
| [vorteil/solve](https://hub.docker.com/r/vorteil/solve) | Solves an array of expressions and returns the answers. | [README](https://github.com/vorteil/direktiv-apps/tree/master/solver) |
| [vorteil/greeting](https://hub.docker.com/r/vorteil/greeting) | Returns a message to greet you | [README](https://github.com/vorteil/direktiv-apps/tree/master/greeter) |
