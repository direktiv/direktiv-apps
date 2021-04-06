+++
{
  "image": "vorteil/tweet",
  "desc": "Tweets a message to a twitter account."
}
+++

# Tweet

Tweets a message to a twitter account.

## Direktiv

An example workflow of using the container on Direktiv

```yaml
id: tweet-message
function: 
- id: tweet
  image: vorteil/tweet
description: "Tweet a message"
states:
- id: postATweet
  type: action
  action:
    function: tweet
    input: .
```

## Input
The input needed to run the above workflow properly is the following:

```json
{
    "consumerKey": ".secrets.CONSUMER_KEY",
    "consumerSecret": ".secrets.CONSUMER_SECRET",
    "tokenKey": ".secrets.TOKEN_KEY",
    "tokenSecret": ".secrets.TOKEN_SECRET",
    "message": "This is getting tweeted :)"
}
```

To obtain a consumer key and secret you need a Twitter developer account which can be found [here](https://developer.twitter.com/en/portal/dashboard). You can also generate your token key and secret when you create your application on the Twitter development account. Make sure you provide the account with read and write.


## Output

If the action is successful, the tweet should appear on the account with the given credentials.

## Error

In the case that an error is encountered, it will present in the following format:

```json
{
    "errorCode": "com.tweet.error",
    "errorMsg": "Something went wrong"
}
```