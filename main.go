// Copyright (C) 2017-2018 Betalo AB - All Rights Reserved

package main

import (
	"flag"
	"log"
	"os"

	"github.com/betalo-org/go-corp-linter/internal/engine"
)

func init() {
	log.SetFlags(0)
	log.SetPrefix("")
}

func main() {
	verboseFlag := flag.Bool("v", false, "verbose output mode")
	verboseLongFlag := flag.Bool("verbose", false, "verbose output mode")
	flag.Parse()
	dirs := flag.Args()

	verboseMode := *verboseFlag || *verboseLongFlag

	if len(dirs) == 0 {
		dirs = []string{"."}
	}

	err := engine.Process(dirs, os.Stdout, verboseMode)
	if err != nil {
		log.Fatalln("Error:", err)
	}
}
