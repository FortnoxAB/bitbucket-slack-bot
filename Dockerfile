FROM alpine:3.8
RUN apk add --no-cache ca-certificates
COPY bitbucket-slack-bot /
ENTRYPOINT ["/bitbucket-slack-bot"]
