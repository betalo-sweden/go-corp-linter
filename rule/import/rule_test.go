// Copyright (C) 2017 Betalo AB - All Rights Reserved

package imports

import (
	"go/parser"
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProcessImportsSuccess(t *testing.T) {
	fset := token.NewFileSet()
	source := `
	package main

	import (
		"flag"
		"log"
		"os"
		"path/filepath"
		"strings"
	)
	func main() {
		var b bytes.Buffer
		b.Write([]byte("Hello "))
		fmt.Fprintf(&b, "world!")
	}`

	f, err := parser.ParseFile(fset, "main.go", source, parser.ImportsOnly)
	require.NoError(t, err)
	hits := processImportRule(f, fset)

	assert.Equal(t, 0, len(hits))
}

func TestProcessImportsWrongOrder(t *testing.T) {
	fset := token.NewFileSet()

	source := `
	package main

	import (
		"flag"
		"os"
		"log"
		"path/filepath"
		"strings"
	)
	func main() {
		var b bytes.Buffer
		b.Write([]byte("Hello "))
		fmt.Fprintf(&b, "world!")
	}`

	f, err := parser.ParseFile(fset, "main.go", source, parser.ImportsOnly)
	require.NoError(t, err)

	hits := processImportRule(f, fset)

	assert.Equal(t, 1, len(hits))
	assert.Equal(t, notSorted, hits[0].problem)
}

func TestProcessImportsWrongGrouping(t *testing.T) {
	fset := token.NewFileSet()

	source := `
	package main

	import (
		"flag"
		"log"
		"os"

		"path/filepath"
		"strings"
	)
	func main() {
		var b bytes.Buffer
		b.Write([]byte("Hello "))
		fmt.Fprintf(&b, "world!")
	}`

	f, err := parser.ParseFile(fset, "main.go", source, parser.ImportsOnly)
	require.NoError(t, err)

	hits := processImportRule(f, fset)

	assert.Equal(t, 1, len(hits))
	assert.Equal(t, notGrouped, hits[0].problem)
}

func TestProcessImportsWrongFirstPackage(t *testing.T) {
	fset := token.NewFileSet()
	source := `
	package main

	import (
		"github.com/stretchr/testify/assert"
		"flag"
		"log"
		"os"
		"path/filepath"
		"strings"
	)
	func main() {
		var b bytes.Buffer
		b.Write([]byte("Hello "))
		fmt.Fprintf(&b, "world!")
	}`

	f, err := parser.ParseFile(fset, "main.go", source, parser.ImportsOnly)
	require.NoError(t, err)

	hits := processImportRule(f, fset)

	assert.Equal(t, 1, len(hits))
	assert.Equal(t, firstPackage, hits[0].problem)
}
