#!/bin/bash

set -eu

BRANCH_NAME="$1"

# Config
# git config --global credential.helper "cache --timeout=120"

# @TODO use bot email here
# git config --global user.email "matthew.padich@gmail.com"
# git config --global user.name "Slack Bot!"

rm -rf gtm-settings-poc

# Clone config repo and copy new config into it
git clone "https://${GH_USER}:${GH_PASSWORD}@github.com/mattyfresh/gtm-settings-poc.git"
cp gtm-config.json ./gtm-settings-poc/gtm-config.json

# cd into config repo, commit and push
cd gtm-settings-poc/
git checkout -b $BRANCH_NAME
git add gtm-config.json
git commit -m "update GTM config via Slack on $(date)"

# Push quietly to prevent showing the token in log
git push -q "https://${GH_USER}:${GH_PASSWORD}@github.com/mattyfresh/gtm-settings-poc.git $BRANCH_NAME"

# Echo out link to create new PR
# @NB '@@@' is for easy text parsing for output in Slack
echo "@@@https://github.com/mattyfresh/gtm-settings-poc/pull/new/$BRANCH_NAME@@@"
