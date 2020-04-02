package rasa

import (
	"github.com/stretchr/testify/assert"
	"github.com/wochinge/go-rasa-sdk/rasa/events"
	"github.com/wochinge/go-rasa-sdk/rasa/responses"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParsedMinimalRequest(t *testing.T) {
	parsed, err := parsedJSON("minimal_request.json")
	assert.Nil(t, err)

	expectedAction := "action_hello_world"
	assert.Equal(t, parsed.ActionToRun, expectedAction)
}

func TestParsedInvalidEvents(t *testing.T) {
	invalidJSON := `{"tracker": {"events": [[]]}}`

	parsed, err := Parsed(strings.NewReader(invalidJSON))

	assert.NotNil(t, err)
	assert.Empty(t, parsed.Tracker.Events)
}

func TestParsedDomainActions(t *testing.T) {
	parsed, err := parsedJSON("request_With_domain.json")
	assert.Nil(t, err)

	domain := parsed.Domain

	expectedForms := [3]string{"sales_form", "subscribe_newsletter_form", "suggestion_form"}
	assert.ElementsMatch(t, domain.Forms, expectedForms)

	expectedActions := []string{"action_chitchat",
		"action_default_ask_affirmation",
		"action_default_fallback",
		"respond_out_of_scope",
		"utter_already_subscribed",
		"utter_also_explain_core"}
	assert.ElementsMatch(t, domain.Actions, expectedActions)
}

func TestParsedDomainConfig(t *testing.T) {
	parsed, err := parsedJSON("request_With_domain.json")
	assert.Nil(t, err)

	domain := parsed.Domain

	expectedSessionConfig := SessionConfig{123.45, true}
	assert.Equal(t, domain.SessionConfig, expectedSessionConfig)

	expectedConfig := Config{true}
	assert.Equal(t, domain.Config, expectedConfig)
}

func TestParsedDomainIntents(t *testing.T) {
	parsed, err := parsedJSON("request_With_domain.json")
	assert.Nil(t, err)

	domain := parsed.Domain

	var actualIntents []string

	for _, intent := range domain.Intents {
		for key := range intent {
			actualIntents = append(actualIntents, key)
		}
	}

	expectedIntents := []string{"affirm", "ask_builder", "enter_data", "out_of_scope"}
	assert.ElementsMatch(t, actualIntents, expectedIntents)
}

func TestParsedDomainSlots(t *testing.T) {
	parsed, err := parsedJSON("request_With_domain.json")
	assert.Nil(t, err)

	domain := parsed.Domain

	expectedSlots := []Slot{
		{Name: "budget", Type: "rasa.core.slots.UnfeaturizedSlot", AutoFill: true},
		{Name: "current_api", Type: "rasa.core.slots.CategoricalSlot", AutoFill: true},
		{Name: "name", Type: "rasa.core.slots.TextSlot", AutoFill: true},
		{Name: "onboarding", Type: "rasa.core.slots.BooleanSlot", AutoFill: true}}
	assert.ElementsMatch(t, domain.Slots, expectedSlots)
}

func TestParsedDomainResponses(t *testing.T) {
	parsed, err := parsedJSON("request_With_domain.json")
	assert.Nil(t, err)

	domain := parsed.Domain

	expectedResponses := map[string][]Response{
		"utter_already_subscribed": {{Text: "spam folder 🗑"}},
		"utter_ask_docs_help": {{"Did that help?", "",
			[]responses.Button{{Title: "👍", PayLoad: `/affirm`}, {Title: "👎", PayLoad: `/deny`}}}},
		"utter_continue_step2": {
			{Text: "Let's continue", Channel: "socketio"},
			{Text: "Let's continue, please click the button below.",
				Buttons: []responses.Button{{Title: "Next step", PayLoad: `/get_started_step2`}}}},
	}

	for responseKey, response := range expectedResponses {
		assert.ElementsMatch(t, domain.Responses[responseKey], response)
	}
}

func TestParsedSmallTracker(t *testing.T) {
	parsed, err := parsedJSON("request_with_small_tracker.json")
	assert.Nil(t, err)

	assert.Equal(t, "wochinge", parsed.Tracker.ConversationID)
	assert.Equal(t, false, parsed.Tracker.Paused)
	assert.Equal(t, "rasa", parsed.Tracker.LatestInputChannel)
	assert.Equal(t, 1584966507.4803030491, parsed.Tracker.LatestEventTime)
	assert.Equal(t, "", parsed.Tracker.FollowUpAction)
	assert.Equal(t, "action_listen", parsed.Tracker.LatestActionName)

	expectedLatestMessage := events.ParseData{
		Intent: events.IntentParseResult{Name: "ask_howold", Confidence: 0.7406903505}, Entities: []events.Entity{},
		IntentRanking: []events.IntentParseResult{{Name: "ask_howold", Confidence: 0.7406903505}}}
	assert.Equal(t, expectedLatestMessage, parsed.Tracker.LatestMessage)

	expectedSlots := map[string]interface{}{
		"job_function": "nurse",
		"use_case":     true,
	}
	assert.Equal(t, expectedSlots, parsed.Tracker.Slots)
}

func TestParseTrackerEvents(t *testing.T) {
	parsed, err := parsedJSON("request_with_tracker_containing_all_events.json")
	assert.Nil(t, err)

	expectedEvents := []events.Event{
		&events.Action{Base: events.Base{Type: "action", Timestamp: 1584966507.4802880287},
			Name: "action_session_start"},
		&events.User{Base: events.Base{Type: "user",
			Timestamp: 1585158505.1458339691,
			Metadata:  map[string]interface{}{"rasa_x_flagged": false, "rasa_x_id": 4.0}},
			Text: "hello", ParseData: events.ParseData{
				Intent: events.IntentParseResult{Name: "greet", Confidence: 0.9908843637},
				Entities: []events.Entity{
					{Start: 0, End: 13, Value: "Windows Linux", Name: "name", Confidence: 0.7906, Extractor: "ner_crf"},
					{Start: 0, End: 13, Value: 5.0, Name: "number", Confidence: 0.7906, Extractor: "ner_crf"},
					{Start: 0, End: 13, Value: true, Name: "isHot", Confidence: 0.7906, Extractor: "ner_crf"}},
				IntentRanking: []events.IntentParseResult{{Name: "greet", Confidence: 0.9908843637},
					{Name: "mood_deny", Confidence: 0.01}}, Text: "hello"},
			MessageID: "c25928b830814f8180336745d9ad29f2", InputChannel: "rasa"},
		&events.Bot{Base: events.Base{Type: "bot", Timestamp: 1234}, Text: "Peace",
			Data: responses.BotMessage{Elements: []interface{}{}, Buttons: []responses.Button{}, Attachment: nil}},
		&events.SessionStarted{Base: events.Base{Type: "session_started", Timestamp: 1584966507.4802930355}},
		&events.SlotSet{Base: events.Base{Type: "slot", Timestamp: 1560425053.3079407215}, Name: "name", Value: "test"},
		&events.ConversationPaused{Base: events.Base{Type: "pause", Timestamp: 99.1}},
		&events.ConversationResumed{Base: events.Base{Type: "resume", Timestamp: 99.1}},
		&events.Form{Base: events.Base{Type: "form", Timestamp: 1556550828.3499741554}},
		&events.FormValidation{Base: events.Base{Type: "form_validation", Timestamp: 12345}, Validate: false},
		&events.FollowUpAction{Base: events.Base{Type: "followup", Timestamp: 99.1}, Name: "next action"},
		&events.StoryExported{Base: events.Base{Type: "export", Timestamp: 99.1}},
		&events.ActionReverted{Base: events.Base{Type: "undo", Timestamp: 99.1}},
		&events.UserUtteranceReverted{Base: events.Base{Type: "rewind", Timestamp: 99.1}},
		&events.Restarted{Base: events.Base{Type: "restart", Timestamp: 1560424318.9264261723}},
		&events.AllSlotsReset{Base: events.Base{Type: "reset_slots", Timestamp: 1560424318.9264261723}},
		&events.ActionExecutionRejected{
			Action: events.Action{
				Base:   events.Base{Type: "action_execution_rejected", Timestamp: 1556550399.6700005531},
				Policy: "policy_3_FormPolicy", Confidence: 1.0, Name: "subscribe_newsletter_form"},
		},
	}

	for i, e := range expectedEvents {
		assert.Equal(t, e, parsed.Tracker.Events[i])
	}

	assert.ElementsMatch(t, parsed.Tracker.Events, expectedEvents)
}

func TestParsedActiveForm(t *testing.T) {
	parsed, err := parsedJSON("request_with_active_form.json")
	assert.Nil(t, err)

	assert.Equal(t, parsed.Tracker.ActiveForm, ActiveForm{Name: "my-form", Validate: true,
		Rejected: false, TriggerMessage: events.ParseData{}})
}

func parsedJSON(path string) (CustomActionRequest, error) {
	const testDataDir = "../test/"
	fullPath := filepath.Join(testDataDir, path)
	reader, _ := os.Open(fullPath)

	return Parsed(reader)
}
