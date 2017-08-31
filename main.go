// Copyright (C) 2017 Betalo AB - All Rights Reserved

package main

import (
	"flag"
	"os"

	"github.com/betalo-sweden/go-corp-linter/internal/engine"
)

func main() {
	flag.Parse()
	dirs := flag.Args()

	if len(dirs) == 0 {
		dirs = []string{"."}
	}

	engine.Process(dirs, os.Stdout)
}
