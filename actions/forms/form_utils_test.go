package forms

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wochinge/go-rasa-sdk/rasa"
	"github.com/wochinge/go-rasa-sdk/rasa/events"
	"github.com/wochinge/go-rasa-sdk/rasa/responses"
)

const testFormName = "test_form"

type ExactMatchValidator struct {
	toMatch string
}

func (v *ExactMatchValidator) IsValid(value interface{}, _ *rasa.Domain, _ *rasa.Tracker,
	_ responses.ResponseDispatcher) (interface{}, bool) {
	return value, value == v.toMatch
}

func TestFormRun(t *testing.T) {
	validators, extractors := make(map[string]SlotValidator), make(map[string]SlotExtractor)
	formValidator := FormValidationAction{
		"test_form", validators, extractors, nil,
	}

	tracker := rasa.Tracker{ActiveLoop: rasa.ActiveLoop{Name: testFormName},
		Slots: map[string]interface{}{requestedSlot: nil}}

	newEvents := formValidator.Run(&tracker, &rasa.Domain{}, responses.NewDispatcher())

	assert.Empty(t, newEvents)
	assert.NotNil(t, newEvents)
	assert.Equal(t, "validate_test_form", formValidator.Name())
}

func TestFormValidSlot(t *testing.T) {
	slotName := "color"
	validators := map[string]SlotValidator{slotName: &ExactMatchValidator{
		"green",
	}}
	extractors := make(map[string]SlotExtractor)
	formName := "test_form"
	formValidator := FormValidationAction{
		formName, validators, extractors, nil,
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
		&events.SlotSet{Name: "another", Value: "bla"},
	}
	assert.ElementsMatch(t, expected, newEvents)
}

func TestFormInValidSlot(t *testing.T) {
	slotName := "color"
	validators := map[string]SlotValidator{slotName: &ExactMatchValidator{
		"green",
	}}
	extractors := make(map[string]SlotExtractor)
	formName := "test_form"
	formValidator := FormValidationAction{
		formName, validators, extractors, nil,
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
		&events.SlotSet{Name: "another", Value: "bla"},
	}
	assert.ElementsMatch(t, expected, newEvents)
}

type EntityExtractor struct {
	entityToMatch string
}

func (v *EntityExtractor) Extract(_ *rasa.Domain, tracker *rasa.Tracker,
	_ responses.ResponseDispatcher) (interface{}, bool) {
	for _, entity := range tracker.LatestMessage.Entities {
		if entity.Name == v.entityToMatch {
			return entity.Value, true
		}
	}

	return nil, false
}

func TestExtractCustomSlot(t *testing.T) {
	slotName := "color"
	validators := make(map[string]SlotValidator)
	extractors := map[string]SlotExtractor{slotName: &EntityExtractor{
		"color",
	}}
	formName := "test_form"
	formValidator := FormValidationAction{
		formName, validators, extractors, nil,
	}
	tracker := rasa.Tracker{ActiveLoop: rasa.ActiveLoop{Name: testFormName},
		Slots: map[string]interface{}{requestedSlot: nil, slotName: "green"},
		Events: []events.Event{
			&events.SlotSet{Name: "bla", Value: 5},
			&events.Action{Name: formName},
		}}
	tracker.LatestMessage.Entities = []events.Entity{{Name: slotName, Value: "green"}}

	newEvents := formValidator.Run(&tracker, &rasa.Domain{}, responses.NewDispatcher())

	expected := []events.Event{
		&events.SlotSet{Name: slotName, Value: "green"},
	}
	assert.ElementsMatch(t, expected, newEvents)
}

