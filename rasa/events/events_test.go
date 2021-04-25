package events

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEventType(t *testing.T) {
	types := []Type{action, sessionStarted, user, bot, userUtteranceReverted, actionReverted, restarted,
		storyExported, followUpAction, conversationPaused, conversationResumed, slotSet, allSlotsReset, activeLoop,
		form, loopInterrupted, formValidation, actionExecutionRejected, reminderScheduled, reminderCancelled}

	for _, eventType := range types {
		eventCreator, found := eventParser(Base{Type: eventType})
		assert.True(t, found)

		event := eventCreator()

		// Nuke type to invalid one if it's correctly re-assigned based on struct type
		event.SetType(action)

		events := WithTypeKeys(event)
		assert.Len(t, events, 1)
		assert.Equal(t, eventType, event.EventType())
	}
}

func TestParseUnknownEvent(t *testing.T) {
	unknownType := "never seen this before"
	unknownEvent := json.RawMessage([]byte(fmt.Sprintf(`{"event": "%s"}`, unknownType)))

	events, err := Parsed([]json.RawMessage{unknownEvent})

	assert.Nil(t, err)
	assert.ElementsMatch(t, []Event{&Base{Type: unknown}}, events)
	assert.Equal(t, unknown, events[0].EventType())
}

func TestParseBasedOnyTypeKeyError(t *testing.T) {
	eventWithUnexpectedFormat := json.RawMessage([]byte(`{"event": "action", "policy": 2}`))
	events, err := Parsed([]json.RawMessage{eventWithUnexpectedFormat})

	assert.Nil(t, err)
	assert.ElementsMatch(t, []Event{&Base{Type: action}}, events)
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
