package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"url-shortener/internal/handler"
	"url-shortener/internal/storage"
)

func main() {
	// Connetti a MongoDB
	mongoStore, err := storage.NewMongoStore(
		"mongodb://192.168.1.115:27017", // IP del tuo container Mongo
		"urlshortener",                  // nome del database
		"urls",                          // nome della collection
	)
	if err != nil {
		log.Fatalf("errore nella connessione a Mongo: %v", err)
	}

	// Inizializza handler
	h := &handler.Handler{Store: mongoStore}

	// Setup router
	r := chi.NewRouter()
	r.Post("/shorten", h.ShortenHandler)
	r.Get("/{code}", h.RedirectHandler)

	log.Println("Server in ascolto su :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
