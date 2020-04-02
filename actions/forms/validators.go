package forms

import (
	"github.com/wochinge/go-rasa-sdk/rasa"
	"github.com/wochinge/go-rasa-sdk/rasa/responses"
)

type SlotValidator interface {
	IsValid(value interface{}, domain *rasa.Domain, tracker *rasa.Tracker,
		dispatcher responses.ResponseDispatcher) (interface{}, bool)
}

type DefaultValidator struct{}

func (v *DefaultValidator) IsValid(value interface{}, _ *rasa.Domain, _ *rasa.Tracker,
	_ responses.ResponseDispatcher) (interface{}, bool) {
	return value, value != nil
}

type MultiValidator struct {
	validators []SlotValidator
}

func (v *MultiValidator) IsValid(value interface{}, domain *rasa.Domain, tracker *rasa.Tracker,
	dispatcher responses.ResponseDispatcher) (interface{}, bool) {
	for _, validator := range v.validators {
		if validated, valid := validator.IsValid(value, domain, tracker, dispatcher); valid {
			return validated, true
		}
	}

	return nil, false
}
