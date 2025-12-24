// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gfile_test

import (
	"testing"

	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/test/gtest"
)

func Test_MatchGlob_Basic(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Basic glob patterns (no **)
		matched, err := gfile.MatchGlob("*.go", "main.go")
		t.AssertNil(err)
		t.Assert(matched, true)

		matched, err = gfile.MatchGlob("*.go", "main.txt")
		t.AssertNil(err)
		t.Assert(matched, false)

		matched, err = gfile.MatchGlob("test_*.go", "test_main.go")
		t.AssertNil(err)
		t.Assert(matched, true)

		matched, err = gfile.MatchGlob("?est.go", "test.go")
		t.AssertNil(err)
		t.Assert(matched, true)

		matched, err = gfile.MatchGlob("[abc].go", "a.go")
		t.AssertNil(err)
		t.Assert(matched, true)

		matched, err = gfile.MatchGlob("[a-z].go", "x.go")
		t.AssertNil(err)
		t.Assert(matched, true)
	})
}

func Test_MatchGlob_Globstar(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// ** matches everything
		matched, err := gfile.MatchGlob("**", "any/path/to/file.go")
		t.AssertNil(err)
		t.Assert(matched, true)

		matched, err = gfile.MatchGlob("**", "file.go")
		t.AssertNil(err)
		t.Assert(matched, true)

		matched, err = gfile.MatchGlob("**", "")
		t.AssertNil(err)
		t.Assert(matched, true)
	})
}

func Test_MatchGlob_GlobstarWithSuffix(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// **/*.go - matches .go files in any directory
		matched, err := gfile.MatchGlob("**/*.go", "main.go")
		t.AssertNil(err)
		t.Assert(matched, true)

		matched, err = gfile.MatchGlob("**/*.go", "src/main.go")
		t.AssertNil(err)
		t.Assert(matched, true)

		matched, err = gfile.MatchGlob("**/*.go", "src/foo/bar/main.go")
		t.AssertNil(err)
		t.Assert(matched, true)

		matched, err = gfile.MatchGlob("**/*.go", "src/main.txt")
		t.AssertNil(err)
		t.Assert(matched, false)
	})
}

func Test_MatchGlob_GlobstarWithPrefix(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// src/** - matches everything under src/
		matched, err := gfile.MatchGlob("src/**", "src/main.go")
		t.AssertNil(err)
		t.Assert(matched, true)

		matched, err = gfile.MatchGlob("src/**", "src/foo/bar/main.go")
		t.AssertNil(err)
		t.Assert(matched, true)

		matched, err = gfile.MatchGlob("src/**", "other/main.go")
		t.AssertNil(err)
		t.Assert(matched, false)
	})
}

func Test_MatchGlob_GlobstarWithPrefixAndSuffix(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// src/**/*.go - matches .go files under src/
		matched, err := gfile.MatchGlob("src/**/*.go", "src/main.go")
		t.AssertNil(err)
		t.Assert(matched, true)

		matched, err = gfile.MatchGlob("src/**/*.go", "src/foo/main.go")
		t.AssertNil(err)
		t.Assert(matched, true)

		matched, err = gfile.MatchGlob("src/**/*.go", "src/foo/bar/baz/main.go")
		t.AssertNil(err)
		t.Assert(matched, true)

		matched, err = gfile.MatchGlob("src/**/*.go", "src/main.txt")
		t.AssertNil(err)
		t.Assert(matched, false)

		matched, err = gfile.MatchGlob("src/**/*.go", "other/main.go")
		t.AssertNil(err)
		t.Assert(matched, false)
	})
}

func Test_MatchGlob_GlobstarMultiple(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Multiple ** in pattern
		matched, err := gfile.MatchGlob("src/**/test/**/*.go", "src/foo/test/bar/main.go")
		t.AssertNil(err)
		t.Assert(matched, true)

		matched, err = gfile.MatchGlob("src/**/test/**/*.go", "src/test/main.go")
		t.AssertNil(err)
		t.Assert(matched, true)

		matched, err = gfile.MatchGlob("src/**/test/**/*.go", "src/a/b/test/c/d/main.go")
		t.AssertNil(err)
		t.Assert(matched, true)
	})
}

func Test_MatchGlob_GlobstarEdgeCases(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// ** at the beginning
		matched, err := gfile.MatchGlob("**/main.go", "main.go")
		t.AssertNil(err)
		t.Assert(matched, true)

		matched, err = gfile.MatchGlob("**/main.go", "src/main.go")
		t.AssertNil(err)
		t.Assert(matched, true)

		matched, err = gfile.MatchGlob("**/main.go", "src/foo/bar/main.go")
		t.AssertNil(err)
		t.Assert(matched, true)

		// Hidden directories
		matched, err = gfile.MatchGlob(".*", ".git")
		t.AssertNil(err)
		t.Assert(matched, true)

		matched, err = gfile.MatchGlob(".*", ".vscode")
		t.AssertNil(err)
		t.Assert(matched, true)

		matched, err = gfile.MatchGlob("_*", "_test")
		t.AssertNil(err)
		t.Assert(matched, true)
	})
}

func Test_MatchGlob_WindowsPath(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Windows-style paths should also work
		matched, err := gfile.MatchGlob("src/**/*.go", "src\\foo\\main.go")
		t.AssertNil(err)
		t.Assert(matched, true)

		matched, err = gfile.MatchGlob("src\\**\\*.go", "src/foo/main.go")
		t.AssertNil(err)
		t.Assert(matched, true)
	})
}
