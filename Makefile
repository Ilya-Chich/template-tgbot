-include .env
GO_BIN := $(GOPATH)/bin
TOOLBIN = go -C tools tool -n
GOIMPORTS := $(shell $(TOOLBIN) goimports)
GOFUMPT := $(shell $(TOOLBIN) gofumpt)
GOLINES := $(shell $(TOOLBIN) golines)
OAPI_CODEGEN := $(GO_BIN)/oapi-codegen
MERGED_OAPI_V1=$(PWD)/api/openapi/v1/merged.json
OAPI_MERGER := $(GO_BIN)/oapi-merger
GOLANGCI := $(shell $(TOOLBIN) golangci-lint)

init: \
	postgres-init \
	up

init-prod: \
	lint \
	up-prod


up: \
	stop \
	docker-up

stop:
	docker-compose stop

.PHONY: clean-volumes
clean-volumes:
	docker-compose down --volumes --remove-orphans --rmi local

docker-up:
	docker-compose up -d --build --remove-orphans

.PHONY: up-prod
up-prod:
	docker compose -f docker-compose.prod.yml up -d --build

.PHONY: down-prod
down-prod:
	docker compose -f docker-compose.prod.yml down

build:
	docker-compose build

.PHONY: postgres-init
postgres-init:
	docker-compose up -d postgres

.PHONY: docs
docs: openapi_merge openapi_http

.PHONY: openapi_merge
openapi_merge: $(OAPI_MERGER)
	$(OAPI_MERGER) -wdir api/openapi/v1 -spec openapi.yaml -o $(MERGED_OAPI_V1)

.PHONY: openapi_http
openapi_http: $(OAPI_CODEGEN)
	$(OAPI_CODEGEN) --old-config-style  -generate types,skip-prune -o ./internal/controllers/v1/view/types.gen.go -package view $(MERGED_OAPI_V1)
	$(OAPI_CODEGEN) --old-config-style -generate spec -o ./internal/controllers/v1/spec.gen.go -package v1 $(MERGED_OAPI_V1)
	rm -f $(MERGED_OAPI_V1)

$(OAPI_MERGER):
	go install github.com/felicson/oapi-merger/cmd/oapi-merger@v0.0.2
$(OAPI_CODEGEN):
	go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.11.0

.PHONY: lint
lint:
	"$(GOLANGCI)" run

.PHONY: test
test:
	go test ./...

generate_models:
	tables-to-go -fn-format s -u ${POSTGRES_USER} -p ${POSTGRES_PASSWORD} -d ${POSTGRES_DB} -v -of ./models

.PHONY: fmt
fmt:
	"$(GOIMPORTS)" -w -l .
	"$(GOFUMPT)" -w -l .
	"$(GOLINES)" -w --no-reformat-tags --max-len=120 .

.PHONY: migrate-create
migrate-create:
	@if [ -z "$(NAME)" ]; then \
		echo "Usage: make migrate-create NAME=your_migration_name"; \
		exit 1; \
	fi
	"$(GOOSEBIN)" -dir database/migrations create $(NAME) sql

GOOSEBIN := $(shell $(TOOLBIN) goose)
GOOSE := "$(GOOSEBIN)" -dir database/migrations postgres

define migrate
	@saved_url="$${POSTGRES_URL}"; \
	set -a; . ./.env.dev; set +a; \
	if [ -n "$$saved_url" ]; then dsn="$$saved_url"; else dsn="$$(go run ./cmd/pgurl)"; fi; \
	$(GOOSE) "$$dsn" $(1)
endef

.PHONY: migrate-up migrate-down migrate-status migrate-version
migrate-up:
	$(call migrate,up)

migrate-down:
	$(call migrate,down)

migrate-status:
	$(call migrate,status)

migrate-version:
	$(call migrate,version)
