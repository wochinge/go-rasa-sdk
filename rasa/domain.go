// Package rasa handles data sent by Rasa Open Source.
package rasa

import (
	"github.com/wochinge/go-rasa-sdk/rasa/responses"
)

// Domain represents the currently loaded bot's domain.
type Domain struct {
	// Entities which are specified in the `domain.yml`.
	Entities []string `json:"entities"`
	// Actions which are specified in the `domain.yml`.
	Actions []string `json:"actions"`
	// Forms which are specified in the `domain.yml`.
	Forms []string `json:"forms"`
	// Intents which are specified in the `domain.yml`.
	Intents []DomainIntent `json:"intents"`
	// Slots which are specified in the `domain.yml`.
	Slots map[string]Slot `json:"slots"`
	// Responses which are specified in the `domain.yml`.
	Responses map[string][]Response `json:"responses"`
	// Config which is specified in the `domain.yml`.
	Config Config `json:"config"`
	// SessionConfig which is specified in the `domain.yml`.
	SessionConfig SessionConfig `json:"session_config"`
}

// DomainIntent specifies an intent description with the domain.
type DomainIntent map[string]interface{}

// Slot specifies a slot declaration in the `domain.yml`.
type Slot struct {
	// Type of the slot.
	Type string `json:"type"`
	// InitialValue of the slot.
	InitialValue interface{} `json:"initial_value"`
	// AutoFill the slot when an entity with the same name was extracted.
	AutoFill bool `json:"auto_fill"`
}

// Response represents a bot response in the `domain.yml`.
type Response struct {
	// Text of the response.
	Text string `json:"text"`
	// Channel which this response applies to.
	Channel string `json:"channel"`
	// Buttons which are part of the response.
	Buttons []responses.Button `json:"buttons"`
}

// Config to specify if entities should be stored as slots.
type Config struct {
	StoreEntitiesAsSlots bool `json:"store_entities_as_slots"`
}

// SessionConfig for https://rasa.com/docs/rasa/core/domains/#session-configuration.
type SessionConfig struct {
	// SessionExpirationTime in minutes.
	SessionExpirationTime float64 `json:"session_expiration_time"`
	// CarryOverSlotsToNewSession specifies if slots are re-applied when a new session is started.
	CarryOverSlotsToNewSession bool `json:"carry_over_slots_to_new_session"`
}
