/*
Package go-rasa-sdk is a Go implementation of the Rasa SDK (https://github.com/rasahq/rasa-sdk).

It can be used to implement custom actions (https://rasa.com/docs/rasa/core/actions/#custom-actions) and
forms (https://rasa.com/docs/rasa/core/forms/) in Go.

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
		tracker *rasa.Tracker,
		domain *rasa.Domain,
		dispatcher responses.ResponseDispatcher) []events.Event {

		// Your action code goes here

		// Sending a message to the user
		dispatcher.Utter(responses.Message{Text: "Hello world from the go-rasa-sdk!!"})

		// Adding events to the conversation, e.g. to restart the conversation:
		return []events.Event{&events.Restarted{}}
	}

	func (action *HelloWorldAction) Name() string {
		// the name of your action which should be used in your stories and in the `domain.yml`
		return "action_hello_world"
	}

	// main runs the custom action server on port 5055.
	func main() {
		server.Serve(server.DefaultPort, &HelloWorldAction{})
	}

For more information and examples visit https://github.com/wochinge/go-rasa-sdk.
*/
package go_rasa_sdk // nolint
