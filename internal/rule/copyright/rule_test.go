// Copyright (C) 2017-2018 Betalo AB - All Rights Reserved

package copyright

import (
	"bufio"
	"bytes"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProcessImports(t *testing.T) {
	currentYear := strconv.Itoa(time.Now().Year())

	var testcases = []struct {
		name     string
		given    []byte
		expected string
	}{
		{
			name:     "Betalo with year",
			given:    []byte(`// Copyright (C) ` + currentYear + ` Betalo AB - All Rights Reserved`),
			expected: "",
		},
		{
			name:     "PFC with period",
			given:    []byte(`// Copyright (C) 2017-` + currentYear + ` P.F.C. AB - All Rights Reserved`),
			expected: "",
		},
		{
			name:     "Simple",
			given:    []byte(`// © 2019 PFC Technology AB`),
			expected: "",
		},
		{
			name:     "Without",
			given:    []byte(`package main`),
			expected: "dummy.go:1:1: Missing copyright header",
		},
		{
			name:     "WrongYear",
			given:    []byte(`// Copyright (C) 2001 Betalo AB - All Rights Reserved`),
			expected: "dummy.go:1:1: Missing copyright header",
		},
		{
			name:     "Invalid",
			given:    []byte(`// Copyright (C) 2017 foo bar baz`),
			expected: "dummy.go:1:1: Missing copyright header",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			var b bytes.Buffer
			w := bufio.NewWriter(&b)
			findMissingCopyrightInHeader("dummy.go", tc.given, w)
			require.NoError(t, w.Flush())
			assert.Equal(t, tc.expected, strings.TrimSpace(b.String()))
		})
	}
}
