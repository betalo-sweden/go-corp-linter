// Copyright (C) 2017 Betalo AB - All Rights Reserved

package sqlstatement

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"log"
	"strings"
)

// ProcessFile checks for SQL statements that contain undesired tabs indentation
// instead of spaces.
func ProcessFile(fp string, out io.Writer) error {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, fp, nil, 0)
	if err != nil {
		log.Println("Error: parsing file:", err)
		return err
	}

	findMalformedSQLStatements(f, fset, out)
	return nil
}

// findMalformedSQLStatements will find variable that are of type string and check if it contain a sql command.
// The linter rule is intentionally made strict so if we get false positive you can change the rule to make it more
// relaxed. One example of this i contain instead of prefix.
// Two concatenated strings are not supported for the sql statement.
func findMalformedSQLStatements(f ast.Node, fset *token.FileSet, out io.Writer) {
	ast.Inspect(f, func(node ast.Node) bool {
		if assignStmt, ok := node.(*ast.AssignStmt); ok {
			if len(assignStmt.Lhs) != 1 {
				return true
			}
			if ident, ok := assignStmt.Lhs[0].(*ast.Ident); ok {
				if ident.Obj == nil || ident.Obj.Kind != ast.Var {
					return true
				}
			}

			if len(assignStmt.Rhs) != 1 {
				return true
			}
			if basicLit, ok := assignStmt.Rhs[0].(*ast.BasicLit); ok {
				if basicLit.Kind != token.STRING {
					return true
				}
				if !strings.HasPrefix(basicLit.Value, "`") ||
					!strings.HasSuffix(basicLit.Value, "`") {
					return true
				}
				sqlStatementFound := false
				for _, sqlStatementPrefix := range sqlStatementPrefixes {
					if strings.Contains(basicLit.Value, sqlStatementPrefix) {
						sqlStatementFound = true
						break
					}
				}
				if !sqlStatementFound {
					return true
				}

				pos := fset.Position(assignStmt.TokPos)
				if strings.Contains(basicLit.Value, "\t") {
					report(out, pos)
				}
			}
		}
		return true
	})
}

func report(out io.Writer, position token.Position) {
	fmt.Fprintf(out, "%s: sql query contain tabs\n", position)
}
