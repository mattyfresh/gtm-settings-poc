#!/bin/bash
REPO_URL="https://github.com/mattyfresh/gtm-settings-poc"
BRANCH_NAME="$1"

# Config
git config credential.helper "cache --timeout=120"
git config user.email "matthew.padich@gmail.com"
git config user.name "Slack Bot!"

# Reset
rm -rf gtm-settings-poc

# Clone Repo
git clone "$REPO_URL.git"

# Copy Config File to repo
cp gtm-config.json ./gtm-settings-poc/gtm-config.json

# cd into config repo, commit and push
cd gtm-settings-poc/
git checkout -b $BRANCH_NAME
git add gtm-config.json
git commit -m "update to gtm config via Slack"

# Push quietly to prevent showing the token in log
git push -q https://${GITHUB_PERSONAL_ACCESS_TOKEN}@github.com/mattyfresh/gtm-settings-poc.git $BRANCH_NAME

# Echo out link to create new PR
# @NB '@@@' is for easy text parsing for output in Slack
echo "@@@$REPO_URL/pull/new/$BRANCH_NAME@@@"