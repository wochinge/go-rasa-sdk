// Package request is responsible for parsing the payload of the POST request which Rasa Open Source sends
// when requesting to execute a custom action.
package request

import (
	"encoding/json"
	"github.com/wochinge/go-rasa-sdk/rasa"
	"github.com/wochinge/go-rasa-sdk/rasa/events"
	"io"
)

// Parsed parses the payload which is sent to us from Rasa Open Source upon a request to execute a custom action.
// It returns the parsed payload or an error in case the json -> struct conversion failed.
func Parsed(requestBody io.Reader) (CustomActionRequest, error) {
	var parsedRequest CustomActionRequest
	parsedRequest.Tracker = *rasa.EmptyTracker()

	decoder := json.NewDecoder(requestBody)

	if err := decoder.Decode(&parsedRequest); err != nil {
		return parsedRequest, err
	}

	// parsedRequest.Domain = rasa.sanitizeDomain(parsedRequest.Domain)

	if parsedRequest.Tracker.RawEvents == nil {
		parsedRequest.Tracker.Events = []events.Event{}
		return parsedRequest, nil
	}

	trackerEvents, err := events.Parsed(parsedRequest.Tracker.RawEvents)

	if err != nil {
		return parsedRequest, err
	}

	parsedRequest.Tracker.Events = trackerEvents

	return parsedRequest, err
}

// CustomActionRequest exposes the data which Rasa Open Source sends as part of the action execution request.
type CustomActionRequest struct {
	// ActionToRun is the action which Rasa Open Source wants to run.
	ActionToRun string `json:"next_action"`
	// Tracker is the representation of the current conversation history of the user.
	Tracker rasa.Tracker `json:"tracker"`
	// Domain is the content of the current domain.yml of the currently running model.
	Domain rasa.Domain `json:"domain"`
}
