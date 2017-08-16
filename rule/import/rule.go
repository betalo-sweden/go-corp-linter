// Copyright (C) 2017 Betalo AB - All Rights Reserved

package imports

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	firstPackage = "first_package_no_stdlib_package"
	notGrouped   = "stdlib_packages_not_grouped_together"
	notSorted    = "stdlib_packages_not_sorted"
)

type hit struct {
	filename string
	line     int
	col      int
	pkg      string
	problem  string
}

// Rule checks for import errors in the go code and reports them
func Rule(fp string) error {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, fp, nil, parser.ImportsOnly)
	if err != nil {
		log.Println("Error: parsing dir:", err)
		return nil
	}
	report(processImportRule(f, fset))
	return nil
}

func report(hits []hit) {
	for _, item := range hits {
		if item.problem == firstPackage {
			fmt.Printf("%s:%d:%d:error: First package no stdlib package: %s\n", item.filename, item.line, item.col, item.pkg)
		} else if item.problem == notGrouped {
			fmt.Printf("%s:%d:%d:error: stdlib packages not grouped together: %s\n", item.filename, item.line, item.col, item.pkg)
		} else if item.problem == notSorted {
			fmt.Printf("%s:%d:%d:error: stdlib packages not sorted: %s\n", item.filename, item.line, item.col, item.pkg)
		}
	}
}

func processImportRule(f *ast.File, fset *token.FileSet) []hit {
	var hits []hit
	var possibleFirstPackageHit *hit

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
			} else {
				possibleFirstPackageHit = &hit{filename: filename, line: line, col: col, pkg: pkg, problem: firstPackage}
			}
		}

		if isStdlib(pkg) {
			hasStdlibPkgs = true

			if stdlibBlockEnd+1 < line && possibleFirstPackageHit == nil {
				hits = append(hits, hit{filename: filename, line: line, col: col, pkg: pkg, problem: notGrouped})
			} else if lastStdlibPkg > pkg {
				hits = append(hits, hit{filename: filename, line: line, col: col, pkg: pkg, problem: notSorted})
			}

			stdlibBlockEnd = line
			lastStdlibPkg = pkg
		}
	}
	if hasStdlibPkgs && !firstStdlibPkg {
		hits = append(hits, *possibleFirstPackageHit)
	}
	return hits
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
