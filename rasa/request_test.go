package rasa

import (
	"github.com/stretchr/testify/assert"
	"github.com/wochinge/go-rasa-sdk/rasa/events"
	"github.com/wochinge/go-rasa-sdk/rasa/responses"
	"os"
	"path/filepath"
	"testing"
)

func TestParsedMinimalRequest(t *testing.T) {
	path := filepath.Join("testdata", "../../test/minimal_request.json") // relative path
	reader, err := os.Open(path)
	parsed, err := Parsed(reader)

	if err != nil {
		t.Fatalf("Parsing the JSON failed with %s", err)
	}

	expectedAction := "action_hello_world"
	assert.Equal(t, parsed.ActionToRun, expectedAction)
	// TODO: Check for empty tracker and domain
}

func TestParsedDomain(t *testing.T) {
	path := filepath.Join("testdata", "../../test/request_with_domain.json") // relative path
	reader, err := os.Open(path)
	parsed, err := Parsed(reader)

	if err != nil {
		t.Fatalf("Parsing the JSON failed with %s", err)
	}

	domain := parsed.Domain

	expectedForms := [3]string{"sales_form", "subscribe_newsletter_form", "suggestion_form"}
	assert.ElementsMatch(t, domain.Forms, expectedForms)

	expectedActions := []string{"action_chitchat",
		"action_default_ask_affirmation",
		"action_default_fallback",
		"respond_out_of_scope",
		"utter_already_subscribed",
		"utter_also_explain_core"}
	assert.ElementsMatch(t, domain.Actions, expectedActions, )

	expectedSessionConfig := SessionConfig{123.45, true}
	assert.Equal(t, domain.SessionConfig, expectedSessionConfig)

	expectedConfig := Config{true}
	assert.Equal(t, domain.Config, expectedConfig)

	var actualIntents []string
	for _, intent := range domain.Intents {
		for key := range intent {
			actualIntents = append(actualIntents, key)
		}
	}
	expectedIntents := []string{"affirm", "ask_builder", "enter_data", "out_of_scope"}
	assert.ElementsMatch(t, actualIntents, expectedIntents)

	expectedSlots := []Slot{
		{"budget", "rasa.core.slots.UnfeaturizedSlot", nil, true},
		{"current_api", "rasa.core.slots.CategoricalSlot", nil, true},
		{"name", "rasa.core.slots.TextSlot", nil, true},
		{"onboarding", "rasa.core.slots.BooleanSlot", nil, true}}
	assert.ElementsMatch(t, domain.Slots, expectedSlots)

	expectedResponses := map[string][]Response{
		"utter_already_subscribed": {{Text: "spam folder 🗑"}},
		"utter_ask_docs_help": {{"Did that help?", "",
			[]responses.Button{{"👍", `/affirm`}, {"👎", `/deny`}}}},
		"utter_continue_step2": {
			{Text: "Let's continue", Channel: "socketio"},
			{"Let's continue, please click the button below.", "", []responses.Button{{"Next step", `/get_started_step2`}}}},
	}

	for responseKey, responses := range expectedResponses {
		assert.ElementsMatch(t, domain.Responses[responseKey], responses)
	}
}

