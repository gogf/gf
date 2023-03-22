// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*"

package gstr_test

import (
	"testing"

	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/text/gstr"
)

func Test_ToLower(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := "abcdEFG乱入的中文abcdefg"
		e1 := "abcdefg乱入的中文abcdefg"
		t.Assert(gstr.ToLower(s1), e1)
	})
}

func Test_ToUpper(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := "abcdEFG乱入的中文abcdefg"
		e1 := "ABCDEFG乱入的中文ABCDEFG"
		t.Assert(gstr.ToUpper(s1), e1)
	})
}

func Test_UcFirst(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := "abcdEFG乱入的中文abcdefg"
		e1 := "AbcdEFG乱入的中文abcdefg"
		t.Assert(gstr.UcFirst(""), "")
		t.Assert(gstr.UcFirst(s1), e1)
		t.Assert(gstr.UcFirst(e1), e1)
	})
}

func Test_LcFirst(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := "AbcdEFG乱入的中文abcdefg"
		e1 := "abcdEFG乱入的中文abcdefg"
		t.Assert(gstr.LcFirst(""), "")
		t.Assert(gstr.LcFirst(s1), e1)
		t.Assert(gstr.LcFirst(e1), e1)
	})
}

func Test_UcWords(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := "我爱GF: i love go frame"
		e1 := "我爱GF: I Love Go Frame"
		t.Assert(gstr.UcWords(s1), e1)
	})
}

func Test_IsLetterLower(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gstr.IsLetterLower('a'), true)
		t.Assert(gstr.IsLetterLower('A'), false)
		t.Assert(gstr.IsLetterLower('1'), false)
	})
}

func Test_IsLetterUpper(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gstr.IsLetterUpper('a'), false)
		t.Assert(gstr.IsLetterUpper('A'), true)
		t.Assert(gstr.IsLetterUpper('1'), false)
	})
}

func Test_IsNumeric(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gstr.IsNumeric("1a我"), false)
		t.Assert(gstr.IsNumeric("0123"), true)
		t.Assert(gstr.IsNumeric("我是中国人"), false)
		t.Assert(gstr.IsNumeric("1.2.3.4"), false)
	})
}

func Test_SubStr(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gstr.SubStr("我爱GoFrame", 0), "我爱GoFrame")
		t.Assert(gstr.SubStr("我爱GoFrame", 6), "GoFrame")
		t.Assert(gstr.SubStr("我爱GoFrame", 6, 2), "Go")
		t.Assert(gstr.SubStr("我爱GoFrame", -1, 30), "e")
		t.Assert(gstr.SubStr("我爱GoFrame", 30, 30), "")
		t.Assert(gstr.SubStr("abcdef", 0, -1), "abcde")
		t.Assert(gstr.SubStr("abcdef", 2, -1), "cde")
		t.Assert(gstr.SubStr("abcdef", 4, -4), "")
		t.Assert(gstr.SubStr("abcdef", -3, -1), "de")
	})
}

func Test_SubStrRune(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gstr.SubStrRune("我爱GoFrame", 0), "我爱GoFrame")
		t.Assert(gstr.SubStrRune("我爱GoFrame", 2), "GoFrame")
		t.Assert(gstr.SubStrRune("我爱GoFrame", 2, 2), "Go")
		t.Assert(gstr.SubStrRune("我爱GoFrame", -1, 30), "e")
		t.Assert(gstr.SubStrRune("我爱GoFrame", 30, 30), "")
		t.Assert(gstr.SubStrRune("abcdef", 0, -1), "abcde")
		t.Assert(gstr.SubStrRune("abcdef", 2, -1), "cde")
		t.Assert(gstr.SubStrRune("abcdef", 4, -4), "")
		t.Assert(gstr.SubStrRune("abcdef", -3, -1), "de")
		t.Assert(gstr.SubStrRune("我爱GoFrame呵呵", -3, 100), "e呵呵")
		t.Assert(gstr.SubStrRune("abcdef哈哈", -3, -1), "f哈")
		t.Assert(gstr.SubStrRune("ab我爱GoFramecdef哈哈", -3, -1), "f哈")
		t.Assert(gstr.SubStrRune("我爱GoFrame", 0, 3), "我爱G")
	})
}

func Test_StrLimit(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gstr.StrLimit("我爱GoFrame", 6), "我爱...")
		t.Assert(gstr.StrLimit("我爱GoFrame", 6, ""), "我爱")
		t.Assert(gstr.StrLimit("我爱GoFrame", 6, "**"), "我爱**")
		t.Assert(gstr.StrLimit("我爱GoFrame", 8, ""), "我爱Go")
		t.Assert(gstr.StrLimit("*", 4, ""), "*")
	})
}

func Test_StrLimitRune(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gstr.StrLimitRune("我爱GoFrame", 2), "我爱...")
		t.Assert(gstr.StrLimitRune("我爱GoFrame", 2, ""), "我爱")
		t.Assert(gstr.StrLimitRune("我爱GoFrame", 2, "**"), "我爱**")
		t.Assert(gstr.StrLimitRune("我爱GoFrame", 4, ""), "我爱Go")
		t.Assert(gstr.StrLimitRune("*", 4, ""), "*")
	})
}

