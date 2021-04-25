package forms

import (
	"github.com/wochinge/go-rasa-sdk/rasa"
	"github.com/wochinge/go-rasa-sdk/rasa/responses"
)

// SlotValidator can be used to validate candidates before filling a slot with them.
type SlotValidator interface {
	// IsValid checks if a slot candidate is valid and can be used to fill a slot.
	// Returns the validated slot value and `true` if the value is valid.
	IsValid(value interface{}, domain *rasa.Domain, tracker *rasa.Tracker,
		dispatcher responses.ResponseDispatcher) (validatedValue interface{}, isValid bool)
}

// DefaultValidator is a validator which only checks that the value is not `nil` before filling a slot.
type DefaultValidator struct{}

func (v *DefaultValidator) IsValid(value interface{}, _ *rasa.Domain, _ *rasa.Tracker,
	_ responses.ResponseDispatcher) (interface{}, bool) {
	return value, value != nil
}

// MultiValidator represents multiple `SlotValidator`s as it would be one.
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

func (form *Form) validatorFor(slotName string, tracker *rasa.Tracker) SlotValidator {
	validators := form.Validators[slotName]

	if tracker.NoFormValidation() || validators == nil {
		return &DefaultValidator{}
	}

	return &MultiValidator{validators}
}
