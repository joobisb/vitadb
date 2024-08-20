.PHONY: fmt lint test vet security run

fmt:
	goimports -w .

lint:
	golangci-lint run

test:
	go test -v -race -cover ./...

vet:
	go vet ./...

build:
	go build ./...

run:
	go run cmd/server/main.go

cli:
	go run cmd/client/main.go

check: fmt lint test vet security