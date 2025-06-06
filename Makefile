include .env

migrate_up:
	migrate -path database/ -database "postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_DB_HOST):$(POSTGRES_DB_PORT)/$(POSTGRES_DB)?sslmode=disable" -verbose up

migrate_down:
	migrate -path database/ -database "postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_DB_HOST):$(POSTGRES_DB_PORT)/$(POSTGRES_DB)?sslmode=disable" -verbose down 1

run_api:
	go run cmd/main/main.go

mocks_gen:
	mockgen -source=internal/usecase/usecase.go -destination=mocks/usecase/usecase.go
	mockgen -source=internal/repository/repository.go -destination=mocks/repository/repository.go

run_test:
	go test -v -coverprofile=coverage.out ./internal/...

run_test_coverage:
	go test ./... -coverprofile=cover.out
	go tool cover -html=cover.out