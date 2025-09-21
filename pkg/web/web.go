package web

import (
	"encoding/json"
	"net/http"
)

func JSON[T any](w http.ResponseWriter, status int, t T) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(t)
}
