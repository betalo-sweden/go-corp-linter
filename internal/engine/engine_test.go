// Copyright (C) 2017 Betalo AB - All Rights Reserved

package engine

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWalk(t *testing.T) {
	// Arrange

	var buf bytes.Buffer
	w := bufio.NewWriter(&buf)

	givenDirs := []string{
		"../../testdata",
		"../../testdata/foo",
		"../../testdata/foo/bar",
		"../../testdata/vendor",
	}

	// Act

	err := Process(givenDirs, w, false)
	require.NoError(t, err)
	require.NoError(t, w.Flush())

	// Assert

	assert.Contains(t, buf.String(), `testdata/foo/bar/main.go`)
	assert.Contains(t, buf.String(), `testdata/foo/bar/nocopyright.go`)
	assert.NotContains(t, buf.String(), `testdata/foo/bar/main.js`)
	assert.NotContains(t, buf.String(), `testdata/foo/bar/auto.go`)
	assert.NotContains(t, buf.String(), `testdata/vendor/foo/bar/main.go`)
}
