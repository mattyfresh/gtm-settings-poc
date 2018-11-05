package main

import (
	"fmt"
	"reflect"
	"strings"

	gtm "google.golang.org/api/tagmanager/v2"
)

// ensure custom events always have the proper prefix
func customEvents(trigger *gtm.Trigger) error {
	customEventPrefix := "Custom Event - "
	if trigger.Type == "customEvent" && !strings.HasPrefix(trigger.Name, customEventPrefix) {
		errMsg := fmt.Sprintf(
			"Trigger `%s` failed validation, all Custom Event trigger names must start with the prefix: `%s`",
			trigger.Name,
			customEventPrefix,
		)
		return validationError(errMsg)
	}
	return nil
}

// ensure CSS selectors always begin with `.js_`
func enforceSelectors(trigger *gtm.Trigger) error {
	if len(trigger.Filter) < 1 {
		return nil
	}

	const jsPrefix = ".js_"
	cssSelector := ""

	// loop over each filter
	for _, filter := range trigger.Filter {
		if filter.Type == "cssSelector" {

			// each filter param
			for _, p := range filter.Parameter {

				// unfortunately this is the very non-descriptive name the GTM api gives us
				// this will be the value, or the actual CSS selector paramater
				if p.Key == "arg1" {
					cssSelector = p.Value
					hasFormattingError := false

					// split the CSS selector string
					selectors := strings.Split(cssSelector, " ")
					for _, selector := range selectors {

						// if any selector does not have the prefix `.js_` or aren't `*`, flag an error and bail
						if !strings.HasPrefix(selector, jsPrefix) && selector != "*" {
							hasFormattingError = true
							break
						}
					}

					if hasFormattingError {
						errMsg := fmt.Sprintf(
							"Trigger `%s` failed validation, Make sure your css selectors begin with `%s`. Right now they look like: `%s`",
							trigger.Name,
							jsPrefix,
							cssSelector,
						)
						return validationError(errMsg)
					}
				}
			}
		}
	}
	return nil
}

var triggerValidators = []func(trigger *gtm.Trigger) error{customEvents, enforceSelectors}

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
		errMsg := fmt.Sprintf(
			"Variable `%s` failed validation, all Data Layer variables must start with the prefix: `%s`",
			variable.Name,
			dataLayerPrefix,
		)
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

func formattedEventParams(tag *gtm.Tag) error {
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

	properlyFormattedMap := make(map[string]string)
	for key, val := range rawEventMap {
		// lower case and replace all `-` with " "
		formattedVal := strings.ToLower(val)
		properlyFormattedMap[key] = strings.Replace(formattedVal, "-", " ", -1)
	}

	if !reflect.DeepEqual(rawEventMap, properlyFormattedMap) {
		errMsg := fmt.Sprintf(
			"Tag `%s` failed validation, ensure values are lower case and formatted correctly! `%v`.  It should probably be: `%v`",
			tag.Name,
			rawEventMap,
			properlyFormattedMap,
		)
		return validationError(errMsg)
	}
	return nil
}

var tagValidators = []func(tag *gtm.Tag) error{formattedEventParams}

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
