# Copyright (C) 2017 Betalo AB - All Rights Reserved

.PHONY: all
all: build

.PHONY: build
build:
	go build

.PHONY: copyright
copyright:
	@find . -type f -name '*.go' -exec grep -H -m 1 . {} \; \
	    | grep -v '/vendor/' \
	    | grep -v '/testdata/' \
	    | (! grep -v "// Copyright (C) .*20[1-9][0-9] \(Betalo AB\|P.F.C. AB\) - All Rights Reserved") \
	    | (! grep -v "// Code generated by .*")

.PHONY: deps
deps:
	go get -u github.com/golang/dep/cmd/dep

	go get -u github.com/betalo-sweden/go-corp-linter
	go get -u github.com/alecthomas/gometalinter
	gometalinter --install

.PHONY: deps-ensure
deps-ensure:
	dep ensure

.PHONY: install
install:
	go install

.PHONY: lint
lint:
	@if [ $$(gofmt -l . | grep -v vendor/ | wc -l) != 0 ]; then \
	    echo "gofmt: code not formatted" \
	    ; gofmt -l . \
	    | grep -v '/vendor/' \
	    | grep -v '/testdata/' \
	    ; exit 1; \
	fi

	@gometalinter \
	    --vendor \
	    --tests \
	    --disable=gocyclo \
	    --disable=dupl \
	    --disable=deadcode \
	    --disable=gotype \
	    --disable=maligned \
	    --disable=interfacer \
	    --disable=varcheck \
	    --disable=gas \
	    --disable=megacheck \
	    --disable=gosec \
	    --disable=errcheck \
	    --enable=go-corp-linter \
	    --linter='go-corp-linter:go-corp-linter:PATH:LINE:COL:MESSAGE' \
	    ./internal/...

.PHONY: test
test:
	go test -v -race ./...
