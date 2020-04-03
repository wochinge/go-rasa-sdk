// Package responses is to deal with responses which should be sent to the user during the execution of custom actions.
package responses

// Interface for dispatching messages to the users.
type ResponseDispatcher interface {
	// Utter sends a message to the user.
	Utter(Message)

	// Responses returns the messages which were dispatched and will be sent back to Rasa Open Source
	// as part of the response body.
	Responses() []Message
}

type responseDispatcher struct {
	responses []Message
}

func (dispatcher *responseDispatcher) Utter(message Message) {
	dispatcher.responses = append(dispatcher.responses, message)
}

func (dispatcher *responseDispatcher) Responses() []Message {
	return dispatcher.responses
}

// NewDispatcher returns a new `ResponseDispatcher` to send messages to the user.
func NewDispatcher() ResponseDispatcher {
	return &responseDispatcher{responses: []Message{}}
}

// Button which should be shown to the user.
type Button struct {
	// Title of the button.
	Title string `json:"title"`
	// Payload which should be sent to Rasa Open Source when this button is triggered.
	PayLoad string `json:"payload"`
}

// Message which should be sent to the user.
type Message struct {
	// Text of the message.
	Text string `json:"text"`
	// Template is the response from the `domain.yml` which should be triggered instead of the hard coding the message
	// content as part of the custom action.
	Template string `json:"template,omitempty"`
	// Elements of the message.
	Elements []interface{} `json:"elements,omitempty"`
	// QuickReplies for the user.
	QuickReplies []interface{} `json:"quick_replies,omitempty"`
	// Buttons which should be presented to the user.
	Buttons []Button `json:"buttons,omitempty"`
	// Attachment for the bot message.
	Attachment interface{} `json:"attachment,omitempty"`
	// ImageURL is the url of an image which should be shown to the user.
	ImageURL string `json:"image,omitempty"`
	// Custom can be used to send custom payloads to the user.
	Custom interface{} `json:"custom,omitempty"`
}
