// Package events contains the representation of Rasa Open Source conversation events
// (https://rasa.com/docs/rasa/api/events/) in Go and tools to work with them.
package events

import (
	"encoding/json"
	"github.com/wochinge/go-rasa-sdk/rasa/responses"
)

// Type represents the type key which is part of each conversation event.
type Type string

const (
	action         Type = "action"
	user           Type = "user"
	bot            Type = "bot"
	sessionStarted Type = "session_started"
	slotSet        Type = "slot"

	conversationPaused  Type = "pause"
	conversationResumed Type = "resume"

	form                    Type = "form"
	formValidation          Type = "form_validation"
	actionExecutionRejected Type = "action_execution_rejected"

	followUpAction Type = "followup"
	storyExported  Type = "export"

	actionReverted        Type = "undo"
	userUtteranceReverted Type = "rewind"
	restarted             Type = "restart"
	allSlotsReset         Type = "reset_slots"
)

// Parsed parses and returns conversation events from JSON to their Go representation.
func Parsed(rawEvents []json.RawMessage) ([]Event, error) {
	var events []Event

	for _, rawEvent := range rawEvents {
		var minimalEvent Base

		if err := json.Unmarshal(rawEvent, &minimalEvent); err != nil {
			return []Event{}, err
		}

		if event := parseBasedOnyTypeKey(minimalEvent, rawEvent); event != nil {
			events = append(events, event)
		}
	}

	return events, nil
}

func parseBasedOnyTypeKey(base Base, raw json.RawMessage) Event {
	eventCreator, ok := eventParser(base)

	if !ok {
		return nil
	}

	event := eventCreator()
	if err := json.Unmarshal(raw, &event); err != nil {
		return nil
	}

	return event
}

func eventParser(base Base) (func() Event, bool) {
	eventParsers := map[Type]func() Event{
		action:         func() Event { return &Action{Base: base} },
		user:           func() Event { return &User{Base: base} },
		bot:            func() Event { return &Bot{Base: base} },
		sessionStarted: func() Event { return &SessionStarted{Base: base} },
		slotSet:        func() Event { return &SlotSet{Base: base} },

		conversationPaused:  func() Event { return &ConversationPaused{Base: base} },
		conversationResumed: func() Event { return &ConversationResumed{Base: base} },

		form:                    func() Event { return &Form{Base: base} },
		formValidation:          func() Event { return &FormValidation{Base: base} },
		actionExecutionRejected: func() Event { return &ActionExecutionRejected{Action: Action{Base: base}} },

		followUpAction: func() Event { return &FollowUpAction{Base: base} },
		storyExported:  func() Event { return &StoryExported{Base: base} },

		actionReverted:        func() Event { return &ActionReverted{Base: base} },
		userUtteranceReverted: func() Event { return &UserUtteranceReverted{Base: base} },
		restarted:             func() Event { return &Restarted{Base: base} },
		allSlotsReset:         func() Event { return &AllSlotsReset{Base: base} },
	}

	eventCreator, found := eventParsers[base.Type]

	return eventCreator, found
}

// WithTypeKeys sets the event type based on their current struct type.
// This is required to make sure that structs initialized like `SessionStarted{}` have the correct type key when
// they are encoded as JSON.
func WithTypeKeys(events ...Event) []Event {
	for _, event := range events {
		event.SetType(event.EventType())
	}

	return events
}

// HasRejection returns true if there is a `ActionExecutionRejected` in the given list of events.
func HasRejection(events []Event) bool {
	for _, event := range events {
		if _, ok := event.(*ActionExecutionRejected); ok {
			return true
		}
	}

	return false
}

// Event is the interface which all conversation events have to suffice.
type Event interface {
	EventType() Type
	SetType(Type)
}

