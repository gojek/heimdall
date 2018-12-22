.PHONY: all
all: build test

ALL_PACKAGES=$(shell go list ./... | grep -v "vendor")

setup:
	mkdir -p $(GOPATH)/bin
	go get -u golang.org/x/lint/golint

compile:
	mkdir -p out/
	env GO111MODULE=on go build -race ./...

build: compile fmt vet lint

fmt:
	env GO111MODULE=on go fmt ./...

vet:
	env GO111MODULE=on go vet ./...

lint:
	env GO111MODULE=on golint -set_exit_status $(ALL_PACKAGES)

test: fmt vet build
	GO111MODULE=on ENVIRONMENT=test go test -race ./...

test-cover-html:
	@echo "mode: count" > coverage-all.out

	$(foreach pkg, $(ALL_PACKAGES),\
	ENVIRONMENT=test go test -coverprofile=coverage.out -covermode=count $(pkg);\
	tail -n +2 coverage.out >> coverage-all.out;)
	go tool cover -html=coverage-all.out -o out/coverage.html
