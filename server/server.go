// Package server provides the functions to run an action server with your given custom actions.
package server

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/wochinge/go-rasa-sdk/actions"
	"github.com/wochinge/go-rasa-sdk/rasa/request"
	"net/http"

	"github.com/gorilla/mux"
)

// DefaultPort is the port which the Rasa Action Servers runs on by convention.
// If you use a different port, remember to adapt the `action_endpoint` configuration in your `endpoints.yml`.
const DefaultPort int = 5055

type healthResponse struct {
	Status string `json:"status"`
}

type errorResponse struct {
	Error      string `json:"error"`
	ActionName string `json:"action_name"`
}

// Serve runs the action server on the provided port.
// Supplied actions will be executed upon request from Rasa Open Source.
func Serve(port int, actions ...actions.Action) {
	setup(actions)

	err := http.ListenAndServe(address(port), GetRouter(actions...))

	tearDown(err)
}

func setup(actions []actions.Action) {
	log.SetLevel(log.InfoLevel)
	logAvailableActions(actions)
}

func address(port int) string {
	log.Infof("Action server running on on port %v", port)
	return fmt.Sprintf(":%v", port)
}

func tearDown(err error) {
	if err != nil {
		log.Error(err)
	}
}

func logAvailableActions(actions []actions.Action) {
	var actionNames []string
	for _, action := range actions {
		actionNames = append(actionNames, action.Name())
	}

	log.Infof("The following actions are loaded: %v", actionNames)
}

// GetRouter returns the routes for which the server accepts requests.
// By default this is the `health` endpoint which can be used for health checks and
// the `/webhook` endpoint which Rasa Open Source calls to execute a custom action.
// There should only be a need to call this if you want to add custom endpoints.
func GetRouter(actions ...actions.Action) http.Handler {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/health", health).Methods("GET", "OPTIONS")
	router.HandleFunc("/webhook", runAction(actions)).Methods("POST")

	return router
}

func health(w http.ResponseWriter, _ *http.Request) {
	responseBody := healthResponse{"ok"}

	sendJSONResponse(w, responseBody, http.StatusOK)
}

func sendJSONResponse(writer http.ResponseWriter, responseBody interface{}, status int) {
	serialized, _ := json.Marshal(responseBody)

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)

	_, err := writer.Write(serialized)
	if err != nil {
		log.Error(err)
	}
}

func runAction(availableActions []actions.Action) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		actionRequest, err := request.Parsed(r.Body)
		if err != nil {
			sendJSONResponse(w, errorResponse{Error: fmt.Sprintf("parsing body failed with error: %v", err)},
				http.StatusBadRequest)
			return
		}

		responseBody, err := actions.ExecuteAction(actionRequest, availableActions)

		if err == nil {
			sendJSONResponse(w, responseBody, http.StatusOK)
			return
		}

		handleExecutionError(w, actionRequest.ActionToRun, err)
	}
}

func handleExecutionError(w http.ResponseWriter, actionName string, err error) {
	switch err.(type) {
	case *actions.NotFoundError:
		sendJSONResponse(w, errorResponse{Error: fmt.Sprintf("Action execution failed with error: %v.", err),
			ActionName: actionName}, http.StatusNotFound)
		return
	case *actions.ExecutionRejectedError:
		sendJSONResponse(w, errorResponse{Error: fmt.Sprintf("Action execution failed with error: %v.", err),
			ActionName: actionName}, http.StatusBadRequest)
		return
	}
}
