package server

import (
	"encoding/json"
	"fmt"
	"github.com/wochinge/go-rasa-sdk/actions"
	"github.com/wochinge/go-rasa-sdk/rasa"
	"net/http"

	"github.com/gorilla/mux"
)

type healthResponse struct {
	Status string `json:"status"`
}

type errorResponse struct {
	Error      string `json:"error"`
	ActionName string `json:"action_name"`
}

func health(w http.ResponseWriter, _ *http.Request) {
	responseBody := healthResponse{"ok"}

	sendJSONResponse(w, responseBody, http.StatusOK)
}

func runAction(availableActions []actions.Action) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		actionRequest, err := rasa.Parsed(r.Body)
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
			ActionName:actionName}, http.StatusNotFound)
		return
	case *actions.ExecutionRejectedError:
		sendJSONResponse(w, errorResponse{Error: fmt.Sprintf("Action execution failed with error: %v.", err),
			ActionName:actionName}, http.StatusBadRequest)
		return
	}
}

func sendJSONResponse(writer http.ResponseWriter, responseBody interface{}, status int) {
	serialized, _ := json.Marshal(responseBody)

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	writer.Write(serialized)
}

// GetRouter returns the routes for which the server accepts requests.
func GetRouter(actions ...actions.Action) http.Handler {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/health", health).Methods("GET", "OPTIONS")
	router.HandleFunc("/webhook", runAction(actions)).Methods("POST")

	return router
}

// Serve runs the action server on port port.
func Serve(port int, actions ...actions.Action) error {
	fmt.Printf("Running Rasa action server on port '%v'.\n", port)
	address := fmt.Sprintf(":%v", port)
	return http.ListenAndServe(address, GetRouter(actions...))
}
