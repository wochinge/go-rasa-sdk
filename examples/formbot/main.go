// Package main contains the Go equivalent of the Rasa example for forms
// (https://github.com/RasaHQ/rasa/tree/master/examples/formbot).
package main

import (
	"github.com/wochinge/go-rasa-sdk/actions/forms"
	"github.com/wochinge/go-rasa-sdk/rasa"
	"github.com/wochinge/go-rasa-sdk/rasa/events"
	"github.com/wochinge/go-rasa-sdk/rasa/responses"
	"github.com/wochinge/go-rasa-sdk/server"
	"strconv"
	"strings"
)

// CuisineValidator validates if the provided cuisine type is a valid choice.
type CuisineValidator struct{}

func (v *CuisineValidator) IsValid(value interface{}, _ *rasa.Domain, _ *rasa.Tracker,
	dispatcher responses.ResponseDispatcher) (interface{}, bool) {
	validCuisines := []string{"caribbean", "chinese", "french", "greek", "indian", "italian", "mexican"}

	cuisine, isString := value.(string)
	if !isString {
		return nil, false
	}

	for _, valid := range validCuisines {
		if strings.ToLower(cuisine) == valid {
			return cuisine, true
		}
	}

	dispatcher.Utter(responses.Message{Template: "utter_wrong_cuisine"})

	return nil, false
}

// NumPeopleValidator validates if the provided number of people for the restaurant reservation is valid.
type NumPeopleValidator struct{}

func (v *NumPeopleValidator) IsValid(value interface{}, _ *rasa.Domain, _ *rasa.Tracker,
	dispatcher responses.ResponseDispatcher) (interface{}, bool) {

	var people int
	var err error

	switch v := value.(type) {
	case string:
		people, err = strconv.Atoi(v)
		if err != nil {
			return nil, false
		}
	case int:
		people = v
	}

	if people < 1 {
		dispatcher.Utter(responses.Message{Template: "utter_wrong_num_people"})
		return nil, false
	}

	return people, true
}

// OutdoorSeatingValidator validates the answer of the user whether they want to sit outside.
type OutdoorSeatingValidator struct{}

func (v *OutdoorSeatingValidator) IsValid(value interface{}, _ *rasa.Domain, _ *rasa.Tracker,
	dispatcher responses.ResponseDispatcher) (interface{}, bool) {

	switch v := value.(type) {
	case bool:
		return true, true
	case string:
		if strings.Contains(v, "out") {
			return true, true
		} else if strings.Contains(v, "in") {
			return false, true
		}
	}

	dispatcher.Utter(responses.Message{Template: "utter_wrong_outdoor_seating"})
	return nil, false
}

func main() {
	form := forms.Form{
		FormName:      "restaurant_form",
		RequiredSlots: []string{"cuisine", "num_people", "outdoor_seating", "preferences", "feedback"},
		SlotMappings: map[string][]forms.SlotMapping{
			"cuisine": {{FromEntity: "cuisine", ExcludedIntents: []string{"chitchat"}}},
			"num_people": {
				{FromEntity: "num_people", Intents: []string{"inform", "request_restaurant"}},
				{FromEntity: "number"}},
			"outdoor_seating": {
				{FromEntity: "seating"},
				{Intents: []string{"affirm"}, Value: true},
				{Intents: []string{"deny"}, Value: false}},
			"preferences": {
				{Intents: []string{"deny"}, Value: "no additional preferences"},
				{FromText: true, ExcludedIntents: []string{"affirm"}}},
			"feedback": {{FromEntity: "feedback"}, {FromText: true}}},
		Validators: map[string][]forms.SlotValidator{
			"cuisine":         {&CuisineValidator{}},
			"num_people":      {&NumPeopleValidator{}},
			"outdoor_seating": {&OutdoorSeatingValidator{}},
		},
		OnSubmit: func(_ *rasa.Tracker, _ *rasa.Domain, dispatcher responses.ResponseDispatcher) []events.Event {
			dispatcher.Utter(responses.Message{Template: "utter_submit"})
			return []events.Event{}
		},
	}

	server.Serve(server.DefaultPort, &form)
}
