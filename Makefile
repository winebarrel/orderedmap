.PHONY: all
all: vet test

.PHONY: vet
vet:
	go vet ./...

.PHONY: test
test:
	go test -race -v ./...

.PHONY: lint
lint:
	golangci-lint run
