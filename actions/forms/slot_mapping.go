package forms

import (
	"github.com/wochinge/go-rasa-sdk/rasa"
)

// SlotMapping specifies which information is used to fill a form slot.
type SlotMapping struct {
	// ApplyToFirstRunOnly specifies if the mapping should only apply when the form is activated.
	ApplyToFirstRunOnly bool
	// FromText fills the requested slot with the user message.
	FromText bool
	// FromEntity fills the slot with an entity of a given name.
	FromEntity string
	// Intents which the latest message has to have that this mapping applies.
	// If `nil` and `ExcludedIntents` is also `nil`, all intents are valid.
	Intents []string
	// ExcludedIntents specifies intents which this mapping should not apply to.
	ExcludedIntents []string
	// Value can be used to hard code the slot value in case the mapping applies.
	Value interface{}
}

func defaultSlotMapping(slotName string) []SlotMapping {
	return []SlotMapping{{FromEntity: slotName}}
}

func (mapping SlotMapping) apply(form *Form, tracker *rasa.Tracker) (interface{}, bool) {
	latestMessage := &tracker.LatestMessage
	if !mapping.allows(latestMessage.Intent.Name) || (mapping.ApplyToFirstRunOnly && form.wasAlreadyActive(tracker)) {
		return nil, false
	}

	if entity, found := latestMessage.EntityFor(mapping.FromEntity); found && mapping.FromEntity != "" {
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
