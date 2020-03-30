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

func health(w http.ResponseWriter, _ *http.Request) {
	responseBody := healthResponse{"ok"}

	sendJSONResponse(w, responseBody, http.StatusOK)
}

func runAction(availableActions []actions.Action) func(http.ResponseWriter, *http.Request) {
	return func (w http.ResponseWriter, r *http.Request) {

		actionRequest, err := rasa.Parsed(r.Body)
		if err != nil {
			return // TODO
		}

		actionToRun, err := actions.ActionFor(actionRequest.ActionToRun, availableActions)

		if err != nil {

		}

		dispatcher := actions.NewDispatcher()
		newEvents := actionToRun.Run(&actionRequest.Tracker, &actionRequest.Domain, dispatcher)
		fmt.Printf("%v", newEvents)
		// parse request
		// Run action
		// marshall events + responses

	}
}


func sendJSONResponse(writer http.ResponseWriter, responseBody interface{}, status int) {
	serialized, _ := json.Marshal(responseBody)

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
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
func Serve(port int) error {
	fmt.Printf("Running Rasa action server on port '%v'.\n", port)
	address := fmt.Sprintf(":%v", port)
	return http.ListenAndServe(address, GetRouter())
}
