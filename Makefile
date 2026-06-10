.PHONY: run tidy test migrate-up migrate-down

# Run the local application entrypoint
run:
	go run cmd/api/main.go

stop:
	@echo "Sending SIGTERM to PayGo application..."
	@powershell -Command "Stop-Process -Id (Get-NetTCPConnection -LocalPort 8080).OwningProcess -Force"

# Tidy up project dependencies
tidy:
	go mod tidy

# Execute all unit tests in the project
test:
# 	go test -v ./...
	go test -v ./internal/service -run GetByPublicID


db-stat:
	pg_ctl -D "$$USERPROFILE/scoop/apps/postgresql/current/data" status

db-start:
	pg_ctl -D "$$USERPROFILE/scoop/apps/postgresql/current/data" start

db-stop:
	pg_ctl -D "$$USERPROFILE/scoop/apps/postgresql/current/data" stop



# Apply SQL migrations to local Postgres instance
migrate-up:
	migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/paygo?sslmode=disable" -verbose up

# Rollback SQL migrations
migrate-down:
	migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/paygo?sslmode=disable" -verbose down
