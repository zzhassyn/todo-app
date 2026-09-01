include .env
export

PROJECT_ROOT ?= $(CURDIR)
export PROJECT_ROOT

env-up:
	@docker compose up -d todoapp-postgres

env-down:
	@docker compose down todoapp-postgres

env-cleanup:
	@read -p "This will remove all containers, volumes, and networks. Are you sure? (y/n): " answer; \
	if [ "$$answer" = "y" ]; then \
		docker compose down todoapp-postgres && \
		rm -rf ${PROJECT_ROOT}/out/pgdata && \
		echo "Cleanup completed."; \
	else \
		echo "Cleanup aborted."; \
	fi

migrate-create:
	@if [ -z "$(seq)" ]; then \
		echo "Error: Migration name is required. Usage: make migrate-create seq=your_migration_name"; \
		exit 1; \
	fi; \
	docker compose run --rm todoapp-postgres-migrate \
		create \
		-ext sql \
		-dir /migrations \
		-seq "$(seq)"

migrate-up:
	@make migrate-action action=up

migrate-down:
	@make migrate-action action=down

migrate-action:
	@if [ -z "$(action)" ]; then \
		echo "Error: Action is required. Usage: make migrate-action action=up|down"; \
		exit 1; \
	fi; \
	docker compose up -d --wait todoapp-postgres; \
	docker compose run --rm todoapp-postgres-migrate \
		-path /migrations \
		-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@todoapp-postgres:5432/${POSTGRES_DB}?sslmode=disable \
		"$(action)"

todoapp-run:
	@export LOGGER_FOLDER=${PROJECT_ROOT}/out/logs && \
	export POSTGRES_HOST=127.0.0.1 && \
	export POSTGRES_PORT=5433 && \
	cd backend && \
	go mod tidy && \
	go run cmd/todoapp/main.go