OAPI_CODEGEN_VERSION = v2.8.0

.PHONY: oapi-codegen golangci-lint unit-tests golangci-lint-fix migrate

oapi-codegen:
	cd internal/api/models && \
	go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@$(OAPI_CODEGEN_VERSION) \
		--config api.yaml ../../../../swagger/components.yaml

	cd internal/api/server && \
	go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@$(OAPI_CODEGEN_VERSION) \
		--config api.yaml ../../../../swagger/openapi.yaml

golangci-lint:
	go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.13.2 run

unit-tests:
	go test ./...

golangci-lint-fix:
	go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.13.2 run --fix

migrate:
	docker compose exec debug sql-migrate up
