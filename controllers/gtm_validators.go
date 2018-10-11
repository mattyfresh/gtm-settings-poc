package controllers

import (
	"fmt"
	"strings"

	gtm "google.golang.org/api/tagmanager/v2"
)

var triggerValidators = []func(trigger *gtm.Trigger) error{customEvents}

func customEvents(trigger *gtm.Trigger) error {
	customEventPrefix := "Custom Event -"
	if trigger.Type == "customEvent" && !strings.HasPrefix(trigger.Name, customEventPrefix) {
		return fmt.Errorf("_Trigger `%s` failed validation, all Custom Event variables must start with the prefix: `%s`_", trigger.Name, customEventPrefix)
	}
	return nil
}

func validateTrigger(trigger *gtm.Trigger) (errors []error) {
	for _, f := range triggerValidators {
		err := f(trigger)
		if err != nil {
			errors = append(errors, err)
		}
	}
	return errors
}

var variableValidators = []func(variable *gtm.Variable) error{dataLayerVariables}

func dataLayerVariables(variable *gtm.Variable) error {
	dataLayerPrefix := "Data Layer - "
	if variable.Type == "v" && !strings.HasPrefix(variable.Name, dataLayerPrefix) {
		return fmt.Errorf("_Variable `%s` failed validation, all Data Layer variables must start with the prefix: `%s`_", variable.Name, dataLayerPrefix)
	}
	return nil
}

func validateVariable(variable *gtm.Variable) (errors []error) {
	for _, f := range variableValidators {
		err := f(variable)
		if err != nil {
			errors = append(errors, err)
		}
	}
	return errors
}

var tagValidators = []func(variable *gtm.Tag) error{}

func validateTag(tag *gtm.Tag) (errors []error) {
	for _, f := range tagValidators {
		err := f(tag)
		if err != nil {
			errors = append(errors, err)
		}
	}
	return errors
}
