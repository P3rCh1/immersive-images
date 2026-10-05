OAPI_CODEGEN_VERSION = v2.8.0
GOLANGCI_LINT_VERSION = v2.13.2

.PHONY: oapi-codegen golangci-lint unit-tests acc-tests golangci-lint-fix migrate

oapi-codegen:
	cd internal/api/models && \
	go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@$(OAPI_CODEGEN_VERSION) \
		--config api.yaml ../../../../swagger/components.yaml

	cd internal/api/server && \
	go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@$(OAPI_CODEGEN_VERSION) \
		--config api.yaml ../../../../swagger/openapi.yaml

golangci-lint:
	go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION) run

golangci-lint-fix:
	go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION) run --fix

unit-tests:
	go test ./...

acc-tests:
	docker compose -f acceptance/docker-compose.yaml up --build --abort-on-container-exit --exit-code-from tests tests

migrate:
	docker compose exec debug sql-migrate up
