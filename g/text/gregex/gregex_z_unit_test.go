// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*"

package gregex_test

import (
	"github.com/gogf/gf/g/test/gtest"
	"github.com/gogf/gf/g/text/gregex"
	"strings"
	"testing"
)

func Test_Quote(t *testing.T) {
	gtest.Case(t, func() {
		s1 := `[foo]` //`\[foo\]`
		gtest.Assert(gregex.Quote(s1), `\[foo\]`)
	})
}

func Test_Validate(t *testing.T) {
	gtest.Case(t, func() {
		var s1 = `(.+):(\d+)`
		gtest.Assert(gregex.Validate(s1), nil)
		s1 = `((.+):(\d+)`
		gtest.Assert(gregex.Validate(s1) == nil, false)
	})
}

func Test_IsMatch(t *testing.T) {
	gtest.Case(t, func() {
		var pattern = `(.+):(\d+)`
		s1 := []byte(`sfs:2323`)
		gtest.Assert(gregex.IsMatch(pattern, s1), true)
		s1 = []byte(`sfs2323`)
		gtest.Assert(gregex.IsMatch(pattern, s1), false)
		s1 = []byte(`sfs:`)
		gtest.Assert(gregex.IsMatch(pattern, s1), false)
	})
}

func Test_IsMatchString(t *testing.T) {
	gtest.Case(t, func() {
		var pattern = `(.+):(\d+)`
		s1 := `sfs:2323`
		gtest.Assert(gregex.IsMatchString(pattern, s1), true)
		s1 = `sfs2323`
		gtest.Assert(gregex.IsMatchString(pattern, s1), false)
		s1 = `sfs:`
		gtest.Assert(gregex.IsMatchString(pattern, s1), false)
	})
}

func Test_Match(t *testing.T) {
	gtest.Case(t, func() {
		re := "a(a+b+)b"
		wantSubs := "aaabb"
		s := "acbb" + wantSubs + "dd"
		subs, err := gregex.Match(re, []byte(s))
		gtest.Assert(err, nil)
		if string(subs[0]) != wantSubs {
			t.Fatalf("regex:%s,Match(%q)[0] = %q; want %q", re, s, subs[0], wantSubs)
		}
		if string(subs[1]) != "aab" {
			t.Fatalf("Match(%q)[1] = %q; want %q", s, subs[1], "aab")
		}
	})
}

func Test_MatchString(t *testing.T) {
	gtest.Case(t, func() {
		re := "a(a+b+)b"
		wantSubs := "aaabb"
		s := "acbb" + wantSubs + "dd"
		subs, err := gregex.MatchString(re, s)
		gtest.Assert(err, nil)
		if string(subs[0]) != wantSubs {
			t.Fatalf("regex:%s,Match(%q)[0] = %q; want %q", re, s, subs[0], wantSubs)
		}
		if string(subs[1]) != "aab" {
			t.Fatalf("Match(%q)[1] = %q; want %q", s, subs[1], "aab")
		}
	})
}

func Test_MatchAll(t *testing.T) {
	gtest.Case(t, func() {
		re := "a(a+b+)b"
		wantSubs := "aaabb"
		s := "acbb" + wantSubs + "dd"
		s = s + `其他的` + s
		subs, err := gregex.MatchAll(re, []byte(s))
		gtest.Assert(err, nil)
		if string(subs[0][0]) != wantSubs {
			t.Fatalf("regex:%s,Match(%q)[0] = %q; want %q", re, s, subs[0][0], wantSubs)
		}
		if string(subs[0][1]) != "aab" {
			t.Fatalf("Match(%q)[1] = %q; want %q", s, subs[0][1], "aab")
		}

		if string(subs[1][0]) != wantSubs {
			t.Fatalf("regex:%s,Match(%q)[0] = %q; want %q", re, s, subs[1][0], wantSubs)
		}
		if string(subs[1][1]) != "aab" {
			t.Fatalf("Match(%q)[1] = %q; want %q", s, subs[1][1], "aab")
		}
	})
}

func Test_MatchAllString(t *testing.T) {
	gtest.Case(t, func() {
		re := "a(a+b+)b"
		wantSubs := "aaabb"
		s := "acbb" + wantSubs + "dd"
		subs, err := gregex.MatchAllString(re, s+`其他的`+s)
		gtest.Assert(err, nil)
		if string(subs[0][0]) != wantSubs {
			t.Fatalf("regex:%s,Match(%q)[0] = %q; want %q", re, s, subs[0][0], wantSubs)
		}
		if string(subs[0][1]) != "aab" {
			t.Fatalf("Match(%q)[1] = %q; want %q", s, subs[0][1], "aab")
		}

		if string(subs[1][0]) != wantSubs {
			t.Fatalf("regex:%s,Match(%q)[0] = %q; want %q", re, s, subs[1][0], wantSubs)
		}
		if string(subs[1][1]) != "aab" {
			t.Fatalf("Match(%q)[1] = %q; want %q", s, subs[1][1], "aab")
		}
	})
}

func Test_Replace(t *testing.T) {
	gtest.Case(t, func() {
		re := "a(a+b+)b"
		wantSubs := "aaabb"
		replace := "12345"
		s := "acbb" + wantSubs + "dd"
		wanted := "acbb" + replace + "dd"
		replacedStr, err := gregex.Replace(re, []byte(replace), []byte(s))
		gtest.Assert(err, nil)
		if string(replacedStr) != wanted {
			t.Fatalf("regex:%s,old:%s; want %q", re, s, wanted)
		}
	})
}

