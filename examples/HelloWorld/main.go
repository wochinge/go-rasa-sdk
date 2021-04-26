// Package main contains implementation of a hello world example to illustrate the usage of the go-rasa-sdk.
package main

import (
	"github.com/wochinge/go-rasa-sdk/v2/rasa"
	"github.com/wochinge/go-rasa-sdk/v2/rasa/events"
	"github.com/wochinge/go-rasa-sdk/v2/rasa/responses"
	"github.com/wochinge/go-rasa-sdk/v2/server"
)

// HelloWorldAction is an action which sends the user the message "Hello world from the go-rasa-sdk!!" when it's
// triggered during a conversation.
type HelloWorldAction struct{}

func (action *HelloWorldAction) Run(
	_ *rasa.Tracker, // the tracker containing the conversation history
	_ *rasa.Domain, // the domain of the currently loaded model in Rasa
	dispatcher responses.ResponseDispatcher, // a dispatcher to send messages to the user
) []events.Event {

	// Your action code goes here

	// Dispatching the message
	dispatcher.Utter(&responses.Message{Text: "Hello world from the go-rasa-sdk!!"})

	// We are not returning any events for this simple action.
	// See all possible events to return in github.com/wochinge/go-rasa-sdk/v2/rasa/events
	return []events.Event{}
}

func (action *HelloWorldAction) Name() string {
	// the name of your action which should be used in your stories and in the `domain.yml`
	return "action_hello_world"
}

func main() {
	server.Serve(server.DefaultPort, &HelloWorldAction{})
}
