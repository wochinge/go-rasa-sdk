package forms

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wochinge/go-rasa-sdk/rasa"
	"github.com/wochinge/go-rasa-sdk/rasa/events"
	"github.com/wochinge/go-rasa-sdk/rasa/responses"
)

const testFormName = "test-form"

type ExactMatchValidator struct {
	toMatch string
}

func (v *ExactMatchValidator) IsValid(value interface{}, _ *rasa.Domain, _ *rasa.Tracker,
	_ responses.ResponseDispatcher) (interface{}, bool) {
	return value, value == v.toMatch
}

func TestActivateFormIfActive(t *testing.T) {
	nextSlot, value := "next slot", "bla"

	tracker := rasa.Tracker{ActiveLoop: rasa.ActiveLoop{Name: testFormName},
		Slots: map[string]interface{}{requestedSlot: nextSlot}}
	tracker.LatestMessage.Entities = []events.Entity{{Name: nextSlot, Value: value}}

	testForm := Form{FormName: testFormName, RequiredSlots: []string{nextSlot}}

	newEvents := testForm.Run(&tracker, &rasa.Domain{}, responses.NewDispatcher())

	expectedEvents := []events.Event{
		&events.SlotSet{Name: nextSlot, Value: value},
		&events.Form{},
		&events.SlotSet{Name: requestedSlot, Value: nil}}

	assert.ElementsMatch(t, expectedEvents, newEvents)
}

func TestActivateFormIfNotActive(t *testing.T) {
	nextSlot, otherSlot := "x", "y"

	tracker := rasa.Tracker{Slots: map[string]interface{}{nextSlot: nil, otherSlot: nil}}

	testForm := Form{FormName: testFormName, RequiredSlots: []string{nextSlot, otherSlot}}

	newEvents := testForm.Run(&tracker, &rasa.Domain{}, responses.NewDispatcher())
	expected := []events.Event{&events.Form{Name: testFormName}, &events.SlotSet{Name: requestedSlot, Value: nextSlot}}

	assert.ElementsMatch(t, expected, newEvents)
}

func TestActivateValidateExistingSlots(t *testing.T) {
	slot1, slot2 := "slot1", "slot2"
	toMatch := "Sara"
	formName := "my cool form"

	requiredSlots := []string{slot1, slot2}
	currentlyFilledSLots := map[string]interface{}{slot1: toMatch, slot2: "tada"}

	tracker := &rasa.Tracker{Slots: currentlyFilledSLots, ActiveLoop: rasa.ActiveLoop{Validate: true},
		LatestActionName: "action_listen"}

	validator := &ExactMatchValidator{toMatch: toMatch}
	testForm := Form{FormName: formName, RequiredSlots: requiredSlots, Validators: map[string][]SlotValidator{
		slot1: {validator}, slot2: {validator}}}

	newEvents := testForm.Run(tracker, &rasa.Domain{}, responses.NewDispatcher())

	expected := []events.Event{&events.Form{Name: formName},
		&events.SlotSet{Name: slot1, Value: toMatch},
		&events.SlotSet{Name: slot2, Value: nil},
		&events.SlotSet{Name: requestedSlot, Value: slot2}}

	assert.ElementsMatch(t, expected, newEvents)
}

func TestActivateFormWithAlreadyFilledSlots(t *testing.T) {
	slot, value := "there-and-valid", "Sara"
	requiredSlots := []string{slot}
	currentlyFilledSLots := map[string]interface{}{slot: value, "other": "not part of the form"}

	tracker := rasa.Tracker{Slots: currentlyFilledSLots}
	testForm := Form{FormName: testFormName, RequiredSlots: requiredSlots}

	newEvents := testForm.Run(&tracker, &rasa.Domain{}, responses.NewDispatcher())
	expected := []events.Event{
		&events.Form{Name: testFormName},
		&events.SlotSet{Name: slot, Value: value},
		&events.Form{Name: ""},
		&events.SlotSet{Name: requestedSlot, Value: nil}}

	assert.ElementsMatch(t, expected, newEvents)
}

func TestFillSlots(t *testing.T) {
	slotName := "this slot should be filled"

	// Prepare tracker
	name, value := "my entity", "test"
	lastMessage := events.ParseData{Entities: []events.Entity{{Name: name, Value: value}}}
	tracker := rasa.Tracker{ActiveLoop: rasa.ActiveLoop{Name: testFormName},
		LatestMessage: lastMessage,
		Slots:         map[string]interface{}{requestedSlot: slotName}}

	// Prepare form
	requiredSlots := []string{slotName}
	slotMapping := SlotMapping{FromEntity: name}
	testForm := Form{FormName: testFormName, RequiredSlots: requiredSlots,
		SlotMappings: map[string][]SlotMapping{slotName: {slotMapping}}}

	newEvents := testForm.Run(&tracker, &rasa.Domain{}, responses.NewDispatcher())

	expected := []events.Event{
		&events.SlotSet{Name: slotName, Value: value}, &events.Form{Name: ""},
		&events.SlotSet{Name: requestedSlot, Value: nil}}
	assert.ElementsMatch(t, expected, newEvents)
}

func TestFillSlotWithoutMapping(t *testing.T) {
	slotName := "this slot should be filled"

	// Prepare tracker
	entityValue := "test"
	lastMessage := events.ParseData{Entities: []events.Entity{{Name: slotName, Value: entityValue}}}
	tracker := rasa.Tracker{ActiveLoop: rasa.ActiveLoop{Name: testFormName},
		LatestMessage: lastMessage,
		Slots:         map[string]interface{}{requestedSlot: slotName}}

	// Prepare form
	requiredSlots := []string{slotName}
	testForm := Form{FormName: testFormName, RequiredSlots: requiredSlots}

	newEvents := testForm.slotCandidates(&tracker)

	expected := []events.SlotSet{{Name: slotName, Value: entityValue}}
	assert.ElementsMatch(t, expected, newEvents)
}

