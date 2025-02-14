package helper

import (
	"encoding/json"
	"net/http"
)

type HTTPError struct {
	StatusCode int
	Error      any
}
type HTTPFunc func(w http.ResponseWriter, r *http.Request) *HTTPError

func MakeHttpHandlerFunc(f HTTPFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if httpError := f(w, r); httpError != nil {

			WriteJSON(w, httpError.StatusCode, map[string]any{"error": httpError.Error})
		}
	}
}

func WriteJSON(w http.ResponseWriter, statusCode int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(body)
}
