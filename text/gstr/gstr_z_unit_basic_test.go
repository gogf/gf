// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*"

package gstr_test

import (
	"testing"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/test/gtest"
	"github.com/gogf/gf/text/gstr"
)

func Test_Replace(t *testing.T) {
	gtest.Case(t, func() {
		s1 := "abcdEFG乱入的中文abcdefg"
		gtest.Assert(gstr.Replace(s1, "ab", "AB"), "ABcdEFG乱入的中文ABcdefg")
		gtest.Assert(gstr.Replace(s1, "EF", "ef"), "abcdefG乱入的中文abcdefg")
		gtest.Assert(gstr.Replace(s1, "MN", "mn"), s1)

		gtest.Assert(gstr.ReplaceByArray(s1, g.ArrayStr{
			"a", "A",
			"A", "-",
			"a",
		}), "-bcdEFG乱入的中文-bcdefg")

		gtest.Assert(gstr.ReplaceByMap(s1, g.MapStrStr{
			"a": "A",
			"G": "g",
		}), "AbcdEFg乱入的中文Abcdefg")
	})
}

func Test_ReplaceI_1(t *testing.T) {
	gtest.Case(t, func() {
		s1 := "abcd乱入的中文ABCD"
		s2 := "a"
		gtest.Assert(gstr.ReplaceI(s1, "ab", "aa"), "aacd乱入的中文aaCD")
		gtest.Assert(gstr.ReplaceI(s1, "ab", "aa", 0), "abcd乱入的中文ABCD")
		gtest.Assert(gstr.ReplaceI(s1, "ab", "aa", 1), "aacd乱入的中文ABCD")

		gtest.Assert(gstr.ReplaceI(s1, "abcd", "-"), "-乱入的中文-")
		gtest.Assert(gstr.ReplaceI(s1, "abcd", "-", 1), "-乱入的中文ABCD")

		gtest.Assert(gstr.ReplaceI(s1, "abcd乱入的", ""), "中文ABCD")
		gtest.Assert(gstr.ReplaceI(s1, "ABCD乱入的", ""), "中文ABCD")

		gtest.Assert(gstr.ReplaceI(s2, "A", "-"), "-")
		gtest.Assert(gstr.ReplaceI(s2, "a", "-"), "-")

		gtest.Assert(gstr.ReplaceIByArray(s1, g.ArrayStr{
			"abcd乱入的", "-",
			"-", "=",
			"a",
		}), "=中文ABCD")

		gtest.Assert(gstr.ReplaceIByMap(s1, g.MapStrStr{
			"ab": "-",
			"CD": "=",
		}), "-=乱入的中文-=")
	})
}

func Test_ToLower(t *testing.T) {
	gtest.Case(t, func() {
		s1 := "abcdEFG乱入的中文abcdefg"
		e1 := "abcdefg乱入的中文abcdefg"
		gtest.Assert(gstr.ToLower(s1), e1)
	})
}

func Test_ToUpper(t *testing.T) {
	gtest.Case(t, func() {
		s1 := "abcdEFG乱入的中文abcdefg"
		e1 := "ABCDEFG乱入的中文ABCDEFG"
		gtest.Assert(gstr.ToUpper(s1), e1)
	})
}

func Test_UcFirst(t *testing.T) {
	gtest.Case(t, func() {
		s1 := "abcdEFG乱入的中文abcdefg"
		e1 := "AbcdEFG乱入的中文abcdefg"
		gtest.Assert(gstr.UcFirst(""), "")
		gtest.Assert(gstr.UcFirst(s1), e1)
		gtest.Assert(gstr.UcFirst(e1), e1)
	})
}

func Test_LcFirst(t *testing.T) {
	gtest.Case(t, func() {
		s1 := "AbcdEFG乱入的中文abcdefg"
		e1 := "abcdEFG乱入的中文abcdefg"
		gtest.Assert(gstr.LcFirst(""), "")
		gtest.Assert(gstr.LcFirst(s1), e1)
		gtest.Assert(gstr.LcFirst(e1), e1)
	})
}

func Test_UcWords(t *testing.T) {
	gtest.Case(t, func() {
		s1 := "我爱GF: i love go frame"
		e1 := "我爱GF: I Love Go Frame"
		gtest.Assert(gstr.UcWords(s1), e1)
	})
}

