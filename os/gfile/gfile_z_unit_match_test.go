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

func Test_MatchGlob_InvalidGlobstar(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// "**" not as complete path component should be treated as two "*"
		// "a**b" should match "ab", "axb", "axxb", etc. (but not "a/b")
		matched, err := gfile.MatchGlob("a**b", "ab")
		t.AssertNil(err)
		t.Assert(matched, true)

		matched, err = gfile.MatchGlob("a**b", "axb")
		t.AssertNil(err)
		t.Assert(matched, true)

		matched, err = gfile.MatchGlob("a**b", "axxb")
		t.AssertNil(err)
		t.Assert(matched, true)

		matched, err = gfile.MatchGlob("a**b", "axxxb")
		t.AssertNil(err)
		t.Assert(matched, true)

		// "a**b" should NOT match paths with separators
		matched, err = gfile.MatchGlob("a**b", "a/b")
		t.AssertNil(err)
		t.Assert(matched, false)

		matched, err = gfile.MatchGlob("a**b", "ax/yb")
		t.AssertNil(err)
		t.Assert(matched, false)

		// "**a" at start (not valid globstar)
		matched, err = gfile.MatchGlob("**a", "a")
		t.AssertNil(err)
		t.Assert(matched, true)

		matched, err = gfile.MatchGlob("**a", "xa")
		t.AssertNil(err)
		t.Assert(matched, true)

		matched, err = gfile.MatchGlob("**a", "xxa")
		t.AssertNil(err)
		t.Assert(matched, true)

		// "a**" at end (not valid globstar)
		matched, err = gfile.MatchGlob("a**", "a")
		t.AssertNil(err)
		t.Assert(matched, true)

		matched, err = gfile.MatchGlob("a**", "ax")
		t.AssertNil(err)
		t.Assert(matched, true)

		matched, err = gfile.MatchGlob("a**", "axx")
		t.AssertNil(err)
		t.Assert(matched, true)

		// Mixed valid and invalid globstars
		// "src/**a" - "**" is valid globstar, "a" is suffix
		matched, err = gfile.MatchGlob("src/**/a", "src/foo/a")
		t.AssertNil(err)
		t.Assert(matched, true)

		matched, err = gfile.MatchGlob("src/**/a", "src/a")
		t.AssertNil(err)
		t.Assert(matched, true)
	})
}

func Test_MatchGlob_PrefixBoundary(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// "abc/**" should NOT match "abcdef/file.go" (prefix must be complete path component)
		matched, err := gfile.MatchGlob("abc/**", "abcdef/file.go")
		t.AssertNil(err)
		t.Assert(matched, false)

		// "abc/**" should match "abc/file.go"
		matched, err = gfile.MatchGlob("abc/**", "abc/file.go")
		t.AssertNil(err)
		t.Assert(matched, true)

		// "abc/**" should match "abc/def/file.go"
		matched, err = gfile.MatchGlob("abc/**", "abc/def/file.go")
		t.AssertNil(err)
		t.Assert(matched, true)

		// "abc/**" should match "abc" (prefix equals name)
		matched, err = gfile.MatchGlob("abc/**", "abc")
		t.AssertNil(err)
		t.Assert(matched, true)

		// "src/foo/**" should NOT match "src/foobar/file.go"
		matched, err = gfile.MatchGlob("src/foo/**", "src/foobar/file.go")
		t.AssertNil(err)
		t.Assert(matched, false)

		// "src/foo/**" should match "src/foo/bar/file.go"
		matched, err = gfile.MatchGlob("src/foo/**", "src/foo/bar/file.go")
		t.AssertNil(err)
		t.Assert(matched, true)
	})
}

func Test_MatchGlob_MultipleGlobstars(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Test with multiple ** operators - this would be slow without memoization
		matched, err := gfile.MatchGlob("a/**/b/**/c/**/d.go", "a/x/y/b/z/c/w/d.go")
		t.AssertNil(err)
		t.Assert(matched, true)

		matched, err = gfile.MatchGlob("a/**/b/**/c/**/d.go", "a/b/c/d.go")
		t.AssertNil(err)
		t.Assert(matched, true)

		matched, err = gfile.MatchGlob("a/**/b/**/c/**/d.go", "a/1/2/3/b/4/5/c/6/d.go")
		t.AssertNil(err)
		t.Assert(matched, true)

		matched, err = gfile.MatchGlob("a/**/b/**/c/**/d.go", "a/b/c/e.go")
		t.AssertNil(err)
		t.Assert(matched, false)

		// Deep nesting test
		matched, err = gfile.MatchGlob("**/*.go", "a/b/c/d/e/f/g/h/i/j/main.go")
		t.AssertNil(err)
		t.Assert(matched, true)
	})
}

func Test_MatchGlob_MalformedPatterns(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Unclosed bracket - should return error
		_, err := gfile.MatchGlob("[", "a")
		t.AssertNE(err, nil)

		_, err = gfile.MatchGlob("[abc", "a")
		t.AssertNE(err, nil)

		_, err = gfile.MatchGlob("[[", "a")
		t.AssertNE(err, nil)

		// Malformed patterns with globstar - errors should propagate
		_, err = gfile.MatchGlob("**/[", "a/b")
		t.AssertNE(err, nil)

		_, err = gfile.MatchGlob("[/**", "a/b")
		t.AssertNE(err, nil)

		_, err = gfile.MatchGlob("a/**/[abc", "a/b/c")
		t.AssertNE(err, nil)

		// Malformed pattern in prefix with wildcards
		_, err = gfile.MatchGlob("[a/**/b", "a/x/b")
		t.AssertNE(err, nil)

		// Invalid escape sequence on non-Windows (backslash at end)
		// Note: behavior may vary by platform
		_, err = gfile.MatchGlob("test\\", "test")
		// On Unix, this might not error but won't match
		// The key is it shouldn't panic

		// Valid patterns should still work
		matched, err := gfile.MatchGlob("[abc]", "a")
		t.AssertNil(err)
		t.Assert(matched, true)

		matched, err = gfile.MatchGlob("[a-z]", "m")
		t.AssertNil(err)
		t.Assert(matched, true)

		// Note: filepath.Match uses [^...] for negation, not [!...]
		matched, err = gfile.MatchGlob("[^abc]", "d")
		t.AssertNil(err)
		t.Assert(matched, true)

		matched, err = gfile.MatchGlob("[^a-z]", "1")
		t.AssertNil(err)
		t.Assert(matched, true)
	})
}
