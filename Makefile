VERSION=$(shell cat version)

install:
	@go get ./...

clean:
	@rm -rf ./bin/*
	@go clean

format:
	@go fmt

test: format
	@go test
	@go build

build: format
	@which gox > /dev/null || go get github.com/mitchellh/gox
	@gox -ldflags="-X main.Version=$(VERSION)" -output="./bin/veracode-cli_{{.OS}}_{{.Arch}}" ./...

version:
	@echo $(VERSION)

.PHONY: install clean format build version
