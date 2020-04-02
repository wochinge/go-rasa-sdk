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

type NotFoundError struct {name string}
func (e *NotFoundError) Error() string { return fmt.Sprintf("action '%s' was not found.", e.name)}

type ExecutionRejectedError struct {name string}
func (e *ExecutionRejectedError) Error() string { return fmt.Sprintf("action '%s' rejected execution.", e.name)}

func ExecuteAction(actionRequest rasa.CustomActionRequest, availableActions []Action) (map[string]interface{}, error) {
	actionToRun := actionFor(actionRequest.ActionToRun, availableActions)

	if actionToRun == nil {
		return nil, &NotFoundError{name:actionRequest.ActionToRun}
	}

	dispatcher := responses.NewDispatcher()
	newEvents := actionToRun.Run(&actionRequest.Tracker, &actionRequest.Domain, dispatcher)

	if events.HasRejection(newEvents) {
		return nil, &ExecutionRejectedError{name:actionRequest.ActionToRun}
	}
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
	return map[string]interface{}{
		"events": events.WithTypeKeys(newEvents...),
		"responses": dispatcher.Responses(),
	}
}