func Test_ReplaceString(t *testing.T) {
	gtest.Case(t, func() {
		re := "a(a+b+)b"
		wantSubs := "aaabb"
		replace := "12345"
		s := "acbb" + wantSubs + "dd"
		wanted := "acbb" + replace + "dd"
		replacedStr, err := gregex.ReplaceString(re, replace, s)
		gtest.Assert(err, nil)
		if replacedStr != wanted {
			t.Fatalf("regex:%s,old:%s; want %q", re, s, wanted)
		}
	})
}

func Test_ReplaceFun(t *testing.T) {
	gtest.Case(t, func() {
		re := "a(a+b+)b"
		wantSubs := "aaabb"
		//replace :="12345"
		s := "acbb" + wantSubs + "dd"
		wanted := "acbb[x" + wantSubs + "y]dd"
		wanted = "acbb" + "3个a" + "dd"
		replacedStr, err := gregex.ReplaceFunc(re, []byte(s), func(s []byte) []byte {
			if strings.Index(string(s), "aaa") >= 0 {
				return []byte("3个a")
			}
			return []byte("[x" + string(s) + "y]")
		})
		gtest.Assert(err, nil)
		if string(replacedStr) != wanted {
			t.Fatalf("regex:%s,old:%s; want %q", re, s, wanted)
		}
	})
}

func Test_ReplaceFuncMatch(t *testing.T) {
	gtest.Case(t, func() {
		s := []byte("1234567890")
		p := `(\d{3})(\d{3})(.+)`
		s0, e0 := gregex.ReplaceFuncMatch(p, s, func(match [][]byte) []byte {
			return match[0]
		})
		gtest.Assert(e0, nil)
		gtest.Assert(s0, s)
		s1, e1 := gregex.ReplaceFuncMatch(p, s, func(match [][]byte) []byte {
			return match[1]
		})
		gtest.Assert(e1, nil)
		gtest.Assert(s1, []byte("123"))
		s2, e2 := gregex.ReplaceFuncMatch(p, s, func(match [][]byte) []byte {
			return match[2]
		})
		gtest.Assert(e2, nil)
		gtest.Assert(s2, []byte("456"))
		s3, e3 := gregex.ReplaceFuncMatch(p, s, func(match [][]byte) []byte {
			return match[3]
		})
		gtest.Assert(e3, nil)
		gtest.Assert(s3, []byte("7890"))
	})
}

func Test_ReplaceStringFunc(t *testing.T) {
	gtest.Case(t, func() {
		re := "a(a+b+)b"
		wantSubs := "aaabb"
		//replace :="12345"
		s := "acbb" + wantSubs + "dd"
		wanted := "acbb[x" + wantSubs + "y]dd"
		wanted = "acbb" + "3个a" + "dd"
		replacedStr, err := gregex.ReplaceStringFunc(re, s, func(s string) string {
			if strings.Index(s, "aaa") >= 0 {
				return "3个a"
			}
			return "[x" + s + "y]"
		})
		gtest.Assert(err, nil)
		if replacedStr != wanted {
			t.Fatalf("regex:%s,old:%s; want %q", re, s, wanted)
		}
	})
}

func Test_ReplaceStringFuncMatch(t *testing.T) {
	gtest.Case(t, func() {
		s := "1234567890"
		p := `(\d{3})(\d{3})(.+)`
		s0, e0 := gregex.ReplaceStringFuncMatch(p, s, func(match []string) string {
			return match[0]
		})
		gtest.Assert(e0, nil)
		gtest.Assert(s0, s)
		s1, e1 := gregex.ReplaceStringFuncMatch(p, s, func(match []string) string {
			return match[1]
		})
		gtest.Assert(e1, nil)
		gtest.Assert(s1, "123")
		s2, e2 := gregex.ReplaceStringFuncMatch(p, s, func(match []string) string {
			return match[2]
		})
		gtest.Assert(e2, nil)
		gtest.Assert(s2, "456")
		s3, e3 := gregex.ReplaceStringFuncMatch(p, s, func(match []string) string {
			return match[3]
		})
		gtest.Assert(e3, nil)
		gtest.Assert(s3, "7890")
	})
}

func Test_Split(t *testing.T) {
	gtest.Case(t, func() {
		re := "a(a+b+)b"
		matched := "aaabb"
		item0 := "acbb"
		item1 := "dd"
		s := item0 + matched + item1
		gtest.Assert(gregex.IsMatchString(re, matched), true)
		items := gregex.Split(re, s) //split string with matched
		if items[0] != item0 {
			t.Fatalf("regex:%s,Split(%q) want %q", re, s, item0)
		}
		if items[1] != item1 {
			t.Fatalf("regex:%s,Split(%q) want %q", re, s, item0)
		}
	})

	gtest.Case(t, func() {
		re := "a(a+b+)b"
		notmatched := "aaxbb"
		item0 := "acbb"
		item1 := "dd"
		s := item0 + notmatched + item1
		gtest.Assert(gregex.IsMatchString(re, notmatched), false)
		items := gregex.Split(re, s) //split string with notmatched then nosplitting
		if items[0] != s {
			t.Fatalf("regex:%s,Split(%q) want %q", re, s, item0)
		}

	})
}
