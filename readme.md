# Direktiv Apps

Simple Containers that run on Direktiv

<em>created by [102272c7](https://github.com/vorteil/direktiv-apps/tree/102272c7a8425fdc8f14475943fdb66bf5f76b19)</em>


## Containers

| Image | Description | How to Use |
| ----- | ----------- | ---------- |
| [vorteil/amazon-sns](https://hub.docker.com/r/vorteil/amazon-sns) | Writes a message to Amazon's Simple Notification Service | [README](https://github.com/vorteil/direktiv-apps/tree/master/amazon-sns) |
| [vorteil/amazon-upload](https://hub.docker.com/r/vorteil/amazon-upload) | Uploads a base64 encoded string to a blob on Amazon. | [README](https://github.com/vorteil/direktiv-apps/tree/master/amazon-upload) |
| [vorteil/awsgo](https://hub.docker.com/r/vorteil/awsgo) | Executes any cli command with aws using a golang wrapper to provide the authentication before executing. | [README](https://github.com/vorteil/direktiv-apps/tree/master/awsgo) |
| [vorteil/awslog](https://hub.docker.com/r/vorteil/awslog) | Writes a log line to aws cloudwatch logs with provided log stream and group. | [README](https://github.com/vorteil/direktiv-apps/tree/master/awslog) |
| [vorteil/azgo](https://hub.docker.com/r/vorteil/azgo) | Executes any azure cli command using a golang wrapper to provide the authentication before executing. | [README](https://github.com/vorteil/direktiv-apps/tree/master/azgo) |
| [vorteil/azinvoke](https://hub.docker.com/r/vorteil/azinvoke) | Executes a cloud function on azure using the function name, function app and function key as authentication to do so. | [README](https://github.com/vorteil/direktiv-apps/tree/master/azinvoke) |
| [vorteil/azlog](https://hub.docker.com/r/vorteil/azlog) | Writes a log line to Azure Log Analytics Workspace. | [README](https://github.com/vorteil/direktiv-apps/tree/master/azlog) |
| [vorteil/azure-servicebus](https://hub.docker.com/r/vorteil/azure-servicebus) | Writes a message to Azure's Servicebus | [README](https://github.com/vorteil/direktiv-apps/tree/master/azure-servicebus) |
| [vorteil/azure-upload](https://hub.docker.com/r/vorteil/azure-upload) | Uploads a base64 encoded string to a blob on Azure. | [README](https://github.com/vorteil/direktiv-apps/tree/master/azure-upload) |
| [vorteil/debug](https://hub.docker.com/r/vorteil/debug) | Prints all the information that direktiv sends to the container | [README](https://github.com/vorteil/direktiv-apps/tree/master/debug) |
| [vorteil/discordmsg](https://hub.docker.com/r/vorteil/discordmsg) | Writes a discord message to a webhook URL. | [README](https://github.com/vorteil/direktiv-apps/tree/master/discordmsg) |
| [vorteil/gcloud](https://hub.docker.com/r/vorteil/gcloud) | Executes a gcloud cli command using a golang wrapper to provide authentication via service account key file. | [README](https://github.com/vorteil/direktiv-apps/tree/master/gcloud) |
| [vorteil/gcloud-instance-create](https://hub.docker.com/r/vorteil/gcloud-instance-create) | Create a compute engine instance on google cloud. | [README](https://github.com/vorteil/direktiv-apps/tree/master/gcloud-instance-create) |
| [vorteil/gcloud-instance-delete](https://hub.docker.com/r/vorteil/gcloud-instance-delete) | Delete a compute engine instance on google cloud. | [README](https://github.com/vorteil/direktiv-apps/tree/master/gcloud-instance-delete) |
| [vorteil/gcplog](https://hub.docker.com/r/vorteil/gcplog) | Writes a log line to a stackdriver logging implementation | [README](https://github.com/vorteil/direktiv-apps/tree/master/gcplog) |
| [vorteil/google-pubsub](https://hub.docker.com/r/vorteil/google-pubsub) | Writes a message to Google's Pubsub service. | [README](https://github.com/vorteil/direktiv-apps/tree/master/google-pubsub) |
| [vorteil/google-sentiment-check](https://hub.docker.com/r/vorteil/google-sentiment-check) | Reads a string and tells you the sentiment of the written text. | [README](https://github.com/vorteil/direktiv-apps/tree/master/google-sentiment-check) |
| [vorteil/google-translator](https://hub.docker.com/r/vorteil/google-translator) | Reads a string and converts it to the targeted language. | [README](https://github.com/vorteil/direktiv-apps/tree/master/google-translator) |
| [vorteil/google-upload](https://hub.docker.com/r/vorteil/google-upload) | Uploads a base64 encoded string to a bucket on Google. | [README](https://github.com/vorteil/direktiv-apps/tree/master/google-upload) |
| [vorteil/googleinvoke](https://hub.docker.com/r/vorteil/googleinvoke) | Executes a cloud function on google using a client authenticated via a service account key. | [README](https://github.com/vorteil/direktiv-apps/tree/master/googleinvoke) |
| [vorteil/googlemsg](https://hub.docker.com/r/vorteil/googlemsg) | Writes a google message to a webhook URL. | [README](https://github.com/vorteil/direktiv-apps/tree/master/googlemsg) |
| [vorteil/greeting](https://hub.docker.com/r/vorteil/greeting) | Outputs a greeting | [README](https://github.com/vorteil/direktiv-apps/tree/master/greeting) |
| [vorteil/imagerecognition](https://hub.docker.com/r/vorteil/imagerecognition) | Checks an image to see if it is safe for work and responds with content is likely to be racy, adult or violence themed. | [README](https://github.com/vorteil/direktiv-apps/tree/master/imagerecognition) |
| [vorteil/influxdb](https://hub.docker.com/r/vorteil/influxdb) | A container that provides the opportunities to write or query data on an InfluxDB instance. | [README](https://github.com/vorteil/direktiv-apps/tree/master/influxdb) |
| [vorteil/lambda](https://hub.docker.com/r/vorteil/lambda) | Executes a cloud function on aws using their golang SDK. | [README](https://github.com/vorteil/direktiv-apps/tree/master/lambda) |
| [vorteil/request](https://hub.docker.com/r/vorteil/request) | Perform a basic HTTP/S request. | [README](https://github.com/vorteil/direktiv-apps/tree/master/request) |
| [vorteil/simplepostgresclient](https://hub.docker.com/r/vorteil/simplepostgresclient) | Perform simple interactions on a PostgreSQL database.  | [README](https://github.com/vorteil/direktiv-apps/tree/master/simplepostgresclient) |
| [vorteil/slackmsg](https://hub.docker.com/r/vorteil/slackmsg) | Writes a message to a slack webhook url | [README](https://github.com/vorteil/direktiv-apps/tree/master/slackmsg) |
| [vorteil/smtp](https://hub.docker.com/r/vorteil/smtp) | A simple smtp client to send an email | [README](https://github.com/vorteil/direktiv-apps/tree/master/smtp) |
| [vorteil/solve](https://hub.docker.com/r/vorteil/solve) | Solves math expressions | [README](https://github.com/vorteil/direktiv-apps/tree/master/solve) |
| [vorteil/store](https://hub.docker.com/r/vorteil/store) | Add a row to a 'Google Sheets' spreadsheet. | [README](https://github.com/vorteil/direktiv-apps/tree/master/store) |
| [vorteil/tweet](https://hub.docker.com/r/vorteil/tweet) | Tweets a message to a twitter account. | [README](https://github.com/vorteil/direktiv-apps/tree/master/tweet) |
| [vorteil/twilio](https://hub.docker.com/r/vorteil/twilio) | Sends an email or SMS message using Twilio. | [README](https://github.com/vorteil/direktiv-apps/tree/master/twilio) |

