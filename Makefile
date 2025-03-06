.PHONY: all test run gofumpt lint swagger

all: test gofumpt lint swagger

test:
	go test -count=1 -race -covermode=atomic ./...


lint:
	golangci-lint run ./...

gofumpt:
	gofumpt -l -w .
