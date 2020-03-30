package rasa

import (
	"github.com/wochinge/go-rasa-sdk/rasa/responses"
)

func sanitizeDomain(domain Domain) Domain {
	var sanitizedSlots []Slot
	for slotName, slot := range domain.RawSlots {
		slot.Name = slotName
		sanitizedSlots = append(sanitizedSlots, slot)
	}

	domain.Slots = sanitizedSlots
	return domain
}

type Domain struct {
	Entities []string `json:"entities"`
	Actions  []string `json:"actions"`
	Forms    []string `json:"forms"`
	Intents  []DomainIntent `json:"intents"`
	RawSlots map[string]Slot `json:"slots"`
	Slots []Slot
	Responses map[string][]Response  `json:"responses"`

	Config        Config        `json:"config"`
	SessionConfig SessionConfig `json:"session_config"`

}

type DomainIntent map[string]interface{}

type Slot struct {
	Name string
	Type string `json:"type"`
	InitialValue interface{} `json:"initial_value"`
	AutoFill bool `json:"auto_fill"`
}

type Response struct {
	Text string                `json:"text"`
	Channel string             `json:"channel"`
	Buttons []responses.Button `json:"buttons"`
}

type Config struct {
	StoreEntitiesAsSlots bool `json:"store_entities_as_slots"`
}

type SessionConfig struct {
	SessionExpirationTime      float64  `json:"session_expiration_time"`
	CarryOverSlotsToNewSession bool `json:"carry_over_slots_to_new_session"`
}
