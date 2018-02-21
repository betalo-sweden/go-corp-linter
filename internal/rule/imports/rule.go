// Copyright (C) 2017-2018 Betalo AB - All Rights Reserved

package imports

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"log"
	"sort"
	"strings"
)

// ProcessFile checks for import statement violations in source code files and
// reports them.
func ProcessFile(fp string, out io.Writer) error {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, fp, nil, parser.ImportsOnly)
	if err != nil {
		log.Println("Error: parsing file:", err)
		return err
	}

	processImports(f, fset, out)
	return nil
}

func processImports(f *ast.File, fset *token.FileSet, out io.Writer) {
	givenImports := parse(f, fset)
	expectedImports := ordered(givenImports)
	idx := compare(givenImports, expectedImports)
	if idx >= 0 {
		report(out, givenImports[idx])
	}
}

type importStmt struct {
	pkg      string
	position token.Position
}

func parse(f *ast.File, fset *token.FileSet) []importStmt {
	var imports []importStmt

	for _, i := range f.Imports {
		pos := i.Path.ValuePos
		position := fset.Position(pos)
		imports = append(imports, importStmt{
			pkg:      strings.Trim(i.Path.Value, `"`),
			position: position,
		})
	}

	return imports
}

func ordered(imports []importStmt) []importStmt {
	var stdlibs, others []importStmt

	// Partiion
	for _, i := range imports {
		if i.pkg == "" {
			continue
		} else if isStdlib(i.pkg) {
			stdlibs = append(stdlibs, i)
		} else {
			others = append(others, i)
		}
	}

	// Sort
	sort.Slice(stdlibs, func(i, j int) bool { return stdlibs[i].pkg < stdlibs[j].pkg })
	sort.Slice(others, func(i, j int) bool { return others[i].pkg < others[j].pkg })

	// Compact
	for i := range stdlibs {
		if i > 0 {
			stdlibs[i].position.Line = stdlibs[i-1].position.Line + 1
		}
	}
	if len(stdlibs) > 0 && len(others) > 0 {
		others[0].position.Line = stdlibs[len(stdlibs)-1].position.Line + 2
	}
	for i := range others {
		if i > 0 {
			others[i].position.Line = others[i-1].position.Line + 1
		}
	}

	// Merge
	return append(stdlibs, others...)
}

func compare(given, expected []importStmt) int {
	for i, g := range given {
		if g != expected[i] {
			return i
		}
	}
	return -1
}

func report(out io.Writer, i importStmt) {
	fmt.Fprintf(out, "%s: incorrectly sorted import package: %s\n", i.position, i.pkg)
}
