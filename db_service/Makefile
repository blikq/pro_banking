.DEFAULT_GOAL := build

# Formatting section
fmt:
	golangci-lint run


.PHONY: fmt

# Build section
build: fmt
	go run ./cmd .
.PHONY: build
