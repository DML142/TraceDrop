DATABASE_URL ?= postgres://tracedrop:tracedrop@localhost:5432/tracedrop?sslmode=disable

.PHONY: migrate-up migrate-down sqlc test

migrate-up:
	go run github.com/pressly/goose/v3/cmd/goose@v3.28.0 -dir db/migrations postgres "$(DATABASE_URL)" up

migrate-down:
	go run github.com/pressly/goose/v3/cmd/goose@v3.28.0 -dir db/migrations postgres "$(DATABASE_URL)" down

sqlc:
	go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.31.1 generate

test:
	go test ./...
