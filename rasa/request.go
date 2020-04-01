package rasa

import (
	"encoding/json"
	"github.com/wochinge/go-rasa-sdk/rasa/events"
	"io"
)

func Parsed(requestBody io.Reader) (CustomActionRequest, error) {
	var parsedRequest CustomActionRequest
	parsedRequest.Tracker = *EmptyTracker()

	decoder := json.NewDecoder(requestBody)

	if err:= decoder.Decode(&parsedRequest); err != nil {
		return parsedRequest, err
	}

	parsedRequest.Domain = sanitizeDomain(parsedRequest.Domain)

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


type CustomActionRequest struct {
	ActionToRun string  `json:"next_action"`
	Tracker     Tracker `json:"tracker"`
	Domain      Domain  `json:"domain"`
}

// https://mholt.github.io/json-to-go/
