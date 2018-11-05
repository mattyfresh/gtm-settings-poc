package main

import (
	"fmt"
	"reflect"
	"strings"

	gtm "google.golang.org/api/tagmanager/v2"
)

func customEvents(trigger *gtm.Trigger) error {
	customEventPrefix := "Custom Event - "
	if trigger.Type == "customEvent" && !strings.HasPrefix(trigger.Name, customEventPrefix) {
		errMsg := fmt.Sprintf("Trigger `%s` failed validation, all Custom Event trigger names must start with the prefix: `%s`", trigger.Name, customEventPrefix)
		return validationError(errMsg)
	}
	return nil
}

var triggerValidators = []func(trigger *gtm.Trigger) error{customEvents}

// ValidateTrigger takes a gtm trigger and runs all of the relevant validation functions
func ValidateTrigger(trigger *gtm.Trigger) (errors []error) {
	// remove fields we don't care about diffing
	trigger.Path = ""
	trigger.Fingerprint = ""
	trigger.TagManagerUrl = ""
	trigger.WorkspaceId = ""

	for _, f := range triggerValidators {
		err := f(trigger)
		if err != nil {
			errors = append(errors, err)
		}
	}
	return errors
}

func dataLayerVariables(variable *gtm.Variable) error {
	dataLayerPrefix := "Data Layer - "
	if variable.Type == "v" && !strings.HasPrefix(variable.Name, dataLayerPrefix) {
		errMsg := fmt.Sprintf("Variable `%s` failed validation, all Data Layer variables must start with the prefix: `%s`", variable.Name, dataLayerPrefix)
		return validationError(errMsg)
	}
	return nil
}

var variableValidators = []func(variable *gtm.Variable) error{dataLayerVariables}

// ValidateVariable takes a gtm variable and runs all of the relevant validation functions
func ValidateVariable(variable *gtm.Variable) (errors []error) {
	// remove fields we don't care about diffing
	variable.Path = ""
	variable.Fingerprint = ""
	variable.TagManagerUrl = ""
	variable.WorkspaceId = ""

	for _, f := range variableValidators {
		err := f(variable)
		if err != nil {
			errors = append(errors, err)
		}
	}
	return errors
}

func lowerCaseEventCategory(tag *gtm.Tag) error {
	rawEventMap := make(map[string]string)
	for _, param := range tag.Parameter {
		// ignore properties that have GTM built in variables,
		// these can be uppercase and are not editable :/
		if strings.Contains(param.Value, "{{") {
			continue
		}

		if param.Key == "eventAction" {
			rawEventMap["action"] = param.Value
		}
		if param.Key == "eventCategory" {
			rawEventMap["category"] = param.Value
		}
		if param.Key == "eventLabel" {
			rawEventMap["label"] = param.Value
		}
	}

	lowerCaseMap := make(map[string]string)
	for key, val := range rawEventMap {
		lowerCaseMap[key] = strings.ToLower(val)
	}

	if !reflect.DeepEqual(rawEventMap, lowerCaseMap) {
		errMsg := fmt.Sprintf(
			"Tag `%s` failed validation, ensure all of these values are lower case! `%v`.  It should probably be: `%v`",
			tag.Name,
			rawEventMap,
			lowerCaseMap,
		)
		return validationError(errMsg)
	}
	return nil
}

var tagValidators = []func(tag *gtm.Tag) error{lowerCaseEventCategory}

// ValidateTag takes a gtm tag and runs all of the relevant validation functions
func ValidateTag(tag *gtm.Tag) (errors []error) {
	// remove fields we don't care about diffing
	tag.Path = ""
	tag.Fingerprint = ""
	tag.TagManagerUrl = ""
	tag.WorkspaceId = ""

	for _, f := range tagValidators {
		err := f(tag)
		if err != nil {
			errors = append(errors, err)
		}
	}
	return errors
}

func validationError(errMsg string) error {
	return fmt.Errorf(":x: _%s_", errMsg)
}
