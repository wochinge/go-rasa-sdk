package forms

import (
	"github.com/stretchr/testify/assert"
	"github.com/wochinge/go-rasa-sdk/rasa"
	"github.com/wochinge/go-rasa-sdk/rasa/events"
	"testing"
)

func TestFillSlotFromEntity(t *testing.T) {
	entityName, expectedValue := "entity FormName", "expectedValue"

	lastMessage := events.ParseData{Entities: []events.Entity{{Name: entityName, Value: expectedValue}}}
	tracker := rasa.Tracker{LatestMessage: lastMessage}

	value, found := SlotMapping{Entity: entityName}.apply(nil, &tracker)
	assert.True(t, found)
	assert.Equal(t, expectedValue, value)
}

func TestFillSlotFromEntityWithIntentSpecified(t *testing.T) {
	entityName, expectedValue := "entity FormName", "expectedValue"
	intentName := "greet"

	lastMessage := events.ParseData{Entities: []events.Entity{{Name: entityName, Value: expectedValue}},
		Intent: events.IntentParseResult{Name: intentName}}
	tracker := rasa.Tracker{LatestMessage: lastMessage}

	value, found := SlotMapping{Entity: entityName, Intents: []string{"bye", intentName}}.apply(nil, &tracker)
	assert.True(t, found)
	assert.Equal(t, expectedValue, value)
}

func TestFillSlotFromEntityWithIntentNotSpecified(t *testing.T) {
	entityName, expectedValue := "entity FormName", "expectedValue"

	lastMessage := events.ParseData{Entities: []events.Entity{{Name: entityName, Value: expectedValue}},
		Intent: events.IntentParseResult{Name: "greet"}}
	tracker := rasa.Tracker{LatestMessage: lastMessage}

	_, found := SlotMapping{Entity: entityName, Intents: []string{"bye"}}.apply(nil, &tracker)
	assert.False(t, found)
}

func TestFillSlotFromEntityWithIntentExcluded(t *testing.T) {
	entityName, expectedValue := "entity FormName", "expectedValue"
	intentName := "greet"

	lastMessage := events.ParseData{Entities: []events.Entity{{Name: entityName, Value: expectedValue}},
		Intent: events.IntentParseResult{Name: intentName}}
	tracker := rasa.Tracker{LatestMessage: lastMessage}

	_, found := SlotMapping{Entity: entityName, ExcludedIntents: []string{intentName}}.apply(nil, &tracker)
	assert.False(t, found)
}

func TestFillSlotFromText(t *testing.T) {
	intentName, text := "greet", "hellooooooo"

	lastMessage := events.ParseData{Entities: []events.Entity{{Name: "dasdas", Value: "Dasds"}},
		Intent: events.IntentParseResult{Name: intentName}, Text: text}
	tracker := rasa.Tracker{LatestMessage: lastMessage}

	value, found := SlotMapping{FromText: true, Intents: []string{"bye", intentName}}.apply(nil, &tracker)
	assert.True(t, found)
	assert.Equal(t, text, value)
}

func TestFillSlotFromValue(t *testing.T) {
	intentName := "greet"
	expectedValue := "free style"

	lastMessage := events.ParseData{Intent: events.IntentParseResult{Name: intentName}}
	tracker := rasa.Tracker{LatestMessage: lastMessage}

	value, found := SlotMapping{Intents: []string{"bye", intentName}, Value: expectedValue}.apply(nil, &tracker)
	assert.True(t, found)
	assert.Equal(t, expectedValue, value)
}

func TestMappingDoesNotApply(t *testing.T) {
	intentName := "greet"

	lastMessage := events.ParseData{Intent: events.IntentParseResult{Name: intentName}}
	tracker := rasa.Tracker{LatestMessage: lastMessage}

	_, found := SlotMapping{Intents: []string{"bye", intentName}}.apply(nil, &tracker)
	assert.False(t, found)
}

func TestMappingOnlyFirstRun(t *testing.T) {
	intentName, text := "greet", "hellooooooo"

	lastMessage := events.ParseData{Entities: []events.Entity{{Name: "dasdas", Value: "Dasds"}},
		Intent: events.IntentParseResult{Name: intentName}, Text: text}
	tracker := rasa.Tracker{LatestMessage: lastMessage}

	value, found := SlotMapping{FromText: true, ApplyToFirstRunOnly: true}.apply(&Form{FormName:"form"}, &tracker)
	assert.True(t, found)
	assert.Equal(t, text, value)
}

func TestMappingOnlyFirstRunIfNotFirstRun(t *testing.T) {
	intentName, text := "greet", "hellooooooo"
	formName := "my form"

	lastMessage := events.ParseData{Entities: []events.Entity{{Name: "dasdas", Value: "Dasds"}},
		Intent: events.IntentParseResult{Name: intentName}, Text: text}
	tracker := rasa.Tracker{LatestMessage: lastMessage, ActiveForm:rasa.ActiveForm{Name:formName}}


	_, found := SlotMapping{FromText: true, ApplyToFirstRunOnly: true}.apply(&Form{FormName:formName}, &tracker)
	assert.False(t, found)
}