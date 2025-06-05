# bitbucket-slack-bot

This bot tries its best to only send messages about pull requests that can be acted upon.
For example notify author if we have at least 2 approvers.

## Setup

Configure webhooks for your repo to go to `http://bot-url/webhook`. All event's Pull request category should be checked.

## Running tests

`go test -v ./...`

## Running

```
Usage of bitbucket-slack-bot:
  -bitbucketpassword
    	Change value of BitbucketPassword.
  -bitbucketurl
    	Change value of BitbucketURL.
  -bitbucketuser
    	Change value of BitbucketUser.
  -log-format
    	Change value of Log-Format. (default text)
  -log-formatter
    	Change value of Log-Formatter. (default <nil>)
  -log-level
    	Change value of Log-Level.
  -port
    	Change value of Port. (default 8080)
  -token
    	Change value of Token.

Generated environment variables:
   CONFIG_BITBUCKETPASSWORD
   CONFIG_BITBUCKETURL
   CONFIG_BITBUCKETUSER
   CONFIG_LOG_FORMAT
   CONFIG_LOG_FORMATTER
   CONFIG_LOG_LEVEL
   CONFIG_PORT
   CONFIG_TOKEN

```
