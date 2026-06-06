.PHONY: run tidy test migrate-up migrate-down

# Run the local application entrypoint
run:
	go run cmd/api/main.go

# Tidy up project dependencies
tidy:
	go mod tidy

# Execute all unit tests in the project
test:
	go test -v ./...

# Apply SQL migrations to local Postgres instance
migrate-up:
	migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/paygo?sslmode=disable" -verbose up

# Rollback SQL migrations
migrate-down:
	migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/paygo?sslmode=disable" -verbose down