#!/bin/bash

docker run -e SLACK_BOT_API_TOKEN=$SLACK_BOT_API_TOKEN \
    -e GITHUB_ACCESS_TOKEN=$GITHUB_ACCESS_TOKEN \
    -e GTM_ACCOUNT_ID=$GTM_ACCOUNT_ID \
    -e GOOGLE_APPLICATION_CREDENTIALS="google-creds.json" \
    mpadich/gobot:latest # CHANGE ME to the tag name you gave during the docker build step!