func Test_IsLetterLower(t *testing.T) {
	gtest.Case(t, func() {
		gtest.Assert(gstr.IsLetterLower('a'), true)
		gtest.Assert(gstr.IsLetterLower('A'), false)
		gtest.Assert(gstr.IsLetterLower('1'), false)
	})
}

func Test_IsLetterUpper(t *testing.T) {
	gtest.Case(t, func() {
		gtest.Assert(gstr.IsLetterUpper('a'), false)
		gtest.Assert(gstr.IsLetterUpper('A'), true)
		gtest.Assert(gstr.IsLetterUpper('1'), false)
	})
}

func Test_IsNumeric(t *testing.T) {
	gtest.Case(t, func() {
		gtest.Assert(gstr.IsNumeric("1a我"), false)
		gtest.Assert(gstr.IsNumeric("0123"), true)
		gtest.Assert(gstr.IsNumeric("我是中国人"), false)
	})
}

func Test_SubStr(t *testing.T) {
	gtest.Case(t, func() {
		gtest.Assert(gstr.SubStr("我爱GoFrame", 0), "我爱GoFrame")
		gtest.Assert(gstr.SubStr("我爱GoFrame", 2), "GoFrame")
		gtest.Assert(gstr.SubStr("我爱GoFrame", 2, 2), "Go")
		gtest.Assert(gstr.SubStr("我爱GoFrame", -1, 30), "我爱GoFrame")
		gtest.Assert(gstr.SubStr("我爱GoFrame", 30, 30), "")
	})
}

func Test_StrLimit(t *testing.T) {
	gtest.Case(t, func() {
		gtest.Assert(gstr.StrLimit("我爱GoFrame", 2), "我爱...")
		gtest.Assert(gstr.StrLimit("我爱GoFrame", 2, ""), "我爱")
		gtest.Assert(gstr.StrLimit("我爱GoFrame", 2, "**"), "我爱**")
		gtest.Assert(gstr.StrLimit("我爱GoFrame", 4, ""), "我爱Go")
		gtest.Assert(gstr.StrLimit("*", 4, ""), "*")
	})
}

func Test_HasPrefix(t *testing.T) {
	gtest.Case(t, func() {
		gtest.Assert(gstr.HasPrefix("我爱GoFrame", "我爱"), true)
		gtest.Assert(gstr.HasPrefix("en我爱GoFrame", "我爱"), false)
		gtest.Assert(gstr.HasPrefix("en我爱GoFrame", "en"), true)
	})
}

func Test_HasSuffix(t *testing.T) {
	gtest.Case(t, func() {
		gtest.Assert(gstr.HasSuffix("我爱GoFrame", "GoFrame"), true)
		gtest.Assert(gstr.HasSuffix("en我爱GoFrame", "a"), false)
		gtest.Assert(gstr.HasSuffix("GoFrame很棒", "棒"), true)
	})
}

func Test_Reverse(t *testing.T) {
	gtest.Case(t, func() {
		gtest.Assert(gstr.Reverse("我爱123"), "321爱我")
	})
}

func Test_NumberFormat(t *testing.T) {
	gtest.Case(t, func() {
		gtest.Assert(gstr.NumberFormat(1234567.8910, 2, ".", ","), "1,234,567.89")
		gtest.Assert(gstr.NumberFormat(1234567.8910, 2, "#", "/"), "1/234/567#89")
		gtest.Assert(gstr.NumberFormat(-1234567.8910, 2, "#", "/"), "-1/234/567#89")
	})
}

func Test_ChunkSplit(t *testing.T) {
	gtest.Case(t, func() {
		gtest.Assert(gstr.ChunkSplit("1234", 1, "#"), "1#2#3#4#")
		gtest.Assert(gstr.ChunkSplit("我爱123", 1, "#"), "我#爱#1#2#3#")
		gtest.Assert(gstr.ChunkSplit("1234", 1, ""), "1\r\n2\r\n3\r\n4\r\n")
	})
}

func Test_SplitAndTrim(t *testing.T) {
	gtest.Case(t, func() {
		s := `

010    

020  

`
		a := gstr.SplitAndTrim(s, "\n", "0")
		gtest.Assert(len(a), 2)
		gtest.Assert(a[0], "1")
		gtest.Assert(a[1], "2")
	})
}

func Test_Fields(t *testing.T) {
	gtest.Case(t, func() {
		gtest.Assert(gstr.Fields("我爱 Go Frame"), []string{
			"我爱", "Go", "Frame",
		})
	})
}

