//go:build e2e

package test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"
)

func TestEndToEnd_ShortenAndRedirect(t *testing.T) {
	baseURL := "http://localhost:8080"

	// 1. POST /shorten
	payload := []byte(`{"url":"https://test.com"}`)
	resp, err := http.Post(baseURL+"/shorten", "application/json", bytes.NewReader(payload))
	if err != nil {
		t.Fatalf("failed to POST /shorten: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", resp.StatusCode)
	}

	var data map[string]string
	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &data); err != nil {
		t.Fatalf("cannot decode response: %v", err)
	}

	code := data["code"]
	if code == "" {
		t.Fatal("code not returned")
	}

	// 2. GET /{code}
	client := &http.Client{
		// prevent automatic redirect
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Timeout: 5 * time.Second,
	}

	resp2, err := client.Get(baseURL + "/" + code)
	if err != nil {
		t.Fatalf("GET /%s failed: %v", code, err)
	}
	defer resp2.Body.Close()

	if resp2.StatusCode != http.StatusFound {
		t.Fatalf("expected 302 redirect, got %d", resp2.StatusCode)
	}

	location := resp2.Header.Get("Location")
	if location != "https://test.com" {
		t.Fatalf("expected redirect to https://example.com, got %s", location)
	}
}
