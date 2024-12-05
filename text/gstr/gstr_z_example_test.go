// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gstr_test

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/text/gstr"
)

func ExampleCount() {
	var (
		str     = `goframe is very, very easy to use`
		substr1 = "goframe"
		substr2 = "very"
		result1 = gstr.Count(str, substr1)
		result2 = gstr.Count(str, substr2)
	)
	fmt.Println(result1)
	fmt.Println(result2)

	// Output:
	// 1
	// 2
}

func ExampleCountI() {
	var (
		str     = `goframe is very, very easy to use`
		substr1 = "GOFRAME"
		substr2 = "VERY"
		result1 = gstr.CountI(str, substr1)
		result2 = gstr.CountI(str, substr2)
	)
	fmt.Println(result1)
	fmt.Println(result2)

	// Output:
	// 1
	// 2
}

func ExampleToLower() {
	var (
		s      = `GOFRAME`
		result = gstr.ToLower(s)
	)
	fmt.Println(result)

	// Output:
	// goframe
}

func ExampleToUpper() {
	var (
		s      = `goframe`
		result = gstr.ToUpper(s)
	)
	fmt.Println(result)

	// Output:
	// GOFRAME
}

func ExampleUcFirst() {
	var (
		s      = `hello`
		result = gstr.UcFirst(s)
	)
	fmt.Println(result)

	// Output:
	// Hello
}

func ExampleLcFirst() {
	var (
		str    = `Goframe`
		result = gstr.LcFirst(str)
	)
	fmt.Println(result)

	// Output:
	// goframe
}

func ExampleUcWords() {
	var (
		str    = `hello world`
		result = gstr.UcWords(str)
	)
	fmt.Println(result)

	// Output:
	// Hello World
}

func ExampleIsLetterLower() {
	fmt.Println(gstr.IsLetterLower('a'))
	fmt.Println(gstr.IsLetterLower('A'))

	// Output:
	// true
	// false
}

func ExampleIsLetterUpper() {
	fmt.Println(gstr.IsLetterUpper('A'))
	fmt.Println(gstr.IsLetterUpper('a'))

	// Output:
	// true
	// false
}

func ExampleIsNumeric() {
	fmt.Println(gstr.IsNumeric("88"))
	fmt.Println(gstr.IsNumeric("3.1415926"))
	fmt.Println(gstr.IsNumeric("abc"))
	// Output:
	// true
	// true
	// false
}

func ExampleReverse() {
	var (
		str    = `123456`
		result = gstr.Reverse(str)
	)
	fmt.Println(result)

	// Output:
	// 654321
}

func ExampleNumberFormat() {
	var (
		number       float64 = 123456
		decimals             = 2
		decPoint             = "."
		thousandsSep         = ","
		result               = gstr.NumberFormat(number, decimals, decPoint, thousandsSep)
	)
	fmt.Println(result)

	// Output:
	// 123,456.00
}

func ExampleChunkSplit() {
	var (
		body     = `1234567890`
		chunkLen = 2
		end      = "#"
		result   = gstr.ChunkSplit(body, chunkLen, end)
	)
	fmt.Println(result)

	// Output:
	// 12#34#56#78#90#
}

func ExampleCompare() {
	fmt.Println(gstr.Compare("c", "c"))
	fmt.Println(gstr.Compare("a", "b"))
	fmt.Println(gstr.Compare("c", "b"))

	// Output:
	// 0
	// -1
	// 1
}

func ExampleEqual() {
	fmt.Println(gstr.Equal(`A`, `a`))
	fmt.Println(gstr.Equal(`A`, `A`))
	fmt.Println(gstr.Equal(`A`, `B`))

	// Output:
	// true
	// true
	// false
}

func ExampleFields() {
	var (
		str    = `Hello World`
		result = gstr.Fields(str)
	)
	fmt.Printf(`%#v`, result)

	// Output:
	// []string{"Hello", "World"}
}

func ExampleHasPrefix() {
	var (
		s      = `Hello World`
		prefix = "Hello"
		result = gstr.HasPrefix(s, prefix)
	)
	fmt.Println(result)

	// Output:
	// true
}

