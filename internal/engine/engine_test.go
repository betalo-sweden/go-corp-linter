package engine

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWalk(t *testing.T) {
	// Arrange

	var buf bytes.Buffer
	w := bufio.NewWriter(&buf)

	// Act

	Process([]string{"testdata", "testdata/foo", "testdata/foo/bar"}, w, false)
	w.Flush()

	// Assert

	assert.Contains(t, buf.String(), `testdata/foo/bar/main.go`)
}
