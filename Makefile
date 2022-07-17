.PHONY: run.docker run.local deps clean test build.docker build.local build.linux.armv8 build.linux.armv7 build.linux build.osx build.windows

ifneq ("$(wildcard .env)","")
ENV_FILE = .env
else
ENV_FILE = .env.example
endif

GOPKGS = $(shell go list ./... | grep -v /vendor/)
BIN_OUTPUT = bin/brcep

default: build.local

run.local: deps
	go run $(GOPKGS) server.go

deps:
	go mod vendor

test:
	go test -v $(GOPKGS) -coverprofile=coverage.txt -covermode=atomic

clean:
	rm -rf vendor
	rm -rf bin
