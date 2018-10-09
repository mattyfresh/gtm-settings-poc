package controllers

import (
	"errors"
	"fmt"
	"go-bot/services"
	"os"
	"regexp"
	"strings"

	"github.com/nlopes/slack"
	gtm "google.golang.org/api/tagmanager/v2"
)

var realTimeMessenger *slack.RTM
var gtmService *gtm.Service

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

func validateVariable(variable *gtm.Variable) (errors []error) {
	// data layer vars
	dataLayerPrefix := "Data Layer - "
	if variable.Type == "v" && !strings.HasPrefix(variable.Name, dataLayerPrefix) {
		errMsg := fmt.Errorf("_Variable `%s` failed validation, all Data Layer variables must start with the prefix: `%s`_", variable.Name, dataLayerPrefix)
		errors = append(errors, errMsg)
	}

	return errors
}

func validateTag(tag *gtm.Tag) (errors []error) {
	// @TODO
	return errors
}

func validateTrigger(trigger *gtm.Trigger) (errors []error) {
	// custom events
	customEventPrefix := "Custom Event -"
	if trigger.Type == "customEvent" && !strings.HasPrefix(trigger.Name, customEventPrefix) {
		errMsg := fmt.Errorf("_Trigger `%s` failed validation, all Custom Event variables must start with the prefix: `%s`_", trigger.Name, customEventPrefix)
		errors = append(errors, errMsg)
	}
	return errors
}

// GtmHandler controller for all GTM related events
func GtmHandler(msg *slack.MessageEvent, rtm *slack.RTM) {
	// expose this globally to the module so we don't have to pass it around as an arg
	realTimeMessenger = rtm

	// get args for gtm command
	pattern := regexp.MustCompile(`gtm\s+(?P<command_type>\w*)\s+(?P<container_name>\S*)`)
	match := pattern.FindStringSubmatch(msg.Text)
	if match == nil {
		sendMessage("gtm requests must be in format `gtm <name_of_command> <container_name>`", msg.Channel)
		return
	}

	captures := make(map[string]string)
	for i, name := range pattern.SubexpNames() {
		// Ignore the whole regexp match and unnamed groups
		if i == 0 || name == "" {
			continue
		}
		captures[name] = match[i]
	}

	commandType, _ := captures["command_type"]
	containerName, _ := captures["container_name"]

	if !(commandType == "publish" || commandType == "validate") {
		sendMessage("There are two commands available, `publish` and `validate`.", msg.Channel)
		return
	}

	// make the GTM api service available
	if initErr := gtmInit(); initErr != nil {
		fmt.Printf("Error initializing GTM api: %s", initErr.Error())
		return
	}

	accountID := os.Getenv("GTM_ACCOUNT_ID")
	accountPath := "accounts/" + accountID

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
	wipMsg := fmt.Sprintf("validating container with ID %s and workspace #%s :ram:", containerID, workspaceID)
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
		if errors := validateVariable(v); errors != nil {
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
		if errors := validateTag(t); errors != nil {
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
		if errors := validateTrigger(t); errors != nil {
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

	sendMessage(":shipit: Validation Succeeded!", msg.Channel)

	// @TODO
	if commandType == "publish" {
		// write current GTM state to file
		sendMessage(":shipit: Publishing...", msg.Channel)

		// create PR in GitHub
		// @TODO
	}
}
