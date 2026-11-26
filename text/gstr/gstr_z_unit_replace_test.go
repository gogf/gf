// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*"

package gstr_test

import (
	"strings"
	"testing"
	"unsafe"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/text/gstr"
)

func Test_Replace(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := "abcdEFG乱入的中文abcdefg"
		t.Assert(gstr.Replace(s1, "ab", "AB"), "ABcdEFG乱入的中文ABcdefg")
		t.Assert(gstr.Replace(s1, "EF", "ef"), "abcdefG乱入的中文abcdefg")
		t.Assert(gstr.Replace(s1, "MN", "mn"), s1)

		t.Assert(gstr.ReplaceByArray(s1, g.ArrayStr{
			"a", "A",
			"A", "-",
			"a",
		}), "-bcdEFG乱入的中文-bcdefg")

		t.Assert(gstr.ReplaceByMap(s1, g.MapStrStr{
			"a": "A",
			"G": "g",
		}), "AbcdEFg乱入的中文Abcdefg")
	})
}

// An empty search string used to loop forever.
func Test_ReplaceI_EmptySearch(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gstr.ReplaceI("abc", "", "x"), "abc")
		t.Assert(gstr.ReplaceI("abc", "", "x", 2), "abc")
		t.Assert(gstr.ReplaceI("", "", "x"), "")
		t.Assert(gstr.ReplaceIByArray("abc", g.ArrayStr{"", "x", "B", "-"}), "a-c")
		t.Assert(gstr.ReplaceIByMap("abc", g.MapStrStr{"": "x"}), "abc")
	})
}

// The lower case of "İ", of the Kelvin sign "K" and of an invalid byte has another byte length.
func Test_ReplaceI_LowerCaseSize(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gstr.ReplaceI("İstanbul abc", "ABC", "x"), "İstanbul x")
		t.Assert(gstr.ReplaceI("İabc İabc", "ABC", "x"), "İx İx")
		t.Assert(gstr.ReplaceI("\xffabc", "ABC", "x"), "\xffx")
		t.Assert(gstr.ReplaceI("İstanbul", "İSTANBUL", "x"), "x")
		t.Assert(gstr.ReplaceI("abk", "K", "x"), "abx")
		t.Assert(gstr.ReplaceI("abK", "k", "x"), "abx")
		// The match is replaced as a whole, even if it is longer or shorter than the search string.
		t.Assert(gstr.ReplaceI("İstanbul İstanbul", "istanbul", "x"), "x x")
		t.Assert(gstr.ReplaceI("istanbul istanbul", "İstanbul", "x"), "x x")
		t.Assert(gstr.ReplaceI("İstanbul İstanbul", "istanbul", "x", 1), "x İstanbul")
		t.Assert(gstr.ReplaceI("Ångström Ωmega", "å", "x"), "xngström Ωmega")
		t.Assert(gstr.ReplaceI("Ångström Ωmega", "ω", "x"), "Ångström xmega")
		// The origin string is returned as it is, not a copy of it, if nothing is replaced.
		origin := strings.Repeat("İstanbul", 2)
		result := gstr.ReplaceI(origin, "abc", "x")
		t.Assert(unsafe.StringData(result) == unsafe.StringData(origin), true)
		t.Assert(gstr.ReplaceIFunc("aİb İ", "i", func(s string) string {
			return "[" + s + "]"
		}), "a[İ]b [İ]")
		t.Assert(gstr.ReplaceIFunc("ab", "İİİ", func(s string) string {
			return "[" + s + "]"
		}), "ab")
		t.Assert(gstr.ReplaceIFunc("İİabc", "ABC", func(s string) string {
			return "[" + s + "]"
		}), "İİ[abc]")
	})
}

func Test_ReplaceI_1(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := "abcd乱入的中文ABCD"
		s2 := "a"
		t.Assert(gstr.ReplaceI(s1, "ab", "aa"), "aacd乱入的中文aaCD")
		t.Assert(gstr.ReplaceI(s1, "ab", "aa", 0), "abcd乱入的中文ABCD")
		t.Assert(gstr.ReplaceI(s1, "ab", "aa", 1), "aacd乱入的中文ABCD")

		t.Assert(gstr.ReplaceI(s1, "abcd", "-"), "-乱入的中文-")
		t.Assert(gstr.ReplaceI(s1, "abcd", "-", 1), "-乱入的中文ABCD")

		t.Assert(gstr.ReplaceI(s1, "abcd乱入的", ""), "中文ABCD")
		t.Assert(gstr.ReplaceI(s1, "ABCD乱入的", ""), "中文ABCD")

		t.Assert(gstr.ReplaceI(s2, "A", "-"), "-")
		t.Assert(gstr.ReplaceI(s2, "a", "-"), "-")

		t.Assert(gstr.ReplaceIByArray(s1, g.ArrayStr{
			"abcd乱入的", "-",
			"-", "=",
			"a",
		}), "=中文ABCD")

		t.Assert(gstr.ReplaceIByMap(s1, g.MapStrStr{
			"ab": "-",
			"CD": "=",
		}), "-=乱入的中文-=")
	})
}

