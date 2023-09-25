SHELL        := /bin/bash
GOPATH       := $(shell go env GOPATH)

.PHONY: proto envoy test build

test:
	@rm -rf profile.cov coverage.html
	go test -mod=vendor -tags test -cover -covermode=atomic \
	        -coverpkg=$(shell go list ./... | grep -v main | grep -v cmd | tr '\n' ',') \
	        -coverprofile=profile.cov -p 1 -v ./...
	go tool cover -func=profile.cov
	go tool cover -html=profile.cov -o coverage.html
	go vet ./...

build:
	CGO_ENABLED=0 go build --mod=vendor -o build/normalize main.go