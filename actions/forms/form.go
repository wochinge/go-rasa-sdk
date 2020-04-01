package forms

import (
	"fmt"
	"github.com/wochinge/go-rasa-sdk/rasa"
	"github.com/wochinge/go-rasa-sdk/rasa/events"
	"github.com/wochinge/go-rasa-sdk/rasa/responses"
)

const requestedSlot = "requested_slot"

type Form struct {
	FormName string

	RequiredSlots []string
	SlotMappings  map[string][]SlotMapping
	Validators    map[string][]SlotValidator
	OnSubmit      func(tracker *rasa.Tracker, domain *rasa.Domain, dispatcher responses.ResponseDispatcher) []events.Event
}

func (form *Form) Name() string { return form.FormName }

func (form *Form) Run(tracker *rasa.Tracker, domain *rasa.Domain, dispatcher responses.ResponseDispatcher) []events.Event {
	tracker.Init()
	var newEvents []events.Event

	if ! form.wasAlreadyActive(tracker) {
		newEvents = append(newEvents, form.activate(tracker, domain, dispatcher)...)
	}

	slotCandidates := form.slotCandidates(tracker)

	if len(slotCandidates) == 0 && form.wasAlreadyActive(tracker) {
		// Reject to execute the form action if some slot was requested but nothing was extracted.
		// This will allow other policies to predict another action.
		// TODO: This has to affect the return code
		return []events.Event{&events.ActionExecutionRejected{}}
	}

	// TODO: Check default behavior if no mapping given for slot
	// TODO: Handle weird trigger intent

	newEvents = append(newEvents, form.validatedSlots(slotCandidates, domain, tracker, dispatcher)...)

	if nextSlot, allSlotsFilled := form.nextSlotToRequest(tracker); allSlotsFilled {
		newEvents = append(newEvents, form.submit(tracker, domain, dispatcher)...)
		newEvents = append(newEvents, form.deactivate()...)
		return newEvents
	} else {
		return append(newEvents, requestSlot(nextSlot, dispatcher)...)
	}
}

func (form *Form) wasAlreadyActive(tracker *rasa.Tracker) bool {
	return tracker.ActiveForm.Name == form.Name()
}

func (form *Form) activate(tracker *rasa.Tracker, domain *rasa.Domain, dispatcher responses.ResponseDispatcher) []events.Event {
	return append([]events.Event{&events.Form{Name: form.Name()}},
		form.candidatesFromExisting(tracker, domain, dispatcher)...)
}

func (form *Form) candidatesFromExisting(tracker *rasa.Tracker, domain *rasa.Domain, dispatcher responses.ResponseDispatcher) []events.Event {
	var candidates []events.SlotSet
	for _, requiredSlot := range form.RequiredSlots {
		if value, found := tracker.Slots[requiredSlot]; found {
			candidates = append(candidates, events.SlotSet{Name: requiredSlot, Value: value})
		}
	}

	return form.validatedSlots(candidates, domain, tracker, dispatcher)
}

func (form *Form) slotCandidates(tracker *rasa.Tracker) []events.SlotSet {
	requestedSlotName, ok := tracker.Slots[requestedSlot].(string)
	if !ok || requestedSlotName == "" {
		return []events.SlotSet{}
		// TODO: Handle error
	}

	otherSlotCandidates := form.fillSlotsWhichWereNotRequestedButFilledAnyway(requestedSlotName, tracker)
	requestedSlotCandidates := form.slotEventsFor(requestedSlotName, form.mappingsFor(requestedSlotName), tracker)

	return append(otherSlotCandidates, requestedSlotCandidates...)
}

func (form *Form) fillSlotsWhichWereNotRequestedButFilledAnyway(requestedSlot string, tracker *rasa.Tracker) []events.SlotSet {
	var newEvents []events.SlotSet
	for _, slotName := range form.RequiredSlots {
		if slotName == requestedSlot {
			continue
		}

		var mappings []SlotMapping

		for _, mapping := range form.mappingsFor(slotName) {
			mappings = append(mappings, SlotMapping{Intents: mapping.Intents, Entity: mapping.Entity})
		}
		// TODO: only if different from already set slot
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
			tracker.Slots[slotName] = found
			return []events.SlotSet{{Name: slotName, Value: value}}
		}
	}

	return []events.SlotSet{}
}

func (form *Form) validatedSlots(candidates []events.SlotSet, domain *rasa.Domain, tracker *rasa.Tracker, dispatcher responses.ResponseDispatcher) []events.Event {
	var slots []events.SlotSet

	if tracker.NoFormValidation() {
		return toEventInterface(candidates)
	}

	for _, candidate := range candidates {
		validator := form.validatorFor(candidate.Name)

		if validator.IsValid(candidate.Value, domain, tracker, dispatcher) {
			tracker.Slots[candidate.Name] = candidate.Value
			slots = append(slots, events.SlotSet{Name: candidate.Name, Value: candidate.Value})
		} else {
			// Nuke slot otherwise
			tracker.Slots[candidate.Name] = nil
			slots = append(slots, events.SlotSet{Name: candidate.Name, Value: nil})
		}
	}

	return toEventInterface(slots)
}

func (form *Form) validatorFor(slotName string) SlotValidator {
	validators := form.Validators[slotName]

	if validators == nil {
		return &DefaultValidator{}
	}

	return &MultiValidator{validators}
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

func (form *Form) deactivate() []events.Event {
	return []events.Event{&events.Form{Name: ""}, &events.SlotSet{Name: requestedSlot, Value: nil}}
}

func requestSlot(slotName string, dispatcher responses.ResponseDispatcher) []events.Event {
	templateNameForSlotRequest := fmt.Sprintf("utter_ask_%s", slotName)
	dispatcher.Utter(responses.BotMessage{Template: templateNameForSlotRequest})

	return []events.Event{&events.SlotSet{Name: requestedSlot, Value: slotName}}
}

func (form *Form) submit(tracker *rasa.Tracker, domain *rasa.Domain, dispatcher responses.ResponseDispatcher) []events.Event {
	if form.OnSubmit != nil {
		return form.OnSubmit(tracker, domain, dispatcher)
	}

	return []events.Event{}
}
