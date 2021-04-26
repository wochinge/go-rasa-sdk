# go-rasa-sdk
![Build Status](https://github.com/wochinge/go-rasa-sdk/workflows/Lint%20and%20Test/badge.svg?branch=main)
[![Coverage Status](https://coveralls.io/repos/github/wochinge/go-rasa-sdk/badge.svg?branch=master)](https://coveralls.io/github/wochinge/go-rasa-sdk?branch=main)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/wochinge/go-rasa-sdk?tab=doc)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat-square)](https://github.com/wochinge/go-rasa-sdk)

Go implementation of the [Rasa Python SDK](https://github.com/rasahq/rasa-sdk). 
Use this SDK to implement [custom actions](https://rasa.com/docs/rasa/core/actions/#custom-actions) for 
Rasa Open Source (>= 2.0). Version 1 of the `go-rasa-sdk` is compatible with Rasa Open Source 1.

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
    Run(_ *rasa.Tracker, // the tracker containing the conversation history
        _ *rasa.Domain,  // the domain of the currently loaded model in Rasa
        dispatcher responses.ResponseDispatcher, // a dispatcher to send messages to the user
        ) []events.Event
    
    // Name returns the name of the custom action.
    Name() string
}
```

E.g. to implement an `Action` which sends a message `Hello` to the user and set a slot `user_was_greeted`:

```go
import (
    "github.com/wochinge/go-rasa-sdk/v2/rasa"
    "github.com/wochinge/go-rasa-sdk/v2/rasa/events"
    "github.com/wochinge/go-rasa-sdk/v2/rasa/responses"
)

type GreetAction struct{}

func (action *GreetAction) Run(
    _ *rasa.Tracker,
    _ *rasa.Domain,
    dispatcher responses.ResponseDispatcher,
) []events.Event {
    
    // Your action code goes here.
    
    // Dispatching the message.
    dispatcher.Utter(&responses.Message{Text: "Hello"})
    
    // See all possible events to return in github.com/wochinge/go-rasa-sdk/rasa/events .
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
    "github.com/wochinge/go-rasa-sdk/v2/server"
)

func main() {
    // Service is variadic function and accepts multiple actions as argument.
	server.Serve(server.DefaultPort, &GreetAction{})
}


```

### Implementing a Form
The `go-rasa-sdk` also provides support for 
[Rasa Open Source forms](https://rasa.com/docs/rasa/forms/). Implement a form using the `FormValidationAction` struct. 
To implement a form which fills an `age` slot:

```go
import (
    "github.com/wochinge/go-rasa-sdk/v2/actions/forms"
    "github.com/wochinge/go-rasa-sdk/v2/rasa"
    "github.com/wochinge/go-rasa-sdk/v2/rasa/events"
    "github.com/wochinge/go-rasa-sdk/v2/rasa/responses"
)

func main() {
    ageForm := forms.FormValidationAction{
        // the name of your form which should be specified in the `forms` section
        // in your `domain.yml`
        FormName: "age_form",
        // Validators for slot candidates
        Validators: map[string][]forms.SlotValidator{
            // AgeValidator will validate that the age is not a negative number.
            "age": {&AgeValidator{}},
        },
        // Extractors specify functions to extract slot candidates.
        Extractors: map[string]forms.SlotExtractor{}
    }
}
```

To run the Go action server with your form loaded:

```go
import (
    "github.com/wochinge/go-rasa-sdk/v2/server"
)

func main() {
	server.Serve(server.DefaultPort, &ageForm)
}
```

#### Slot Validators
You can provide multiple `Validators` for each slot. To implement a `Validator` which validates that a given value is
greater 0:

```go
import (
    "github.com/wochinge/go-rasa-sdk/v2/rasa"
    "github.com/wochinge/go-rasa-sdk/v2/rasa/events"
    "github.com/wochinge/go-rasa-sdk/v2/rasa/responses"
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

#### Extracting Custom Slots
You can provide multiple `Extractors` to extract custom slots. This is in line with what's described in the Rasa
documentation [here](https://rasa.com/docs/rasa/forms#custom-slot-mappings). To implement an `Extractor` which extracts
a slot `age` based on an entity `age`: 

```go
import (
    "github.com/wochinge/go-rasa-sdk/v2/actions/forms"
    "github.com/wochinge/go-rasa-sdk/v2/rasa"
    "github.com/wochinge/go-rasa-sdk/v2/rasa/responses"
    "github.com/wochinge/go-rasa-sdk/v2/server"
)

type AgeExtractor struct{}

func (v *AgeExtractor) Extract(_ *rasa.Domain, tracker *rasa.Tracker,
    _ responses.ResponseDispatcher) (extractedValue interface{}, valueFound bool) {

    for _, entity := range tracker.LatestMessage.Entities {
        if entity.Name == "age" {
            return entity.Value, true
        }
    }

    return nil, false
}
```

You can combine this the usage of `Validators`.

## Docker Usage

Please see the [HelloWorld example](https://github.com/wochinge/go-rasa-sdk/tree/master/examples/HelloWorld) for an
example how to build a Docker image of your Go action server.

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
