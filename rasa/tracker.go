package rasa

import (
	"encoding/json"
	"github.com/wochinge/go-rasa-sdk/rasa/events"
)

// Tracker represents a conversation history of a user.
type Tracker struct {
	// ConversationID is a unique ID for the conversation.
	ConversationID string `json:"sender_id"`
	// Slots and their values.
	Slots map[string]interface{} `json:"slots"`
	// LatestMessage contains the NLU parsing result for the last user message.
	LatestMessage events.ParseData `json:"latest_message"`
	// LatestEventTime as unix timestamp.
	LatestEventTime float64 `json:"latest_event_time"`
	// RawEvents are the unparsed conversation events as JSON.
	RawEvents []json.RawMessage `json:"events"`
	// Events within the conversation.
	Events []events.Event
	// Paused is true if the bot is currently not allowed to send messages to the user.
	Paused bool `json:"paused"`
	// FollowUpAction is the name of an action which the bot is forced to execute next.
	FollowUpAction string `json:"followup_action"`
	// ActiveForm describes whether and which form is currently active.
	ActiveForm ActiveForm `json:"active_form"`
	// LatestActionName is the name of the last action the bot executed.
	LatestActionName string `json:"latest_action_name"`
	// LatestInputChannel is the name of the last channel (e.g. Slack, Telegram) which the user used to speak to the
	// assistant.
	LatestInputChannel string `json:"latest_input_channel"`
}

// NoFormValidation checks if the form should validate candidates before filling the slots.
func (tracker *Tracker) NoFormValidation() bool {
	return !tracker.ActiveForm.Validate || tracker.LatestActionName != "action_listen"
}

// EmptyTracker returns a new tracker with its default default values set.
func EmptyTracker() *Tracker {
	tracker := &Tracker{ActiveForm: ActiveForm{Validate: true}}

	return tracker.Init()
}

func (tracker *Tracker) Init() *Tracker {
	if tracker.Slots == nil {
		tracker.Slots = map[string]interface{}{}
	}

	return tracker
}

// ActiveForm describes a potentially active form.
type ActiveForm struct {
	// Name of the currently active form.
	Name string `json:"name"`
	// Validate is `true` if the slot candidates should be validated before filling the slot.
	Validate bool `json:"validate"`
	// Rejected specifies if the form rejected its execution.
	Rejected bool `json:"rejected"`
	// TriggerMessage is the first message which started the form.
	TriggerMessage events.ParseData `json:"trigger_message"`
}
