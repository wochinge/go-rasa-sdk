package forms

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/wochinge/go-rasa-sdk/logging"
	"github.com/wochinge/go-rasa-sdk/rasa"
	"github.com/wochinge/go-rasa-sdk/rasa/events"
	"github.com/wochinge/go-rasa-sdk/rasa/responses"
)

type FormValidationAction struct {
	// FormName is the name of the form.
	FormName string
	// Validators specify functions to validate slot candidates.
	Validators map[string]SlotValidator
	// Extractors specify functions to extract slot candidates.
	Extractors map[string]SlotValidator
}

// Run is executed whenever Rasa Open Source sends a request to validate this form.
func (action *FormValidationAction) Run(tracker *rasa.Tracker, domain *rasa.Domain,
	dispatcher responses.ResponseDispatcher) []events.Event {
	tracker.Init()
	log.WithFields(
		log.Fields{logging.FormNameKey: action.FormName, logging.FormValidationKey: tracker.ActiveLoop.Validate}).Debug(
		"Validating form.")

	var newEvents []events.Event

	slotsToValidate := tracker.SlotsToValidate()
	for slotName, slotValue := range *slotsToValidate {
		if validator, ok := action.Validators[slotName]; ok {
			if validatedValue, isValid := validator.IsValid(slotValue, domain, tracker, dispatcher); isValid {
				newEvents = append(newEvents, &events.SlotSet{Name: slotName, Value: validatedValue})
			} else {
				newEvents = append(newEvents, &events.SlotSet{Name: slotName, Value: nil})
			}
		}
	}

	return newEvents
}

func (action *FormValidationAction) Name() string { return fmt.Sprintf("validate_%v", action.FormName) }
