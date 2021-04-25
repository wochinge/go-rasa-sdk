// Package form implements the logic for Rasa Open Source forms (https://rasa.com/docs/rasa/core/forms/).
package forms

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/wochinge/go-rasa-sdk/logging"
	"github.com/wochinge/go-rasa-sdk/rasa"
	"github.com/wochinge/go-rasa-sdk/rasa/events"
	"github.com/wochinge/go-rasa-sdk/rasa/responses"
)

const requestedSlot = "requested_slot"

// Form to validate and map user input to slots (https://rasa.com/docs/rasa/core/forms/).
// Special implementation of a custom action which implements the logic required by Rasa Open Source.
type Form struct {
	// FormName is the name of the form.
	FormName string

	// RequiredSlots specifies the slots which have to be filled by this form.
	RequiredSlots []string
	// SlotMappings define how the slots are filled.
	SlotMappings map[string][]SlotMapping
	// Validators specify functions to validate slot candidates.
	Validators map[string][]SlotValidator
	// OnSubmit is a function which is run when all slots are filled and the form is complete.
	OnSubmit func(*rasa.Tracker, *rasa.Domain, responses.ResponseDispatcher) []events.Event
}

func (form *Form) Name() string { return form.FormName }

// Run is executed whenever Rasa Open Source sends a request to execute this form.
func (form *Form) Run(tracker *rasa.Tracker, domain *rasa.Domain,
	dispatcher responses.ResponseDispatcher) []events.Event {
	tracker.Init()
	log.WithFields(
		log.Fields{logging.FormNameKey: form.FormName, logging.FormValidationKey: tracker.ActiveForm.Validate}).Debug(
		"Running form.")

	var newEvents []events.Event

	if !form.wasAlreadyActive(tracker) {
		newEvents = append(newEvents, form.activate(tracker, domain, dispatcher)...)
	}

	newSlots, reject := form.newSlotEvents(tracker, domain, dispatcher)
	newEvents = append(newEvents, newSlots...)

	if reject {
		// Reject to execute the form action if some slot was requested but nothing was extracted.
		// This will allow other policies to predict another action.
		return []events.Event{&events.ActionExecutionRejected{}}
	}

	nextSlot, allSlotsFilled := form.nextSlotToRequest(tracker)
	if allSlotsFilled {
		return append(newEvents, form.deactivate(tracker, domain, dispatcher)...)
	}

	return append(newEvents, requestSlot(nextSlot, dispatcher)...)
}

func (form *Form) newSlotEvents(tracker *rasa.Tracker, domain *rasa.Domain,
	dispatcher responses.ResponseDispatcher) ([]events.Event, bool) {
	slotCandidates := form.slotCandidates(tracker)
	if len(slotCandidates) == 0 && form.wasAlreadyActive(tracker) {
		return nil, true
	}

	return form.validatedSlots(slotCandidates, domain, tracker, dispatcher), false
}

func (form *Form) wasAlreadyActive(tracker *rasa.Tracker) bool {
	return tracker.ActiveForm.Name == form.Name()
}

func (form *Form) activate(tracker *rasa.Tracker, domain *rasa.Domain,
	dispatcher responses.ResponseDispatcher) []events.Event {
	log.WithField(logging.FormNameKey, form.FormName).Debug("Activating form.")

	return append([]events.Event{&events.Form{Name: form.Name()}},
		form.candidatesFromExisting(tracker, domain, dispatcher)...)
}

func (form *Form) candidatesFromExisting(tracker *rasa.Tracker, domain *rasa.Domain,
	dispatcher responses.ResponseDispatcher) []events.Event {
	var candidates []events.SlotSet

	for _, requiredSlot := range form.RequiredSlots {
		if value, found := tracker.Slots[requiredSlot]; found && value != nil {
			candidates = append(candidates, events.SlotSet{Name: requiredSlot, Value: value})
		}
	}

	return form.validatedSlots(candidates, domain, tracker, dispatcher)
}

