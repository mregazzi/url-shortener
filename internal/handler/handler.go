package handler

import (
	"encoding/json"
	"net/http"
)

func ShortenHandler(w http.ResponseWriter, r *http.Request) {
	resp := map[string]string{"code": "abc123"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
