package handler

import (
	"bytes"
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
