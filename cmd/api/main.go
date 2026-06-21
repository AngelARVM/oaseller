package main

import (
	"log"
	"net/http"
	"oaseller/internal/health"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	port := os.Getenv("PORT")

	if port == "" {
		port = "7009"
	}

	r := chi.NewRouter()

	healthService := health.NewService()
	healthHandler := health.NewHandler(healthService)

	r.Get("/health/live", healthHandler.Live)
	r.Get("/health/ready", healthHandler.Ready)

	log.Printf("server running on port %s", port)

	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Fatal(err)
	}
}