func Test_ReplaceI_2(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gstr.ReplaceI("aaa", "A", "-a-"), `-a--a--a-`)
		t.Assert(gstr.ReplaceI("aaaa", "AA", "-"), `--`)
		t.Assert(gstr.ReplaceI("a a a", "A", "b"), `b b b`)
		t.Assert(gstr.ReplaceI("aaaaaa", "aa", "a"), `aaa`)
		t.Assert(gstr.ReplaceI("aaaaaa", "AA", "A"), `AAA`)
		t.Assert(gstr.ReplaceI("aaa", "A", "AA"), `AAAAAA`)
		t.Assert(gstr.ReplaceI("aaa", "A", "AA"), `AAAAAA`)
		t.Assert(gstr.ReplaceI("a duration", "duration", "recordduration"), `a recordduration`)
	})
	// With count parameter.
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gstr.ReplaceI("aaaaaa", "aa", "a", 2), `aaaa`)
		t.Assert(gstr.ReplaceI("aaaaaa", "AA", "A", 1), `Aaaaa`)
		t.Assert(gstr.ReplaceI("aaaaaa", "AA", "A", 3), `AAA`)
		t.Assert(gstr.ReplaceI("aaaaaa", "AA", "A", 4), `AAA`)
		t.Assert(gstr.ReplaceI("aaa", "A", "AA", 2), `AAAAa`)
		t.Assert(gstr.ReplaceI("aaa", "A", "AA", 3), `AAAAAA`)
		t.Assert(gstr.ReplaceI("aaa", "A", "AA", 4), `AAAAAA`)
	})
}

func Test_ReplaceIFunc(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			origin = "hello GF 2018~2020!"
			search = "gf"
		)
		// Simple replacement
		result := gstr.ReplaceIFunc(origin, search, func(s string) string {
			return "GoFrame"
		})
		t.Assert(result, "hello GoFrame 2018~2020!")

		// Replace with original string
		result = gstr.ReplaceIFunc(origin, search, func(s string) string {
			return s
		})
		t.Assert(result, origin)

		// Replace with empty string
		result = gstr.ReplaceIFunc(origin, search, func(s string) string {
			return ""
		})
		t.Assert(result, "hello  2018~2020!")

		// Replace multiple occurrences with different cases
		origin = "GF is best, gf is nice, Gf is excellent"
		result = gstr.ReplaceIFunc(origin, search, func(s string) string {
			return "GoFrame"
		})
		t.Assert(result, "GoFrame is best, GoFrame is nice, GoFrame is excellent")

		// Empty search string
		result = gstr.ReplaceIFunc(origin, "", func(s string) string {
			return "GoFrame"
		})
		t.Assert(result, origin)

		// Empty origin string
		result = gstr.ReplaceIFunc("", search, func(s string) string {
			return "GoFrame"
		})
		t.Assert(result, "")

		// Replace with longer string
		result = gstr.ReplaceIFunc("GF", search, func(s string) string {
			return "GoFrame"
		})
		t.Assert(result, "GoFrame")

		// Replace with shorter string
		result = gstr.ReplaceIFunc("GF", search, func(s string) string {
			return "g"
		})
		t.Assert(result, "g")

		// Replace with mixed case patterns
		origin = "gf GF Gf gF"
		result = gstr.ReplaceIFunc(origin, search, func(s string) string {
			return strings.ToUpper(s)
		})
		t.Assert(result, "GF GF GF GF")
	})
}

func Test_ReplaceFunc(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			origin = "hello gf 2018~2020!"
			search = "gf"
		)
		// Simple replacement
		result := gstr.ReplaceFunc(origin, search, func(s string) string {
			return "GoFrame"
		})
		t.Assert(result, "hello GoFrame 2018~2020!")

		// Replace with original string
		result = gstr.ReplaceFunc(origin, search, func(s string) string {
			return s
		})
		t.Assert(result, origin)

		// Replace with empty string
		result = gstr.ReplaceFunc(origin, search, func(s string) string {
			return ""
		})
		t.Assert(result, "hello  2018~2020!")

		// Replace multiple occurrences
		origin = "gf is best, gf is nice"
		result = gstr.ReplaceFunc(origin, search, func(s string) string {
			return "GoFrame"
		})
		t.Assert(result, "GoFrame is best, GoFrame is nice")

		// Empty search string
		result = gstr.ReplaceFunc(origin, "", func(s string) string {
			return "GoFrame"
		})
		t.Assert(result, origin)

		// Empty origin string
		result = gstr.ReplaceFunc("", search, func(s string) string {
			return "GoFrame"
		})
		t.Assert(result, "")

		// Case sensitive
		origin = "GF is best, gf is nice"
		result = gstr.ReplaceFunc(origin, search, func(s string) string {
			return "GoFrame"
		})
		t.Assert(result, "GF is best, GoFrame is nice")

		// Replace with longer string
		result = gstr.ReplaceFunc("gf", search, func(s string) string {
			return "GoFrame"
		})
		t.Assert(result, "GoFrame")

		// Replace with shorter string
		result = gstr.ReplaceFunc("gf", search, func(s string) string {
			return "g"
		})
		t.Assert(result, "g")
	})
	gtest.C(t, func(t *gtest.T) {
		var (
			origin  = "gggg"
			search  = "g"
			replace = "gg"
		)
		// Simple replacement
		result := gstr.ReplaceFunc(origin, search, func(s string) string {
			return replace
		})
		t.Assert(result, "gggggggg")
	})
}
