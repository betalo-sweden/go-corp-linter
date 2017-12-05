// Copyright (C) 2017 Betalo AB - All Rights Reserved

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
	pkg string
	loc location
}

type location struct {
	filename string
	line     int
	col      int
}

func parse(f *ast.File, fset *token.FileSet) []importStmt {
	var imports []importStmt

	for _, i := range f.Imports {
		pos := i.Path.ValuePos
		loc := location{
			filename: fset.Position(pos).Filename,
			line:     fset.Position(pos).Line,
			col:      fset.Position(pos).Column,
		}
		imports = append(imports, importStmt{
			pkg: strings.Trim(i.Path.Value, `"`),
			loc: loc,
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
			stdlibs[i].loc.line = stdlibs[i-1].loc.line + 1
		}
	}
	if len(stdlibs) > 0 && len(others) > 0 {
		others[0].loc.line = stdlibs[len(stdlibs)-1].loc.line + 2
	}
	for i := range others {
		if i > 0 {
			others[i].loc.line = others[i-1].loc.line + 1
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
	fmt.Fprintf(out, "%s:%d:%d: incorrectly sorted import package: %s\n",
		i.loc.filename, i.loc.line, i.loc.col, i.pkg)
}
