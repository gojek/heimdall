.PHONY: all
all: build test

ALL_PACKAGES=$(shell go list ./... | grep -v "vendor")

setup:
	mkdir -p $(GOPATH)/bin
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.8.0
	go install github.com/mattn/goveralls@v0.0.12

compile:
	mkdir -p out/
	go build -race ./...

build: compile fmt lint

fmt:
	go fmt ./...

lint:
	golangci-lint run ./...

test: fmt build
	ENVIRONMENT=test go test -race -covermode=atomic -coverprofile=coverage.out ./...

test-cover-html:
	@echo "mode: count" > coverage-all.out

	$(foreach pkg, $(ALL_PACKAGES),\
		ENVIRONMENT=test go test -race -coverprofile=coverage.out -covermode=count $(pkg);\
		tail -n +2 coverage.out >> coverage-all.out;)
	go tool cover -html=coverage-all.out -o out/coverage.html
