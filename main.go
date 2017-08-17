// Copyright (C) 2017 Betalo AB - All Rights Reserved

package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/betalo-sweden/go-corp-linter/rule/import"
)

func main() {
	flag.Parse()
	dirs := flag.Args()

	if len(dirs) == 0 {
		dirs = []string{"."}
	}

	for _, dir := range dirs {
		filepath.Walk(dir, findImportViolations)
	}
}

func findImportViolations(fp string, fi os.FileInfo, err error) error {
	if err != nil {
		log.Println("Error:", err)
		return nil
	}
	if fi.IsDir() {
		if fp == "." {
			return nil
		}
			return filepath.SkipDir
		}
	}
	if strings.Contains(fp, "vendor/") {
		return nil
	}
	if !strings.HasSuffix(fi.Name(), ".go") {
		return nil
	}

	if err = imports.Rule(fp); err != nil {
		return nil
	}

	return nil
}
