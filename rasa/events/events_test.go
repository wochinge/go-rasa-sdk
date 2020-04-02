package events

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEventType(t *testing.T) {
	types := []Type{action, sessionStarted, user, bot, userUtteranceReverted, actionReverted, restarted,
		storyExported, followUpAction, conversationPaused, conversationResumed, slotSet, allSlotsReset, form,
		formValidation, actionExecutionRejected}

	for _, eventType := range types {
		eventCreator, found := eventParser(Base{Type: eventType})
		assert.True(t, found)

		event := eventCreator()

		// Nuke type to invalid one if it's correctly re-assigned based on struct type
		event.SetType(action)

		events := WithTypeKeys(event)
		assert.Len(t, events, 1)
		assert.True(t, eventType == event.EventType())
	}
}

func TestParsedDataEntityFor(t *testing.T) {
	name, expectedValue := "user name", "Maria"

	parsed := ParseData{Entities: []Entity{{Name: "other", Value: "doesnt matter"}, {Name: name, Value: expectedValue}}}

	actualValue, found := parsed.EntityFor(name)
	assert.True(t, found)
	assert.Equal(t, expectedValue, actualValue)
}

func TestParsedDataEntityForNotExisting(t *testing.T) {
	parsed := ParseData{Entities: []Entity{{Name: "other", Value: "doesnt matter"}}}

	_, found := parsed.EntityFor("not there")
	assert.False(t, found)
}
