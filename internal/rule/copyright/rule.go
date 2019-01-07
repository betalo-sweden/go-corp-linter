// Copyright (C) 2017-2018 Betalo AB - All Rights Reserved

package copyright

import (
	"bytes"
	"fmt"
	"go/token"
	"io"
	"os"
	"regexp"
)

var (
	betalo                      = "Betalo AB"
	pfc                         = "P.F.C. AB"
	copyrightHeaderPrefixRegexp = regexp.MustCompile(fmt.Sprintf(`// Copyright \([cC]\) .*20[1-9][0-9] %s|%s - All Rights Reserved`, betalo, pfc))
)

// ProcessFile checks for copyright headers in Go files.
func ProcessFile(fp string, out io.Writer) error {
	return findMissingCopyright(fp, out)
}

func findMissingCopyright(filepath string, out io.Writer) error {
	b, err := head(filepath, 100) // Arbitrary, just has to be large enough
	if err != nil {
		return err
	}

	findMissingCopyrightInHeader(filepath, bytes.TrimSpace(b), out)
	return nil
}

func findMissingCopyrightInHeader(filepath string, header []byte, out io.Writer) {
	if !copyrightHeaderPrefixRegexp.Match(header) {
		report(out, token.Position{Filename: filepath, Line: 1, Column: 1})
	}
}

func head(filepath string, size int64) ([]byte, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()

	b := make([]byte, size)
	_, err = f.Read(b)
	if err != io.EOF {
		err = nil
	}
	return b, err
}

func report(out io.Writer, position token.Position) {
	fmt.Fprintf(out, "%s: Missing copyright header\n", position)
}
