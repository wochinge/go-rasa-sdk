package actions

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/wochinge/go-rasa-sdk/rasa"
	"github.com/wochinge/go-rasa-sdk/rasa/events"
	"github.com/wochinge/go-rasa-sdk/rasa/responses"
	"testing"
)

type SimpleTestAction struct{}
func (action *SimpleTestAction) Run(_ *rasa.Tracker, _ *rasa.Domain, _ responses.ResponseDispatcher) []events.Event {
	return []events.Event{events.Slot{Name: "test-slot", Value: "test-value"}}
}
func (action *SimpleTestAction) Name() string {return "test-action"}

func TestActionReturningEvents(t *testing.T) {
	var action Action
	action = &SimpleTestAction{}

	newEvents := action.Run(&rasa.Tracker{}, &rasa.Domain{}, responses.NewDispatcher())

	expectedEvents := []events.Event{events.Slot{Name: "test-slot", Value: "test-value"}}
	assert.ElementsMatch(t, expectedEvents, newEvents)
}

type ActionDispatchingResponses struct{}
func (action *ActionDispatchingResponses) Run(_ *rasa.Tracker, _ *rasa.Domain, dispatcher responses.ResponseDispatcher) []events.Event {
	message := responses.BotMessage{Text:"Hello World"}
	dispatcher.Utter(message)
	return []events.Event{}
}
func (action *ActionDispatchingResponses) Name() string {return "action-dispatching-responses"}

func TestActionDispatchingResponses(t *testing.T) {
	var action Action
	action = &ActionDispatchingResponses{}

	dispatcher := responses.NewDispatcher()
	action.Run(&rasa.Tracker{}, &rasa.Domain{}, dispatcher)

	expectedResponses := []responses.BotMessage{{Text: "Hello World"}}
	assert.ElementsMatch(t, expectedResponses, dispatcher.Responses())
}

func TestActionResponseEmpty(t *testing.T) {
	response := ActionResponse([]events.Event{}, responses.NewDispatcher())
	actualAsJson, err := json.Marshal(response)

	assert.Nil(t, err)


	expectedResponse := `{"events":[],"responses":[]}`
	assert.Equal(t, expectedResponse, string(actualAsJson))
}

func TestActionResponseWithMultipleResponses(t *testing.T) {
	dispatcher := responses.NewDispatcher()
	dispatcher.Utter(responses.BotMessage{Text:"hi"})
	dispatcher.Utter(responses.BotMessage{Template: "utter_ask"})

	response := ActionResponse([]events.Event{}, dispatcher)
	actualAsJson, err := json.Marshal(response)

	assert.Nil(t, err)

	expectedResponse := `{"events":[],"responses":[{"text":"hi"},{"text":"","template":"utter_ask"}]}`
	assert.Equal(t, expectedResponse, string(actualAsJson))
}

func TestActionResponseWithEvents(t *testing.T) {
	newEvents := []events.Event{events.Restart(), events.SetSlot("my cool slot", "best value")}

	response  := ActionResponse(newEvents,responses. NewDispatcher())
	actualAsJson, err := json.Marshal(response)

	assert.Nil(t, err)

	expectedResponse := `{"events":[{"event":"restart"},{"event":"slot","name":"my cool slot","value":"best value"}],"responses":[]}`
	assert.Equal(t, expectedResponse, string(actualAsJson))
}
