// Copyright (C) 2017 Betalo AB - All Rights Reserved

package sqlstatement

import (
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

func findMalformedSQLStatements(f ast.Node, fset *token.FileSet, out io.Writer) {
	ast.Inspect(f, func(node ast.Node) bool {
		if assignStmt, ok := node.(*ast.AssignStmt); ok {
			if len(assignStmt.Lhs) != 1 {
				return true
			}
			if ident, ok := assignStmt.Lhs[0].(*ast.Ident); ok {
				if ident.Obj.Kind != ast.Var || ident.Obj.Name != "stmt" {
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
			}
		}
		return true

		//*ast.AssignStmt {
		//	Lhs: []ast.Expr (len = 1) {
		//		0: *ast.Ident {
		//			NamePos: 4:5
		//			Name: "stmt"
		//			Obj: *ast.Object {
		//				Kind: var
		//				Name: "stmt"
		//				Decl: *(obj @ 27)
		//			}
		//		}
		//	}
		//	TokPos: 4:10
		//	Tok: :=
		//	Rhs: []ast.Expr (len = 1) {
		//		0: *ast.BasicLit {
		//			ValuePos: 4:13
		//			Kind: STRING
		//			Value: "`\nSELECT\n    foo\nFROM\n   bar\n`"
		//		}
		//	}
		//}
	})
}
