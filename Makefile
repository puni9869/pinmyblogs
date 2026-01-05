BINARY_NAME=main
MAIN_PATH=./cmd/main.go
VERSION=$(shell git rev-parse --short HEAD)
PKG=github.com/puni9869/pinmyblogs/cmd/command
GOLANGCI_LINT_VERSION=v2.1.6

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
lint: install-linter
	golangci-lint run --config .golangci.yml ./... && djlint templates/**/* --lint

.PHONY: install-linter
install-linter:
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

.PHONY: golangci-ci
golangci-ci:
	@command -v golangci-lint >/dev/null || (echo "golangci-lint not found; please install it in CI image"; exit 1)
	@golangci-lint run --config .golangci.yml ./...

.PHONY: format
format: install-linter
	golangci-lint run --config .golangci.yml ./... --fix && djlint templates/**/* --reformat

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
