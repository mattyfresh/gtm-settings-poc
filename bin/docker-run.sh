#!/bin/bash

set -e
set -u

docker run -e SLACK_BOT_API_TOKEN=$SLACK_BOT_API_TOKEN \
    -e GITHUB_ACCESS_TOKEN=$GITHUB_ACCESS_TOKEN \
    -e GTM_ACCOUNT_ID=$GTM_ACCOUNT_ID \
    -e GH_USER=$GH_USER \
    -e GTM_PASSWORD=$GH_PASSWORD \
    -e GOOGLE_APPLICATION_CREDENTIALS="google-creds.json" \
    artnet/gobot:latest