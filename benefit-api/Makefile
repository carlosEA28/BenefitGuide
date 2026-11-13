.PHONY: run build test fmt tidy migration-up migration-down migration-rollback

GOOSE ?= goose
MIGRATIONS_DIR := migrations

# Export the same DB_* settings used by the application for Goose/PostgreSQL.
# Keep the password in the environment instead of passing it on the command line.
define run_migration
	@set -a; \
	if [ -f .env ]; then . ./.env; fi; \
	set +a; \
	: "$${DB_HOST:?DB_HOST is required}"; \
	: "$${DB_USER:?DB_USER is required}"; \
	: "$${DB_NAME:?DB_NAME is required}"; \
	PGHOST="$$DB_HOST" PGPORT="$${DB_PORT:-5432}" \
	PGUSER="$$DB_USER" PGPASSWORD="$${DB_PASSWORD:-}" \
	PGDATABASE="$$DB_NAME" PGSSLMODE="$${DB_SSL_MODE:-disable}" \
	$(GOOSE) -dir $(MIGRATIONS_DIR) postgres "sslmode=$${DB_SSL_MODE:-disable}" $(1)
endef

run:
	go run .

build:
	go build -o bin/benefit-api .

test:
	go test ./...

fmt:
	go fmt ./...

tidy:
	go mod tidy

migration-up:
	$(call run_migration,up)

# Revert every applied migration.
migration-down:
	$(call run_migration,reset)

# Revert only the latest applied migration.
migration-rollback:
	$(call run_migration,down)
