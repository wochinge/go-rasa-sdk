// Package main contains implementation of a hello world example to illustrate the usage of the go-rasa-sdk.
package main

import (
	"github.com/wochinge/go-rasa-sdk/rasa"
	"github.com/wochinge/go-rasa-sdk/rasa/events"
	"github.com/wochinge/go-rasa-sdk/rasa/responses"
	"github.com/wochinge/go-rasa-sdk/server"
)

// HelloWorldAction is an action which sends the user the message "Hello world from the go-rasa-sdk!!" when it's
// triggered during a conversation.
type HelloWorldAction struct{}

func (action *HelloWorldAction) Run(_ *rasa.Tracker,
	_ *rasa.Domain,
	dispatcher responses.ResponseDispatcher) []events.Event {

	// Your action code goes here

	// Dispatching the message
	dispatcher.Utter(responses.Message{Text: "Hello world from the go-rasa-sdk!!"})

	// We are not returning any events for this simple action.
	// See all possible events to return in github.com/wochinge/go-rasa-sdk/rasa/events
	return []events.Event{}
}

func (action *HelloWorldAction) Name() string {
	// the name of your action which should be used in your stories and in the `domain.yml`
	return "action_hello_world"
}

func main() {
	server.Serve(server.DefaultPort, &HelloWorldAction{})
}