func ExampleHasSuffix() {
	var (
		s      = `my best love is goframe`
		prefix = "goframe"
		result = gstr.HasSuffix(s, prefix)
	)
	fmt.Println(result)

	// Output:
	// true
}

func ExampleCountWords() {
	var (
		str    = `goframe is very, very easy to use!`
		result = gstr.CountWords(str)
	)
	fmt.Printf(`%#v`, result)

	// Output:
	// map[string]int{"easy":1, "goframe":1, "is":1, "to":1, "use!":1, "very":1, "very,":1}
}

func ExampleCountChars() {
	var (
		str    = `goframe`
		result = gstr.CountChars(str)
	)
	fmt.Println(result)

	// May Output:
	// map[a:1 e:1 f:1 g:1 m:1 o:1 r:1]
}

func ExampleWordWrap() {
	{
		var (
			str    = `A very long woooooooooooooooooord. and something`
			width  = 8
			br     = "\n"
			result = gstr.WordWrap(str, width, br)
		)
		fmt.Println(result)
	}
	{
		var (
			str    = `The quick brown fox jumped over the lazy dog.`
			width  = 20
			br     = "<br />\n"
			result = gstr.WordWrap(str, width, br)
		)
		fmt.Printf("%v", result)
	}

	// Output:
	// A very
	// long
	// woooooooooooooooooord.
	// and
	// something
	// The quick brown fox<br />
	// jumped over the lazy<br />
	// dog.
}

func ExampleLenRune() {
	var (
		str    = `GoFrame框架`
		result = gstr.LenRune(str)
	)
	fmt.Println(result)

	// Output:
	// 9
}

func ExampleRepeat() {
	var (
		input      = `goframe `
		multiplier = 3
		result     = gstr.Repeat(input, multiplier)
	)
	fmt.Println(result)

	// Output:
	// goframe goframe goframe
}

func ExampleShuffle() {
	var (
		str    = `123456`
		result = gstr.Shuffle(str)
	)
	fmt.Println(result)

	// May Output:
	// 563214
}

func ExampleSplit() {
	var (
		str       = `a|b|c|d`
		delimiter = `|`
		result    = gstr.Split(str, delimiter)
	)
	fmt.Printf(`%#v`, result)

	// Output:
	// []string{"a", "b", "c", "d"}
}

func ExampleSplitAndTrim() {
	var (
		str       = `a|b|||||c|d`
		delimiter = `|`
		result    = gstr.SplitAndTrim(str, delimiter)
	)
	fmt.Printf(`%#v`, result)

	// Output:
	// []string{"a", "b", "c", "d"}
}

func ExampleJoin() {
	var (
		array  = []string{"goframe", "is", "very", "easy", "to", "use"}
		sep    = ` `
		result = gstr.Join(array, sep)
	)
	fmt.Println(result)

	// Output:
	// goframe is very easy to use
}

func ExampleJoinAny() {
	var (
		sep    = `,`
		arr2   = []int{99, 73, 85, 66}
		result = gstr.JoinAny(arr2, sep)
	)
	fmt.Println(result)

	// Output:
	// 99,73,85,66
}

func ExampleExplode() {
	var (
		str       = `Hello World`
		delimiter = " "
		result    = gstr.Explode(delimiter, str)
	)
	fmt.Printf(`%#v`, result)

	// Output:
	// []string{"Hello", "World"}
}

func ExampleImplode() {
	var (
		pieces = []string{"goframe", "is", "very", "easy", "to", "use"}
		glue   = " "
		result = gstr.Implode(glue, pieces)
	)
	fmt.Println(result)

	// Output:
	// goframe is very easy to use
}

func ExampleChr() {
	var (
		ascii  = 65 // A
		result = gstr.Chr(ascii)
	)
	fmt.Println(result)

	// Output:
	// A
}

// '103' is the 'g' in ASCII
func ExampleOrd() {
	var (
		str    = `goframe`
		result = gstr.Ord(str)
	)

	fmt.Println(result)

	// Output:
	// 103
}

