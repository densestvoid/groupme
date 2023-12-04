package groupme

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// HTTPMessageCallback is a function that acts on new messages sent from
// the GroupMe server to a callback URL.
type HTTPMessageCallback func(Message)

// HTTPHandlerFunc creates an http.HandlerFunc that executes callback functions
// on each received message. Function should be registered on an http.Handler
// route for use in a callback URL server.
func HTTPHandlerFunc(callbacks ...HTTPMessageCallback) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			http.Error(w, fmt.Sprintf("unsupported Content-Type: %s", ct), http.StatusUnsupportedMediaType)
			return
		}

		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()

		var msg Message
		if err := decoder.Decode(&msg); err != nil {
			http.Error(w, fmt.Sprintf("unable to decode message: %v", err), http.StatusBadRequest)
			return
		}

		for _, callback := range callbacks {
			callback(msg)
		}

		w.WriteHeader(http.StatusOK)
	}
}
