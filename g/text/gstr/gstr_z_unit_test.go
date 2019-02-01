// Copyright 2019 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// go test *.go -bench=".*"

package gstr_test

import (
    "gitee.com/johng/gf/g"
    "gitee.com/johng/gf/g/text/gstr"
    "gitee.com/johng/gf/g/test/gtest"
    "testing"
)

func Test_Replace(t *testing.T) {
    gtest.Case(t, func() {
        s1 := "abcdEFG乱入的中文abcdefg"
        gtest.Assert(gstr.Replace(s1, "ab", "AB"), "ABcdEFG乱入的中文ABcdefg")
        gtest.Assert(gstr.Replace(s1, "EF", "ef"), "abcdefG乱入的中文abcdefg")
        gtest.Assert(gstr.Replace(s1, "MN", "mn"), s1)
        gtest.Assert(gstr.ReplaceByMap(s1, g.MapStrStr{
            "a" : "A",
            "G" : "g",
        }), "AbcdEFg乱入的中文Abcdefg")
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
        gtest.Assert(gstr.UcFirst(s1), e1)
    })
}

func Test_LcFirst(t *testing.T) {
    gtest.Case(t, func() {
        s1 := "AbcdEFG乱入的中文abcdefg"
        e1 := "abcdEFG乱入的中文abcdefg"
        gtest.Assert(gstr.LcFirst(s1), e1)
    })
}

func Test_UcWords(t *testing.T) {
    gtest.Case(t, func() {
        s1 := "我爱GF: i love go frame"
        e1 := "我爱GF: I Love Go Frame"
        gtest.Assert(gstr.UcWords(s1), e1)
    })
}

func Test_SearchArray(t *testing.T) {
    gtest.Case(t, func() {
        array := []string{"a", "b", "c"}
        gtest.Assert(gstr.SearchArray(array, "a"),  0)
        gtest.Assert(gstr.SearchArray(array, "b"),  1)
        gtest.Assert(gstr.SearchArray(array, "c"),  2)
        gtest.Assert(gstr.SearchArray(array, "d"), -1)
    })
}

func Test_InArray(t *testing.T) {
    gtest.Case(t, func() {
        array := []string{"a", "b", "c"}
        gtest.Assert(gstr.InArray(array, "a"), true)
        gtest.Assert(gstr.InArray(array, "b"), true)
        gtest.Assert(gstr.InArray(array, "c"), true)
        gtest.Assert(gstr.InArray(array, "d"), false)
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
        gtest.Assert(gstr.SubStr("我爱GoFrame", 0),    "我爱GoFrame")
        gtest.Assert(gstr.SubStr("我爱GoFrame", 2),    "GoFrame")
        gtest.Assert(gstr.SubStr("我爱GoFrame", 2, 2), "Go")
    })
}

func Test_StrLimit(t *testing.T) {
    gtest.Case(t, func() {
        gtest.Assert(gstr.StrLimit("我爱GoFrame", 2),        "我爱...")
        gtest.Assert(gstr.StrLimit("我爱GoFrame", 2, ""),    "我爱")
        gtest.Assert(gstr.StrLimit("我爱GoFrame", 2, "**"),  "我爱**")
        gtest.Assert(gstr.StrLimit("我爱GoFrame", 4, ""),    "我爱Go")
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
    })
}

func Test_ChunkSplit(t *testing.T) {
    gtest.Case(t, func() {
        gtest.Assert(gstr.ChunkSplit("1234",   1, "#"), "1#2#3#4#")
        gtest.Assert(gstr.ChunkSplit("我爱123", 1, "#"), "我#爱#1#2#3#")
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
            "Go"  : 3,
            "我爱" : 1,
        })
    })
}

func Test_CountChars(t *testing.T) {
    gtest.Case(t, func() {
        gtest.Assert(gstr.CountChars("我爱 Go Go Go"), map[string]int{
            " "  : 3,
            "G"  : 3,
            "o"  : 3,
            "我" : 1,
            "爱" : 1,
        })
        gtest.Assert(gstr.CountChars("我爱 Go Go Go", true), map[string]int{
            "G"  : 3,
            "o"  : 3,
            "我" : 1,
            "爱" : 1,
        })
    })
}

func Test_WordWrap(t *testing.T) {
    gtest.Case(t, func() {
        gtest.Assert(gstr.WordWrap("12 34", 2, "<br>"), "12<br>34")
        gtest.Assert(gstr.WordWrap("A very long woooooooooooooooooord. and something", 7, "<br>"),
            "A very<br>long<br>woooooooooooooooooord.<br>and<br>something")
    })
}

func Test_RuneLen(t *testing.T) {
    gtest.Case(t, func() {
        gtest.Assert(gstr.RuneLen("1234"),       4)
        gtest.Assert(gstr.RuneLen("我爱GoFrame"), 9)
    })
}

func Test_Repeat(t *testing.T) {
    gtest.Case(t, func() {
        gtest.Assert(gstr.Repeat("go",  3), "gogogo")
        gtest.Assert(gstr.Repeat("好的", 3), "好的好的好的")
    })
}

func Test_Str(t *testing.T) {
    gtest.Case(t, func() {
        gtest.Assert(gstr.Str("name@example.com", "@"), "@example.com")
    })
}

func Test_Shuffle(t *testing.T) {
    gtest.Case(t, func() {
        gtest.Assert(len(gstr.Shuffle("123456")), 6)
    })
}

func Test_Trim(t *testing.T) {
    gtest.Case(t, func() {
        gtest.Assert(gstr.Trim(" 123456\n "), "123456")
        gtest.Assert(gstr.Trim("#123456#;", "#;"), "123456")
    })
}

func Test_TrimRight(t *testing.T) {
    gtest.Case(t, func() {
        gtest.Assert(gstr.TrimRight(" 123456\n "), " 123456")
        gtest.Assert(gstr.TrimRight("#123456#;", "#;"), "#123456")
    })
}

func Test_TrimRightStr(t *testing.T) {
    gtest.Case(t, func() {
        gtest.Assert(gstr.TrimRightStr("gogo我爱gogo", "go"), "gogo我爱")
    })
}

func Test_TrimLeft(t *testing.T) {
    gtest.Case(t, func() {
        gtest.Assert(gstr.TrimLeft(" \r123456\n "), "123456\n ")
        gtest.Assert(gstr.TrimLeft("#;123456#;", "#;"), "123456#;")
    })
}

func Test_TrimLeftStr(t *testing.T) {
    gtest.Case(t, func() {
        gtest.Assert(gstr.TrimLeftStr("gogo我爱gogo", "go"), "我爱gogo")
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
    })
}