func ExampleHideStr() {
	var (
		str     = `13800138000`
		percent = 40
		hide    = `*`
		result  = gstr.HideStr(str, percent, hide)
	)
	fmt.Println(result)

	// Output:
	// 138****8000
}

func ExampleNl2Br() {
	var (
		str = `goframe
is
very
easy
to
use`
		result = gstr.Nl2Br(str)
	)

	fmt.Println(result)

	// Output:
	// goframe<br>is<br>very<br>easy<br>to<br>use
}

func ExampleAddSlashes() {
	var (
		str    = `'aa'"bb"cc\r\n\d\t`
		result = gstr.AddSlashes(str)
	)

	fmt.Println(result)

	// Output:
	// \'aa\'\"bb\"cc\\r\\n\\d\\t
}

func ExampleStripSlashes() {
	var (
		str    = `C:\\windows\\GoFrame\\test`
		result = gstr.StripSlashes(str)
	)
	fmt.Println(result)

	// Output:
	// C:\windows\GoFrame\test
}

func ExampleQuoteMeta() {
	{
		var (
			str    = `.\+?[^]()`
			result = gstr.QuoteMeta(str)
		)
		fmt.Println(result)
	}
	{
		var (
			str    = `https://goframe.org/pages/viewpage.action?pageId=1114327`
			result = gstr.QuoteMeta(str)
		)
		fmt.Println(result)
	}

	// Output:
	// \.\\\+\?\[\^\]\(\)
	// https://goframe\.org/pages/viewpage\.action\?pageId=1114327

}

// array
func ExampleSearchArray() {
	var (
		array  = []string{"goframe", "is", "very", "nice"}
		str    = `goframe`
		result = gstr.SearchArray(array, str)
	)
	fmt.Println(result)

	// Output:
	// 0
}

func ExampleInArray() {
	var (
		a      = []string{"goframe", "is", "very", "easy", "to", "use"}
		s      = "goframe"
		result = gstr.InArray(a, s)
	)
	fmt.Println(result)

	// Output:
	// true
}

func ExamplePrefixArray() {
	var (
		strArray = []string{"tom", "lily", "john"}
	)

	gstr.PrefixArray(strArray, "classA_")

	fmt.Println(strArray)

	// Output:
	// [classA_tom classA_lily classA_john]
}

// case
func ExampleCaseCamel() {
	var (
		str    = `hello world`
		result = gstr.CaseCamel(str)
	)
	fmt.Println(result)

	// Output:
	// HelloWorld
}

func ExampleCaseCamelLower() {
	var (
		str    = `hello world`
		result = gstr.CaseCamelLower(str)
	)
	fmt.Println(result)

	// Output:
	// helloWorld
}

func ExampleCaseSnake() {
	var (
		str    = `hello world`
		result = gstr.CaseSnake(str)
	)
	fmt.Println(result)

	// Output:
	// hello_world
}

func ExampleCaseSnakeScreaming() {
	var (
		str    = `hello world`
		result = gstr.CaseSnakeScreaming(str)
	)
	fmt.Println(result)

	// Output:
	// HELLO_WORLD
}

func ExampleCaseSnakeFirstUpper() {
	var (
		str    = `RGBCodeMd5`
		result = gstr.CaseSnakeFirstUpper(str)
	)
	fmt.Println(result)

	// Output:
	// rgb_code_md5
}

func ExampleCaseKebab() {
	var (
		str    = `hello world`
		result = gstr.CaseKebab(str)
	)
	fmt.Println(result)

	// Output:
	// hello-world
}

func ExampleCaseKebabScreaming() {
	var (
		str    = `hello world`
		result = gstr.CaseKebabScreaming(str)
	)
	fmt.Println(result)

	// Output:
	// HELLO-WORLD
}

func ExampleCaseDelimited() {
	var (
		str    = `hello world`
		del    = byte('-')
		result = gstr.CaseDelimited(str, del)
	)
	fmt.Println(result)

	// Output:
	// hello-world
}

