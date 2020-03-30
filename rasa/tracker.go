package rasa

import (
	"encoding/json"
	"github.com/wochinge/go-rasa-sdk/rasa/events"
)

type Tracker struct {
	ConversationId     string                 `json:"sender_id"`
	Slots              map[string]interface{} `json:"slots"`
	LatestMessage      events.ParseData       `json:"latest_message"`
	LatestEventTime    float64                `json:"latest_event_time"`
	RawEvents          []json.RawMessage      `json:"events"`
	Events             []events.Event
	Paused             bool                   `json:"paused"`
	FollowUpAction     string                 `json:"followup_action"`
	ActiveForm         map[string]interface{} `json:"action_form"` // TODO: Maybe this is a string?
	LatestActionName   string                 `json:"latest_action_name"`
	LatestInputChannel string                 `json:"latest_input_channel"`
}
