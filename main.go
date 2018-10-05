package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/nlopes/slack"
)

func respond(rtm *slack.RTM, msg *slack.MessageEvent, prefix string) {
	text := msg.Text

	// get username of the person who mentioed the bot
	user, err := rtm.GetUserInfo(msg.User)
	if err != nil {
		fmt.Println("error retrieving user name", err)
		return
	}
	userName := user.Name

	// clean up the msg text
	text = strings.TrimPrefix(text, prefix)
	text = strings.TrimSpace(text)
	text = strings.ToLower(text)

	acceptedGreetings := map[string]bool{
		"hey!":       true,
		"hello":      true,
		"what's up?": true,
		"yo":         true,
	}
	acceptedHowAreYou := map[string]bool{
		"how's it going?": true,
		"how are ya?":     true,
		"feeling okay?":   true,
	}

	if acceptedGreetings[text] {
		response := fmt.Sprintf("What's up %s?", userName)
		rtm.SendMessage(rtm.NewOutgoingMessage(response, msg.Channel))
	} else if acceptedHowAreYou[text] {
		response := fmt.Sprintf("I'm doing well, how are you %s?", userName)
		rtm.SendMessage(rtm.NewOutgoingMessage(response, msg.Channel))
	}
}

func main() {
	token := os.Getenv("SLACK_BOT_API_TOKEN")

	// slack API service
	service := slack.New(token)
	service.SetDebug(true)

	// log to stdout so we can see what's going on
	logger := log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)
	slack.SetLogger(logger)

	// real time messaging service
	rtm := service.NewRTM()

	// open the websocket connection
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		fmt.Print("Event Received: ")

		switch slackEvent := msg.Data.(type) {
		case *slack.HelloEvent:
			// Ignore hello
			// fmt.Println("HELLO!")
		case *slack.ConnectedEvent:
			// fmt.Println("Infos:", eventType.Info)
			// fmt.Println("Connection counter:", eventType.ConnectionCount)

		case *slack.MessageEvent:
			info := rtm.GetInfo()
			gtmPrefix := fmt.Sprintf("<@%s> gtm", info.User.ID)

			// @TODO parse request for GTM and build API calls
			if slackEvent.User != info.User.ID && strings.HasPrefix(slackEvent.Text, gtmPrefix) {
				rtm.SendMessage(rtm.NewOutgoingMessage("google tag manager integration coming soon!", slackEvent.Channel))
				break
			}

			// fallback to generic responses
			if slackEvent.User != info.User.ID && strings.HasPrefix(slackEvent.Text, fmt.Sprintf("<@%s> ", info.User.ID)) {
				respond(rtm, slackEvent, fmt.Sprintf("<@%s> ", info.User.ID))
			}

		case *slack.PresenceChangeEvent:
			fmt.Printf("Presence Change: %v\n", slackEvent)

		case *slack.LatencyReport:
			fmt.Printf("Current latency: %v\n", slackEvent.Value)

		case *slack.RTMError:
			fmt.Printf("Error: %s\n", slackEvent.Error())

		case *slack.InvalidAuthEvent:
			fmt.Printf("Invalid credentials")
			return

		default:
			// Ignore other events..
			// fmt.Printf("Unexpected: %v\n", msg.Data)
		}
	}
}