func TestExtractCustomSlotIfNotFound(t *testing.T) {
	slotName := "color"
	validators := make(map[string]SlotValidator)
	extractors := map[string]SlotExtractor{slotName: &EntityExtractor{
		"color",
	}}
	formName := "test_form"
	formValidator := FormValidationAction{
		formName, validators, extractors, nil,
	}
	tracker := rasa.Tracker{ActiveLoop: rasa.ActiveLoop{Name: testFormName},
		Slots: map[string]interface{}{requestedSlot: nil, slotName: "green"},
		Events: []events.Event{
			&events.SlotSet{Name: "bla", Value: 5},
			&events.Action{Name: formName},
		}}

	newEvents := formValidator.Run(&tracker, &rasa.Domain{}, responses.NewDispatcher())

	assert.Empty(t, newEvents)
}

func TestExtractCustomSlotAndValidate(t *testing.T) {
	slotName := "color"
	validators := map[string]SlotValidator{slotName: &ExactMatchValidator{
		"green",
	}}
	extractors := map[string]SlotExtractor{slotName: &EntityExtractor{
		"color",
	}}
	formName := "test_form"
	formValidator := FormValidationAction{
		formName, validators, extractors, nil,
	}
	tracker := rasa.Tracker{ActiveLoop: rasa.ActiveLoop{Name: testFormName},
		Slots: map[string]interface{}{requestedSlot: nil, slotName: "green"},
		Events: []events.Event{
			&events.SlotSet{Name: "bla", Value: 5},
			&events.Action{Name: formName},
		}}
	tracker.LatestMessage.Entities = []events.Entity{{Name: slotName, Value: "green"}}

	newEvents := formValidator.Run(&tracker, &rasa.Domain{}, responses.NewDispatcher())

	expected := []events.Event{
		&events.SlotSet{Name: slotName, Value: "green"},
	}
	assert.ElementsMatch(t, expected, newEvents)
}

type ConstantSlotRequester struct {
	slotToRequest interface{}
}

func (v *ConstantSlotRequester) NextSlot(_ *rasa.Domain, tracker *rasa.Tracker,
	_ responses.ResponseDispatcher) (string, bool) {
	if v.slotToRequest == nil {
		return "", false
	}

	return v.slotToRequest.(string), true
}

func TestRequestNextSlot(t *testing.T) {
	slotName := "color"

	validators, extractors := make(map[string]SlotValidator), make(map[string]SlotExtractor)
	formName := "test_form"
	formValidator := FormValidationAction{
		formName, validators, extractors, &ConstantSlotRequester{"color"},
	}
	tracker := rasa.Tracker{ActiveLoop: rasa.ActiveLoop{Name: testFormName},
		Slots: map[string]interface{}{requestedSlot: nil, slotName: "green"},
		Events: []events.Event{
			&events.SlotSet{Name: "bla", Value: 5},
			&events.Action{Name: formName},
		}}
	tracker.LatestMessage.Entities = []events.Entity{{Name: slotName, Value: "green"}}

	newEvents := formValidator.Run(&tracker, &rasa.Domain{}, responses.NewDispatcher())

	expected := []events.Event{
		&events.SlotSet{Name: requestedSlot, Value: "color"},
	}
	assert.ElementsMatch(t, expected, newEvents)
}

func TestRequestNoNextSlot(t *testing.T) {
	slotName := "color"

	validators, extractors := make(map[string]SlotValidator), make(map[string]SlotExtractor)
	formName := "test_form"
	formValidator := FormValidationAction{
		formName, validators, extractors, &ConstantSlotRequester{nil},
	}
	tracker := rasa.Tracker{ActiveLoop: rasa.ActiveLoop{Name: testFormName},
		Slots: map[string]interface{}{requestedSlot: nil, slotName: "green"},
		Events: []events.Event{
			&events.SlotSet{Name: "bla", Value: 5},
			&events.Action{Name: formName},
		}}
	tracker.LatestMessage.Entities = []events.Entity{{Name: slotName, Value: "green"}}

	newEvents := formValidator.Run(&tracker, &rasa.Domain{}, responses.NewDispatcher())

	expected := []events.Event{
		&events.SlotSet{Name: requestedSlot, Value: nil},
	}
	assert.ElementsMatch(t, expected, newEvents)
}
