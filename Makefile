.PHONY: test
test:
	go test ./... -cover

.PHONY: lint
lint:
	golangci-lint run

.PHONY: govulncheck
govulncheck:
	go install golang.org/x/vuln/cmd/govulncheck@latest
	govulncheck ./...

.PHONY: vet
vet:
	go vet ./...
