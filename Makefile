.PHONY: build test integration-test release

build:
	go build ./cmd/ittconv

test:
	go test ./...

integration-test:
	@echo "TODO: add integration tests when explicitly requested."

release:
	mkdir -p dist
	go build -o dist/ittconv ./cmd/ittconv
