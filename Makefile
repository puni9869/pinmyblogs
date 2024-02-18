.PHONY: build
build:
	go build ./cmd/main.go

.PHONY: server
server:
	air

.PHONY: test
test:
	go test ./... -cover

.PHONY: lint
lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.56.2
	golangci-lint run ./...

.PHONY: govulncheck
govulncheck:
	go install golang.org/x/vuln/cmd/govulncheck@latest
	govulncheck ./...

.PHONY: vet
vet:
	go vet ./...
