// Copyright (C) 2017 Betalo AB - All Rights Reserved

package engine

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/betalo-sweden/go-corp-linter/internal/rule/imports"
	"github.com/betalo-sweden/go-corp-linter/internal/rule/writeerror"
)

var rules = []func(fp string, out io.Writer) error{
	imports.ProcessFile,
	writeerror.ProcessFile,
}

// Process walks a given sequence of directories and tries to identify rule
// violations.
func Process(dirs []string, out io.Writer, verboseMode bool) error {
	for _, dir := range dirs {
		err := filepath.Walk(dir, process(dir, out, verboseMode))
		if err != nil {
			return err
		}
	}
	return nil
}

func process(root string, out io.Writer, verbose bool) func(fp string, fi os.FileInfo, err error) error {
	return func(fp string, fi os.FileInfo, err error) error {
		if verbose {
			log.Println("Debug: processing", fp)
		}

		if err != nil {
			log.Println("Error:", err)
			return nil
		}

		if fi.IsDir() {
			if fp == root {
				return nil
			}
			return filepath.SkipDir
		}

		// Excludes

		if strings.Contains(fp, "vendor/") {
			return nil
		}

		if !strings.HasSuffix(fi.Name(), ".go") {
			return nil
		}

		// Rules

		for _, rule := range rules {
			if err = rule(fp, out); err != nil {
				return err
			}
		}

		return nil
	}
}
