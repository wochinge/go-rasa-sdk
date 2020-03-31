package actions

import (
	"fmt"
	"github.com/wochinge/go-rasa-sdk/rasa"
	"github.com/wochinge/go-rasa-sdk/rasa/events"
	"github.com/wochinge/go-rasa-sdk/rasa/responses"
)

type Action interface {
	Run(tracker *rasa.Tracker, domain *rasa.Domain, dispatcher responses.ResponseDispatcher) []events.Event
	Name() string
}

func ExecuteAction(actionRequest rasa.CustomActionRequest, availableActions []Action) (map[string]interface{}, error) {
	actionToRun := actionFor(actionRequest.ActionToRun, availableActions)

	if actionToRun == nil {
		return nil, fmt.Errorf("action with this name not found")
	}

	dispatcher := responses.NewDispatcher()
	newEvents := actionToRun.Run(&actionRequest.Tracker, &actionRequest.Domain, dispatcher)

	responseBody := ActionResponse(newEvents, dispatcher)
	return responseBody, nil
}

func actionFor(name string, actions []Action) Action {
	for _, action := range actions {
		if action.Name() == name {
			return action
		}
	}
	return nil
}

func ActionResponse(newEvents []events.Event, dispatcher responses.ResponseDispatcher) map[string]interface{} {
	for _, event := range newEvents {
		event.SetType(event.EventType())
	}

	return map[string]interface{}{
		"events": newEvents,
		"responses": dispatcher.Responses(),
	}
}
