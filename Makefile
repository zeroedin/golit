# golit Makefile
# Build, test, and release targets

BINARY_NAME := golit
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-s -w -X main.version=$(VERSION)"
PLATFORMS := linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64 windows/arm64

.PHONY: all build test clean install cross-compile lint fmt help

## help: Print this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@sed -n 's/^## //p' $(MAKEFILE_LIST) | column -t -s ':' | sed 's/^/  /'

## build: Build the golit binary for the current platform
build:
	@mkdir -p dist
	CGO_ENABLED=0 go build $(LDFLAGS) -o dist/$(BINARY_NAME) ./cmd/golit/

## test: Run all tests
test:
	go test ./... -count=1

## test-verbose: Run all tests with verbose output
test-verbose:
	go test ./... -v -count=1

## lint: Run go vet
lint:
	go vet ./...

## fmt: Format all Go source files
fmt:
	gofmt -w -s .

## clean: Remove build artifacts
clean:
	rm -rf dist/

## install: Install golit to $GOPATH/bin
install:
	CGO_ENABLED=0 go install $(LDFLAGS) ./cmd/golit/

## cross-compile: Build for all platforms (output to dist/)
cross-compile: clean
	@mkdir -p dist
	@for platform in $(PLATFORMS); do \
		os=$${platform%/*}; \
		arch=$${platform#*/}; \
		output=dist/$(BINARY_NAME)-$${os}-$${arch}; \
		if [ "$$os" = "windows" ]; then output=$${output}.exe; fi; \
		echo "Building $$output..."; \
		GOOS=$$os GOARCH=$$arch CGO_ENABLED=0 go build $(LDFLAGS) -o $$output ./cmd/golit/ || exit 1; \
	done
	@echo "Built binaries:"
	@ls -lh dist/

## deps: Download and tidy Go dependencies
deps:
	go mod download
	go mod tidy

## npm-deps: Install npm dependencies (for component bundling tests)
npm-deps:
	npm install

## all: Build, test, and lint
all: fmt lint test build
