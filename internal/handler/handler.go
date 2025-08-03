package handler

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"time"
)

type shortenRequest struct {
	URL string `json:"url"`
}

type shortenResponse struct {
	Code string `json:"code"`
}

func ShortenHandler(w http.ResponseWriter, r *http.Request) {
	var req shortenRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.URL == "" {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	code := generateCode(6)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(shortenResponse{Code: code})

}

var randGen = rand.New(rand.NewSource(time.Now().UnixNano()))

func generateCode(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[randGen.Intn(len(letters))]
	}
	return string(b)
}