func TestParsedSmallTracker(t *testing.T) {
	path := filepath.Join("testdata", "../../test/request_with_small_tracker.json") // relative path
	reader, err := os.Open(path)
	parsed, err := Parsed(reader)

	if err != nil {
		t.Fatalf("Parsing the JSON failed with %s", err)
	}

	assert.Equal(t, "wochinge", parsed.Tracker.ConversationId,)
	assert.Equal(t, false, parsed.Tracker.Paused, )
	assert.Equal(t, "rasa", parsed.Tracker.LatestInputChannel)
	assert.Equal(t, 1584966507.4803030491, parsed.Tracker.LatestEventTime)
	assert.Equal(t, "", parsed.Tracker.FollowUpAction, )
	assert.Equal(t, "action_listen", parsed.Tracker.LatestActionName,)
	expectedLatestMessage := events.ParseData{
		Intent: events.IntentParseResult{Name: "ask_howold", Confidence: 0.7406903505}, Entities: []events.Entity{},
		IntentRanking: []events.IntentParseResult{{"ask_howold", 0.7406903505}},}
	assert.Equal(t, expectedLatestMessage, parsed.Tracker.LatestMessage)

	expectedSlots := map[string]interface{}{
		"job_function": "nurse",
		"use_case": true,
	}
	assert.Equal(t, expectedSlots, parsed.Tracker.Slots)
	expectedEvents := []events.Event{
		&events.Action{Base: events.Base{Type: "action", Timestamp: 1584966507.4802880287}, Name: "action_session_start"},
		&events.Base{Type: "session_started", Timestamp: 1584966507.4802930355},
		&events.Action{Base: events.Base{Type: "action", Timestamp: 1584966507.4803030491}, Name: "action_listen"},
		&events.User{Base: events.Base{Type: "user",
			Timestamp: 1585158505.1458339691,
			Metadata:  map[string]interface{}{"rasa_x_flagged": false, "rasa_x_id": 4.0}},
			Text: "hello", ParseData: events.ParseData{
				Intent: events.IntentParseResult{Name: "greet", Confidence: 0.9908843637}, Entities: []events.Entity{},
				IntentRanking: []events.IntentParseResult{{"greet", 0.9908843637},
					{"mood_deny", 0.004441225}}, Text: "hello"}, MessageId: "c25928b830814f8180336745d9ad29f2", InputChannel: "rasa"},
	}

	assert.ElementsMatch(t, parsed.Tracker.Events, expectedEvents)
}

func TestParseTrackerEvents(t *testing.T) {
	path := filepath.Join("testdata", "../../test/request_with_tracker_containing_all_events.json") // relative path
	reader, err := os.Open(path)
	parsed, err := Parsed(reader)

	if err != nil {
		t.Fatalf("Parsing the JSON failed with %s", err)
	}

	expectedEvents := []events.Event{
		&events.Action{Base: events.Base{Type: "action", Timestamp: 1584966507.4802880287}, Name: "action_session_start"},
		&events.User{Base: events.Base{Type: "user",
			Timestamp: 1585158505.1458339691,
			Metadata:  map[string]interface{}{"rasa_x_flagged": false, "rasa_x_id": 4.0}},
			Text: "hello", ParseData: events.ParseData{
				Intent: events.IntentParseResult{Name: "greet", Confidence: 0.9908843637},
				Entities: []events.Entity{{Start: 0, End: 13, Value: "Windows Linux", Name: "name", Confidence: 0.7906183672, Extractor: "ner_crf"}},
				IntentRanking: []events.IntentParseResult{{"greet", 0.9908843637},
					{"mood_deny", 0.004441225}}, Text: "hello"}, MessageId: "c25928b830814f8180336745d9ad29f2", InputChannel: "rasa"},
		&events.Bot{Base: events.Base{Type: "bot", Timestamp: 1545048302.4110603333}, Text: "Peace", Data: responses.BotMessage{Elements: []interface{}{}, Buttons: []responses.Button{}, Attachment: nil}},
		&events.Base{Type: "session_started", Timestamp: 1584966507.4802930355},
		&events.Slot{Base: events.Base{Type: "slot", Timestamp: 1560425053.3079407215}, Name: "name", Value: "test"},
		&events.Base{Type: "pause", Timestamp: 1560425075.5327758789},
		&events.Base{Type: "resume", Timestamp: 1560425075.5327758789},
		&events.Form{Base: events.Base{Type: "form", Timestamp: 1556550828.3499741554}},
		&events.FormValidation{Base: events.Base{Type: "form_validation", Timestamp: 1556550812.2503328323}, Validate: false},
		&events.FollowUpAction{Base: events.Base{Type: "followup", Timestamp: 1560425075.5327758789}, Name: "next action"},
		&events.Base{Type: "export", Timestamp: 1560425075.5327758789},
		&events.Base{Type: "undo", Timestamp: 1560425075.5327758789},
		&events.Base{Type: "rewind", Timestamp: 1560425075.5327758789},
		&events.Base{Type: "restart", Timestamp: 1560424318.9264261723},
		&events.Base{Type: "reset_slots", Timestamp: 1560424318.9264261723},
		&events.ActionExecutionRejected{
			Action: events.Action{
				Base:   events.Base{Type: "action_execution_rejected", Timestamp: 1556550399.6700005531},
				Policy: "policy_3_FormPolicy", Confidence: 1.0, Name: "subscribe_newsletter_form"},
		},
	}

	assert.ElementsMatch(t, parsed.Tracker.Events, expectedEvents)
}

// TODO Test active form field
