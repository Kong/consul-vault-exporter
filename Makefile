.PHONY: all unit e2e bin tests

bin:
	go build -ldflags="-X main.version=$(shell git describe --always --long --dirty)"

unit:
	go test -race $(shell go list ./... | grep -v e2e)

e2e:
	go test -race $(shell go list ./... | grep e2e)

vet:
	go vet ./...

tests: vet unit e2e

all: unit e2e bin
