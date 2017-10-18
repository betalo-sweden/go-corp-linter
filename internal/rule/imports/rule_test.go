// Copyright (C) 2017 Betalo AB - All Rights Reserved

package imports

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

var testcases = []struct {
	name     string
	given    string
	expected string
}{
	{
		name: "valid",
		given: `
				"flag"
				"log"
				"os"
				"path/filepath"
				"strings"
				`,
		expected: "",
	},
	{
		name:     "empty",
		given:    ``,
		expected: "",
	},
	{
		name:     "noStdlib",
		given:    `"github.com/stretchr/testify/assert"`,
		expected: "",
	},
	{
		name:     "onlyStdlib",
		given:    `"io"`,
		expected: "",
	},
	{
		name: "unsortedStdlib",
		given: `
				"flag"
				"os"
				"log"
				"path/filepath"
				"strings"
				`,
		expected: "main.go:5:5:error: incorrectly sorted import package: os",
	},
	{
		name: "redundantGrouping",
		given: `
				"flag"
				"log"
				"os"

				"path/filepath"
				"strings"
				`,
		expected: "main.go:8:5:error: incorrectly sorted import package: path/filepath",
	},
	{
		name: "stdlibAndOthersMixed",
		given: `
				"github.com/stretchr/testify/assert"
				"flag"
				"log"
				"os"
				"path/filepath"
				"strings"
				`,
		expected: "main.go:4:5:error: incorrectly sorted import package: github.com/stretchr/testify/assert",
	},
	{
		name: "othersNotSorted",
		given: `
				"context"
				"crypto/md5"
				"encoding/hex"
				"fmt"
				"net/http"
				"net/url"
				"strings"
				"time"

				"github.com/betalo-sweden/pkg/log"
				"go.uber.org/zap"
				"github.com/betalo-sweden/pkg/router/middleware"
				`,
		expected: "main.go:14:5:error: incorrectly sorted import package: go.uber.org/zap",
	},
}

func TestProcessImports(t *testing.T) {
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			fset := token.NewFileSet()
			source := "package main\nimport (\n" + tc.given + ")\nfunc main() {}"
			f, err := parser.ParseFile(fset, "main.go", source, parser.ImportsOnly)
			require.NoError(t, err)

			var b bytes.Buffer
			w := bufio.NewWriter(&b)
			processImports(f, fset, w)
			require.NoError(t, w.Flush())
			assert.Equal(t, tc.expected, strings.TrimSpace(b.String()))
		})
	}
}
