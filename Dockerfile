FROM gcr.io/distroless/static-debian11:nonroot
COPY bitbucket-slack-bot /
USER nonroot
ENTRYPOINT ["/bitbucket-slack-bot"]
