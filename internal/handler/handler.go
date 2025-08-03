package handler

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"math/rand"
	"net/http"
	"sync"
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
	saveURL(code, req.URL)

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

var (
	store = make(map[string]string)
	mu    = sync.RWMutex{}
)

func saveURL(code, url string) {
	mu.Lock()
	defer mu.Unlock()
	store[code] = url
}

func GetURL(code string) (string, bool) {
	mu.RLock()
	defer mu.RUnlock()
	url, ok := store[code]
	return url, ok
}

func ResetStore() {
	mu.Lock()
	defer mu.Unlock()
	store = make(map[string]string)
}

func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")

	url, ok := GetURL(code)
	if !ok {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, url, http.StatusFound)
}