func ExampleCaseDelimitedScreaming() {
	{
		var (
			str    = `hello world`
			del    = byte('-')
			result = gstr.CaseDelimitedScreaming(str, del, true)
		)
		fmt.Println(result)
	}
	{
		var (
			str    = `hello world`
			del    = byte('-')
			result = gstr.CaseDelimitedScreaming(str, del, false)
		)
		fmt.Println(result)
	}

	// Output:
	// HELLO-WORLD
	// hello-world
}

// contain
func ExampleContains() {
	{
		var (
			str    = `Hello World`
			substr = `Hello`
			result = gstr.Contains(str, substr)
		)
		fmt.Println(result)
	}
	{
		var (
			str    = `Hello World`
			substr = `hello`
			result = gstr.Contains(str, substr)
		)
		fmt.Println(result)
	}

	// Output:
	// true
	// false
}

func ExampleContainsI() {
	var (
		str     = `Hello World`
		substr  = "hello"
		result1 = gstr.Contains(str, substr)
		result2 = gstr.ContainsI(str, substr)
	)
	fmt.Println(result1)
	fmt.Println(result2)

	// Output:
	// false
	// true
}

func ExampleContainsAny() {
	{
		var (
			s      = `goframe`
			chars  = "g"
			result = gstr.ContainsAny(s, chars)
		)
		fmt.Println(result)
	}
	{
		var (
			s      = `goframe`
			chars  = "G"
			result = gstr.ContainsAny(s, chars)
		)
		fmt.Println(result)
	}

	// Output:
	// true
	// false
}

// convert
func ExampleOctStr() {
	var (
		str    = `\346\200\241`
		result = gstr.OctStr(str)
	)
	fmt.Println(result)

	// Output:
	// 怡
}

// domain
func ExampleIsSubDomain() {
	var (
		subDomain  = `s.goframe.org`
		mainDomain = `goframe.org`
		result     = gstr.IsSubDomain(subDomain, mainDomain)
	)
	fmt.Println(result)

	// Output:
	// true
}

// levenshtein
func ExampleLevenshtein() {
	var (
		str1    = "Hello World"
		str2    = "hallo World"
		costIns = 1
		costRep = 1
		costDel = 1
		result  = gstr.Levenshtein(str1, str2, costIns, costRep, costDel)
	)
	fmt.Println(result)

	// Output:
	// 2
}

// parse
func ExampleParse() {
	{
		var (
			str       = `v1=m&v2=n`
			result, _ = gstr.Parse(str)
		)
		fmt.Println(result)
	}
	{
		var (
			str       = `v[a][a]=m&v[a][b]=n`
			result, _ = gstr.Parse(str)
		)
		fmt.Println(result)
	}
	{
		// The form of nested Slice is not yet supported.
		var str = `v[][]=m&v[][]=n`
		result, err := gstr.Parse(str)
		if err != nil {
			panic(err)
		}
		fmt.Println(result)
	}
	{
		// This will produce an error.
		var str = `v=m&v[a]=n`
		result, err := gstr.Parse(str)
		if err != nil {
			println(err)
		}
		fmt.Println(result)
	}
	{
		var (
			str       = `a .[[b=c`
			result, _ = gstr.Parse(str)
		)
		fmt.Println(result)
	}

	// May Output:
	// map[v1:m v2:n]
	// map[v:map[a:map[a:m b:n]]]
	// map[v:map[]]
	// Error: expected type 'map[string]interface{}' for key 'v', but got 'string'
	// map[]
	// map[a___[b:c]
}

// pos
func ExamplePos() {
	var (
		haystack = `Hello World`
		needle   = `World`
		result   = gstr.Pos(haystack, needle)
	)
	fmt.Println(result)

	// Output:
	// 6
}

func ExamplePosRune() {
	var (
		haystack = `GoFrame是一款模块化、高性能、企业级的Go基础开发框架`
		needle   = `Go`
		posI     = gstr.PosRune(haystack, needle)
		posR     = gstr.PosRRune(haystack, needle)
	)
	fmt.Println(posI)
	fmt.Println(posR)

	// Output:
	// 0
	// 22
}

