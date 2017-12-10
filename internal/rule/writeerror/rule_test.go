// Copyright (C) 2017 Betalo AB - All Rights Reserved

package writeerror

import (
	"bufio"
	"bytes"
	"go/parser"
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_findMissingReturn(t *testing.T) {
	var cases = []struct {
		name     string
		given    string
		expected string
	}{
		{
			name: "Valid",
			given: `
				package main

				import (
					"errors"
					"fmt"
				)

				func main() {
					err := errors.New("dummy error")
					if true {
						rest.WriteError(w, r, err)
						return
					}
					fmt.Println("dummy text")
				}
				`,
			expected: "",
		}, {
			name: "Violation",
			given: `
				package main

				import (
					"errors"
					"fmt"
				)

				func main() {
					err := errors.New("dummy error")
					if true {
						writeError(w, r, err)
					}
					fmt.Println("dummy text")
				}
				`,
			expected: "main.go:13:7: missing return statement after writeError call\n",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fset := token.NewFileSet()
			f, err := parser.ParseFile(fset, "main.go", tc.given, 0)
			require.NoError(t, err)

			var b bytes.Buffer
			w := bufio.NewWriter(&b)
			findMissingReturn(f, fset, w)
			require.NoError(t, w.Flush())
			assert.Equal(t, tc.expected, b.String())
		})
	}
}
