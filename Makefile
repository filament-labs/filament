# Define variables for directories and binaries
BIN_DIR := bin
FRONTEND_DIR := web/frontend
WEB_DIR := web
DESKTOP_BIN := $(BIN_DIR)/desktop
WEB_BIN := $(BIN_DIR)/web-server

PROTO_DIR := proto               # your .proto files location
GO_OUT := internal/pb           # generated Go code
TS_OUT := web/src/pb            # generated TS code

# Plugins
PROTOC_GEN_GO := $(shell which protoc-gen-go)
PROTOC_GEN_GO_GRPC := $(shell which protoc-gen-go-grpc)
PROTOC_GEN_TS := $(shell which protoc-gen-ts_proto)

# Phony targets (not associated with files)
.PHONY: all build-desktop build-web dev-desktop dev-web test docker-build clean proto

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
	rm -rf $(GO_OUT) $(TS_OUT)

# Generate Go and TS protobuf files
proto:
	@echo "Generating Go protobuf files..."
	# Recursively find all .proto files
	protoc -I=$(PROTO_DIR) \
		--go_out=. \
		--go-grpc_out=. \
		$(shell find $(PROTO_DIR) -name '*.proto')

	@echo "Generating TypeScript protobuf files..."
	@mkdir -p $(TS_OUT)
	protoc -I=$(PROTO_DIR) \
		--plugin=protoc-gen-ts_proto=$(PROTOC_GEN_TS) \
		--ts_proto_out=$(TS_OUT) \
		--ts_proto_opt=esModuleInterop=true,forceLong=long,useOptionals=messages \
		$(shell find $(PROTO_DIR) -name '*.proto')