func ExamplePosI() {
	var (
		haystack = `goframe is very, very easy to use`
		needle   = `very`
		posI     = gstr.PosI(haystack, needle)
		posR     = gstr.PosR(haystack, needle)
	)
	fmt.Println(posI)
	fmt.Println(posR)

	// Output:
	// 11
	// 17
}

func ExamplePosIRune() {
	{
		var (
			haystack    = `GoFrame是一款模块化、高性能、企业级的Go基础开发框架`
			needle      = `高性能`
			startOffset = 10
			result      = gstr.PosIRune(haystack, needle, startOffset)
		)
		fmt.Println(result)
	}
	{
		var (
			haystack    = `GoFrame是一款模块化、高性能、企业级的Go基础开发框架`
			needle      = `高性能`
			startOffset = 30
			result      = gstr.PosIRune(haystack, needle, startOffset)
		)
		fmt.Println(result)
	}

	// Output:
	// 14
	// -1
}

func ExamplePosR() {
	var (
		haystack = `goframe is very, very easy to use`
		needle   = `very`
		posI     = gstr.PosI(haystack, needle)
		posR     = gstr.PosR(haystack, needle)
	)
	fmt.Println(posI)
	fmt.Println(posR)

	// Output:
	// 11
	// 17
}

func ExamplePosRRune() {
	var (
		haystack = `GoFrame是一款模块化、高性能、企业级的Go基础开发框架`
		needle   = `Go`
		posI     = gstr.PosIRune(haystack, needle)
		posR     = gstr.PosRRune(haystack, needle)
	)
	fmt.Println(posI)
	fmt.Println(posR)

	// Output:
	// 0
	// 22
}

func ExamplePosRI() {
	var (
		haystack = `goframe is very, very easy to use`
		needle   = `VERY`
		posI     = gstr.PosI(haystack, needle)
		posR     = gstr.PosRI(haystack, needle)
	)
	fmt.Println(posI)
	fmt.Println(posR)

	// Output:
	// 11
	// 17
}

func ExamplePosRIRune() {
	var (
		haystack = `GoFrame是一款模块化、高性能、企业级的Go基础开发框架`
		needle   = `GO`
		posI     = gstr.PosIRune(haystack, needle)
		posR     = gstr.PosRIRune(haystack, needle)
	)
	fmt.Println(posI)
	fmt.Println(posR)

	// Output:
	// 0
	// 22
}

// replace
func ExampleReplace() {
	var (
		origin  = `golang is very nice!`
		search  = `golang`
		replace = `goframe`
		result  = gstr.Replace(origin, search, replace)
	)
	fmt.Println(result)

	// Output:
	// goframe is very nice!
}

func ExampleReplaceI() {
	var (
		origin  = `golang is very nice!`
		search  = `GOLANG`
		replace = `goframe`
		result  = gstr.ReplaceI(origin, search, replace)
	)
	fmt.Println(result)

	// Output:
	// goframe is very nice!
}

func ExampleReplaceByArray() {
	{
		var (
			origin = `golang is very nice`
			array  = []string{"lang", "frame"}
			result = gstr.ReplaceByArray(origin, array)
		)
		fmt.Println(result)
	}
	{
		var (
			origin = `golang is very good`
			array  = []string{"golang", "goframe", "good", "nice"}
			result = gstr.ReplaceByArray(origin, array)
		)
		fmt.Println(result)
	}

	// Output:
	// goframe is very nice
	// goframe is very nice
}

func ExampleReplaceIByArray() {
	var (
		origin = `golang is very Good`
		array  = []string{"Golang", "goframe", "GOOD", "nice"}
		result = gstr.ReplaceIByArray(origin, array)
	)

	fmt.Println(result)

	// Output:
	// goframe is very nice
}

func ExampleReplaceByMap() {
	{
		var (
			origin   = `golang is very nice`
			replaces = map[string]string{
				"lang": "frame",
			}
			result = gstr.ReplaceByMap(origin, replaces)
		)
		fmt.Println(result)
	}
	{
		var (
			origin   = `golang is very good`
			replaces = map[string]string{
				"golang": "goframe",
				"good":   "nice",
			}
			result = gstr.ReplaceByMap(origin, replaces)
		)
		fmt.Println(result)
	}

	// Output:
	// goframe is very nice
	// goframe is very nice
}

