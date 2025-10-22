.PHONY: build-desktop build-web test

build-desktop:
    @go build -o bin/desktop ./cmd/desktop

build-web:
    @go build -o bin/web-server ./cmd/web
    @cd web/frontend && pnpm build

dev-desktop:
    @go run ./cmd/desktop

dev-web:
    @go run ./cmd/web-server & 
    @cd web/frontend && pnpm dev

test:
    @go test ./...

docker-build:
    @docker-compose build