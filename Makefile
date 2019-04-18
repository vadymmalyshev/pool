ADMIN_BIN_NAME = bin/hadmin

OS := $(shell uname | tr '[:upper:]' '[:lower:]')

admin:
	@echo "Building ADMIN binary..."
	@go build -o ${ADMIN_BIN_NAME} ./cmd/hadmin
	@echo "You can now use ./${ADMIN_BIN_NAME}"
