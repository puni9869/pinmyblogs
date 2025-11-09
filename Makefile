BINARY_NAME=main
MAIN_PATH=./cmd/main.go
VERSION=$(shell git rev-parse --short HEAD)
PKG=github.com/puni9869/pinmyblogs/cmd/command

.PHONY: build
build:
	go build -ldflags "-w -s -X $(PKG).BuildVersion=$(VERSION)" -o $(BINARY_NAME) $(MAIN_PATH)

.PHONY: server
server:
	air

.PHONY: test
test:
	go test ./... -cover

.PHONY: lint
lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	golangci-lint run ./...

.PHONY: govulncheck
govulncheck:
	go install golang.org/x/vuln/cmd/govulncheck@latest
	govulncheck ./...

.PHONY: vet
vet:
	go vet ./...


# Static Linux binary build (for Ubuntu deployment)
build-linux:
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 \
	CC=x86_64-linux-musl-gcc \
	go build -tags 'osusergo netgo static_build' \
	-ldflags="-linkmode external -extldflags '-static' -s -w -X $(PKG).BuildVersion=$(VERSION)" \
	-o $(BINARY_NAME)-linux $(MAIN_PATH)