func ExampleReplaceIByMap() {
	var (
		origin   = `golang is very nice`
		replaces = map[string]string{
			"Lang": "frame",
		}
		result = gstr.ReplaceIByMap(origin, replaces)
	)
	fmt.Println(result)

	// Output:
	// goframe is very nice
}

func ExampleReplaceFunc() {
	str := "hello gf 2018~2020!"
	// Replace "gf" with a custom function that returns "GoFrame"
	result := gstr.ReplaceFunc(str, "gf", func(s string) string {
		return "GoFrame"
	})
	fmt.Println(result)

	// Replace numbers with their doubled values
	result = gstr.ReplaceFunc("1 2 3", "2", func(s string) string {
		n, _ := strconv.Atoi(s)
		return strconv.Itoa(n * 2)
	})
	fmt.Println(result)

	// Output:
	// hello GoFrame 2018~2020!
	// 1 4 3
}

func ExampleReplaceIFunc() {
	str := "Hello GF, hello gf, HELLO Gf!"
	// Replace any case variation of "gf" with "GoFrame"
	result := gstr.ReplaceIFunc(str, "gf", func(s string) string {
		return "GoFrame"
	})
	fmt.Println(result)

	// Preserve the original case of each match
	result = gstr.ReplaceIFunc(str, "gf", func(s string) string {
		if s == strings.ToUpper(s) {
			return "GOFRAME"
		}
		if s == strings.ToLower(s) {
			return "goframe"
		}
		return "GoFrame"
	})
	fmt.Println(result)

	// Output:
	// Hello GoFrame, hello GoFrame, HELLO GoFrame!
	// Hello GOFRAME, hello goframe, HELLO GoFrame!
}

// similartext
func ExampleSimilarText() {
	var (
		first   = `AaBbCcDd`
		second  = `ad`
		percent = 0.80
		result  = gstr.SimilarText(first, second, &percent)
	)
	fmt.Println(result)

	// Output:
	// 2
}

// soundex
func ExampleSoundex() {
	var (
		str1    = `Hello`
		str2    = `Hallo`
		result1 = gstr.Soundex(str1)
		result2 = gstr.Soundex(str2)
	)
	fmt.Println(result1, result2)

	// Output:
	// H400 H400
}

// str
func ExampleStr() {
	var (
		haystack = `xxx.jpg`
		needle   = `.`
		result   = gstr.Str(haystack, needle)
	)
	fmt.Println(result)

	// Output:
	// .jpg
}

func ExampleStrEx() {
	var (
		haystack = `https://goframe.org/index.html?a=1&b=2`
		needle   = `?`
		result   = gstr.StrEx(haystack, needle)
	)
	fmt.Println(result)

	// Output:
	// a=1&b=2
}

func ExampleStrTill() {
	var (
		haystack = `https://goframe.org/index.html?test=123456`
		needle   = `?`
		result   = gstr.StrTill(haystack, needle)
	)
	fmt.Println(result)

	// Output:
	// https://goframe.org/index.html?
}

func ExampleStrTillEx() {
	var (
		haystack = `https://goframe.org/index.html?test=123456`
		needle   = `?`
		result   = gstr.StrTillEx(haystack, needle)
	)
	fmt.Println(result)

	// Output:
	// https://goframe.org/index.html
}

// substr
func ExampleSubStr() {
	var (
		str    = `1234567890`
		start  = 0
		length = 4
		subStr = gstr.SubStr(str, start, length)
	)
	fmt.Println(subStr)

	// Output:
	// 1234
}

func ExampleSubStrRune() {
	var (
		str    = `GoFrame是一款模块化、高性能、企业级的Go基础开发框架。`
		start  = 14
		length = 3
		subStr = gstr.SubStrRune(str, start, length)
	)
	fmt.Println(subStr)

	// Output:
	// 高性能
}

