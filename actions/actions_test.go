package actions

import (
	"github.com/stretchr/testify/assert"
	"github.com/wochinge/go-rasa-sdk/rasa"
	"github.com/wochinge/go-rasa-sdk/rasa/events"
	"github.com/wochinge/go-rasa-sdk/rasa/responses"
	"testing"
)

type SimpleTestAction struct{}
func (action *SimpleTestAction) Run(_ *rasa.Tracker, _ *rasa.Domain, _ ResponseDispatcher) []events.Event {
	return []events.Event{events.SlotSet{Name: "test-slot", Value: "test-value"}}
}
func (action *SimpleTestAction) Name() string {return "test-action"}

func TestActionReturningEvents(t *testing.T) {
	var action Action
	action = &SimpleTestAction{}

	newEvents := action.Run(&rasa.Tracker{}, &rasa.Domain{}, &responseDispatcher{})

	expectedEvents := []events.Event{events.SlotSet{Name: "test-slot", Value: "test-value"}}
	assert.ElementsMatch(t, expectedEvents, newEvents)
}

type ActionDispatchingResponses struct{}
func (action *ActionDispatchingResponses) Run(_ *rasa.Tracker, _ *rasa.Domain, dispatcher ResponseDispatcher) []events.Event {
	message := responses.BotMessage{Text:"Hello World"}
	dispatcher.Utter(message)
	return []events.Event{}
}
func (action *ActionDispatchingResponses) Name() string {return "action-dispatching-responses"}

func TestActionDispatchingResponses(t *testing.T) {
	var action Action
	action = &ActionDispatchingResponses{}

	dispatcher := NewDispatcher()
	action.Run(&rasa.Tracker{}, &rasa.Domain{}, dispatcher)

	expectedResponses := []responses.BotMessage{{Text: "Hello World"}}
	assert.ElementsMatch(t, expectedResponses, dispatcher.Responses())
}

func TestActionResponseEmpty(t *testing.T) {
	response, err := ActionResponse([]events.Event{}, NewDispatcher())

	assert.Nil(t, err)

	expectedResponse := `{"events":[],"responses":[]}`
	assert.Equal(t, expectedResponse, string(response))
}

func TestActionResponseWithMultipleResponses(t *testing.T) {
	dispatcher := NewDispatcher()
	dispatcher.Utter(responses.BotMessage{Text:"hi"})
	dispatcher.Utter(responses.BotMessage{Template: "utter_ask"})

	response, err := ActionResponse([]events.Event{}, dispatcher)

	assert.Nil(t, err)

	expectedResponse := `{"events":[],"responses":[{"text":"hi"},{"text":"","template":"utter_ask"}]}`
	assert.Equal(t, expectedResponse, string(response))
}

func TestActionResponseWithEvents(t *testing.T) {
	newEvents := []events.Event{events.Restarted{}, events.SlotSet{Name:"my cool slot", Value: "best value"}}

	response, err := ActionResponse(newEvents, NewDispatcher())

	assert.Nil(t, err)

	expectedResponse := `{"events":[{},{"name":"my cool slot","value":"best value"}],"responses":[]}`
	assert.Equal(t, expectedResponse, string(response))
}

