// Package httpx contains small, transport-level helpers shared across HTTP
// handlers (JSON encoding/decoding, error envelopes). Keeping these here
// avoids repeating boilerplate in every handler and keeps handlers focused
// on request/response mapping rather than plumbing.
package httpx

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse is the JSON envelope returned for any error response.
type ErrorResponse struct {
	Error string `json:"error"`
}

// WriteJSON writes v as a JSON response body with the given status code.
// It always sets the Content-Type header, even if encoding later fails.
func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	// Encoding errors here would mean v itself is unencodable, which is a
	// programmer error (e.g. a channel or func field), not a request-time
	// failure. There's nothing meaningful left to do but drop it; the
	// status code and headers have already been written.
	_ = json.NewEncoder(w).Encode(v)
}

// WriteError writes a JSON error envelope: {"error": message}.
func WriteError(w http.ResponseWriter, status int, message string) {
	WriteJSON(w, status, ErrorResponse{Error: message})
}
