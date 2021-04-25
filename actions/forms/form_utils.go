// Package forms implements helpers for Rasa Open Source forms (https://rasa.com/docs/rasa/forms/).
package forms

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/wochinge/go-rasa-sdk/logging"
	"github.com/wochinge/go-rasa-sdk/rasa"
	"github.com/wochinge/go-rasa-sdk/rasa/events"
	"github.com/wochinge/go-rasa-sdk/rasa/responses"
)

const requestedSlot = "requested_slot"

type FormValidationAction struct {
	// FormName is the name of the form.
	FormName string
	// Validators specify functions to validate slot candidates.
	Validators map[string]SlotValidator
	// Extractors specify functions to extract slot candidates.
	Extractors        map[string]SlotExtractor
	NextSlotRequester NextSlotRequester
}

// SlotValidator can be used to validate candidates before filling a slot with them.
type SlotValidator interface {
	// IsValid checks if a slot candidate is valid and can be used to fill a slot.
	// Returns the validated slot value and `true` if the value is valid.
	IsValid(value interface{}, domain *rasa.Domain, tracker *rasa.Tracker,
		dispatcher responses.ResponseDispatcher) (validatedValue interface{}, isValid bool)
}

// SlotExtractor can be used to extract custom slots.
type SlotExtractor interface {
	// Extract can be used to extract custom slots.
	Extract(domain *rasa.Domain, tracker *rasa.Tracker,
		dispatcher responses.ResponseDispatcher) (extractedValue interface{}, valueFound bool)
}

// NextSlotRequester can be used to TODO.
type NextSlotRequester interface {
	// NextSlot be used to TODO.
	NextSlot(domain *rasa.Domain, tracker *rasa.Tracker,
		dispatcher responses.ResponseDispatcher) (nextSlot string, shouldRequestNextSlot bool)
}

// Run is executed whenever Rasa Open Source sends a request to validate this form.
func (action *FormValidationAction) Run(tracker *rasa.Tracker, domain *rasa.Domain,
	dispatcher responses.ResponseDispatcher) []events.Event {
	tracker.Init()
	log.WithFields(
		log.Fields{logging.FormNameKey: action.FormName, logging.FormValidationKey: tracker.ActiveLoop.Validate}).Debug(
		"Validating form.")

	var newEvents []events.Event

	for slotName, extractor := range action.Extractors {
		if extractedValue, valueFound := extractor.Extract(domain, tracker, dispatcher); valueFound {
			tracker.Slots[slotName] = extractedValue
			tracker.Events = append(tracker.Events, &events.SlotSet{Name: slotName, Value: extractedValue})
		}
	}

	slotsToValidate := tracker.SlotsToValidate()
	for slotName, slotValue := range slotsToValidate {
		if validator, ok := action.Validators[slotName]; ok {
			if validatedValue, isValid := validator.IsValid(slotValue, domain, tracker, dispatcher); isValid {
				newEvents = append(newEvents, &events.SlotSet{Name: slotName, Value: validatedValue})
			} else {
				newEvents = append(newEvents, &events.SlotSet{Name: slotName, Value: nil})
			}
		} else {
			// no validator function provided
			newEvents = append(newEvents, &events.SlotSet{Name: slotName, Value: slotValue})
		}
	}

	if action.NextSlotRequester != nil {
		if nextSlot, shouldRequestNextSlot := action.NextSlotRequester.NextSlot(
			domain, tracker, dispatcher); shouldRequestNextSlot {
			newEvents = append(newEvents, &events.SlotSet{Name: requestedSlot, Value: nextSlot})
		} else {
			newEvents = append(newEvents, &events.SlotSet{Name: requestedSlot, Value: nil})
		}
	}

	return newEvents
}

func (action *FormValidationAction) Name() string { return fmt.Sprintf("validate_%v", action.FormName) }
