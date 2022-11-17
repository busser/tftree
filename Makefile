.DEFAULT_TARGET=help
VERSION:=$(shell cat VERSION)

## help: Display list of commands
.PHONY: help
help: Makefile
	@sed -n 's|^##||p' $< | column -t -s ':' | sed -e 's|^| |'

## build: Build tftree binary
.PHONY: build
build: fmt vet
	go build -o bin/tftree

## fmt: Format source code
.PHONY: fmt
fmt:
	go fmt ./...

## vet: Vet source code
.PHONY: vet
vet:
	go vet ./...

## test: Run unit tests
.PHONY: test
test:
	go test ./...

## release: Release a new version
.PHONY: release
release: test
	git tag -a "$(VERSION)" -m "$(VERSION)"
	git push origin "$(VERSION)"
	goreleaser release --rm-dist
