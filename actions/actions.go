package actions

import (
	"encoding/json"
	"github.com/wochinge/go-rasa-sdk/rasa"
	"github.com/wochinge/go-rasa-sdk/rasa/events"
	"github.com/wochinge/go-rasa-sdk/rasa/responses"
)

type Action interface {
	Run(tracker *rasa.Tracker, domain *rasa.Domain, dispatcher ResponseDispatcher) []events.Event
	Name() string
}

func ActionFor(name string, actions []Action) (Action, error) {
	for _, action := range actions {
		if action.Name() == name {
			return action, nil
		}
	}
	return nil, nil
}

type ResponseDispatcher interface {
	Utter(responses.BotMessage)

	Responses() []responses.BotMessage
}

type responseDispatcher struct {
	responses []responses.BotMessage
}

func (dispatcher *responseDispatcher) Utter(message responses.BotMessage) {
	dispatcher.responses = append(dispatcher.responses, message)
}

func (dispatcher *responseDispatcher) Responses() []responses.BotMessage {
	return dispatcher.responses
}

func NewDispatcher() ResponseDispatcher {
	return &responseDispatcher{responses: []responses.BotMessage{}}
}

func ActionResponse(newEvents []events.Event, dispatcher ResponseDispatcher) ([]byte, error) {
	response := map[string]interface{}{
		"events": newEvents,
		"responses": dispatcher.Responses(),
	}

	return json.Marshal(response)
}