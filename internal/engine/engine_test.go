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

	// Act

	err := Process([]string{"../../testdata", "../../testdata/foo", "../../testdata/foo/bar"}, w, false)
	require.NoError(t, err)
	require.NoError(t, w.Flush())

	// Assert

	assert.Contains(t, buf.String(), `testdata/foo/bar/main.go`)
}
