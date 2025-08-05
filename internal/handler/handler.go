package handler

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"math/rand"
	"net/http"
	"sync"
	"time"
	"url-shortener/internal/storage"
)

type Handler struct {
	Store storage.Store
}
type shortenRequest struct {
	URL string `json:"url"`
}

type shortenResponse struct {
	Code string `json:"code"`
}

func (h *Handler) ShortenHandler(w http.ResponseWriter, r *http.Request) {
	var req shortenRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.URL == "" {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	const maxAttempts = 5
	var code string

	for i := 0; i < maxAttempts; i++ {
		candidate := codeGenerator()
		if _, exists := GetURL(candidate); !exists {
			code = candidate
			break
		}
	}
	if code == "" {
		http.Error(w, "unable to generate unique code", http.StatusInternalServerError)
		return
	}

	err := h.Store.Save(code, req.URL)
	if err != nil {
		http.Error(w, "failed to save", http.StatusInternalServerError)
		return
	}

	//saveURL(code, req.URL)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(shortenResponse{Code: code})

}

var randGen = rand.New(rand.NewSource(time.Now().UnixNano()))

var codeGenerator = func() string {
	return generateCode(6)
}

func SetCodeGenerator(f func() string) (restore func()) {
	old := codeGenerator
	codeGenerator = f
	return func() {
		codeGenerator = old
	}
}

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

func (h *Handler) RedirectHandler(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")

	url, found, err := h.Store.Get(code)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	if !found {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, url, http.StatusFound)
}
