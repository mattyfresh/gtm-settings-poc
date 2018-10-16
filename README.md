# Go Bot

A Slack bot for Chat Ops

# Required Env Vars

- `GOOGLE_APPLICATION_CREDENTIALS`: required for authenticating a service account for GTM
  - see https://developers.google.com/accounts/docs/application-default-credentials
- `GTM_ACCOUNT_ID`: the google tag manager account ID you want to validate / publish to.
  - e.g. `https://tagmanager.google.com/?authuser=0#/container/accounts/${GTM_ACCOUNT_ID}/containers/${CONTAINER_ID}/workspaces/${WORKSPACE_ID}`
- `SLACK_BOT_API_TOKEN`: an access token that allows your bot to communicate with our slack instance via the API
  - follow the _Getting Started_ section here: https://api.slack.com/bot-users
- `GITHUB_ACCESS_TOKEN`: used to create and push commits to a github repo with the GTM configuration.
  - see https://github.com/settings/tokens

# Getting Started

1. Set up your bot using the instructions at https://api.slack.com/bot-users and set the `SLACK_BOT_API_TOKEN` once your bot user is created. Remember to invite your bot to whichever channel you would like it to be active in. There is a Slack shortcut for inviting members: `/invite @user_name`.
1. ensure you are in the project root, then run the bot: `go run main.go`.
1. Give it a go! For example, if you named your both `gobot`: type `@gobot hello` into Slack. You should get a response back.
1. Congrats! You are up and running.

# Google Tag Manager Commands

_@NB this is a WIP, currently points to a public repository just as a POC. If you want this to work for you, you need to generate a `GITHUB_ACCESS_TOKEN` and change the `GTM_CONFIG_REPO_URL` this scripts references_

- `@bot gtm validate ${name_of_container}` will pull the latest workspace, run it against our validation spec (@TODO this is just hard-coded into the app for now), and print out any possible errors.
- `@bot gtm publish ${name_of_container}` will run all of the `gtm validate` commands, plus run a bash script to pull down the GTM config repo, add your changes, and push a new branch to that repo with the proposed changes. @TODO this is just hard-coded into `github-commit.sh` right now, these should come from the env somehow.

# Docker

1. ensure all ENV variables are set
1. copy the contents of the file path that the `GOOGLE_APPLICATION_CREDENTIALS` env var points to and add it to a file called `google-creds.json` in the project root.
1. run `docker build -t some/alias .` to build a docker container with the tag `some/alias` (eg `mpadich/gobot`)
1. edit the `run.sh` to reflect the tag name you chose for your docker container.
1. run `./run.sh` to run your docker container.
