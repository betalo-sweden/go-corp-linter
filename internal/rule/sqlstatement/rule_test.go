// Copyright (C) 2017 Betalo AB - All Rights Reserved

package sqlstatement

import (
	"bufio"
	"bytes"
	"go/parser"
	"go/token"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProcessImports(t *testing.T) {
	var testcases = []struct {
		name     string
		given    string
		expected string
	}{
		{
			name: "tab",
			given: `
			SELECT
			id,
			WHERE id=$1`,
			expected: "main.go:4:7: sql query contain tabs",
		},
		{
			name: "tab on last line",
			given: `
            SELECT
            id,
			WHERE id=$1`,
			expected: "main.go:4:7: sql query contain tabs",
		},
		{
			name: "tab on last line 2",
			given: `
            SELECT
            id,
            WHERE id=$1
			`,
			expected: "main.go:4:7: sql query contain tabs",
		},
		{
			name: "tab on first line",
			given: `
			SELECT
            id,
            WHERE id=$1`,
			expected: "main.go:4:7: sql query contain tabs",
		},
		{
			name: "space",
			given: `
            SELECT
            id,
            WHERE id=$1`,
			expected: "",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			fset := token.NewFileSet()
			source := `package main

func main() {
	stmt := ` + "`" + tc.given + "`" + `
}`
			f, err := parser.ParseFile(fset, "main.go", source, 0)
			require.NoError(t, err)

			var b bytes.Buffer
			w := bufio.NewWriter(&b)
			findMalformedSQLStatements(f, fset, w)
			require.NoError(t, w.Flush())
			assert.Equal(t, tc.expected, strings.TrimSpace(b.String()))
		})
	}
}
