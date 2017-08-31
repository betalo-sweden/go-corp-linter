// Copyright (C) 2017 Betalo AB - All Rights Reserved

package engine

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/betalo-sweden/go-corp-linter/internal/rule/imports"
)

// Process walks a given sequence of directories and tries to identify rule
// violations.
func Process(dirs []string, out io.Writer) {
	for _, dir := range dirs {
		filepath.Walk(dir, findImportViolations(dirs[0], out))
	}
}

func findImportViolations(root string, out io.Writer) func(fp string, fi os.FileInfo, err error) error {
	return func(fp string, fi os.FileInfo, err error) error {
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
		if strings.Contains(fp, "vendor/") {
			return nil
		}
		if !strings.HasSuffix(fi.Name(), ".go") {
			return nil
		}

		if err = imports.ProcessFile(fp, out); err != nil {
			return nil
		}

		return nil
	}
}
