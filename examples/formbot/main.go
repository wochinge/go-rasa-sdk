// Package main contains the Go equivalent of the Rasa example for forms
// (https://github.com/RasaHQ/rasa/tree/master/examples/formbot).
package main

import (
	"strconv"
	"strings"

	"github.com/wochinge/go-rasa-sdk/v2/actions/forms"
	"github.com/wochinge/go-rasa-sdk/v2/rasa"
	"github.com/wochinge/go-rasa-sdk/v2/rasa/responses"
	"github.com/wochinge/go-rasa-sdk/v2/server"
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

	dispatcher.Utter(&responses.Message{Template: "utter_wrong_cuisine"})

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
		dispatcher.Utter(&responses.Message{Template: "utter_wrong_num_people"})
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

	dispatcher.Utter(&responses.Message{Template: "utter_wrong_outdoor_seating"})
	return nil, false
}

func main() {
	form := forms.FormValidationAction{
		FormName: "restaurant_form",
		Validators: map[string]forms.SlotValidator{
			"cuisine":         &CuisineValidator{},
			"num_people":      &NumPeopleValidator{},
			"outdoor_seating": &OutdoorSeatingValidator{},
		},
	}

	server.Serve(server.DefaultPort, &form)
}
