// Copyright (C) 2017-2018 Betalo AB - All Rights Reserved

package sqlstatement

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"log"
	"regexp"
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

var regexpSelectAsterisk = regexp.MustCompile(`(?i)SELECT\s+\*`)

// findMalformedSQLStatements will find variable that are of type string and
// check if it contain a sql command. The linter rule is intentionally made
// strict so if we get false positive you can change the rule to make it more
// relaxed. One example of this that is uses contains instead of prefix.
//
// In addition the rule will check that no SELECT statement selects all columns,
// i.e. `SELECT *`.
//
// Caveats: Two concatenated strings are not supported for the sql statement.
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
					report(out, pos, "sql query contain tabs")
				}
				if regexpSelectAsterisk.FindStringIndex(basicLit.Value) != nil {
					report(out, pos, "sql query selects '*'")
				}
			}
		}
		return true
	})
}

func report(out io.Writer, position token.Position, offense string) {
	fmt.Fprintf(out, "%s: %s\n", position, offense)
}
