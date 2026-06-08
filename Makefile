include .env
export

export PROJECT_ROOT=${shell pwd}

env-up:
	@docker compose up -d chat-postgres

env-down:
	@docker compose down chat-postgres

migrate-create:
	@if [ -z "${seq}"]; then \
		echo "pls, try again with seq=value"; \
		exit 1; \
	fi; \
	docker compose run --rm chat-migrate \
	create -ext sql -dir ./migrations -seq ${seq}

migrate-up:
	@docker compose run --rm chat-migrate \
	-path ./migrations \
	-database "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@chat-postgres:5432/${POSTGRES_DB}?sslmode=disable" \
	up

migrate-down:
	@docker compose run --rm chat-migrate \
	-path ./migrations \
	-database "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@chat-postgres:5432/${POSTGRES_DB}?sslmode=disable" \
	down

migrate-force:
	@if [ -z "${version}"]; then \
		echo "pls, try again with version=value"; \
		exit 1; \
	fi; \
	docker compose run --rm chat-migrate \
	-path ./migrations \
	-database "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@chat-postgres:5432/${POSTGRES_DB}?sslmode=disable" \
	force ${version}

app-run:
	@go mod tidy && \
	go run ${PROJECT_ROOT}/cmd/chat/main.go