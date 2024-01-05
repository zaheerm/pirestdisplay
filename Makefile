# Makefile to build Go application for Linux ARM64

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOGET=$(GOCMD) get
BINARY_NAME=activatedisplay
BINARY_LINUX_ARM64=$(BINARY_NAME)_linux_arm64

# Build for Linux ARM64
all: build-linux-arm64

build-linux-arm64:
	@CGO_ENABLED=0 GOOS=linux GOARCH=arm64 $(GOBUILD) -o $(BINARY_LINUX_ARM64) -v

local:
	@$(GOBUILD) -o $(BINARY_NAME) -v

clean:
	@$(GOCLEAN)
	rm -f $(BINARY_LINUX_ARM64)
	rm -f $(BINARY_NAME)

run:
	@./$(BINARY_NAME)
