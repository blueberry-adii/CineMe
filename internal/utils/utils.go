package utils

import (
	"encoding/json"
	"net/http"
)

/*
* Helper function to set content json, encode data and return over http
 */
func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

/*
* Helper function to parse json body into go readable format
 */
func ParseJSON(r *http.Request, v any) error {
	return json.NewDecoder(r.Body).Decode(v)
}
