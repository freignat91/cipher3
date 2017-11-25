.PHONY: all clean install fmt check version build run test

SHELL := /bin/sh
BASEDIR := $(shell echo $${PWD})

# build variables (provided to binaries by linker LDFLAGS below)
VERSION := 0.0.1

LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"

# ignore vendor directory for go files
SRC := $(shell find . -type f -name '*.go' -not -path './vendor/*' -not -path './.git/*')

# for walking directory tree (like for proto rule)
DIRS = $(shell find . -type d -not -path '.' -not -path './vendor' -not -path './vendor/*' -not -path './.git' -not -path './.git/*')

# generated files that can be cleaned
GENERATED := $(shell find . -type f -name '*.pb.go' -not -path './vendor/*' -not -path './.git/*')

# ignore generated files when formatting/linting/vetting
CHECKSRC := $(shell find . -type f -name '*.go' -not -name '*.pb.go' -not -path './vendor/*' -not -path './.git/*')

OWNER := freignat91
PROJECT := cipher3
NAME :=  cipher3

REPO := github.com/$(OWNER)/$(PROJECT)

all: version check install

version:
	@echo "version: $(VERSION)"

clean:
	@rm -rf $(GENERATED)

install:
	@go install $(LDFLAGS) $(REPO)/$(NAME)

fmt:
	@gofmt -s -l -w $(CHECKSRC)

check:
	@for d in $$(go list ./... | grep -v /vendor/); do golint $${d} | sed '/pb\.go/d'; done
	@go tool vet ${CHECKSRC}

test:   
	@go test ./tests -v
