.PHONYthub.com/mattn/goveralls: all
all: build test coverage

ALL_PACKAGES=$(shell go list ./... | grep -v "vendor")

setup:
	mkdir -p $(GOPATH)/bin
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.8
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
	ENVIRONMENT=test go test -race ./...

coverage:
	ENVIRONMENT=test goveralls -service=github-actions

test-cover-html:
	@echo "mode: count" > coverage-all.out

	$(foreach pkg, $(ALL_PACKAGES),\
	ENVIRONMENT=test go test -coverprofile=coverage.out -covermode=count $(pkg);\
	tail -n +2 coverage.out >> coverage-all.out;)
	go tool cover -html=coverage-all.out -o out/coverage.html
