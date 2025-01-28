# Makefile for Nox

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOCLEAN=$(GOCMD) clean
GOFMT=$(GOCMD) fmt
GOVET=$(GOCMD) vet

# Binary name
BINARY_NAME=nox

# Build flags
VERSION=1.0.0
BUILD_TIME=$(shell date +%F_%T)
GIT_COMMIT=$(shell git rev-parse HEAD)
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT)"

# Default target
.PHONY: all
all: clean fmt vet test build

# Build the binary
.PHONY: build
build:
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) cmd/nox.go

# Run tests
.PHONY: test
test:
	$(GOTEST) -v ./...

# Format code
.PHONY: fmt
fmt:
	$(GOFMT) ./...

# Run go vet
.PHONY: vet
vet:
	$(GOVET) ./...

# Clean build files
.PHONY: clean
clean:
	rm -f $(BINARY_NAME)
	$(GOCLEAN)

# Cross compilation
.PHONY: build-linux
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME)_linux_amd64 cmd/nox.go

.PHONY: build-windows
build-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME)_windows_amd64.exe cmd/nox.go

.PHONY: build-darwin
build-darwin:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME)_darwin_amd64 cmd/nox.go

# Build all platforms
.PHONY: build-all
build-all: build-linux build-windows build-darwin

# Install binary
.PHONY: install
install: build
	mv $(BINARY_NAME) $(GOPATH)/bin/

# Help target
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all        : Clean, format, test and build"
	@echo "  build      : Build binary"
	@echo "  test       : Run tests"
	@echo "  fmt        : Format code"
	@echo "  vet        : Run go vet"
	@echo "  clean      : Clean build files"
	@echo "  build-all  : Build for all platforms"
	@echo "  install    : Install binary"