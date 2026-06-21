package main

import (
	"context"
	"log"
	"net/http"
	"oaseller/internal/app/merchants"
	"oaseller/internal/health"
	"oaseller/internal/platform/postgres"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	pgPool, err := postgres.NewPool(ctx)
	if err != nil {
		log.Fatal(err)
	}

	defer pgPool.Close()

	merchantRepository := merchants.NewRepository(pgPool)
	merchantService := merchants.NewService(merchantRepository)
	merchantHandler := merchants.NewHandler(merchantService)

	port := os.Getenv("PORT")

	if port == "" {
		port = "7009"
	}

	r := chi.NewRouter()

	postgresCheck := postgres.NewHealthChecker(pgPool)

	healthService := health.NewService(postgresCheck)
	healthHandler := health.NewHandler(healthService)

	r.Get("/health/live", healthHandler.Live)
	r.Get("/health/ready", healthHandler.Ready)

	r.Get("/merchants", merchantHandler.ListMerchants)
	r.Post("/merchants", merchantHandler.CreateMerchant)

	log.Printf("server running on port %s", port)

	err = http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Printf("%v", err)
		return
	}
}
