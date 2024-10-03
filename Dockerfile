FROM gcr.io/distroless/static-debian12:nonroot
COPY bitbucket-slack-bot /
USER nonroot
ENTRYPOINT ["/bitbucket-slack-bot"]
