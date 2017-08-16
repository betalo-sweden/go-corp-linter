# Copyright (C) 2017 Betalo AB - All Rights Reserved

.PHONY: all
all: build

.PHONY: build
build:
	go build ./

.PHONY: copyright
copyright:
	find ./cmd ./pkg -type f -name '*.go' -exec grep -H -m 1 . {} \; | \
	    grep -v '/vendor/' | \
	    (! grep -v "// Copyright (C) .*$$(date +%Y) Betalo AB - All Rights Reserved")
