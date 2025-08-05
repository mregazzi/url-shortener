package handler

import (
	"bytes"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
	"net/http/httptest"
	"testing"
	"url-shortener/internal/storage"
)

func TestShortenHandler(t *testing.T) {
	body := []byte(`{"url":"https://example.com"}`)
	req := httptest.NewRequest(http.MethodPost, "/shorten", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	h := &Handler{Store: storage.NewInMemoryStore()}
	rr := httptest.NewRecorder()
	h.ShortenHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rr.Code)
	}

	if !bytes.Contains(rr.Body.Bytes(), []byte("code")) {
		t.Fatalf("expected response to contain 'code', got %s", rr.Body.String())
	}
}

func TestShortenHandler_returnsDifferentCodes(t *testing.T) {
	body := []byte(`{"url":"https://example.com"}`)
	req := httptest.NewRequest(http.MethodPost, "/shorten", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	h := &Handler{Store: storage.NewInMemoryStore()}
	rr := httptest.NewRecorder()
	h.ShortenHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rr.Code)
	}

	var resp map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("cannot decode response: %v", err)
	}

	code, ok := resp["code"]
	if !ok || len(code) != 6 {
		t.Fatalf("expected a code of length 6, got: %v", code)
	}

}

func TestShortenHandler_savesCodeInMemory(t *testing.T) {
	store := storage.NewInMemoryStore() // crea e tieni il riferimento
	h := &Handler{Store: store}

	body := []byte(`{"url":"https://example.com"}`)
	req := httptest.NewRequest(http.MethodPost, "/shorten", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.ShortenHandler(rr, req)

	var resp map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("cannot decode response: %v", err)
	}

	code := resp["code"]

	url, ok, err := store.Get(code) // accedi allo stesso store usato dall’handler
	if err != nil {
		t.Fatalf("unexpected error getting url: %v", err)
	}
	if !ok || url != "https://example.com" {
		t.Fatalf("expected to find a saved url, got: %v (ok=%v)", url, ok)
	}
}

func TestShortenHandler_redirectsToOriginalURL(t *testing.T) {
	store := storage.NewInMemoryStore()
	_ = store.Save("xyz789", "https://golang.org")

	h := &Handler{Store: store}

	req := httptest.NewRequest(http.MethodGet, "/xyz789", nil)
	rr := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Get("/{code}", h.RedirectHandler)
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusFound {
		t.Fatalf("expected status code to be %d, got %d", http.StatusFound, rr.Code)
	}

	location := rr.Header().Get("Location")
	if location != "https://golang.org" {
		t.Fatalf("expected redirect to 'https://golang.org', got %s", location)
	}
}

func TestShortenHandler_404IfCodeNotFound(t *testing.T) {
	ResetStore()
	req := httptest.NewRequest(http.MethodGet, "/nope123", nil)
	rr := httptest.NewRecorder()
	h := &Handler{Store: storage.NewInMemoryStore()}
	r := chi.NewRouter()
	r.Get("/{code}", h.RedirectHandler)
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected status code to be %d, got %d", http.StatusNotFound, rr.Code)
	}
}

func TestShortenHandler_avoidDuplicateCodes(t *testing.T) {
	ResetStore()
	const already_present_code = "abc123"
	saveURL(already_present_code, "https://already-present.org")

	restore := SetCodeGenerator(func() string { return already_present_code })
	defer restore()

	body := []byte(`{"url":"https://new-url.org"}`)
	req := httptest.NewRequest(http.MethodPost, "/shorten", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	h := &Handler{Store: storage.NewInMemoryStore()}
	rr := httptest.NewRecorder()
	h.ShortenHandler(rr, req)

	if rr.Code != http.StatusInternalServerError {

	}
}

func TestShortenHandler_failsAfterTooManyDuplicates(t *testing.T) {
	ResetStore()
	usedCodes := []string{"a1", "a2", "a3", "a4", "a5"}
	for _, c := range usedCodes {
		saveURL(c, "https://old.com"+c)
	}

	attempts := 0
	restore := SetCodeGenerator(func() string {
		attempts++
		return usedCodes[(attempts-1)%len(usedCodes)]
	})
	defer restore()

	body := []byte(`{"url":"https://new.com"}`)
	req := httptest.NewRequest(http.MethodPost, "/shorten", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	h := &Handler{Store: storage.NewInMemoryStore()}
	rr := httptest.NewRecorder()
	h.ShortenHandler(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected status code to be %d, got %d", http.StatusInternalServerError, rr.Code)
	}

	if attempts != 5 {
		t.Fatalf("expected 5 attempts, got %d", attempts)
	}
}
