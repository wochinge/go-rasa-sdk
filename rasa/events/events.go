package events

import (
	"encoding/json"
)


// TODO: consider merging this with the constants declarations above
func eventParser(base Base) (func() Event, bool) {
	eventParsers := map[Type]func() Event{
		action:         func() Event { return &Action{Base: base} },
		user:           func() Event { return &User{Base: base} },
		bot:            func() Event { return &Bot{Base: base} },
		slotSet:        func() Event { return &Slot{Base: base} },
		actionExecutionRejected: func() Event { return &ActionExecutionRejected{Action: Action{Base: base}} },

		form:                    func() Event { return &Form{Base: base} },
		formValidation:          func() Event { return &FormValidation{Base: base} },

		followUpAction: func() Event { return &FollowUpAction{Base: base} },
	}

	value, ok := eventParsers[base.Type]
	return value, ok
}

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


func StartNewSession() Event {
	return &Base{Type:sessionStarted}
}


func RevertUserUtterance() Event {
	return &Base{Type:userUtteranceReverted}
}

func RevertActions() Event {
	return &Base{Type:actionReverted}
}

func Restart() Event {
	return &Base{Type:restarted}
}

func ExportStory() Event {
	return &Base{Type:storyExported}
}


func Pause() Event {
	return &Base{Type:conversationPaused}
}

func Resume() Event {
	return &Base{Type:conversationResumed}
}

func SetSlot(name string, value interface{}) Event {
	return &Slot{Base: Base{Type: slotSet}, Name:name, Value:value}
}

func ResetSlots() Event {
	return &Base{Type:allSlotsReset}
}
