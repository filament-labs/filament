# Define variables for directories and binaries
BIN_DIR := bin
FRONTEND_DIR := web/frontend
WEB_DIR := web
DESKTOP_BIN := $(BIN_DIR)/desktop
WEB_BIN := $(BIN_DIR)/web-server

# Phony targets (not associated with files)
.PHONY: all build-desktop build-web dev-desktop dev-web test docker-build clean

# Default target
all: build-desktop build-web

# Build desktop application
build-desktop:
	@mkdir -p $(BIN_DIR)
	go build -o $(DESKTOP_BIN) ./cmd/desktop

# Build web application (backend and frontend)
build-web:
	@mkdir -p $(BIN_DIR)
	go build -o $(WEB_BIN) ./cmd/web
	@if [ -d "$(FRONTEND_DIR)" ]; then \
		cd $(FRONTEND_DIR) && pnpm build; \
	else \
		echo "Error: $(FRONTEND_DIR) not found"; exit 1; \
	fi

# Run desktop app in development mode
dev-desktop:
	go run ./cmd/desktop

generate-constants:
	go generate ./...

# Run web app in development mode
dev-web:
	@if [ -d "$(WEB_DIR)" ]; then \
		cd $(WEB_DIR) && pnpm dev & go run ./cmd/web; \
	else \
		echo "Error: $(WEB_DIR) not found"; exit 1; \
	fi

# Run tests
test:
	go test ./...

# Build Docker images
docker-build:
	docker compose build

# Clean build artifacts
clean:
	rm -rf $(BIN_DIR)
	@if [ -d "$(FRONTEND_DIR)" ]; then \
		cd $(FRONTEND_DIR) && pnpm clean || true; \
	fi
