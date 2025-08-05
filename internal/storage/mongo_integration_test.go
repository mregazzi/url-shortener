//go:build integration

package storage

import (
	"os"
	"testing"
)

func TestMongoStore_SaveAndGet(t *testing.T) {
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		uri = "mongodb://localhost:27017"
	}
	return
	store, err := NewMongoStore(uri, "testdb", "testurls")
	if err != nil {
		t.Fatalf("failed to connect to mongo %s: %v", uri, err)
	}

	//Clean up collection before and after
	_ = store.collection.Drop(nil)

	code := "test123"
	url := "https://example.com"

	err = store.Save(code, url)
	if err != nil {
		t.Fatalf("failed to save url: %v", err)
	}

	found, ok, err := store.Get(code)
	if err != nil {
		t.Fatalf("failed to retrieve url: %v", err)
	}
	if !ok {
		t.Fatalf("expected url to be found")
	}
	if found != url {
		t.Fatalf("expected %s, got%s", url, found)
	}
}
