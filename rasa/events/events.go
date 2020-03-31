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

func Parsed(rawEvents []json.RawMessage) ([]Event, error) {
	var events []Event
	for _, rawEvent := range rawEvents {
		var minimalEvent Base

		if err := json.Unmarshal(rawEvent, &minimalEvent); err != nil {
			return []Event{}, err
		}

		eventCreator, ok := eventParser(minimalEvent)

		if ok {
			event := eventCreator()
			if err := json.Unmarshal(rawEvent, &event); err != nil {
				return []Event{}, err
			}
			events = append(events, event)
		} else {
			events = append(events, &minimalEvent)
		}
	}

	return events, nil
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

type Event interface {
	EventType() Type
	SetType(Type)
}

type Base struct {
	Type      Type                   `json:"event,omitempty"`
	Timestamp float64                `json:"timestamp,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}
func (_ *Base) EventType() Type {return "unknown"}
func (base *Base) SetType(eventType Type) {base.Type = eventType}

type Action struct {
	Base
	Policy     string  `json:"policy"`
	Confidence float64 `json:"confidence"`
	Name       string  `json:"name"`
}
func (_ *Action) EventType() Type {return action}

type SessionStarted struct {
	Base
}
func (_ *SessionStarted) EventType() Type {return sessionStarted}


type User struct {
	Base
	Text         string    `json:"text"`
	InputChannel string    `json:"input_channel"`
	MessageId    string    `json:"message_id"`
	ParseData    ParseData `json:"parse_data"`
}
func (_ *User) EventType() Type {return user}


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
func (_ *Bot) EventType() Type {return bot}

type UserUtteranceReverted struct {
	Base
}
func (_ *UserUtteranceReverted) EventType() Type {return userUtteranceReverted}

type ActionReverted struct {
	Base
}
func (_ *ActionReverted) EventType() Type {return actionReverted}

type Restarted struct {
	Base
}
func (_ *Restarted) EventType() Type {return restarted}

type StoryExported struct {
	Base
}
func (_ *StoryExported) EventType() Type {return storyExported}

type FollowUpAction struct {
	Base
	Name string `json:"name"`
}
func (_ *FollowUpAction) EventType() Type {return followUpAction}

type ConversationPaused struct {
	Base
}
func (_ *ConversationPaused) EventType() Type {return conversationPaused}


type ConversationResumed struct {
	Base
}
func (_ *ConversationResumed) EventType() Type {return conversationResumed}

type SlotSet struct {
	Base
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}
func (_ *SlotSet) EventType() Type {return slotSet}

type AllSlotsReset struct {
	Base
}
func (_ *AllSlotsReset) EventType() Type {return allSlotsReset}

type Form struct {
	Base
	Name string `json:"name"`
}
func (_ *Form) EventType() Type {return form}

type FormValidation struct {
	Base
	Validate bool `json:"validate"`
}
func (_ *FormValidation) EventType() Type {return formValidation}


type ActionExecutionRejected struct {
	Action
}
func (_ *ActionExecutionRejected) EventType() Type {return actionExecutionRejected}
