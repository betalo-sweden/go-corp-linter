// Copyright (C) 2017 Betalo AB - All Rights Reserved

package main

import (
	"flag"
	"os"

	"github.com/betalo-sweden/go-corp-linter/internal/engine"
)

func main() {
	verboseFlag := flag.Bool("v", false, "verbose output mode")
	verboseLongFlag := flag.Bool("verbose", false, "verbose output mode")
	flag.Parse()
	dirs := flag.Args()

	verboseMode := *verboseFlag || *verboseLongFlag

	if len(dirs) == 0 {
		dirs = []string{"."}
	}

	engine.Process(dirs, os.Stdout, verboseMode)
}
