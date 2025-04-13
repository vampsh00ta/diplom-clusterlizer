include .env
export

.PHONY: help
help: ## Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

compose-up: ### Run docker-compose
	docker-compose up --build -d && docker-compose logs -f
.PHONY: compose-up

compose-down: ### Down docker-compose
	docker-compose down --remove-orphans
.PHONY: compose-down

docker-rm-volume: ### remove docker volume
	docker volume rm pg-data
.PHONY: docker-rm-volume

linter-golangci: ### check by golangci linter
	golangci-lint run
.PHONY: linter-golangci

linter-hadolint: ### check by hadolint linter
	git ls-files --exclude='Dockerfile*' -c --ignored | xargs hadolint
.PHONY: linter-hadolint

gen-proto: ### generate proto interfaces
	protoc   ./internal/handler/grpc/pb/matchmaking.proto --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative
.PHONY: gen-proto
cover: ### run test with coverage
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out
	rm coverage.out
.PHONY: coverage

migrate-create:  ### create new migration
	migrate create -ext sql -dir migration -seq  $(NAME)
.PHONY: migrate-create

migrate-force:  ### create new migration
	 migrate -path migration  -database  '$(POSTGRES_URL)' force $(SEQ)
.PHONY: migrate-force

migrate-up: ### migration up
	migrate -path migration -database '$(POSTGRES_URL)' up
.PHONY: migrate-up

migrate-down: ### migration up
	migrate -path migration -database '$(POSTGRES_URL)' down $(SEQ)
.PHONY: migrate-down

