package actions

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wochinge/go-rasa-sdk/rasa"
	"github.com/wochinge/go-rasa-sdk/rasa/events"
	"github.com/wochinge/go-rasa-sdk/rasa/request"
	"github.com/wochinge/go-rasa-sdk/rasa/responses"
)

type SimpleTestAction struct{}

func (action *SimpleTestAction) Run(_ *rasa.Tracker, _ *rasa.Domain, _ responses.ResponseDispatcher) []events.Event {
	return []events.Event{&events.SlotSet{Name: "test-slot", Value: "test-value"}}
}
func (action *SimpleTestAction) Name() string { return "test-action" }

func TestActionReturningEvents(t *testing.T) {
	action := &SimpleTestAction{}

	newEvents := action.Run(&rasa.Tracker{}, &rasa.Domain{}, responses.NewDispatcher())

	expectedEvents := []events.Event{&events.SlotSet{Name: "test-slot", Value: "test-value"}}
	assert.ElementsMatch(t, expectedEvents, newEvents)
}

type ActionDispatchingResponses struct{}

func (action *ActionDispatchingResponses) Run(_ *rasa.Tracker, _ *rasa.Domain,
	dispatcher responses.ResponseDispatcher) []events.Event {
	message := responses.Message{Text: "Hello World"}
	dispatcher.Utter(message)

	return []events.Event{}
}
func (action *ActionDispatchingResponses) Name() string { return "action-dispatching-responses" }

func TestActionDispatchingResponses(t *testing.T) {
	action := &ActionDispatchingResponses{}

	dispatcher := responses.NewDispatcher()
	action.Run(&rasa.Tracker{}, &rasa.Domain{}, dispatcher)

	expectedResponses := []responses.Message{{Text: "Hello World"}}
	assert.ElementsMatch(t, expectedResponses, dispatcher.Responses())
}

func TestActionResponseEmpty(t *testing.T) {
	response := actionResponse([]events.Event{}, responses.NewDispatcher())
	actualAsJSON, err := json.Marshal(response)

	assert.Nil(t, err)

	expectedResponse := `{"events":[],"responses":[]}`
	assert.Equal(t, expectedResponse, string(actualAsJSON))
}

func TestActionResponseWithMultipleResponses(t *testing.T) {
	dispatcher := responses.NewDispatcher()
	dispatcher.Utter(responses.Message{Text: "hi"})
	dispatcher.Utter(responses.Message{Template: "utter_ask"})

	response := actionResponse([]events.Event{}, dispatcher)
	actualAsJSON, err := json.Marshal(response)

	assert.Nil(t, err)

	expectedResponse := `{"events":[],"responses":[{"text":"hi"},{"text":"","template":"utter_ask"}]}`
	assert.Equal(t, expectedResponse, string(actualAsJSON))
}

func TestActionResponseWithEvents(t *testing.T) {
	newEvents := []events.Event{&events.Restarted{}, &events.SlotSet{Name: "my cool slot", Value: "best value"}}

	response := actionResponse(newEvents, responses.NewDispatcher())
	actualAsJSON, err := json.Marshal(response)

	assert.Nil(t, err)

	expectedResponse := `{"events":[{"event":"restart"},` +
		`{"event":"slot","name":"my cool slot","value":"best value"}],"responses":[]}`
	assert.Equal(t, expectedResponse, string(actualAsJSON))
}

type RejectingAction struct{}

func (action *RejectingAction) Run(_ *rasa.Tracker, _ *rasa.Domain, _ responses.ResponseDispatcher) []events.Event {
	return []events.Event{&events.ActionExecutionRejected{}}
}
func (action *RejectingAction) Name() string { return "test-reject" }

func TestActionRejectingExecution(t *testing.T) {
	actionRequest := request.CustomActionRequest{ActionToRun: "test-reject"}

	_, err := ExecuteAction(actionRequest, []Action{&RejectingAction{}})

	assert.IsType(t, &ExecutionRejectedError{}, err)
}
