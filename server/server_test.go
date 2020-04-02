package server

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/wochinge/go-rasa-sdk/rasa"
	"github.com/wochinge/go-rasa-sdk/rasa/events"
	"github.com/wochinge/go-rasa-sdk/rasa/responses"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealth(t *testing.T) {
	request, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	response := httptest.NewRecorder()
	handler := GetRouter()

	handler.ServeHTTP(response, request)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, `{"status":"ok"}`, response.Body.String())
}

type TestAction struct{}

func (action *TestAction) Run(_ *rasa.Tracker, _ *rasa.Domain, _ responses.ResponseDispatcher) []events.Event {
	return []events.Event{&events.Restarted{}}
}
func (action *TestAction) Name() string { return "test-action" }

func TestRunAction(t *testing.T) {
	body := []byte(`{"next_action": "test-action"}`)
	request, err := http.NewRequest("POST", "/webhook", bytes.NewBuffer(body))

	if err != nil {
		t.Fatal(err)
	}

	response := httptest.NewRecorder()
	handler := GetRouter(&TestAction{})

	handler.ServeHTTP(response, request)

	assert.Equal(t, http.StatusOK, response.Code)

	expectedResponse := `{"events":[{"event":"restart"}],"responses":[]}`
	assert.Equal(t, expectedResponse, response.Body.String())
}

func TestRunActionNotFound(t *testing.T) {
	body := []byte(`{"next_action": "test-action"}`)
	request, err := http.NewRequest("POST", "/webhook", bytes.NewBuffer(body))

	if err != nil {
		t.Fatal(err)
	}

	response := httptest.NewRecorder()
	handler := GetRouter()

	handler.ServeHTTP(response, request)

	assert.Equal(t, http.StatusNotFound, response.Code)
}

func TestRunActionInvalidPayload(t *testing.T) {
	body := []byte(`{"}`)
	request, err := http.NewRequest("POST", "/webhook", bytes.NewBuffer(body))

	if err != nil {
		t.Fatal(err)
	}

	response := httptest.NewRecorder()
	handler := GetRouter()

	handler.ServeHTTP(response, request)

	assert.Equal(t, http.StatusBadRequest, response.Code)
}

type RejectingAction struct{}

func (action *RejectingAction) Run(_ *rasa.Tracker, _ *rasa.Domain, _ responses.ResponseDispatcher) []events.Event {
	return []events.Event{&events.ActionExecutionRejected{}}
}
func (action *RejectingAction) Name() string { return "test-reject" }

func TestActionRejectsExecution(t *testing.T) {
	body := []byte(`{"next_action": "test-reject""}`)
	request, err := http.NewRequest("POST", "/webhook", bytes.NewBuffer(body))

	if err != nil {
		t.Fatal(err)
	}

	response := httptest.NewRecorder()
	handler := GetRouter(&RejectingAction{})

	handler.ServeHTTP(response, request)

	assert.Equal(t, http.StatusBadRequest, response.Code)
}
