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

	mapping := SlotMapping{FromEntity: testEntity}
	value, found := mapping.apply(nil, &tracker)
	assert.True(t, found)
	assert.Equal(t, testEntityValue, value)
}

func TestFillSlotFromEntityWithIntentSpecified(t *testing.T) {
	lastMessage := events.ParseData{Entities: []events.Entity{{Name: testEntity, Value: testEntityValue}},
		Intent: events.IntentParseResult{Name: intentName}}
	tracker := rasa.Tracker{LatestMessage: lastMessage}

	mapping := SlotMapping{FromEntity: testEntity, Intents: []string{"bye", intentName}}
	value, found := mapping.apply(nil, &tracker)
	assert.True(t, found)
	assert.Equal(t, testEntityValue, value)
}

func TestFillSlotFromEntityWithIntentNotSpecified(t *testing.T) {
	lastMessage := events.ParseData{Entities: []events.Entity{{Name: testEntity, Value: testEntityValue}},
		Intent: events.IntentParseResult{Name: intentName}}
	tracker := rasa.Tracker{LatestMessage: lastMessage}

	mapping := SlotMapping{FromEntity: testEntity, Intents: []string{"bye"}}
	_, found := mapping.apply(nil, &tracker)
	assert.False(t, found)
}

func TestFillSlotFromEntityWithIntentExcluded(t *testing.T) {
	lastMessage := events.ParseData{Entities: []events.Entity{{Name: testEntity, Value: testEntityValue}},
		Intent: events.IntentParseResult{Name: intentName}}
	tracker := rasa.Tracker{LatestMessage: lastMessage}

	mapping := SlotMapping{FromEntity: testEntity, ExcludedIntents: []string{intentName}}
	_, found := mapping.apply(nil, &tracker)
	assert.False(t, found)
}

func TestFillSlotFromEntityWithIntentExcludedButAllowedWasGiven(t *testing.T) {
	lastMessage := events.ParseData{Entities: []events.Entity{{Name: testEntity, Value: testEntityValue}},
		Intent: events.IntentParseResult{Name: intentName}}
	tracker := rasa.Tracker{LatestMessage: lastMessage}

	mapping := SlotMapping{FromEntity: testEntity, ExcludedIntents: []string{"someother"}}
	_, found := mapping.apply(nil, &tracker)
	assert.True(t, found)
}

func TestFillSlotFromText(t *testing.T) {
	text := "TestFillSlotFromText"

	lastMessage := events.ParseData{Entities: []events.Entity{{Name: "dasdas", Value: "Dasds"}},
		Intent: events.IntentParseResult{Name: intentName}, Text: text}
	tracker := rasa.Tracker{LatestMessage: lastMessage}

	mapping := SlotMapping{FromText: true, Intents: []string{"bye", intentName}}
	value, found := mapping.apply(nil, &tracker)
	assert.True(t, found)
	assert.Equal(t, text, value)
}

func TestFillSlotFromValue(t *testing.T) {
	expectedValue := "free style"

	lastMessage := events.ParseData{Intent: events.IntentParseResult{Name: intentName}}
	tracker := rasa.Tracker{LatestMessage: lastMessage}

	mapping := SlotMapping{Intents: []string{"bye", intentName}, Value: expectedValue}
	value, found := mapping.apply(nil, &tracker)
	assert.True(t, found)
	assert.Equal(t, expectedValue, value)
}

func TestMappingDoesNotApply(t *testing.T) {
	lastMessage := events.ParseData{Intent: events.IntentParseResult{Name: intentName}}
	tracker := rasa.Tracker{LatestMessage: lastMessage}

	mapping := SlotMapping{Intents: []string{"bye", intentName}}
	_, found := mapping.apply(nil, &tracker)
	assert.False(t, found)
}

func TestMappingOnlyFirstRun(t *testing.T) {
	text := "TestMappingOnlyFirstRun"

	lastMessage := events.ParseData{Entities: []events.Entity{{Name: "dasdas", Value: "Dasds"}},
		Intent: events.IntentParseResult{Name: intentName}, Text: text}
	tracker := rasa.Tracker{LatestMessage: lastMessage}

	mapping := SlotMapping{FromText: true, ApplyToFirstRunOnly: true}
	value, found := mapping.apply(&Form{FormName: "form"}, &tracker)
	assert.True(t, found)
	assert.Equal(t, text, value)
}

func TestMappingOnlyFirstRunIfNotFirstRun(t *testing.T) {
	text := "TestMappingOnlyFirstRunIfNotFirstRun"
	formName := "my form"

	lastMessage := events.ParseData{Entities: []events.Entity{{Name: "dasdas", Value: "Dasds"}},
		Intent: events.IntentParseResult{Name: intentName}, Text: text}
	tracker := rasa.Tracker{LatestMessage: lastMessage, ActiveLoop: rasa.ActiveLoop{Name: formName}}

	mapping := SlotMapping{FromText: true, ApplyToFirstRunOnly: true}
	_, found := mapping.apply(&Form{FormName: formName}, &tracker)
	assert.False(t, found)
}
