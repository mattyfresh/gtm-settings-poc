package controllers

import (
	"fmt"
	"strings"

	gtm "google.golang.org/api/tagmanager/v2"
)

func customEvents(trigger *gtm.Trigger) error {
	customEventPrefix := "Custom Event -"
	if trigger.Type == "customEvent" && !strings.HasPrefix(trigger.Name, customEventPrefix) {
		errMsg := fmt.Sprintf("Trigger `%s` failed validation, all Custom Event variables must start with the prefix: `%s`", trigger.Name, customEventPrefix)
		return validationError(errMsg)
	}
	return nil
}

var triggerValidators = []func(trigger *gtm.Trigger) error{customEvents}

// ValidateTrigger takes a gtm trigger and runs all of the relevant validation functions
func ValidateTrigger(trigger *gtm.Trigger) (errors []error) {
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
	for _, f := range variableValidators {
		err := f(variable)
		if err != nil {
			errors = append(errors, err)
		}
	}
	return errors
}

func lowerCaseEventCategory(tag *gtm.Tag) error {
	eventCategory := ""
	for _, param := range tag.Parameter {
		if param.Key == "eventCategory" {
			eventCategory = param.Value
			break
		}
	}
	lowerCaseCategory := strings.ToLower(eventCategory)
	if lowerCaseCategory != eventCategory {
		errMsg := fmt.Sprintf("Tag `%s` failed validation, category `%s` must be all lower case.  It should probably be: `%s`", tag.Name, eventCategory, lowerCaseCategory)
		return validationError(errMsg)
	}
	return nil
}

var tagValidators = []func(variable *gtm.Tag) error{lowerCaseEventCategory}

// ValidateTag takes a gtm tag and runs all of the relevant validation functions
func ValidateTag(tag *gtm.Tag) (errors []error) {
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
