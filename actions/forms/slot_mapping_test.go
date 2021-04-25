package forms

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wochinge/go-rasa-sdk/rasa"
	"github.com/wochinge/go-rasa-sdk/rasa/events"
)

const (
	testEntity      = "my entity"
	testEntityValue = "value"
	intentName      = "greet"
)

func TestFillSlotFromEntity(t *testing.T) {
	lastMessage := events.ParseData{Entities: []events.Entity{{Name: testEntity, Value: testEntityValue}}}
	tracker := rasa.Tracker{LatestMessage: lastMessage}

	value, found := SlotMapping{FromEntity: testEntity}.apply(nil, &tracker)
	assert.True(t, found)
	assert.Equal(t, testEntityValue, value)
}

func TestFillSlotFromEntityWithIntentSpecified(t *testing.T) {
	lastMessage := events.ParseData{Entities: []events.Entity{{Name: testEntity, Value: testEntityValue}},
		Intent: events.IntentParseResult{Name: intentName}}
	tracker := rasa.Tracker{LatestMessage: lastMessage}

	value, found := SlotMapping{FromEntity: testEntity, Intents: []string{"bye", intentName}}.apply(nil, &tracker)
	assert.True(t, found)
	assert.Equal(t, testEntityValue, value)
}

func TestFillSlotFromEntityWithIntentNotSpecified(t *testing.T) {
	lastMessage := events.ParseData{Entities: []events.Entity{{Name: testEntity, Value: testEntityValue}},
		Intent: events.IntentParseResult{Name: intentName}}
	tracker := rasa.Tracker{LatestMessage: lastMessage}

	_, found := SlotMapping{FromEntity: testEntity, Intents: []string{"bye"}}.apply(nil, &tracker)
	assert.False(t, found)
}

func TestFillSlotFromEntityWithIntentExcluded(t *testing.T) {
	lastMessage := events.ParseData{Entities: []events.Entity{{Name: testEntity, Value: testEntityValue}},
		Intent: events.IntentParseResult{Name: intentName}}
	tracker := rasa.Tracker{LatestMessage: lastMessage}

	_, found := SlotMapping{FromEntity: testEntity, ExcludedIntents: []string{intentName}}.apply(nil, &tracker)
	assert.False(t, found)
}

func TestFillSlotFromEntityWithIntentExcludedButAllowedWasGiven(t *testing.T) {
	lastMessage := events.ParseData{Entities: []events.Entity{{Name: testEntity, Value: testEntityValue}},
		Intent: events.IntentParseResult{Name: intentName}}
	tracker := rasa.Tracker{LatestMessage: lastMessage}

	_, found := SlotMapping{FromEntity: testEntity, ExcludedIntents: []string{"someother"}}.apply(nil, &tracker)
	assert.True(t, found)
}

func TestFillSlotFromText(t *testing.T) {
	text := "TestFillSlotFromText"

	lastMessage := events.ParseData{Entities: []events.Entity{{Name: "dasdas", Value: "Dasds"}},
		Intent: events.IntentParseResult{Name: intentName}, Text: text}
	tracker := rasa.Tracker{LatestMessage: lastMessage}

	value, found := SlotMapping{FromText: true, Intents: []string{"bye", intentName}}.apply(nil, &tracker)
	assert.True(t, found)
	assert.Equal(t, text, value)
}

func TestFillSlotFromValue(t *testing.T) {
	expectedValue := "free style"

	lastMessage := events.ParseData{Intent: events.IntentParseResult{Name: intentName}}
	tracker := rasa.Tracker{LatestMessage: lastMessage}

	value, found := SlotMapping{Intents: []string{"bye", intentName}, Value: expectedValue}.apply(nil, &tracker)
	assert.True(t, found)
	assert.Equal(t, expectedValue, value)
}

func TestMappingDoesNotApply(t *testing.T) {
	lastMessage := events.ParseData{Intent: events.IntentParseResult{Name: intentName}}
	tracker := rasa.Tracker{LatestMessage: lastMessage}

	_, found := SlotMapping{Intents: []string{"bye", intentName}}.apply(nil, &tracker)
	assert.False(t, found)
}

func TestMappingOnlyFirstRun(t *testing.T) {
	text := "TestMappingOnlyFirstRun"

	lastMessage := events.ParseData{Entities: []events.Entity{{Name: "dasdas", Value: "Dasds"}},
		Intent: events.IntentParseResult{Name: intentName}, Text: text}
	tracker := rasa.Tracker{LatestMessage: lastMessage}

	value, found := SlotMapping{FromText: true, ApplyToFirstRunOnly: true}.apply(&Form{FormName: "form"}, &tracker)
	assert.True(t, found)
	assert.Equal(t, text, value)
}

func TestMappingOnlyFirstRunIfNotFirstRun(t *testing.T) {
	text := "TestMappingOnlyFirstRunIfNotFirstRun"
	formName := "my form"

	lastMessage := events.ParseData{Entities: []events.Entity{{Name: "dasdas", Value: "Dasds"}},
		Intent: events.IntentParseResult{Name: intentName}, Text: text}
	tracker := rasa.Tracker{LatestMessage: lastMessage, ActiveForm: rasa.ActiveForm{Name: formName}}

	_, found := SlotMapping{FromText: true, ApplyToFirstRunOnly: true}.apply(&Form{FormName: formName}, &tracker)
	assert.False(t, found)
}
