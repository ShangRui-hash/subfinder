# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOMOD=$(GOCMD) mod
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
    
all: build
build:
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -v -ldflags="-extldflags=-static" -o "subfinder" cmd/subfinder/main.go
test: 
		$(GOTEST) -v ./...
tidy:
		$(GOMOD) tidy
