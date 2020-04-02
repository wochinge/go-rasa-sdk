package forms

import (
	"github.com/wochinge/go-rasa-sdk/rasa"
)

type SlotMapping struct {
	ApplyToFirstRunOnly bool
	FromText            bool
	Entity              string
	Intents             []string
	ExcludedIntents     []string
	Value               interface{}
}

func defaultSlotMapping(slotName string) []SlotMapping {
	return []SlotMapping{{Entity: slotName}}
}

func (mapping SlotMapping) apply(form *Form, tracker *rasa.Tracker) (interface{}, bool) {
	latestMessage := &tracker.LatestMessage
	if !mapping.allows(latestMessage.Intent.Name) || (mapping.ApplyToFirstRunOnly && form.wasAlreadyActive(tracker)) {
		return nil, false
	}

	if entity, found := latestMessage.EntityFor(mapping.Entity); found && mapping.Entity != "" {
		return entity, true
	}

	if mapping.FromText {
		return latestMessage.Text, true
	}

	if mapping.Value != nil {
		return mapping.Value, true
	}

	return nil, false
}

func (mapping SlotMapping) allows(intentName string) bool {
	if mapping.Intents == nil && mapping.ExcludedIntents == nil {
		return true
	}

	for _, excludedIntent := range mapping.ExcludedIntents {
		if excludedIntent == intentName {
			return false
		}
	}

	if mapping.Intents == nil {
		return true
	}

	for _, allowedIntent := range mapping.Intents {
		if allowedIntent == intentName {
			return true
		}
	}

	return false
}
