GIT_VERSION=$(shell git rev-parse HEAD)
TAG_VERSION=$(shell git tag -l --contains $$GIT_VERSION | tail -1)

BUILD_DIR=bin

GOLIST=go list ./...

default: test

test:
	go test -cover -v ./...

build:
	go build -i

.PHONY: build test

deps:
	go get -u github.com/golang/lint/golint
    #dep should also be installed, but globally.

lint:
	for pkg in $$($(GOLIST) ./...) ; do \
		golint $$pkg ; \
	done

cover:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out

ineffassign:
	go get -u github.com/gordonklaus/ineffassign/...
	ineffassign .

ci-lint:
	golangci-lint run

