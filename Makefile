.PHONY: test

GIT_COMMIT ?= $(shell git rev-parse --verify HEAD)
GIT_VERSION ?= $(shell git describe --tags --always --dirty="-dev")
DATE ?= $(shell date -u '+%Y-%m-%d %H:%M UTC')
BUILDER ?= Makefile
VERSION_FLAGS := -X "github.com/chavacava/changelog-lint/main.version=$(GIT_VERSION)" -X "github.com/chavacava/changelog-lint/main.date=$(DATE)" -X "github.com/chavacava/changelog-lint/main.commit=$(GIT_COMMIT)" -X "github.com/chavacava/changelog-lint/main.builtBy=$(BUILDER)"

build:
	@go build -ldflags='$(VERSION_FLAGS)'

test:
	@go test -v -race ./...