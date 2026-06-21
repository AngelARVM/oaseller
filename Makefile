.PHONY: dev run psql

dev:
	air

run:
	go run ./cmd/api

psql:
	docker compose exec postgres psql -U postgres -d postgres
