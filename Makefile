.PHONY: build test get

VERSION := $(shell git rev-parse --short HEAD)
BRANCH := $(shell git rev-parse --abbrev-ref HEAD)

build:
	go build -ldflags "-X pkg.Version=$(VERSION) -X pkg.Branch=$(BRANCH)" \
	-o steam-gameserver
	
test:
	go test -v ./pkg/...

get:
	go get ./pkg/...
	go mod verify
