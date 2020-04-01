package forms

import (
	"github.com/wochinge/go-rasa-sdk/rasa"
	"github.com/wochinge/go-rasa-sdk/rasa/responses"
)

type SlotValidator interface {
	IsValid(value interface{}, domain *rasa.Domain, tracker *rasa.Tracker,
		dispatcher responses.ResponseDispatcher) bool
}

type DefaultValidator struct{}

func (v *DefaultValidator) IsValid(value interface{}, _ *rasa.Domain, _ *rasa.Tracker,
	_ responses.ResponseDispatcher) bool {
	return value != nil
}

type MultiValidator struct {
	validators []SlotValidator
}

func (v *MultiValidator) IsValid(value interface{}, domain *rasa.Domain, tracker *rasa.Tracker,
	dispatcher responses.ResponseDispatcher) bool {
	for _, validator := range v.validators {
		if ! validator.IsValid(value, domain, tracker, dispatcher) {
			return false
		}
	}
	return true
}
