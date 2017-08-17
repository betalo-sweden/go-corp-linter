# Copyright (C) 2017 Betalo AB - All Rights Reserved

.PHONY: all
all: build

.PHONY: build
build:
	go build ./

deps:
	go get -u github.com/golang/dep/cmd/dep

	go get -u github.com/betalo-sweden/go-corp-linter
	go get -u github.com/alecthomas/gometalinter
	gometalinter --install --update

deps-ensure:
	dep ensure
	git checkout vendor/vendor.json

test:
	go test -v -race . ./rule/...

lint:
	gometalinter --vendor --tests --disable=gocyclo --disable=dupl --disable=deadcode --disable=gotype --disable=errcheck --disable=aligncheck --disable=unconvert --disable=interfacer --disable=varcheck --disable=gas --disable=megacheck --enable=go-corp-linter --linter='go-corp-linter:go-corp-linter:PATH:LINE:MESSAGE' ./...

.PHONY: copyright
copyright:
	find ./cmd ./pkg -type f -name '*.go' -exec grep -H -m 1 . {} \; | \
	    grep -v '/vendor/' | \
	    (! grep -v "// Copyright (C) .*$$(date +%Y) Betalo AB - All Rights Reserved")
