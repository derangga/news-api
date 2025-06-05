include .env

migrate_up:
	migrate -path database/ -database "postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_DB_HOST):$(POSTGRES_DB_PORT)/$(POSTGRES_DB)?sslmode=disable" -verbose up

migrate_down:
	migrate -path database/ -database "postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_DB_HOST):$(POSTGRES_DB_PORT)/$(POSTGRES_DB)?sslmode=disable" -verbose down 1

run_api:
	go run cmd/main/main.go