// Base contains the data of an event which is common to all events.
type Base struct {
	Type      Type                   `json:"event,omitempty"`
	Timestamp float64                `json:"timestamp,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// EventType returns the string identifier of the type of event.
func (*Base) EventType() Type { return "unknown" }

// SetType sets the type of an event.
func (base *Base) SetType(eventType Type) { base.Type = eventType }

// Action is an event which represents that the assistant executed an action during the conversation.
type Action struct {
	Base
	// Policy is the policy which decided to run this action.
	Policy string `json:"policy"`
	// Confidence of the policy decision.
	Confidence float64 `json:"confidence"`
	// Name of the action which was run.
	Name string `json:"name"`
}

func (*Action) EventType() Type { return action }

// SessionStarted represents that a new conversation session
// (https://rasa.com/docs/rasa/core/domains/#session-configuration) was started.
type SessionStarted struct {
	Base
}

func (*SessionStarted) EventType() Type { return sessionStarted }

// User represents user messages in the conversation history.
type User struct {
	Base
	// Text of the user message.
	Text string `json:"text"`
	// InputChannel the user used to send their message (e.g. Slack, REST, Telegram).
	InputChannel string `json:"input_channel"`
	// MessageID is a unique ID of the message.
	MessageID string `json:"message_id"`
	// ParseData contains the result of the NLU prediction.
	ParseData ParseData `json:"parse_data"`
}

func (*User) EventType() Type { return user }

// ParseData represents the NLU prediction result.
type ParseData struct {
	// Intent is the predicted intent of the message.
	Intent IntentParseResult `json:"intent"`
	// Entities which were extracted from the message.
	Entities []Entity `json:"entities"`
	// IntentRanking shows the likeliness for other intents.
	IntentRanking []IntentParseResult `json:"intent_ranking"`
	// Text of the message.
	Text string `json:"text"`
}

// EntityFor returns the entity for a given entity name. Returns `nil` in case no entity with this name was found.
func (data ParseData) EntityFor(name string) (interface{}, bool) {
	for _, entity := range data.Entities {
		if entity.Name == name {
			return entity.Value, true
		}
	}

	return "", false
}

// IntentParseResult of the NLU prediction.
type IntentParseResult struct {
	// Name of the intent.
	Name string `json:"name"`
	// Confidence that the message has this intent.
	Confidence float64 `json:"confidence"`
}

// FromEntity represents entities (e.g. names, numbers) which were extracted from the message.
type Entity struct {
	// Start index of the entity in the message.
	Start int `json:"start"`
	// End index of the entity in the message.
	End int `json:"end"`
	// Value is the extracted value for this entity.
	Value interface{} `json:"value"`
	// Name of the extracted entity.
	Name string `json:"entity"`
	// Confidence of the entity extractory.
	Confidence float64 `json:"confidence"`
	// Extractor is the name of the extractor which extracted the entity.
	Extractor string `json:"extractor"`
}

// Bot represents bot messages to the user within a conversation.
type Bot struct {
	Base
	// Text of the message.
	Text string `json:"text"`
	// Data which is part of the message.
	Data responses.Message `json:"data"`
}

func (*Bot) EventType() Type { return bot }

// UserUtteranceReverted is an event which reverts the last user message in the conversation history.
type UserUtteranceReverted struct {
	Base
}

func (*UserUtteranceReverted) EventType() Type { return userUtteranceReverted }

// UserUtteranceReverted is an event which reverts the last bot actions in the conversation history until the last
// user message.
type ActionReverted struct {
	Base
}

func (*ActionReverted) EventType() Type { return actionReverted }

// Restarted symbolizes a conversation restart.
type Restarted struct {
	Base
}

func (*Restarted) EventType() Type { return restarted }

// StoryExported instructs Rasa Open Source to dump the current conversation to a file.
type StoryExported struct {
	Base
}

func (*StoryExported) EventType() Type { return storyExported }

// FollowUpAction forces Rasa Open Source to execute a specific action next.
type FollowUpAction struct {
	Base
	Name string `json:"name"`
}

func (*FollowUpAction) EventType() Type { return followUpAction }

// ConversationPaused pauses the conversation until there is another user message.
type ConversationPaused struct {
	Base
}

func (*ConversationPaused) EventType() Type { return conversationPaused }

// ConversationResumed resumes the conversation.
type ConversationResumed struct {
	Base
}

func (*ConversationResumed) EventType() Type { return conversationResumed }

// SlotSet saves information in the conversation history and can be used to direct the story flow.
type SlotSet struct {
	Base
	// Name of the slot.
	Name string `json:"name"`
	// Value of the slot.
	Value interface{} `json:"value"`
}

func (*SlotSet) EventType() Type { return slotSet }

// AllotSlotsReset sets all slots to `nil`.
type AllSlotsReset struct {
	Base
}

func (*AllSlotsReset) EventType() Type { return allSlotsReset }

// Form is an event which states that a form (https://rasa.com/docs/rasa/core/forms/) was activated or deactivated.
type Form struct {
	Base
	// Name of the form if activated. Empty if the currently active form was deactivated.
	Name string `json:"name,omitempty"`
}

func (*Form) EventType() Type { return form }

// FormValidation instructs the form to validate or not.
type FormValidation struct {
	Base
	// Validate if potential slot candidates. If `false` don't validate slot candidates..
	Validate bool `json:"validate"`
}

func (*FormValidation) EventType() Type { return formValidation }

// ActionExecutionReject tells Rasa Open Source that the execution of the action failed so that other policies can
// be chosen to execute a different action.
type ActionExecutionRejected struct {
	Action
}

func (*ActionExecutionRejected) EventType() Type { return actionExecutionRejected }
