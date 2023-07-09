# Load environment variables from .env file
include .env
export

.PHONY: migrate-up
migrate-up:
	migrate -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" -path migrations up

.PHONY: migrate-down
migrate-down:
	migrate -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" -path migrations down

.PHONY: stop-db
stop-db:
	docker stop $(DOCKER_NAME) && docker rm $(DOCKER_NAME)

.PHONY: start-db
start-db:
	docker run --name $(DOCKER_NAME) -e POSTGRES_PASSWORD=$(DB_PASSWORD) -p $(DB_PORT):$(DB_HOST) -d postgres

.PHONY: run
run:
	go run main.go

.PHONY: test
test:
	go test -v ./...

.PHONY: lint
lint:
	golangci-lint run ./...