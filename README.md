# bitbucket-slack-bot
A bot to send alertmanager notifications to slack to reduce webhook configuration overhead. 
Just setup the api_url in slack_config in alertmanager to your running instance with path: /webhook

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