func Test_HasPrefix(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gstr.HasPrefix("我爱GoFrame", "我爱"), true)
		t.Assert(gstr.HasPrefix("en我爱GoFrame", "我爱"), false)
		t.Assert(gstr.HasPrefix("en我爱GoFrame", "en"), true)
	})
}

func Test_HasSuffix(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gstr.HasSuffix("我爱GoFrame", "GoFrame"), true)
		t.Assert(gstr.HasSuffix("en我爱GoFrame", "a"), false)
		t.Assert(gstr.HasSuffix("GoFrame很棒", "棒"), true)
	})
}

func Test_Reverse(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gstr.Reverse("我爱123"), "321爱我")
	})
}

func Test_NumberFormat(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gstr.NumberFormat(1234567.8910, 2, ".", ","), "1,234,567.89")
		t.Assert(gstr.NumberFormat(1234567.8910, 2, "#", "/"), "1/234/567#89")
		t.Assert(gstr.NumberFormat(-1234567.8910, 2, "#", "/"), "-1/234/567#89")
	})
}

func Test_ChunkSplit(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gstr.ChunkSplit("1234", 1, "#"), "1#2#3#4#")
		t.Assert(gstr.ChunkSplit("我爱123", 1, "#"), "我#爱#1#2#3#")
		t.Assert(gstr.ChunkSplit("1234", 1, ""), "1\r\n2\r\n3\r\n4\r\n")
	})
}

func Test_SplitAndTrim(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := `

010    

020  

`
		a := gstr.SplitAndTrim(s, "\n", "0")
		t.Assert(len(a), 2)
		t.Assert(a[0], "1")
		t.Assert(a[1], "2")
	})
}

func Test_Fields(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gstr.Fields("我爱 Go Frame"), []string{
			"我爱", "Go", "Frame",
		})
	})
}

func Test_CountWords(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gstr.CountWords("我爱 Go Go Go"), map[string]int{
			"Go": 3,
			"我爱": 1,
		})
	})
}

func Test_CountChars(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gstr.CountChars("我爱 Go Go Go"), map[string]int{
			" ": 3,
			"G": 3,
			"o": 3,
			"我": 1,
			"爱": 1,
		})
		t.Assert(gstr.CountChars("我爱 Go Go Go", true), map[string]int{
			"G": 3,
			"o": 3,
			"我": 1,
			"爱": 1,
		})
	})
}

func Test_LenRune(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gstr.LenRune("1234"), 4)
		t.Assert(gstr.LenRune("我爱GoFrame"), 9)
	})
}

func Test_Repeat(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gstr.Repeat("go", 3), "gogogo")
		t.Assert(gstr.Repeat("好的", 3), "好的好的好的")
	})
}

func Test_Str(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gstr.Str("name@example.com", "@"), "@example.com")
		t.Assert(gstr.Str("name@example.com", ""), "")
		t.Assert(gstr.Str("name@example.com", "z"), "")
	})
}

func Test_StrEx(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gstr.StrEx("name@example.com", "@"), "example.com")
		t.Assert(gstr.StrEx("name@example.com", ""), "")
		t.Assert(gstr.StrEx("name@example.com", "z"), "")
	})
}

func Test_StrTill(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gstr.StrTill("name@example.com", "@"), "name@")
		t.Assert(gstr.StrTill("name@example.com", ""), "")
		t.Assert(gstr.StrTill("name@example.com", "z"), "")
	})
}

func Test_StrTillEx(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gstr.StrTillEx("name@example.com", "@"), "name")
		t.Assert(gstr.StrTillEx("name@example.com", ""), "")
		t.Assert(gstr.StrTillEx("name@example.com", "z"), "")
	})
}

func Test_Shuffle(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(len(gstr.Shuffle("123456")), 6)
	})
}

func Test_Split(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gstr.Split("1.2", "."), []string{"1", "2"})
		t.Assert(gstr.Split("我爱 - GoFrame", " - "), []string{"我爱", "GoFrame"})
	})
}

func Test_Join(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gstr.Join([]string{"我爱", "GoFrame"}, " - "), "我爱 - GoFrame")
	})
}

func Test_Explode(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gstr.Explode(" - ", "我爱 - GoFrame"), []string{"我爱", "GoFrame"})
	})
}

func Test_Implode(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gstr.Implode(" - ", []string{"我爱", "GoFrame"}), "我爱 - GoFrame")
	})
}

func Test_Chr(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gstr.Chr(65), "A")
	})
}

func Test_Ord(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gstr.Ord("A"), 65)
	})
}

func Test_HideStr(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gstr.HideStr("15928008611", 40, "*"), "159****8611")
		t.Assert(gstr.HideStr("john@kohg.cn", 40, "*"), "jo*n@kohg.cn")
		t.Assert(gstr.HideStr("张三", 50, "*"), "张*")
		t.Assert(gstr.HideStr("张小三", 50, "*"), "张*三")
		t.Assert(gstr.HideStr("欧阳小三", 50, "*"), "欧**三")
	})
}

