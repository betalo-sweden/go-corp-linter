// Copyright (C) 2017-2018 Betalo AB - All Rights Reserved

package writeerror

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"log"
)

// ProcessFile checks for missing return statements immediately after a call to
// `writeError` function.
func ProcessFile(fp string, out io.Writer) error {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, fp, nil, 0)
	if err != nil {
		log.Println("Error: parsing file:", err)
		return err
	}

	findMissingReturn(f, fset, out)
	return nil
}

func findMissingReturn(f ast.Node, fset *token.FileSet, out io.Writer) {
	var v NodeVisitor = newFindMissingReturnsVisitor(fset, out)
	ast.Inspect(f, v.Visit)
}

// NodeVisitor visits nodes in an AST.
type NodeVisitor interface {
	Visit(n ast.Node) bool
}

func newFindMissingReturnsVisitor(fset *token.FileSet, out io.Writer) *FindMissingReturnsVisitor {
	return &FindMissingReturnsVisitor{
		fset:        fset,
		out:         out,
		stateLevels: map[int]int{},
	}
}

var _ NodeVisitor = (*FindMissingReturnsVisitor)(nil)

// FindMissingReturnsVisitor visits nodes in an AST trying to find a call to a
// `writeError` function with missing subsequent return statement.
type FindMissingReturnsVisitor struct {
	fset        *token.FileSet
	out         io.Writer
	stateLevels map[int]int
	level       int
	state       int
	lastFound   token.Position
}

// Visit implements the NodeVisitor interface.
//
// State machine:
//
//   0 found nothing
//   1 found *ast.ExprStmt
//   2   found *ast.CallExpr
//   3     found *ast.Ident: "writeError"
//   4     found first  parameter: http.ResponseWriter
//   5     found second parameter: *http.Request
//   6     found third  parameter: error
//   7   found *ast.ReturnStmt with no results/arguments
//
// Currently, for simplicity and sufficiency, only states 0, 2, and 7 are
// implemented.
func (v *FindMissingReturnsVisitor) Visit(n ast.Node) bool {
	// fmt.Printf("Debug: %d %T %+v %#v\n", level, n, n, n)

	if n == nil { // Leaf node
		v.level--
		return false
	}
	v.level++

	switch v.state {
	case 0, 1:
		// fmt.Println("Debug: Test for function call")
		if callExpr, ok := n.(*ast.CallExpr); ok && len(callExpr.Args) == 3 {
			if ident, ok := callExpr.Fun.(*ast.Ident); ok && ident.Name == "writeError" {
				// fmt.Println("Debug: Found CallExpr")
				v.state = 7
				v.stateLevels[v.state] = v.level
				v.lastFound = v.fset.Position(n.Pos())
				return false
			}
		}
	case 7:
		// fmt.Println("Debug: Test for return")
		if v.level == v.stateLevels[v.state] {
			if returnStmt, ok := n.(*ast.ReturnStmt); ok && len(returnStmt.Results) == 0 {
				// fmt.Println("Info: Found return statement after writeError call")
				v.state = 0
				return false
			}
		}

		position := v.lastFound
		position.Line++ // Report the next line
		fmt.Fprintf(v.out, "%s: missing return statement after writeError call\n", position)
	}

	v.state = 0 // Reset
	return true
}
