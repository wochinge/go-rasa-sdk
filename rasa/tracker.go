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
	Paused             bool       `json:"paused"`
	FollowUpAction     string     `json:"followup_action"`
	ActiveForm         ActiveForm `json:"active_form"`
	LatestActionName   string     `json:"latest_action_name"`
	LatestInputChannel string     `json:"latest_input_channel"`
}

func (tracker *Tracker) NoFormValidation() bool {
	return !tracker.ActiveForm.Validate || tracker.LatestActionName == "action_listen"
}

func EmptyTracker() *Tracker {
	tracker := &Tracker{ActiveForm: ActiveForm{Validate:true}}
	return tracker.Init()
}

func (tracker *Tracker) Init() *Tracker{
	if tracker.Slots == nil {
		tracker.Slots = map[string]interface{}{}
	}
	return tracker
}


type ActiveForm struct {
	Name           string           `json:"name"`
	Validate       bool             `json:"validate"`
	Rejected       bool             `json:"rejected"`
	TriggerMessage events.ParseData `json:"trigger_message"`
}
