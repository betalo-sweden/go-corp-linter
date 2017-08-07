// Copyright (C) 2017 Betalo AB - All Rights Reserved

package main

import (
	"flag"
	"fmt"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func main() {
	flag.Parse()
	dir := flag.Arg(0)

	filepath.Walk(dir, findImportViolations)
}

func findImportViolations(fp string, fi os.FileInfo, err error) error {
	if err != nil {
		log.Println("Error:", err)
		return nil
	}
	if fi.IsDir() {
		return nil
	}
	if strings.Contains(fp, "/vendor/") {
		return nil
	}
	if !strings.HasSuffix(fi.Name(), ".go") {
		return nil
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, fp, nil, parser.ImportsOnly)
	if err != nil {
		log.Println("Error: parsing dir:", err)
		return nil
	}

	var hasStdlibPkgs bool
	var firstImport bool
	var firstStdlibPkg bool
	var stdlibBlockEnd int
	var lastStdlibPkg string

	for _, i := range f.Imports {
		pkg := strings.Trim(i.Path.Value, `"`)
		pos := i.Path.ValuePos
		line := fset.Position(pos).Line
		col := fset.Position(pos).Column
		filename := fset.Position(pos).Filename

		if !firstImport {
			firstImport = true
			if isStdlib(pkg) {
				firstStdlibPkg = true
				stdlibBlockEnd = line
			}
		}

		if isStdlib(pkg) {
			hasStdlibPkgs = true

			if hasStdlibPkgs && !firstStdlibPkg {
				fmt.Printf("%s:%d:%d:error: First package no stdlib package: %s\n", filename, line, col, pkg)
			} else if stdlibBlockEnd+1 < line {
				fmt.Printf("%s:%d:%d:error: stdlib packages not grouped together: %s\n", filename, line, col, pkg)
			} else if lastStdlibPkg > pkg {
				fmt.Printf("%s:%d:%d:error: stdlib packages not sorted: %s\n", filename, line, col, pkg)
			}

			stdlibBlockEnd = line
			lastStdlibPkg = pkg
		}
	}

	return nil
}

func isStdlib(path string) bool {
	for _, p := range stdlibPkgs {
		if p == path {
			return true
		}
	}
	return false
}

func init() {
	findStdlibPkgs()
}

var stdlibPkgs []string

// findStdlibPkgs tries to get list of stdlib packages by reading GOROOT
//
// Courtesy: https://github.com/divan/depscheck/blob/master/package.go#L72
//
// This approach is used by "go list std" tool
// Based on go/cmd function matchPackages (https://golang.org/src/cmd/go/main.go#L553)
// and listStdPkgs function from https://golang.org/src/go/build/deps_test.go#L420
func findStdlibPkgs() {
	goroot := runtime.GOROOT()

	src := filepath.Join(goroot, "src") + string(filepath.Separator)
	walkFn := func(path string, fi os.FileInfo, err error) error {
		if err != nil || !fi.IsDir() || path == src {
			return nil
		}

		base := filepath.Base(path)
		if strings.HasPrefix(base, ".") || strings.HasPrefix(base, "_") || base == "testdata" {
			return filepath.SkipDir
		}

		name := filepath.ToSlash(path[len(src):])
		if name == "builtin" || name == "cmd" || strings.Contains(name, ".") {
			return filepath.SkipDir
		}

		stdlibPkgs = append(stdlibPkgs, name)
		return nil
	}

	if err := filepath.Walk(src, walkFn); err != nil {
		log.Fatalln("Error:", err)
	}
}
