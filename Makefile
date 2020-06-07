.PHONY: build test get

VERSION := $(shell git rev-parse --short HEAD)
BRANCH := $(shell git rev-parse --abbrev-ref HEAD)

build:
	go build -ldflags "-X cmd.Version=$(VERSION) -X cmd.Branch=$(BRANCH)" \
	-o rcon-cli ./cmd
	
test:
	go test -v ./pkg/...

get:
	go get ./pkg/...
	go mod verify