func ExampleStrLimit() {
	var (
		str    = `123456789`
		length = 3
		suffix = `...`
		result = gstr.StrLimit(str, length, suffix)
	)
	fmt.Println(result)

	// Output:
	// 123...
}

func ExampleStrLimitRune() {
	var (
		str    = `GoFrame是一款模块化、高性能、企业级的Go基础开发框架。`
		length = 17
		suffix = "..."
		result = gstr.StrLimitRune(str, length, suffix)
	)
	fmt.Println(result)

	// Output:
	// GoFrame是一款模块化、高性能...
}

func ExampleSubStrFrom() {
	var (
		str  = "我爱GoFrameGood"
		need = `爱`
	)

	fmt.Println(gstr.SubStrFrom(str, need))

	// Output:
	// 爱GoFrameGood
}

func ExampleSubStrFromEx() {
	var (
		str  = "我爱GoFrameGood"
		need = `爱`
	)

	fmt.Println(gstr.SubStrFromEx(str, need))

	// Output:
	// GoFrameGood
}

func ExampleSubStrFromR() {
	var (
		str  = "我爱GoFrameGood"
		need = `Go`
	)

	fmt.Println(gstr.SubStrFromR(str, need))

	// Output:
	// Good
}

func ExampleSubStrFromREx() {
	var (
		str  = "我爱GoFrameGood"
		need = `Go`
	)

	fmt.Println(gstr.SubStrFromREx(str, need))

	// Output:
	// od
}

// trim
func ExampleTrim() {
	var (
		str           = `*Hello World*`
		characterMask = "*"
		result        = gstr.Trim(str, characterMask)
	)
	fmt.Println(result)

	// Output:
	// Hello World
}

func ExampleTrimStr() {
	var (
		str    = `Hello World`
		cut    = "World"
		count  = -1
		result = gstr.TrimStr(str, cut, count)
	)
	fmt.Println(result)

	// Output:
	// Hello
}

func ExampleTrimLeft() {
	var (
		str           = `*Hello World*`
		characterMask = "*"
		result        = gstr.TrimLeft(str, characterMask)
	)
	fmt.Println(result)

	// Output:
	// Hello World*
}

func ExampleTrimLeftStr() {
	var (
		str    = `**Hello World**`
		cut    = "*"
		count  = 1
		result = gstr.TrimLeftStr(str, cut, count)
	)
	fmt.Println(result)

	// Output:
	// *Hello World**
}

func ExampleTrimRight() {
	var (
		str           = `**Hello World**`
		characterMask = "*def" // []byte{"*", "d", "e", "f"}
		result        = gstr.TrimRight(str, characterMask)
	)
	fmt.Println(result)

	// Output:
	// **Hello Worl
}

func ExampleTrimRightStr() {
	var (
		str    = `Hello World!`
		cut    = "!"
		count  = -1
		result = gstr.TrimRightStr(str, cut, count)
	)
	fmt.Println(result)

	// Output:
	// Hello World
}

func ExampleTrimAll() {
	var (
		str           = `*Hello World*`
		characterMask = "*"
		result        = gstr.TrimAll(str, characterMask)
	)
	fmt.Println(result)

	// Output:
	// HelloWorld
}

// version
func ExampleCompareVersion() {
	fmt.Println(gstr.CompareVersion("v2.11.9", "v2.10.8"))
	fmt.Println(gstr.CompareVersion("1.10.8", "1.19.7"))
	fmt.Println(gstr.CompareVersion("2.8.beta", "2.8"))

	// Output:
	// 1
	// -1
	// 0
}

func ExampleCompareVersionGo() {
	fmt.Println(gstr.CompareVersionGo("v2.11.9", "v2.10.8"))
	fmt.Println(gstr.CompareVersionGo("v4.20.1", "v4.20.1+incompatible"))
	fmt.Println(gstr.CompareVersionGo(
		"v0.0.2-20180626092158-b2ccc119800e",
		"v1.0.1-20190626092158-b2ccc519800e",
	))

	// Output:
	// 1
	// 1
	// -1
}