func (form *Form) slotCandidates(tracker *rasa.Tracker) []events.SlotSet {
	requestedSlotName, ok := tracker.Slots[requestedSlot].(string)

	candidates := form.fillProvidedButNotRequested(requestedSlotName, tracker)

	if ok && requestedSlotName != "" {
		requestedSlotCandidates := form.slotEventsFor(requestedSlotName, form.mappingsFor(requestedSlotName), tracker)
		candidates = append(candidates, requestedSlotCandidates...)
	}

	return candidates
}

func (form *Form) fillProvidedButNotRequested(requestedSlot string, tracker *rasa.Tracker) []events.SlotSet {
	var newEvents []events.SlotSet

	for _, slotName := range form.RequiredSlots {
		if slotName == requestedSlot {
			continue
		}

		var mappings []SlotMapping

		for _, mapping := range form.mappingsFor(slotName) {
			mappings = append(mappings, SlotMapping{Intents: mapping.Intents, FromEntity: mapping.FromEntity})
		}

		newEvents = append(newEvents, form.slotEventsFor(slotName, mappings, tracker)...)
	}

	return newEvents
}

func (form *Form) mappingsFor(slotName string) []SlotMapping {
	slotMappings := form.SlotMappings[slotName]

	if slotMappings == nil {
		slotMappings = defaultSlotMapping(slotName)
	}

	return slotMappings
}

func (form *Form) slotEventsFor(slotName string, mappings []SlotMapping, tracker *rasa.Tracker) []events.SlotSet {
	for _, mapping := range mappings {
		if value, found := mapping.apply(form, tracker); found {
			return []events.SlotSet{{Name: slotName, Value: value}}
		}
	}

	return []events.SlotSet{}
}

func (form *Form) validatedSlots(candidates []events.SlotSet, domain *rasa.Domain, tracker *rasa.Tracker,
	dispatcher responses.ResponseDispatcher) []events.Event {
	var slots []events.SlotSet

	for _, candidate := range candidates {
		validator := form.validatorFor(candidate.Name, tracker)

		if validated, isValid := validator.IsValid(candidate.Value, domain, tracker, dispatcher); isValid {
			tracker.Slots[candidate.Name] = candidate.Value

			slots = append(slots, events.SlotSet{Name: candidate.Name, Value: validated})
		} else {
			log.WithFields(log.Fields{logging.FormNameKey: form.FormName, logging.FormValidatedSlotKey: candidate.Name}).Debug(
				"Slot validation failed")
			tracker.Slots[candidate.Name] = nil
			slots = append(slots, events.SlotSet{Name: candidate.Name, Value: nil})
		}
	}

	return toEventInterface(slots)
}

func toEventInterface(slots []events.SlotSet) []events.Event {
	var wrapped []events.Event

	for _, slot := range slots {
		copied := slot
		wrapped = append(wrapped, &copied)
	}

	return wrapped
}

func (form *Form) nextSlotToRequest(tracker *rasa.Tracker) (string, bool) {
	for _, slot := range form.RequiredSlots {
		currentValue, found := tracker.Slots[slot]
		if !found || currentValue == nil {
			return slot, false
		}
	}

	return "", true
}

func (form *Form) deactivate(tracker *rasa.Tracker, domain *rasa.Domain,
	dispatcher responses.ResponseDispatcher) []events.Event {
	log.WithField(logging.FormNameKey, form.FormName).Debug("Deactivating form.")

	return append(
		form.submit(tracker, domain, dispatcher),
		&events.Form{Name: ""},
		&events.SlotSet{Name: requestedSlot, Value: nil})
}

func requestSlot(slotName string, dispatcher responses.ResponseDispatcher) []events.Event {
	log.WithField(requestedSlot, slotName).Debug("Requesting next slot.")

	templateNameForSlotRequest := fmt.Sprintf("utter_ask_%s", slotName)

	dispatcher.Utter(&responses.Message{Template: templateNameForSlotRequest})

	return []events.Event{&events.SlotSet{Name: requestedSlot, Value: slotName}}
}

func (form *Form) submit(tracker *rasa.Tracker, domain *rasa.Domain,
	dispatcher responses.ResponseDispatcher) []events.Event {
	if form.OnSubmit != nil {
		return form.OnSubmit(tracker, domain, dispatcher)
	}

	return []events.Event{}
}
