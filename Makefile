SHELL := /bin/bash
include .env
export

.PHONY: help
help: ## display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: update
update: ## go modules update
	@go get -u -t ./...
	@go mod tidy
	@go mod vendor

.PHONY: goredis
goredis: ## run the goredis demo
	@go run cmd/kv/goredis/main.go

.PHONY: ruedis
ruedis: ## run the ruedis demo
	@go run cmd/kv/ruedis/main.go

.PHONY: pgx
pgx: ## run the pgx demo
	@go run cmd/postgres/pgx/main.go

.PHONY: pq
pq: ## run the pq demo
	@go run cmd/postgres/pq/main.go

.PHONY: blob
blob: ## run the blob demo
	@go run cmd/blob/main.go

.PHONY: config
config: ## run the config demo
	@go run cmd/config/main.go
