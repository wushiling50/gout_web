PROJECT="gout_web"
GO ?= go
GOFMT ?= gofmt "-s"
GOFILES := $(shell find . -name "*.go")
BINARY := web
VERSION := "1.20.5"

default:
	@echo ${PROJECT}

.PHONY:env
env:
	@${GO} env

.PHONY:ver
ver:
	@echo ${VERSION}

.PHONY: benchmark
benchmark:
	@${GO} test -bench .

.PHONY: fmt
fmt:
	@${GOFMT} -w $(GOFILES)

.PHONY: build
build:
	@$(GO) build -o ${BINARY} -tags=latest

.PHONY: clean
clean:
	@rm -rf .