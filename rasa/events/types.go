package events

import "github.com/wochinge/go-rasa-sdk/rasa/responses"

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

type Event interface {}

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

type Slot struct {
	Base
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

type FollowUpAction struct {
	Base
	Name string `json:"name"`
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
