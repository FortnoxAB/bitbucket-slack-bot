# bitbucket-slack-bot

This bot tries its best to only send messages about pull requests that can be act apon.
For example notify author if we have at least 2 approvers. 

TODO:
Implement canmerge check using read only api user ([merge.go](models/merge.go) )

# Setup

Configure webhooks for your repo to go to `http://bot-url/webhook`. All event's Pull request category should be checked.

# Running

```
Usage of bitbucket-slack-bot:
  -log-format
    	Change value of Log-Format. (default <nil>)
  -log-level
    	Change value of Log-Level.
  -port
    	Change value of Port. (default 8080)
  -token
    	Change value of Token.

Generated environment variables:
   CONFIG_LOG_FORMAT
   CONFIG_LOG_LEVEL
   CONFIG_PORT
   CONFIG_TOKEN

```
