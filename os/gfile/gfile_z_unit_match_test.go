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

func Test_MatchGlob_MemoizationCache(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Test cases that exercise memoization cache hits
		// Multiple ** with same suffix patterns will trigger cache reuse
		matched, err := gfile.MatchGlob("a/**/b/**/c", "a/x/b/y/c")
		t.AssertNil(err)
		t.Assert(matched, true)

		// This pattern creates multiple paths that converge to same subproblems
		matched, err = gfile.MatchGlob("**/a/**/a", "x/a/y/a")
		t.AssertNil(err)
		t.Assert(matched, true)

		// Deep recursion with cache hits
		matched, err = gfile.MatchGlob("**/**/**", "a/b/c")
		t.AssertNil(err)
		t.Assert(matched, true)
	})
}

func Test_MatchGlob_InvalidGlobstarAtEnd(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Pattern where "**" appears at the very end of string (idx >= len(pattern) after pos+2)
		// "x**" - invalid globstar at end, should be treated as two "*"
		matched, err := gfile.MatchGlob("x**", "x")
		t.AssertNil(err)
		t.Assert(matched, true)

		matched, err = gfile.MatchGlob("x**", "xyz")
		t.AssertNil(err)
		t.Assert(matched, true)

		// Pattern ending with invalid globstar that exhausts the string
		matched, err = gfile.MatchGlob("abc**", "abc")
		t.AssertNil(err)
		t.Assert(matched, true)

		matched, err = gfile.MatchGlob("abc**", "abcdef")
		t.AssertNil(err)
		t.Assert(matched, true)
	})
}

func Test_MatchGlob_PrefixWithWildcards(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Prefix contains wildcards - tests lines 220-236
		// Pattern: "s*c/**/file.go" - prefix "s*c" contains wildcard
		matched, err := gfile.MatchGlob("s*c/**/*.go", "src/foo/main.go")
		t.AssertNil(err)
		t.Assert(matched, true)

		matched, err = gfile.MatchGlob("s?c/**/*.go", "src/foo/main.go")
		t.AssertNil(err)
		t.Assert(matched, true)

		// Test line 223-225: name has fewer segments than prefix
		matched, err = gfile.MatchGlob("a/b/c/**", "a/b")
		t.AssertNil(err)
		t.Assert(matched, false)

		matched, err = gfile.MatchGlob("a/b/c/**/d", "a")
		t.AssertNil(err)
		t.Assert(matched, false)

		// Test line 232-234: wildcard prefix doesn't match
		matched, err = gfile.MatchGlob("x*c/**/*.go", "src/foo/main.go")
		t.AssertNil(err)
		t.Assert(matched, false)

		matched, err = gfile.MatchGlob("s?x/**/*.go", "src/foo/main.go")
		t.AssertNil(err)
		t.Assert(matched, false)

		// Test line 236: name update after prefix match
		matched, err = gfile.MatchGlob("a*/b*/**/*.go", "abc/bcd/efg/main.go")
		t.AssertNil(err)
		t.Assert(matched, true)
	})
}

func Test_MatchGlob_EmptyNameWithSuffix(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Test line 246-249: name becomes empty after prefix match, check if suffix can match empty
		// "abc/**" with name "abc" - after prefix match, name is empty
		matched, err := gfile.MatchGlob("abc/**/", "abc")
		t.AssertNil(err)
		t.Assert(matched, true)

		// "abc/**/d" with name "abc" - after prefix match, name is empty but suffix is "d"
		matched, err = gfile.MatchGlob("abc/**/d", "abc")
		t.AssertNil(err)
		t.Assert(matched, false)

		// Test with wildcard prefix that exactly matches
		matched, err = gfile.MatchGlob("a*c/**/x", "abc")
		t.AssertNil(err)
		t.Assert(matched, false)
	})
}

func Test_MatchGlob_FindValidGlobstarExhaust(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Test lines 147-152: findValidGlobstar exhausts pattern without finding valid globstar
		// Pattern with multiple invalid "**" that ends exactly at pattern length
		matched, err := gfile.MatchGlob("a**b**", "ab")
		t.AssertNil(err)
		t.Assert(matched, true)

		matched, err = gfile.MatchGlob("x**y**z", "xyz")
		t.AssertNil(err)
		t.Assert(matched, true)

		// Pattern where last "**" is at the very end but invalid
		matched, err = gfile.MatchGlob("test**", "test")
		t.AssertNil(err)
		t.Assert(matched, true)

		matched, err = gfile.MatchGlob("test**", "testing")
		t.AssertNil(err)
		t.Assert(matched, true)
	})
}

