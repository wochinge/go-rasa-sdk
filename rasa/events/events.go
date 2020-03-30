package events

import (
	"encoding/json"
	"github.com/wochinge/go-rasa-sdk/rasa/responses"
)


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

// TODO: consider merging this with the constants declarations above
func eventParser(base Base) func() Event {
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

	return eventParsers[base.Type]
}

func Parsed(rawEvents []json.RawMessage) ([]Event, error) {
	var events []Event
	for _, rawEvent := range rawEvents {
		var minimalEvent Base

		if err := json.Unmarshal(rawEvent, &minimalEvent); err != nil {
			return []Event{}, err
		}

		event := eventParser(minimalEvent)()

		if event != nil {
			if err := json.Unmarshal(rawEvent, &event); err != nil {
				return []Event{}, err
			}
			events = append(events, event)
		} else {
			events = append(events, minimalEvent)
		}
	}

	return events, nil
}

type Event interface {
}

type Base struct {
	Type      Type                   `json:"event,omitempty"`
	Timestamp float64                `json:"timestamp,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

type Action struct {
	Base
	Policy     string  `json:"policy"`
	Confidence float64 `json:"confidence"`
	Name       string  `json:"name"`
}

type SessionStarted struct {
	Base
}

func StartSession() *SessionStarted {
	return &SessionStarted{Base{Type:sessionStarted}}
}

type User struct {
	Base
	Text         string    `json:"text"`
	InputChannel string    `json:"input_channel"`
	MessageId    string    `json:"message_id"`
	ParseData    ParseData `json:"parse_data"`
}

type ParseData struct {
	Intent        IntentParseResult   `json:"intent"`
	Entities      []Entity      `json:"entities"`
	IntentRanking []IntentParseResult `json:"intent_ranking"`
	Text          string              `json:"text"`
}

type IntentParseResult struct {
	Name       string  `json:"name"`
	Confidence float64 `json:"confidence"`
}

type Entity struct {
	Start int `json:"start"`
	End int `json:"end"`
	Value string `json:"value"`
	Name string `json:"entity"`
	Confidence float64 `json:"confidence"`
	Extractor string `json:"extractor"`
}

type Bot struct {
	Base
	Text string               `json:"text"`
	Data responses.BotMessage `json:"data"`
}

type UserUtteranceReverted struct {
	Base
}

type ActionReverted struct {
	Base
}

type Restarted struct {
	Base
}

type StoryExported struct {
	Base
}

type FollowUpAction struct {
	Base
	Name string `json:"name"`
}

type ConversationPaused struct {
	Base
}

type ConversationResumed struct {
	Base
}

type SlotSet struct {
	Base
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

type AllSlotsReset struct {
	Base
}

type Form struct {
	Base
	Name string `json:"name"`
}

type FormValidation struct {
	Base
	Validate bool `json:"validate"`
}

type ActionExecutionRejected struct {
	Action
}
