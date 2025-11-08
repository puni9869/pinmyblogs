BINARY_NAME=main
MAIN_PATH=./cmd/main.go

.PHONY: build
build:
	go build -ldflags "-w -s" $(MAIN_PATH)

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

build-linux:
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 \
	CC=x86_64-linux-musl-gcc \
	go build -tags 'osusergo netgo static_build' \
	-ldflags="-linkmode external -extldflags '-static' -s -w" \
	-o $(BINARY_NAME)-linux $(MAIN_PATH)
