#!/bin/bash

GTM_CONFIG_REPO_URL="https://github.com/artnetworldwide/automation-googletagmanager-config"
BRANCH_NAME="$1"

# Config
git config --global credential.helper "cache --timeout=120"

# @TODO use bot email here
git config --global user.email "matthew.padich@gmail.com"
git config --global user.name "gobot"

rm -rf automation-googletagmanager-config

# Clone config repo and copy new config into it
git clone "$GTM_CONFIG_REPO_URL.git"
cp gtm-config.json ./automation-googletagmanager-config/gtm-config.json

# cd into config repo, commit and push
cd automation-googletagmanager-config/
git checkout -b $BRANCH_NAME
git add gtm-config.json
git commit -m "update GTM config via Slack on $(date)"

# Push quietly to prevent showing the token in log
git push -q https://${GITHUB_ACCESS_TOKEN}@github.com/artnetworldwide/automation-googletagmanager-config.git $BRANCH_NAME

# Echo out link to create new PR
# @NB '@@@' is for easy text parsing for output in Slack
echo "@@@$GTM_CONFIG_REPO_URL/pull/new/$BRANCH_NAME@@@"
