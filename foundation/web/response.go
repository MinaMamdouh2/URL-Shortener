package web

import (
	"context"
	"encoding/json"
	"net/http"
)

// Respond converts a Go value to JSON and sends it to the client.
// We added Respond that knows how to do the http related stuff for responding
func Respond(ctx context.Context, w http.ResponseWriter, data any, statusCode int) error {
	setStatusCode(ctx, statusCode)

	// If the response is 204 No Content, the HTTP spec says you must not include a response body
	// So we only set the status code and exit
	if statusCode == http.StatusNoContent {
		w.WriteHeader(statusCode)
		return nil
	}

	// Converts the data value into JSON
	// If marshaling fails (e.g., unsupported type or recursive structure), returns the error.
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	// Write the JSON encoded data to the HTTP response body
	// If this fails (e.g., broken pipe, client closed connection), the error is returned
	if _, err := w.Write(jsonData); err != nil {
		return err
	}

	return nil
}