func TestFillOtherSlotsIfEntitiesGiven(t *testing.T) {
	otherSlot, expectedValue := "age", "15"
	requested, expectedText := "FormName", "Hello from the other side"

	lastMessage := events.ParseData{Entities: []events.Entity{{Name: otherSlot, Value: expectedValue}},
		Intent: events.IntentParseResult{Name: "some intent"}, Text: expectedText}
	tracker := rasa.Tracker{LatestMessage: lastMessage, Slots: map[string]interface{}{requestedSlot: requested}}

	mappings := map[string][]SlotMapping{requested: {{FromText: true}},
		otherSlot: {{FromEntity: otherSlot, Value: expectedValue}}}
	testForm := Form{FormName: "bla", RequiredSlots: []string{otherSlot, requested}, SlotMappings: mappings}

	newEvents := testForm.slotCandidates(&tracker)

	expectedEvents := []events.SlotSet{{Name: otherSlot, Value: expectedValue}, {Name: requested, Value: expectedText}}

	assert.ElementsMatch(t, expectedEvents, newEvents)
}

func TestDefaultValidation(t *testing.T) {
	invalidSlot := "invalid!!!!!!"
	candidates := []events.SlotSet{{Name: invalidSlot, Value: nil}, {Name: "valid", Value: "so valid"}}

	form := Form{}
	validSlots := form.validatedSlots(candidates, nil, rasa.EmptyTracker(), nil)
	expected := []events.Event{&events.SlotSet{Name: invalidSlot}, &candidates[1]}
	assert.ElementsMatch(t, expected, validSlots)
}

func TestValidation(t *testing.T) {
	validSlot, invalidSlot := "valid", "also invalid"
	toMatch := "exact match!!"
	candidates := []events.SlotSet{{Name: invalidSlot, Value: nil}, {Name: validSlot, Value: toMatch}}

	form := Form{Validators: map[string][]SlotValidator{invalidSlot: {&ExactMatchValidator{}},
		validSlot: {&ExactMatchValidator{toMatch: toMatch}}}}
	validSlots := form.validatedSlots(candidates, nil, rasa.EmptyTracker(), nil)
	expected := []events.Event{&events.SlotSet{Name: invalidSlot}, &candidates[1]}
	assert.ElementsMatch(t, expected, validSlots)
}

func TestValidationDisabled(t *testing.T) {
	validSlot, validValue, invalidSlot, invalidValue := "valid", "valid value", "invalid", "invalid value"
	toMatch := validValue

	tracker := &rasa.Tracker{ActiveLoop: rasa.ActiveLoop{Validate: false}, Slots: map[string]interface{}{}}

	form := Form{Validators: map[string][]SlotValidator{invalidSlot: {&ExactMatchValidator{}},
		validSlot: {&ExactMatchValidator{toMatch: toMatch}}}}

	candidates := []events.SlotSet{{Name: invalidSlot, Value: invalidValue}, {Name: validSlot, Value: validValue}}
	validSlots := form.validatedSlots(candidates, nil, tracker, nil)

	assert.ElementsMatch(t, toEventInterface(candidates), validSlots)
}

func TestSubmit(t *testing.T) {
	submitEvents := []events.Event{&events.Restarted{}, &events.AllSlotsReset{}}
	onSubmit := func(_ *rasa.Tracker, _ *rasa.Domain, _ responses.ResponseDispatcher) []events.Event {
		return submitEvents
	}

	requiredSlot, value := "age", "bla"
	form := Form{FormName: testFormName, OnSubmit: onSubmit, RequiredSlots: []string{requiredSlot}}

	tracker := &rasa.Tracker{ActiveLoop: rasa.ActiveLoop{Name: form.FormName},
		Slots: map[string]interface{}{requestedSlot: requiredSlot}}
	tracker.LatestMessage.Entities = []events.Entity{{Name: requiredSlot, Value: value}}

	newEvents := form.Run(tracker, nil, nil)

	expected := []events.Event{
		&events.SlotSet{Name: requiredSlot, Value: value}}
	expected = append(expected, submitEvents...)
	deactivationEvents := []events.Event{&events.Form{Name: ""}, &events.SlotSet{Name: requestedSlot, Value: nil}}
	expected = append(expected, deactivationEvents...)

	assert.ElementsMatch(t, expected, newEvents)
}

func TestFormExecutionRejected(t *testing.T) {
	form := Form{}

	newEvents := form.Run(rasa.EmptyTracker(), nil, nil)

	assert.ElementsMatch(t, []events.Event{&events.ActionExecutionRejected{}}, newEvents)
}

func TestRequestNextSlot(t *testing.T) {
	requiredSlot := "my slot"
	dispatcher := responses.NewDispatcher()
	form := Form{FormName: testFormName, RequiredSlots: []string{requiredSlot}}

	newEvents := form.Run(rasa.EmptyTracker(), nil, dispatcher)

	expected := []events.Event{
		&events.Form{Name: testFormName},
		&events.SlotSet{Name: requestedSlot, Value: requiredSlot}}
	assert.ElementsMatch(t, expected, newEvents)

	expectedResponses := []*responses.Message{{Template: fmt.Sprintf("utter_ask_%s", requiredSlot)}}
	assert.ElementsMatch(t, expectedResponses, dispatcher.Responses())
}
