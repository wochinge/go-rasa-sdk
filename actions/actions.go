// Package actions contains everything to implement custom actions using the go-rasa-sdk.
package actions

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/wochinge/go-rasa-sdk/logging"
	"github.com/wochinge/go-rasa-sdk/rasa"
	"github.com/wochinge/go-rasa-sdk/rasa/events"
	"github.com/wochinge/go-rasa-sdk/rasa/request"
	"github.com/wochinge/go-rasa-sdk/rasa/responses"
)

// Action is the interface for all custom action implementations.
type Action interface {
	// Run runs the custom action in the given context and returns new conversation events.
	// Any messages dispatched will be sent to the user.
	Run(tracker *rasa.Tracker, domain *rasa.Domain, dispatcher responses.ResponseDispatcher) []events.Event
	// Name returns the name of the custom action.
	Name() string
}

// NotFoundError happens when no action was found for the given name.
type NotFoundError struct{ name string }

func (e *NotFoundError) Error() string { return fmt.Sprintf("action '%s' was not found.", e.name) }

// ExecutionRejectedError happens when the action rejected its execution.
type ExecutionRejectedError struct{ name string }

func (e *ExecutionRejectedError) Error() string {
	return fmt.Sprintf("action '%s' rejected execution.", e.name)
}

// ExecuteAction executes the custom action which was requested by Rasa Open Source.
func ExecuteAction(actionRequest request.CustomActionRequest,
	availableActions []Action) (map[string]interface{}, error) {
	actionToRun := actionFor(actionRequest.ActionToRun, availableActions)

	if actionToRun == nil {
		log.WithFields(log.Fields{logging.ActionNameKey: actionRequest.ActionToRun}).Warn("Requested action not found.")
		return nil, &NotFoundError{name: actionRequest.ActionToRun}
	}

	log.WithFields(
		log.Fields{logging.ActionNameKey: actionToRun,
			logging.ConversationIDKey: actionRequest.Tracker.ConversationID}).Debug("Received request to run action.")

	dispatcher := responses.NewDispatcher()
	newEvents := actionToRun.Run(&actionRequest.Tracker, &actionRequest.Domain, dispatcher)

	if events.HasRejection(newEvents) {
		log.WithFields(log.Fields{logging.ActionNameKey: actionToRun}).Debug("Action rejected execution.")
		return nil, &ExecutionRejectedError{name: actionRequest.ActionToRun}
	}

	log.WithFields(
		log.Fields{logging.ActionNameKey: actionToRun, logging.EventKeys: newEvents}).Debug("Action execution finished.")

	return actionResponse(newEvents, dispatcher), nil
}

func actionFor(name string, actions []Action) Action {
	for _, action := range actions {
		if action.Name() == name {
			return action
		}
	}

	return nil
}

func actionResponse(newEvents []events.Event, dispatcher responses.ResponseDispatcher) map[string]interface{} {
	return map[string]interface{}{
		"events":    events.WithTypeKeys(newEvents...),
		"responses": dispatcher.Responses(),
	}
}
