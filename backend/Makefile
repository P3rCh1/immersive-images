OAPI_CODEGEN_VERSION = v2.8.0

.PHONY: oapi-codegen

oapi-codegen:
	cd internal/api/models && \
	go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@$(OAPI_CODEGEN_VERSION) \
		--config api.yaml ../../../../swagger/components.yaml

	cd internal/api/server && \
	go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@$(OAPI_CODEGEN_VERSION) \
		--config api.yaml ../../../../swagger/openapi.yaml
