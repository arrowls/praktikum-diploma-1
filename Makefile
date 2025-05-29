include .env
MIGRATIONS_DIR:=migrations

create-migration:
	@read -p "Enter migrate name: " migrate_name; \
	migrate create -seq -ext sql -dir $(MIGRATIONS_DIR) $$migrate_name
migrate-up:
	@migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URI)" up

migrate-down:
	@migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URI)" down

install-migrate:
	@go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

run:
	go run cmd/gophermart/main.go