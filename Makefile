.PHONY: all
all: build test

GLIDE_NOVENDOR=$(shell glide novendor)
ALL_PACKAGES=$(shell go list ./... | grep -v "vendor")

setup:
	mkdir -p $(GOPATH)/bin
	curl https://glide.sh/get | sh
	go get -u github.com/golang/lint/golint

build-deps:
	glide install

update-deps:
	glide update

compile:
	mkdir -p out/
	go build -race $(GLIDE_NOVENDOR)

build: build-deps compile fmt vet lint

fmt:
	go fmt $(GLIDE_NOVENDOR)

vet:
	go vet $(GLIDE_NOVENDOR)

lint:
	@for p in $(GLIDE_NOVENDOR); do \
		echo "==> Linting $$p"; \
		golint -set_exit_status $$p; \
	done

test: build-deps fmt vet build
	ENVIRONMENT=test go test -race $(GLIDE_NOVENDOR)

test-cover-html:
	@echo "mode: count" > coverage-all.out

	$(foreach pkg, $(ALL_PACKAGES),\
	ENVIRONMENT=test go test -coverprofile=coverage.out -covermode=count $(pkg);\
	tail -n +2 coverage.out >> coverage-all.out;)
	go tool cover -html=coverage-all.out -o out/coverage.html
