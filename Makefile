.PHONY: all
all: vet test

.PHONY: vet
vet:
	go vet ./...

.PHONY: test
test:
	go test -v ./... $(TEST_OPTS)

.PHONY: bench
bench:
	go test -run='0^' -bench . -benchmem

.PHONY: lint
lint:
	golangci-lint run
