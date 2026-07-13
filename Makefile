.PHONY: help infra-init infra-plan infra-apply infra-destroy

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

infra-init:
	@cd infra && terraform init

infra-plan:
	@cd infra && terraform plan

infra-apply:
	@cd infra && terraform apply -auto-approve

infra-destroy:
	@cd infra && terraform destroy

# Development helpers
.PHONY: sqlc migrate migrate-up migrate-down migrate-create build run fmt test

sqlc: ## Run sqlc generate to produce DB code
	@sqlc generate

migrate-create: ## Create a new migration: make migrate-create NAME=create_thing
	@if [ -z "$(NAME)" ]; then echo "Please provide NAME, e.g. make migrate-create NAME=create_orders" && exit 1; fi
	@migrate create -ext sql -dir migrations -seq $(NAME)

migrate-up: ## Apply all up migrations (requires DATABASE_URL)
	@if [ -z "$(DATABASE_URL)" ]; then echo "Please set DATABASE_URL environment variable" && exit 1; fi
	@migrate -path ./migrations -database "$(DATABASE_URL)" up

migrate-down: ## Revert last migration (requires DATABASE_URL)
	@if [ -z "$(DATABASE_URL)" ]; then echo "Please set DATABASE_URL environment variable" && exit 1; fi
	@migrate -path ./migrations -database "$(DATABASE_URL)" down 1

build: ## Build the project
	@go build ./...

run: ## Run the main binary (example: make run PKG=cmd/api)
	@if [ -z "$(PKG)" ]; then echo "Please set PKG, e.g. make run PKG=./cmd/api" && exit 1; fi
	@go run $(PKG)

fmt: ## gofmt the repository
	@gofmt -s -w .

test: ## Run unit tests
	@go test ./...
