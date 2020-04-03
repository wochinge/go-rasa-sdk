# go-rasa-sdk

Go implementation of [Rasa Python SDK](https://github.com/rasahq/rasa-sdk). 
Use this SDK to implement [custom actions](https://rasa.com/docs/rasa/core/actions/#custom-actions) for 
Rasa Open Source (>= 1.0).

## Installation

To install the SDK run

```bash
go get github.com/wochinge/go-rasa-sdk
```

## Usage

See the `examples` directory for an `action_hello_world` example as well as the Go implementation of the 
[Rasa formbot example](https://github.com/RasaHQ/rasa/tree/master/examples/formbot).

### Implementing an Action

To implement a custom action, you have to implement two functions in order to suffice the `Action` interface:
```go
type Action interface {
    
    // Run runs the custom action in the given context and returns new conversation events.
    // Any messages dispatched will be sent to the user.
    Run(tracker *rasa.Tracker, 
        domain *rasa.Domain, 
        dispatcher responses.ResponseDispatcher,
        ) []events.Event
    
    // Name returns the name of the custom action.
    Name() string
}
```

E.g. to implement an `Action` which sends a message `Hello` to the user and set a slot `user_was_greeted`:

```go
import (
    "github.com/wochinge/go-rasa-sdk/rasa"
    "github.com/wochinge/go-rasa-sdk/rasa/events"
    "github.com/wochinge/go-rasa-sdk/rasa/responses"
)

type GreetAction struct{}

func (action *GreetAction) Run(
    _ *rasa.Tracker, // the tracker containing the conversation history
    _ *rasa.Domain,  // the domain of the currently loaded model in Rasa
    dispatcher responses.ResponseDispatcher, // a dispatcher to send messages to the user
) []events.Event {
    
    // Your action code goes here
    
    // Dispatching the message
    dispatcher.Utter(responses.Message{Text: "Hello"})
    
    // See all possible events to return in github.com/wochinge/go-rasa-sdk/rasa/events
    return []events.Event{&events.SlotSet{Name: "user_was_greeted", Value: true}}
}

func (action *GreetAction) Name() string {
	// the name of your action which should be used in your stories and in the `domain.yml`
	return "action_hello_world"
}
```

To run the action server on port `5055` with your implemented action:

```go
import (
    "github.com/wochinge/go-rasa-sdk/server"
)

func main() {
    // Service is variadic function and accepts multiple actions as argument
	server.Serve(server.DefaultPort, &GreetAction{})
}


```

### Implementing a Form
The `go-rasa-sdk` also provides support for 
[Rasa Open Source forms](https://rasa.com/docs/rasa/core/forms/). Implement a form using `Form` struct. To implement
a form which fills an `age` slot:

```go
import (
    "github.com/wochinge/go-rasa-sdk/actions/forms"
    "github.com/wochinge/go-rasa-sdk/rasa"
    "github.com/wochinge/go-rasa-sdk/rasa/events"
    "github.com/wochinge/go-rasa-sdk/rasa/responses"
)

func main() {
    ageForm := forms.Form{
        // the name of your form which should be specified in the `forms` section
        // in your `domain.yml`
        FormName: "age_form",

        // Slots your form should fill
        RequiredSlots: []string{"age"},

        // Defines how slots are filled
        SlotMappings: map[string][]forms.SlotMapping{
            // the age slot is filled by an entity with the name `number`
            "age": {{FromEntity: "number"}}},

        // Validators for slot candidates
        Validators: map[string][]forms.SlotValidator{
            // AgeValidator will validate that the age is not a negative number
            "age": {&AgeValidator{}},
        },

        // OnSubmit is triggered when the form filled all required slots are filled
        OnSubmit: func(_ *rasa.Tracker, _ *rasa.Domain, dispatcher responses.ResponseDispatcher) []events.Event {
            // We tell the user that the age was successfully provided
            dispatcher.Utter(responses.Message{Template: "utter_age_provided"})
            return []events.Event{}
        },
    }
}
```

You can provide multiple `Validators` for each slot. If no `Validator` is given, candidates will only be required to 
be not `nil`. To implement a `Validator` which validates that a given value is greater 0:

```go
import (
    "github.com/wochinge/go-rasa-sdk/rasa"
    "github.com/wochinge/go-rasa-sdk/rasa/events"
    "github.com/wochinge/go-rasa-sdk/rasa/responses"
)

type AgeValidator struct{}

func (v *AgeValidator) IsValid(value interface{}, _ *rasa.Domain, _ *rasa.Tracker,
    dispatcher responses.ResponseDispatcher) (interface{}, bool) {
    
    if age, isInt := value.(int); ! isInt || age <= 0 {
        return nil, false
    }
    
    return value, true
}
```  

To run the Go action server with your form loaded:

```go
import (
    "github.com/wochinge/go-rasa-sdk/server"
)

func main() {
	server.Serve(server.DefaultPort, &ageForm)
}
```

## Code Style

### Formatting

The built-in tool `gofmt` is used to format the code. Use `gofmt -s` to check if formatting is required. To format
and apply the changes run `gofmt -s -w`.  

### Linting
This repository uses [golangci-lint](https://github.com/golangci/golangci-lint) to lint the code base.
This will run a bundle of Go linters.

To lint the code:

1. Follow the installation instructions of [golangci-lint](https://github.com/golangci/golangci-lint)
2. Run `golangci-lint run` in the cloned repository
