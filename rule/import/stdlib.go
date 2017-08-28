// Copyright (C) 2017 Betalo AB - All Rights Reserved

package imports

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func isStdlib(path string) bool {

	for _, p := range stdlibPkgs {
		if p == path {
			return true
		}
	}
	return false
}

func init() {
	findStdlibPkgs()
}

var stdlibPkgs []string

// findStdlibPkgs tries to get list of stdlib packages by reading GOROOT
//
// Courtesy: https://github.com/divan/depscheck/blob/master/package.go#L72
//
// This approach is used by "go list std" tool
// Based on go/cmd function matchPackages (https://golang.org/src/cmd/go/main.go#L553)
// and listStdPkgs function from https://golang.org/src/go/build/deps_test.go#L420
func findStdlibPkgs() {
	goroot := runtime.GOROOT()

	src := filepath.Join(goroot, "src") + string(filepath.Separator)
	walkFn := func(path string, fi os.FileInfo, err error) error {
		if err != nil || !fi.IsDir() || path == src {
			return nil
		}

		base := filepath.Base(path)
		if strings.HasPrefix(base, ".") || strings.HasPrefix(base, "_") || base == "testdata" {
			return filepath.SkipDir
		}

		name := filepath.ToSlash(path[len(src):])
		if name == "builtin" || name == "cmd" || strings.Contains(name, ".") {
			return filepath.SkipDir
		}

		stdlibPkgs = append(stdlibPkgs, name)
		return nil
	}

	if err := filepath.Walk(src, walkFn); err != nil {
		log.Fatalln("Error:", err)
	}
}
