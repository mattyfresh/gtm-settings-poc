# Go Bot

Run a bot using the Slack RTM api

# Getting Started

1. Ensure you have an env variable `GOOGLE_APPLICATION_CREDENTIALS` set: https://developers.google.com/accounts/docs/application-default-credentials
1. Set a `GTM_ACCOUNT_ID` env variable. This the account ID from google tag manager `https://tagmanager.google.com/?authuser=0#/container/accounts/${account_id_here}/containers/${container_id_here}/workspaces/${workspace_number_here}`. So `account_id_here` will be your `GTM_ACCOUNT_ID`.
1. Get and set `SLACK_BOT_API_TOKEN` as an env variable: follow the _Getting Started_ section here https://api.slack.com/bot-users. Whatever you decide to name your bot, remember the name and invite your bot to the slack channel you would like it to be active in.
1. run the bot: `go run main.go`
1. Go to Slack and use the `/invite @name_of_your_bot` to bring the bot into a channel (might want to do this in a channel that won't be bothered by a little bit of noise).
1. Give it a go! For example: enter `@gobot hello` into Slack. You should get a response back.
1. Congrats! You are up and running

# Google Tag Manager Commands

_@NB this is a WIP, currently points to a public repository just as a POC. If you want this to work for you, you need to generate a `GITHUB_PERSONAL_ACCESS_TOKEN` and change the `REPO_URL` this scripts references_

- `@bot gtm validate ${name_of_container}` will pull the latest workspace, run it against our validation spec (@TODO this is just hard-coded into the app for now), and print out any possible errors.
- `@bot gtm publish ${name_of_container}` will run all of the `gtm validate` commands, plus run a bash script to pull down the GTM config repo, add your changes, and push a new branch to that repo with the proposed changes. @TODO this is just hard-coded into `github-commit.sh` right now, these should come from the env somehow.
