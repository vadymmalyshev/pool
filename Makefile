ADMIN_BIN_NAME = bin/hadmin

OS := $(shell uname | tr '[:upper:]' '[:lower:]')

admin:
	@echo "Building hAdmin binary..."
	@go build -o ${ADMIN_BIN_NAME} ./cmd/hadmin
	@echo "You can now use ./${ADMIN_BIN_NAME}"

admin-run:
	./${ADMIN_BIN_NAME} -c /config/config.yaml

admin-migrate:
	./${ADMIN_BIN_NAME} -c /config/config.yaml migrate