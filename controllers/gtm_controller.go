package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-bot/services"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/nlopes/slack"
	gtm "google.golang.org/api/tagmanager/v2"
)

var realTimeMessenger *slack.RTM
var gtmService *gtm.Service

// convenience shorthand for sending a slack msg
func sendMessage(msg string, channelID string) {
	realTimeMessenger.SendMessage(realTimeMessenger.NewOutgoingMessage(msg, channelID))
}

func gtmInit() error {
	if gtmService == nil {
		s, err := services.ConfigureGTMService()
		if err != nil {
			return err
		}
		gtmService = s
	}
	return nil
}

func getContainerByName(accountPath string, containerName string) (string, error) {
	containers, err := gtmService.Accounts.Containers.List(accountPath).Do()
	if err != nil {
		return "", err
	}

	var containerID string
	for _, c := range containers.Container {
		if c.Name == containerName {
			containerID = c.ContainerId
			break
		}
	}
	return containerID, nil
}

func getDefaultWorkspaceID(containerPath string) (string, error) {
	workspaceResp, workspaceErr := gtmService.Accounts.Containers.Workspaces.List(containerPath).Do()
	if workspaceErr != nil {
		return "", workspaceErr
	}
	defaultWorkspaceName := "Default Workspace"
	var defaultWorkspaceID string
	for _, ws := range workspaceResp.Workspace {
		if ws.Name == defaultWorkspaceName {
			defaultWorkspaceID = ws.WorkspaceId
			break
		}
	}
	if defaultWorkspaceID == "" {
		return "", errors.New("Workspace ID for '" + defaultWorkspaceName + "' could not be found")
	}
	return defaultWorkspaceID, nil
}

func parseCommand(msg *slack.MessageEvent) (commandType, commandName string) {
	pattern := regexp.MustCompile(`gtm\s+(?P<command_type>\w*)\s+(?P<container_name>\S*)?`)
	match := pattern.FindStringSubmatch(msg.Text)
	if match == nil {
		return
	}

	captures := make(map[string]string)
	for i, name := range pattern.SubexpNames() {
		// ignore the whole regexp match and unnamed groups
		if i == 0 || name == "" {
			continue
		}
		captures[name] = match[i]
	}

	return captures["command_type"], captures["container_name"]
}

// GtmHandler controller for all GTM related events
func GtmHandler(msg *slack.MessageEvent, rtm *slack.RTM) {
	// expose rtm as a singleton
	realTimeMessenger = rtm

	// parse the slack command
	commandType, containerName := parseCommand(msg)
	if !(commandType == "publish" || commandType == "validate") {
		sendMessage("There are two commands available, `publish` and `validate`. For example, '@' the bot and try `gtm validate ${name_of_container}`", msg.Channel)
		return
	}

	accountID := os.Getenv("GTM_ACCOUNT_ID")
	if accountID == "" {
		sendMessage(":x: The `GTM_ACCOUNT_ID` has not been set, please contact tech support", msg.Channel)
		return
	}

	// initialize GTM api service
	if initErr := gtmInit(); initErr != nil {
		fmt.Printf("Error initializing GTM api: %s", initErr.Error())
		return
	}

	accountPath := fmt.Sprintf("accounts/%s", accountID)

	var containerID string
	cID, err := getContainerByName(accountPath, containerName)
	if err != nil {
		noValidContainerMsg := fmt.Sprintf(":crying_cat_face: No valid container with name `%s` found... double check the name of your container!", containerName)
		sendMessage(noValidContainerMsg, msg.Channel)
		return
	}
	containerID = cID

	containerPath := fmt.Sprintf("accounts/%s/containers/%s", accountID, containerID)

	// get current active workspace
	workspaceID, workspaceErr := getDefaultWorkspaceID(containerPath)
	if workspaceErr != nil {
		sendMessage("There was an error getting the default workspace ID: "+workspaceErr.Error(), msg.Channel)
		return
	}

	// let the user know we're validating
	wipMsg := fmt.Sprintf(":ram: validating workspace #%s for container with ID %s", workspaceID, containerID)
	sendMessage(wipMsg, msg.Channel)

	var validationErrors []string
	workspacePath := fmt.Sprintf("%s/workspaces/%s", containerPath, workspaceID)

	// validate variables
	allVars, allVarsErr := gtmService.Accounts.Containers.Workspaces.Variables.List(workspacePath).Do()
	if allVarsErr != nil {
		fmt.Printf(allVarsErr.Error())
		return
	}

	for _, v := range allVars.Variable {
		if errors := ValidateVariable(v); errors != nil {
			for _, e := range errors {
				validationErrors = append(validationErrors, e.Error())
			}
		}
	}

	// validate tags
	allTags, tagsErr := gtmService.Accounts.Containers.Workspaces.Tags.List(workspacePath).Do()
	if tagsErr != nil {
		fmt.Printf(tagsErr.Error())
		return
	}

	for _, t := range allTags.Tag {
		if errors := ValidateTag(t); errors != nil {
			for _, e := range errors {
				validationErrors = append(validationErrors, e.Error())
			}
		}
	}

	// validate triggers
	allTriggers, triggersErr := gtmService.Accounts.Containers.Workspaces.Triggers.List(workspacePath).Do()
	if triggersErr != nil {
		fmt.Printf(triggersErr.Error())
		return
	}

	for _, t := range allTriggers.Trigger {
		if errors := ValidateTrigger(t); errors != nil {
			for _, e := range errors {
				validationErrors = append(validationErrors, e.Error())
			}
		}
	}

	if len(validationErrors) > 0 {
		sendMessage(strings.Join(validationErrors, "\n"), msg.Channel)
		sendMessage(":crying_cat_face: Validation Failed! Please fix the above errors and try again...", msg.Channel)
		return
	}

	sendMessage(":thumbsup: Validation Succeeded!", msg.Channel)

	if commandType == "publish" {
		// build and write JSON to file
		allOutput := []interface{}{allTriggers, allTags, allVars}
		file, err := os.Create("gtm-config.json")
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		defer file.Close()

		outputJSON, jsonErr := json.MarshalIndent(allOutput, "", "    ")
		if jsonErr != nil {
			fmt.Println(jsonErr.Error())
			return
		}
		fmt.Fprint(file, string(outputJSON))

		// create and push a commit with new GTM config file to github
		branchName := fmt.Sprintf("workspace-%s-%d", workspaceID, time.Now().Unix())
		command := exec.Command("/bin/bash", "github-commit.sh", branchName)
		absPath, absPathErr := filepath.Abs(".")
		if absPathErr != nil {
			fmt.Println(absPathErr.Error())
			return
		}
		command.Dir = absPath

		// get output from bash script
		out, err := command.Output()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		pullRequestLink := strings.Split(string(out), "@@@")[1]
		sendMessage(fmt.Sprintf(":shipit: Publish success! Click below to create a PR:\n\n %s", pullRequestLink), msg.Channel)
	}
}
