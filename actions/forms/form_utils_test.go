package forms

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wochinge/go-rasa-sdk/rasa"
	"github.com/wochinge/go-rasa-sdk/rasa/events"
	"github.com/wochinge/go-rasa-sdk/rasa/responses"
)

func TestFormRun(t *testing.T) {
	validators, extractors := make(map[string]SlotValidator), make(map[string]SlotValidator)
	formValidator := FormValidationAction{
		"test_form", validators, extractors,
	}

	tracker := rasa.Tracker{ActiveLoop: rasa.ActiveLoop{Name: testFormName},
		Slots: map[string]interface{}{requestedSlot: nil}}

	newEvents := formValidator.Run(&tracker, &rasa.Domain{}, responses.NewDispatcher())

	assert.Empty(t, newEvents)
	assert.Equal(t, "validate_test_form", formValidator.Name())
}

func TestFormValidSlot(t *testing.T) {
	slotName := "color"
	validators := map[string]SlotValidator{slotName: &ExactMatchValidator{
		"green",
	}}
	extractors := make(map[string]SlotValidator)
	formName := "test_form"
	formValidator := FormValidationAction{
		formName, validators, extractors,
	}
	tracker := rasa.Tracker{ActiveLoop: rasa.ActiveLoop{Name: testFormName},
		Slots: map[string]interface{}{requestedSlot: nil, slotName: "green"},
		Events: []events.Event{
			&events.SlotSet{Name: "bla", Value: 5},
			&events.Action{Name: formName},
			&events.SlotSet{Name: slotName, Value: "green"},
			&events.SlotSet{Name: "another", Value: "bla"},
		}}

	newEvents := formValidator.Run(&tracker, &rasa.Domain{}, responses.NewDispatcher())

	expected := []events.Event{
		&events.SlotSet{Name: slotName, Value: "green"},
	}
	assert.ElementsMatch(t, expected, newEvents)
}


func TestFormInValidSlot(t *testing.T) {
	slotName := "color"
	validators := map[string]SlotValidator{slotName: &ExactMatchValidator{
		"green",
	}}
	extractors := make(map[string]SlotValidator)
	formName := "test_form"
	formValidator := FormValidationAction{
		formName, validators, extractors,
	}
	tracker := rasa.Tracker{ActiveLoop: rasa.ActiveLoop{Name: testFormName},
		Slots: map[string]interface{}{requestedSlot: nil, slotName: "green"},
		Events: []events.Event{
			&events.SlotSet{Name: "bla", Value: 5},
			&events.Action{Name: formName},
			&events.SlotSet{Name: slotName, Value: "blue"},
			&events.SlotSet{Name: "another", Value: "bla"},
		}}

	newEvents := formValidator.Run(&tracker, &rasa.Domain{}, responses.NewDispatcher())

	expected := []events.Event{
		&events.SlotSet{Name: slotName, Value: nil},
	}
	assert.ElementsMatch(t, expected, newEvents)
}
