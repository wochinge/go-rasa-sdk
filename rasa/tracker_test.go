package rasa

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wochinge/go-rasa-sdk/rasa/events"
)

func TestSlotsToValidate(t *testing.T) {
	tracker := Tracker{Events: []events.Event{
		&events.SlotSet{Name: "not_a_candidate", Value: "some value"},
		&events.Action{Name: "my_form"},
		&events.SlotSet{Name: "candidate", Value: "interesting value"},
		&events.SlotSet{Name: "candidate_two", Value: "interesting value2"},
	}}

	slotCandidates := tracker.SlotsToValidate()

	expectedCandidates := map[string]interface{}{
		"candidate":     "interesting value",
		"candidate_two": "interesting value2",
	}
	assert.Equal(t, expectedCandidates, slotCandidates)
}

func TestSlotsToValidateWithNoCandidates(t *testing.T) {
	tracker := Tracker{Events: []events.Event{
		&events.SlotSet{Name: "not_a_candidate", Value: "some value"},
		&events.Action{Name: "my_form"},
	}}

	slotCandidates := tracker.SlotsToValidate()

	expectedCandidates := map[string]interface{}{}
	assert.Equal(t, expectedCandidates, slotCandidates)
}