func Test_CountWords(t *testing.T) {
	gtest.Case(t, func() {
		gtest.Assert(gstr.CountWords("我爱 Go Go Go"), map[string]int{
			"Go": 3,
			"我爱": 1,
		})
	})
}

func Test_CountChars(t *testing.T) {
	gtest.Case(t, func() {
		gtest.Assert(gstr.CountChars("我爱 Go Go Go"), map[string]int{
			" ": 3,
			"G": 3,
			"o": 3,
			"我": 1,
			"爱": 1,
		})
		gtest.Assert(gstr.CountChars("我爱 Go Go Go", true), map[string]int{
			"G": 3,
			"o": 3,
			"我": 1,
			"爱": 1,
		})
	})
}

func Test_WordWrap(t *testing.T) {
	gtest.Case(t, func() {
		gtest.Assert(gstr.WordWrap("12 34", 2, "<br>"), "12<br>34")
		gtest.Assert(gstr.WordWrap("12 34", 2, "\n"), "12\n34")
		gtest.Assert(gstr.WordWrap("A very long woooooooooooooooooord. and something", 7, "<br>"),
			"A very<br>long<br>woooooooooooooooooord.<br>and<br>something")
	})
}

func Test_RuneLen(t *testing.T) {
	gtest.Case(t, func() {
		gtest.Assert(gstr.RuneLen("1234"), 4)
		gtest.Assert(gstr.RuneLen("我爱GoFrame"), 9)
	})
}

func Test_Repeat(t *testing.T) {
	gtest.Case(t, func() {
		gtest.Assert(gstr.Repeat("go", 3), "gogogo")
		gtest.Assert(gstr.Repeat("好的", 3), "好的好的好的")
	})
}

func Test_Str(t *testing.T) {
	gtest.Case(t, func() {
		gtest.Assert(gstr.Str("name@example.com", "@"), "@example.com")
		gtest.Assert(gstr.Str("name@example.com", ""), "")
		gtest.Assert(gstr.Str("name@example.com", "z"), "")
	})
}

func Test_Shuffle(t *testing.T) {
	gtest.Case(t, func() {
		gtest.Assert(len(gstr.Shuffle("123456")), 6)
	})
}

func Test_Split(t *testing.T) {
	gtest.Case(t, func() {
		gtest.Assert(gstr.Split("1.2", "."), []string{"1", "2"})
		gtest.Assert(gstr.Split("我爱 - GoFrame", " - "), []string{"我爱", "GoFrame"})
	})
}

func Test_Join(t *testing.T) {
	gtest.Case(t, func() {
		gtest.Assert(gstr.Join([]string{"我爱", "GoFrame"}, " - "), "我爱 - GoFrame")
	})
}

func Test_Explode(t *testing.T) {
	gtest.Case(t, func() {
		gtest.Assert(gstr.Explode(" - ", "我爱 - GoFrame"), []string{"我爱", "GoFrame"})
	})
}

func Test_Implode(t *testing.T) {
	gtest.Case(t, func() {
		gtest.Assert(gstr.Implode(" - ", []string{"我爱", "GoFrame"}), "我爱 - GoFrame")
	})
}

func Test_Chr(t *testing.T) {
	gtest.Case(t, func() {
		gtest.Assert(gstr.Chr(65), "A")
	})
}

func Test_Ord(t *testing.T) {
	gtest.Case(t, func() {
		gtest.Assert(gstr.Ord("A"), 65)
	})
}

func Test_HideStr(t *testing.T) {
	gtest.Case(t, func() {
		gtest.Assert(gstr.HideStr("15928008611", 40, "*"), "159****8611")
		gtest.Assert(gstr.HideStr("john@kohg.cn", 40, "*"), "jo*n@kohg.cn")
	})
}

func Test_Nl2Br(t *testing.T) {
	gtest.Case(t, func() {
		gtest.Assert(gstr.Nl2Br("1\n2"), "1<br>2")
		gtest.Assert(gstr.Nl2Br("1\r\n2"), "1<br>2")
		gtest.Assert(gstr.Nl2Br("1\r\n2", true), "1<br />2")
	})
}

