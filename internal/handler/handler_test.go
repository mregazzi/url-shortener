package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestShortenHandler(t *testing.T) {
	body := []byte(`{"url":"https://example.com"}`)
	req := httptest.NewRequest(http.MethodPost, "/shorten", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	ShortenHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rr.Code)
	}

	if !bytes.Contains(rr.Body.Bytes(), []byte("code")) {
		t.Fatalf("expected response to contain 'code', got %s", rr.Body.String())
	}
}

func TestShortenHandler_returnsDifferentCOdes(t *testing.T) {
	body := []byte(`{"url":"https://example.com"}`)
	req := httptest.NewRequest(http.MethodPost, "/shorten", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	ShortenHandler(rr, req)

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
