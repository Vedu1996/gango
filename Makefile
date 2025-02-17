BINARY_NAME=gango
BINARY_UNIX=$(BINARY_NAME)
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOLIST=$(GOCMD) list
GOTEST=$(GOCMD) test
GOMOD=GO111MODULE=on $(GOCMD) mod
GOGET=GO111MODULE=on $(GOCMD) get

.PHONY: all
all: clean dependencies build

.PHONY: linux
linux: clean dependencies build-linux

.PHONY: clean
clean:
	@echo "Cleaning..."
	@rm -f $(BINARY_NAME)
	@rm -f $(BINARY_UNIX)
	$(GOMOD) tidy
	$(GOMOD) vendor
	$(GOCLEAN) -i
	@echo  "Done cleaning."

.PHONY: dependencies
dependencies:
	$(GOMOD) tidy
	$(GOMOD) download
	$(GOMOD) vendor

.PHONY: build
build:
	$(GOBUILD) ./...
	$(GOBUILD) -o $(BINARY_NAME) -v cmd/*.go

.PHONY: build-linux
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v  cmd/*.go

.PHONY: build-release
build-release:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -ldflags "-s -w" -o $(BINARY_UNIX) -v cmd/*.go

.PHONY: go-lint
go-lint:
	@echo  "running go-lint ..."
	golangci-lint run -v

.PHONY: build-docker
build-docker:
	make linux
	docker build . -t gango:latest

.PHONY: test-project
test-project:
	@rm -rf ./project
	make all
	./gango generate project