func Test_MatchGlob_CacheHit(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Test line 166-168: cache hit scenario
		// Pattern that creates overlapping subproblems triggering cache hits
		// "**/**" with multiple segments will have cache hits
		matched, err := gfile.MatchGlob("**/x/**/x", "a/x/b/x")
		t.AssertNil(err)
		t.Assert(matched, true)

		// This pattern specifically creates cache hits due to overlapping subproblems
		// when trying different combinations of ** matching
		matched, err = gfile.MatchGlob("**/a/**/b/**/a", "x/a/y/b/z/a")
		t.AssertNil(err)
		t.Assert(matched, true)

		// Pattern with repeated suffix that will be checked multiple times
		matched, err = gfile.MatchGlob("**/**/test", "a/b/c/test")
		t.AssertNil(err)
		t.Assert(matched, true)

		// Pattern that will cause same subproblem to be solved multiple times
		// "**/**/**" matching "a/b/c/d" will have many overlapping subproblems
		matched, err = gfile.MatchGlob("**/**/**/**", "a/b/c/d/e")
		t.AssertNil(err)
		t.Assert(matched, true)
	})
}

func Test_MatchGlob_WildcardPrefixShortName(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Test line 223-225: prefix with wildcards, name has fewer segments
		// Pattern: "a*/b*/**/c" - prefix "a*/b*" has 2 segments
		// Name: "ax" - only 1 segment
		matched, err := gfile.MatchGlob("a*/b*/**/c", "ax")
		t.AssertNil(err)
		t.Assert(matched, false)

		// Pattern: "?/b/c/**/d" - prefix "?/b/c" has 3 segments
		// Name: "x/y" - only 2 segments
		matched, err = gfile.MatchGlob("?/b/c/**/d", "x/y")
		t.AssertNil(err)
		t.Assert(matched, false)

		// Pattern: "[abc]/[def]/**/x" - prefix has 2 segments with brackets
		// Name: "a" - only 1 segment
		matched, err = gfile.MatchGlob("[abc]/[def]/**/x", "a")
		t.AssertNil(err)
		t.Assert(matched, false)
	})
}

func Test_MatchGlob_InvalidGlobstarInSuffix(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Test lines 147-152: findValidGlobstar exhausts pattern in recursive call
		// Pattern "a/**/b**" - first "**" is valid, suffix "b**" has invalid "**" at end
		// When matching suffix "b**", findValidGlobstar will iterate and find "**" is invalid,
		// then idx = pos + 2 = 3, len("b**") = 3, so idx >= len(pattern) triggers break
		matched, err := gfile.MatchGlob("a/**/b**", "a/x/bcd")
		t.AssertNil(err)
		t.Assert(matched, true)

		matched, err = gfile.MatchGlob("a/**/b**", "a/x/b")
		t.AssertNil(err)
		t.Assert(matched, true)

		// Pattern with valid globstar followed by suffix with invalid globstar at end
		matched, err = gfile.MatchGlob("x/**/y**z", "x/a/yabcz")
		t.AssertNil(err)
		t.Assert(matched, true)

		// Multiple invalid globstars in suffix
		matched, err = gfile.MatchGlob("a/**/x**y**", "a/b/xcy")
		t.AssertNil(err)
		t.Assert(matched, true)
	})
}

func Test_MatchGlob_MemoizationCacheHit(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Test line 166-168: cache hit scenario
		// To trigger cache hit, we need:
		// 1. Same (pattern, name) pair called twice
		// 2. First call must complete (not return early)
		// 3. This happens when matching FAILS and we try all combinations

		// Pattern "**/**/z" with name "a/b/c/d" (no match)
		// First ** tries 0,1,2,3,4 segments
		// For each, second ** tries all remaining combinations
		// This creates overlapping subproblems that fail:
		// - ("**/z", "a/b/c/d"), ("**/z", "b/c/d"), ("**/z", "c/d"), ("**/z", "d"), ("**/z", "")
		// - ("z", "a/b/c/d"), ("z", "b/c/d"), ("z", "c/d"), ("z", "d"), ("z", "")
		// When first ** matches 0: check ("**/z", "a/b/c/d")
		//   -> second ** matches 0: check ("z", "a/b/c/d") - false, cached
		//   -> second ** matches 1: check ("z", "b/c/d") - false, cached
		//   -> second ** matches 2: check ("z", "c/d") - false, cached
		//   -> second ** matches 3: check ("z", "d") - false, cached
		//   -> second ** matches 4: check ("z", "") - false, cached
		// When first ** matches 1: check ("**/z", "b/c/d")
		//   -> second ** matches 0: check ("z", "b/c/d") - CACHE HIT!
		matched, err := gfile.MatchGlob("**/**/z", "a/b/c/d")
		t.AssertNil(err)
		t.Assert(matched, false)

		// Another failing pattern that creates cache hits
		matched, err = gfile.MatchGlob("**/**/**/notexist", "a/b/c/d/e")
		t.AssertNil(err)
		t.Assert(matched, false)

		// Pattern with same suffix appearing multiple times in recursion (failing case)
		matched, err = gfile.MatchGlob("**/x/**/x/**/x", "a/b/c/d/e/f")
		t.AssertNil(err)
		t.Assert(matched, false)
	})
}