func Test_AddSlashes(t *testing.T) {
	gtest.Case(t, func() {
		gtest.Assert(gstr.AddSlashes(`1'2"3\`), `1\'2\"3\\`)
	})
}

func Test_StripSlashes(t *testing.T) {
	gtest.Case(t, func() {
		gtest.Assert(gstr.StripSlashes(`1\'2\"3\\`), `1'2"3\`)
	})
}

func Test_QuoteMeta(t *testing.T) {
	gtest.Case(t, func() {
		gtest.Assert(gstr.QuoteMeta(`.\+*?[^]($)`), `\.\\\+\*\?\[\^\]\(\$\)`)
		gtest.Assert(gstr.QuoteMeta(`.\+*中国?[^]($)`), `\.\\\+\*中国\?\[\^\]\(\$\)`)
		gtest.Assert(gstr.QuoteMeta(`.''`, `'`), `.\'\'`)
		gtest.Assert(gstr.QuoteMeta(`中国.''`, `'`), `中国.\'\'`)
	})
}

func Test_Count(t *testing.T) {
	gtest.Case(t, func() {
		s := "abcdaAD"
		gtest.Assert(gstr.Count(s, "0"), 0)
		gtest.Assert(gstr.Count(s, "a"), 2)
		gtest.Assert(gstr.Count(s, "b"), 1)
		gtest.Assert(gstr.Count(s, "d"), 1)
	})
}

func Test_CountI(t *testing.T) {
	gtest.Case(t, func() {
		s := "abcdaAD"
		gtest.Assert(gstr.CountI(s, "0"), 0)
		gtest.Assert(gstr.CountI(s, "a"), 3)
		gtest.Assert(gstr.CountI(s, "b"), 1)
		gtest.Assert(gstr.CountI(s, "d"), 2)
	})
}

func Test_Compare(t *testing.T) {
	gtest.Case(t, func() {
		gtest.Assert(gstr.Compare("a", "b"), -1)
		gtest.Assert(gstr.Compare("a", "a"), 0)
		gtest.Assert(gstr.Compare("b", "a"), 1)
	})
}

func Test_Equal(t *testing.T) {
	gtest.Case(t, func() {
		gtest.Assert(gstr.Equal("a", "A"), true)
		gtest.Assert(gstr.Equal("a", "a"), true)
		gtest.Assert(gstr.Equal("b", "a"), false)
	})
}

func Test_Contains(t *testing.T) {
	gtest.Case(t, func() {
		gtest.Assert(gstr.Contains("abc", "a"), true)
		gtest.Assert(gstr.Contains("abc", "A"), false)
		gtest.Assert(gstr.Contains("abc", "ab"), true)
		gtest.Assert(gstr.Contains("abc", "abc"), true)
	})
}

func Test_ContainsI(t *testing.T) {
	gtest.Case(t, func() {
		gtest.Assert(gstr.ContainsI("abc", "a"), true)
		gtest.Assert(gstr.ContainsI("abc", "A"), true)
		gtest.Assert(gstr.ContainsI("abc", "Ab"), true)
		gtest.Assert(gstr.ContainsI("abc", "ABC"), true)
		gtest.Assert(gstr.ContainsI("abc", "ABCD"), false)
		gtest.Assert(gstr.ContainsI("abc", "D"), false)
	})
}

func Test_ContainsAny(t *testing.T) {
	gtest.Case(t, func() {
		gtest.Assert(gstr.ContainsAny("abc", "a"), true)
		gtest.Assert(gstr.ContainsAny("abc", "cd"), true)
		gtest.Assert(gstr.ContainsAny("abc", "de"), false)
		gtest.Assert(gstr.ContainsAny("abc", "A"), false)
	})
}

func Test_SearchArray(t *testing.T) {
	gtest.Case(t, func() {
		a := g.SliceStr{"a", "b", "c"}
		gtest.AssertEQ(gstr.SearchArray(a, "a"), 0)
		gtest.AssertEQ(gstr.SearchArray(a, "b"), 1)
		gtest.AssertEQ(gstr.SearchArray(a, "c"), 2)
		gtest.AssertEQ(gstr.SearchArray(a, "d"), -1)
	})
}

func Test_InArray(t *testing.T) {
	gtest.Case(t, func() {
		a := g.SliceStr{"a", "b", "c"}
		gtest.AssertEQ(gstr.InArray(a, "a"), true)
		gtest.AssertEQ(gstr.InArray(a, "b"), true)
		gtest.AssertEQ(gstr.InArray(a, "c"), true)
		gtest.AssertEQ(gstr.InArray(a, "d"), false)
	})
}