func Test_Nl2Br(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gstr.Nl2Br("1\n2"), "1<br>2")
		t.Assert(gstr.Nl2Br("1\r\n2"), "1<br>2")
		t.Assert(gstr.Nl2Br("1\r\n2", true), "1<br />2")
	})
}

func Test_AddSlashes(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gstr.AddSlashes(`1'2"3\`), `1\'2\"3\\`)
	})
}

func Test_StripSlashes(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gstr.StripSlashes(`1\'2\"3\\`), `1'2"3\`)
	})
}

func Test_QuoteMeta(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gstr.QuoteMeta(`.\+*?[^]($)`), `\.\\\+\*\?\[\^\]\(\$\)`)
		t.Assert(gstr.QuoteMeta(`.\+*中国?[^]($)`), `\.\\\+\*中国\?\[\^\]\(\$\)`)
		t.Assert(gstr.QuoteMeta(`.''`, `'`), `.\'\'`)
		t.Assert(gstr.QuoteMeta(`中国.''`, `'`), `中国.\'\'`)
	})
}

func Test_Count(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := "abcdaAD"
		t.Assert(gstr.Count(s, "0"), 0)
		t.Assert(gstr.Count(s, "a"), 2)
		t.Assert(gstr.Count(s, "b"), 1)
		t.Assert(gstr.Count(s, "d"), 1)
	})
}

func Test_CountI(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := "abcdaAD"
		t.Assert(gstr.CountI(s, "0"), 0)
		t.Assert(gstr.CountI(s, "a"), 3)
		t.Assert(gstr.CountI(s, "b"), 1)
		t.Assert(gstr.CountI(s, "d"), 2)
	})
}

func Test_Compare(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gstr.Compare("a", "b"), -1)
		t.Assert(gstr.Compare("a", "a"), 0)
		t.Assert(gstr.Compare("b", "a"), 1)
	})
}

func Test_Equal(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gstr.Equal("a", "A"), true)
		t.Assert(gstr.Equal("a", "a"), true)
		t.Assert(gstr.Equal("b", "a"), false)
	})
}

func Test_Contains(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gstr.Contains("abc", "a"), true)
		t.Assert(gstr.Contains("abc", "A"), false)
		t.Assert(gstr.Contains("abc", "ab"), true)
		t.Assert(gstr.Contains("abc", "abc"), true)
	})
}

func Test_ContainsI(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gstr.ContainsI("abc", "a"), true)
		t.Assert(gstr.ContainsI("abc", "A"), true)
		t.Assert(gstr.ContainsI("abc", "Ab"), true)
		t.Assert(gstr.ContainsI("abc", "ABC"), true)
		t.Assert(gstr.ContainsI("abc", "ABCD"), false)
		t.Assert(gstr.ContainsI("abc", "D"), false)
	})
}

func Test_ContainsAny(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gstr.ContainsAny("abc", "a"), true)
		t.Assert(gstr.ContainsAny("abc", "cd"), true)
		t.Assert(gstr.ContainsAny("abc", "de"), false)
		t.Assert(gstr.ContainsAny("abc", "A"), false)
	})
}

func Test_SubStrFrom(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gstr.SubStrFrom("我爱GoFrameGood", `G`), "GoFrameGood")
		t.Assert(gstr.SubStrFrom("我爱GoFrameGood", `GG`), "")
		t.Assert(gstr.SubStrFrom("我爱GoFrameGood", `我`), "我爱GoFrameGood")
		t.Assert(gstr.SubStrFrom("我爱GoFrameGood", `Frame`), "FrameGood")
	})
}

func Test_SubStrFromEx(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gstr.SubStrFromEx("我爱GoFrameGood", `Go`), "FrameGood")
		t.Assert(gstr.SubStrFromEx("我爱GoFrameGood", `GG`), "")
		t.Assert(gstr.SubStrFromEx("我爱GoFrameGood", `我`), "爱GoFrameGood")
		t.Assert(gstr.SubStrFromEx("我爱GoFrameGood", `Frame`), `Good`)
	})
}

func Test_SubStrFromR(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gstr.SubStrFromR("我爱GoFrameGood", `G`), "Good")
		t.Assert(gstr.SubStrFromR("我爱GoFrameGood", `GG`), "")
		t.Assert(gstr.SubStrFromR("我爱GoFrameGood", `我`), "我爱GoFrameGood")
		t.Assert(gstr.SubStrFromR("我爱GoFrameGood", `Frame`), "FrameGood")
	})
}

func Test_SubStrFromREx(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gstr.SubStrFromREx("我爱GoFrameGood", `G`), "ood")
		t.Assert(gstr.SubStrFromREx("我爱GoFrameGood", `GG`), "")
		t.Assert(gstr.SubStrFromREx("我爱GoFrameGood", `我`), "爱GoFrameGood")
		t.Assert(gstr.SubStrFromREx("我爱GoFrameGood", `Frame`), `Good`)
	})
}
