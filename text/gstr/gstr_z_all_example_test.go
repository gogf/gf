package gstr_test

import (
	"fmt"

	"github.com/gogf/gf/v2/text/gstr"
)

func ExampleAddSlashes() {
	var (
		str   = `'aa'"bb"cc\r\n\d\t`
		rsStr = gstr.AddSlashes(str)
	)

	fmt.Println(rsStr)

	// Output:
	// \'aa\'\"bb\"cc\\r\\n\\d\\t
}

func ExampleCaseCamel() {
	var (
		str   = `goframe_is very nice.to-use`
		rsStr = gstr.CaseCamel(str)
	)
	fmt.Println(rsStr)

	// Output:
	// GoframeIsVeryNiceToUse
}

func ExampleCaseCamelLower() {
	var (
		str   = `goframe_is very nice.to-use`
		rsStr = gstr.CaseCamelLower(str)
	)

	fmt.Println(rsStr)

	// Output:
	// goframeIsVeryNiceToUse
}

func ExampleCaseDelimited() {
	var (
		str   = `goframe_is_very-nice.to-use`
		del   = uint8(35)
		rsStr = gstr.CaseDelimited(str, del)
	)
	fmt.Println(rsStr)

	// Output:
	// goframe#is#very#nice#to#use
}

func ExampleCaseDelimitedScreaming() {
	{
		var (
			str   = `goframe_is Very Nice.to-use`
			del   = uint8(35)
			rsStr = gstr.CaseDelimitedScreaming(str, del, true)
		)
		fmt.Println(rsStr)
	}

	{
		var (
			str   = `GOFRAME_IS VERY NICE.TO-USE`
			del   = uint8(35)
			rsStr = gstr.CaseDelimitedScreaming(str, del, false)
		)
		fmt.Println(rsStr)
	}

	// Output:
	// GOFRAME#IS#VERY#NICE#TO#USE
	// goframe#is#very#nice#to#use
}

func ExampleCaseKebab() {
	var (
		str   = `goframe_is Very Nice.to-use`
		rsStr = gstr.CaseKebab(str)
	)
	fmt.Println(rsStr)

	// Output:
	// goframe-is-very-nice-to-use
}

func ExampleCaseKebabScreaming() {
	var (
		str   = `goframe_is Very Nice.to-use`
		rsStr = gstr.CaseKebabScreaming(str)
	)
	fmt.Println(rsStr)

	// Output:
	// GOFRAME-IS-VERY-NICE-TO-USE
}

func ExampleCaseSnake() {
	var (
		str   = `goframe_is Very Nice.to-use`
		rsStr = gstr.CaseSnake(str)
	)

	fmt.Println(rsStr)

	// Output:
	// goframe_is_very_nice_to_use
}

func ExampleCaseSnakeFirstUpper() {
	var (
		str   = `GoframeIsVeryNiceToUse`
		rsStr = gstr.CaseSnakeFirstUpper(str)
	)

	fmt.Println(rsStr)

	// Output:
	// goframe_is_very_nice_to_use
}

func ExampleCaseSnakeScreaming() {
	var (
		str   = `goframe_is Very Nice.to-use`
		rsStr = gstr.CaseSnakeScreaming(str)
	)

	fmt.Println(rsStr)

	// Output:
	// GOFRAME_IS_VERY_NICE_TO_USE
}

func ExampleChr() {
	var (
		ascii = 65
		rsStr = gstr.Chr(ascii)
	)
	fmt.Println(rsStr)

	// Output:
	// A
}

func ExampleChunkSplit() {
	{
		var (
			body     = `1234`
			chunkLen = 2
			end      = "#"
			rsStr    = gstr.ChunkSplit(body, chunkLen, end)
		)

		fmt.Println(rsStr)
	}

	{
		var (
			body     = `我爱Goframe`
			chunkLen = 1
			end      = "-"
			rsStr    = gstr.ChunkSplit(body, chunkLen, end)
		)
		fmt.Println(rsStr)
	}

	{
		var (
			body     = `1234`
			chunkLen = 1
			end      = ""
			rsStr    = gstr.ChunkSplit(body, chunkLen, end)
		)
		fmt.Println(rsStr)
	}

	// May Output:
	// 12#34#
	// 我-爱-G-o-f-r-a-m-e-
	// "1\r\n2\r\n3\r\n4\r\n"
}

func ExampleCompare() {
	{
		var (
			a     = "C"
			b     = "C"
			rsStr = gstr.Compare(a, b)
		)

		fmt.Println(rsStr)
	}

	{
		var (
			a     = "A"
			b     = "B"
			rsStr = gstr.Compare(a, b)
		)
		fmt.Println(rsStr)
	}

	{
		var (
			a     = "C"
			b     = "B"
			rsStr = gstr.Compare(a, b)
		)
		fmt.Println(rsStr)
	}

	// Output:
	// 0
	// -1
	// 1
}

func ExampleCompareVersion() {
	{
		var (
			a = "v2.11.9"
			b = "v2.10.8"

			rsStr = gstr.CompareVersion(a, b)
		)

		fmt.Println(rsStr)
	}

	{
		var (
			a     = "A.10.8"
			b     = "A.19.7"
			rsStr = gstr.CompareVersion(a, b)
		)
		fmt.Println(rsStr)
	}

	{
		var (
			a     = "2.8.beta"
			b     = "2.8"
			rsStr = gstr.CompareVersion(a, b)
		)
		fmt.Println(rsStr)
	}

	// Output:
	// 1
	// -1
	// 0
}

func ExampleCompareVersionGo() {

	{
		var (
			a     = "v2.11.9"
			b     = "v2.10.8"
			rsStr = gstr.CompareVersionGo(a, b)
		)
		fmt.Println(rsStr)
	}

	{
		var (
			a     = "v0.0.2-20180626092158-b2ccc119800e"
			b     = "v1.0.1-20190626092158-b2ccc519800e"
			rsStr = gstr.CompareVersionGo(a, b)
		)
		fmt.Println(rsStr)
	}

	{
		var (
			a     = "v4.20.1"
			b     = "v4.20.1+incompatible"
			rsStr = gstr.CompareVersionGo(a, b)
		)
		fmt.Println(rsStr)
	}

	// Output:
	// 1
	// -1
	// 0
}

func ExampleContains() {
	{
		var (
			str    = `goframe_is Very Nice.to-use`
			substr = `goframe`
			rsStr  = gstr.Contains(str, substr)
		)

		fmt.Println(rsStr)
	}

	{
		var (
			str    = `goframe_is Very Nice.to-use`
			substr = `Goframe`
			rsStr  = gstr.Contains(str, substr)
		)
		fmt.Println(rsStr)
	}

	// Output:
	// true
	// false
}

func ExampleContainsAny() {
	{
		var (
			s     = `goframe_is Very Nice.to-use`
			chars = "g"
			rsStr = gstr.ContainsAny(s, chars)
		)

		fmt.Println(rsStr)
	}

	{
		var (
			s     = `goframe_is Very Nice.to-use`
			chars = "G"
			rsStr = gstr.ContainsAny(s, chars)
		)
		fmt.Println(rsStr)
	}

	{
		var (
			s     = `goframe_is Very Nice.to-use`
			chars = "Nice"
			rsStr = gstr.ContainsAny(s, chars)
		)
		fmt.Println(rsStr)
	}

	// Output:
	// true
	// false
	// true
}

func ExampleContainsI() {
	{
		var (
			str    = `goframe_is Very Nice.to-use`
			substr = "go"
			rsStr  = gstr.ContainsI(str, substr)
		)

		fmt.Println(rsStr)
	}

	{
		var (
			str    = `goframe_is Very Nice.to-use`
			substr = "Go"
			rsStr  = gstr.ContainsI(str, substr)
		)
		fmt.Println(rsStr)
	}

	{
		var (
			str    = `goframe_is Very Nice.to-use`
			substr = "Golang"
			rsStr  = gstr.ContainsI(str, substr)
		)
		fmt.Println(rsStr)
	}

	// Output:
	// true
	// true
	// false
}

func ExampleCount() {
	{
		var (
			str    = `goframe_is Very Nice.to-use`
			substr = "goframe"
			rsStr  = gstr.Count(str, substr)
		)
		fmt.Println(rsStr)
	}

	{
		var (
			str    = `goframe_is Very Nice.to-use`
			substr = "golang"
			rsStr  = gstr.Count(str, substr)
		)
		fmt.Println(rsStr)
	}

	// Output:
	// 1
	// 0
}

func ExampleCountChars() {
	{
		var (
			str   = `goframe_is Very Nice.to-use`
			rsStr = gstr.CountChars(str)
		)
		fmt.Println(rsStr)
	}

	{
		var (
			str   = `goframe_is Very Nice.to-use`
			rsStr = gstr.CountChars(str, true)
		)
		fmt.Println(rsStr)
	}

	// May Output:
	// map[ :2 -:1 .:1 N:1 V:1 _:1 a:1 c:1 e:4 f:1 g:1 i:2 m:1 o:2 r:2 s:2 t:1 u:1 y:1]
	// map[-:1 .:1 N:1 V:1 _:1 a:1 c:1 e:4 f:1 g:1 i:2 m:1 o:2 r:2 s:2 t:1 u:1 y:1]

}

func ExampleCountI() {
	{
		var (
			s      = `goframe_is Very Nice.to-use`
			substr = "goframe"
			rsStr  = gstr.CountI(s, substr)
		)

		fmt.Println(rsStr)
	}

	{
		var (
			s      = `goframe_is Very Nice.to-use`
			substr = "Goframe"
			rsStr  = gstr.CountI(s, substr)
		)

		fmt.Println(rsStr)
	}

	{
		var (
			s      = `Goframe_is Very Nice.to-use`
			substr = "golang"
			rsStr  = gstr.CountI(s, substr)
		)
		fmt.Println(rsStr)
	}

	// Output:
	// 1
	// 1
	// 0
}

func ExampleCountWords() {
	var (
		str   = `goframe is Very Nice to use ! goframe Very Nice !`
		rsStr = gstr.CountWords(str)
	)

	fmt.Println(rsStr)

	// Output:
	// map[!:2 Nice:2 Very:2 goframe:2 is:1 to:1 use:1]
}

func ExampleEqual() {
	{
		var (
			a     = `A`
			b     = `a`
			rsStr = gstr.Equal(a, b)
		)

		fmt.Println(rsStr)
	}

	{
		var (
			a     = `A`
			b     = `A`
			rsStr = gstr.Equal(a, b)
		)
		fmt.Println(rsStr)
	}

	{
		var (
			a     = `C`
			b     = `B`
			rsStr = gstr.Equal(a, b)
		)
		fmt.Println(rsStr)
	}

	// Output:
	// true
	// true
	// false
}

func ExampleExplode() {
	var (
		str       = `goframe_is_Very_Nice_to_use`
		delimiter = "_"
		rsStr     = gstr.Explode(delimiter, str)
	)

	fmt.Printf(`%#v`, rsStr)

	// Output:
	// []string{"goframe", "is", "Very", "Nice", "to", "use"}
}

func ExampleFields() {
	var (
		str   = `goframe is Very Nice to-use`
		rsStr = gstr.Fields(str)
	)

	fmt.Printf(`%#v`, rsStr)

	// Output:
	// []string{"goframe", "is", "Very", "Nice", "to-use"}
}

func ExampleHasPrefix() {
	{
		var (
			s      = `goframe_is Very Nice.to-use`
			prefix = "goframe"
			rsStr  = gstr.HasPrefix(s, prefix)
		)

		fmt.Println(rsStr)
	}

	{
		var (
			s      = `goframe_is Very Nice.to-use`
			prefix = "Goframe"
			rsStr  = gstr.HasPrefix(s, prefix)
		)
		fmt.Println(rsStr)
	}

	// Output:
	// true
	// false
}

func ExampleHasSuffix() {
	{
		var (
			s      = `goframe_is Very Nice.to-use`
			prefix = "use"
			rsStr  = gstr.HasSuffix(s, prefix)
		)
		fmt.Println(rsStr)
	}

	{
		var (
			s      = `goframe_is Very Nice.to-use`
			prefix = "Use"
			rsStr  = gstr.HasSuffix(s, prefix)
		)
		fmt.Println(rsStr)
	}

	// Output:
	// true
	// false
}

func ExampleHideStr() {
	{
		var (
			str     = `13800138000`
			percent = 40
			hide    = `*`
			rsStr   = gstr.HideStr(str, percent, hide)
		)
		fmt.Println(rsStr)
	}

	{
		var (
			str     = `AAAAAAAAAA`
			percent = 60
			hide    = `@`
			rsStr   = gstr.HideStr(str, percent, hide)
		)
		fmt.Println(rsStr)
	}

	// Output:
	// 138****8000
	// AA@@@@@@AA
}

func ExampleImplode() {
	var (
		pieces = []string{"goframe", "is", "Very", "Nice", "to", "use"}
		glue   = "-"
		rsStr  = gstr.Implode(glue, pieces)
	)
	fmt.Println(rsStr)

	// Output:
	// goframe-is-Very-Nice-to-use
}

func ExampleInArray() {
	{
		var (
			a     = []string{"goframe", "is", "Very", "Nice", "to", "use"}
			s     = "goframe"
			rsStr = gstr.InArray(a, s)
		)
		fmt.Println(rsStr)
	}

	{
		var (
			a     = []string{"goframe", "is", "Very", "Nice", "to", "use"}
			s     = "golang"
			rsStr = gstr.InArray(a, s)
		)
		fmt.Println(rsStr)
	}

	// Output:
	// true
	// false
}

func ExampleIsLetterLower() {
	{
		var (
			b     byte = 65 // In ASCII is "a"
			rsStr      = gstr.IsLetterLower(b)
		)
		fmt.Println(rsStr)
	}

	{
		var (
			b     byte = 97 // In ASCII is "A"
			rsStr      = gstr.IsLetterLower(b)
		)
		fmt.Println(rsStr)
	}

	// Output:
	// false
	// true
}

func ExampleIsLetterUpper() {
	{
		var (
			b     byte = 65 // In ASCII is "a"
			rsStr      = gstr.IsLetterUpper(b)
		)
		fmt.Println(rsStr)
	}

	{
		var (
			b     byte = 97 // In ASCII is "A"
			rsStr      = gstr.IsLetterUpper(b)
		)
		fmt.Println(rsStr)
	}

	// Output:
	// true
	// false
}

func ExampleIsNumeric() {
	{
		var (
			s     = "88"
			rsStr = gstr.IsNumeric(s)
		)
		fmt.Println(rsStr)
	}

	{
		var (
			s     = "aa66bb88"
			rsStr = gstr.IsNumeric(s)
		)
		fmt.Println(rsStr)
	}

	{
		var (
			s     = "3.1415926"
			rsStr = gstr.IsNumeric(s)
		)
		fmt.Println(rsStr)
	}

	// Output:
	// true
	// false
	// true
}

func ExampleIsSubDomain() {
	{
		var (
			subDomain  = `s.goframe.org`
			mainDomain = `goframe.org`
			rsStr      = gstr.IsSubDomain(subDomain, mainDomain)
		)
		fmt.Println(rsStr)
	}

	{
		var (
			subDomain  = `s.s.goframe.org`
			mainDomain = `goframe.org`
			rsStr      = gstr.IsSubDomain(subDomain, mainDomain)
		)
		fmt.Println(rsStr)
	}

	{
		var (
			subDomain  = `s.s.goframe.org`
			mainDomain = `*goframe.org`
			rsStr      = gstr.IsSubDomain(subDomain, mainDomain)
		)
		fmt.Println(rsStr)
	}

	// Output:
	// true
	// true
	// false
}

func ExampleJoin() {
	var (
		array = []string{"goframe", "is", "Very", "Nice", "to", "use"}
		sep   = `,`
		rsStr = gstr.Join(array, sep)
	)
	fmt.Println(rsStr)

	// Output:
	// goframe,is,Very,Nice,to,use
}

func ExampleJoinAny() {
	{
		var (
			sep   = `@`
			arr1  = []string{"goframe", "is", "Very", "Nice", "to", "use"}
			rsStr = gstr.JoinAny(arr1, sep)
		)
		fmt.Println(rsStr)
	}

	{
		var (
			sep    = `,`
			arr2   = []int{99, 73, 85, 66}
			rsStr2 = gstr.JoinAny(arr2, sep)
		)
		fmt.Println(rsStr2)
	}

	{
		var (
			sep  = `,`
			arr3 = []interface{}{
				"Mary",
				18,
				99.5,
				"<br>",
				"Jack",
				19,
				66.5,
			}
			rsStr3 = gstr.JoinAny(arr3, sep)
		)
		fmt.Println(rsStr3)
	}

	{
		type StructA struct {
			Name string
			Age  int
		}
		var (
			sep = `|`

			arr4 = []StructA{
				{"Mary", 18},
				{"Jack", 18},
				{"Lucy", 18},
			}
			rsStr = gstr.JoinAny(arr4, sep)
		)
		fmt.Println(rsStr)
	}

	// Output:
	// goframe@is@Very@Nice@to@use
	// 99,73,85,66
	// Mary,18,99.5,<br>,Jack,19,66.5
	// {"Name":"Mary","Age":18}|{"Name":"Jack","Age":18}|{"Name":"Lucy","Age":18}
}

func ExampleLcFirst() {
	{
		var (
			str   = `Goframe`
			rsStr = gstr.LcFirst(str)
		)
		fmt.Println(rsStr)
	}

	{
		var (
			str   = `Goframe is Very Nice to use.`
			rsStr = gstr.LcFirst(str)
		)
		fmt.Println(rsStr)
	}

	// Output:
	// goframe
	// goframe is Very Nice to use.
}

func ExampleLenRune() {
	{
		var (
			str   = `goframe is Very Nice to use`
			rsStr = gstr.LenRune(str)
		)
		fmt.Println(rsStr)
	}

	{
		var (
			str   = `123 4567 890`
			rsStr = gstr.LenRune(str)
		)
		fmt.Println(rsStr)
	}

	{
		var (
			str   = `Goframe是一个非常好用的Go语言框架!`
			rsStr = gstr.LenRune(str)
		)
		fmt.Println(rsStr)
	}

	// Output:
	// 27
	// 12
	// 22
}

func ExampleLevenshtein() {
	{
		var (
			str1    = "Hello World"
			str2    = "ello World"
			costIns = 1
			costRep = 1
			costDel = 1
			rsStr   = gstr.Levenshtein(str1, str2, costIns, costRep, costDel)
		)

		fmt.Println(rsStr)
	}

	{
		var (
			str1    = "Hello World"
			str2    = "ello Worles"
			costIns = 10
			costRep = 20
			costDel = 30
			rsStr   = gstr.Levenshtein(str1, str2, costIns, costRep, costDel)
		)
		fmt.Println(rsStr)
	}

	// Output:
	// 1
	// 60
}

func ExampleNl2Br() {
	var (
		str = `goframe
is
Very
Nice
to
use`
		rsStr = gstr.Nl2Br(str)
	)

	fmt.Println(rsStr)

	// Output:
	// goframe<br>is<br>Very<br>Nice<br>to<br>use
}

func ExampleNumberFormat() {

	{
		var (
			number       float64 = 123456
			decimals             = 2
			decPoint             = "."
			thousandsSep         = ","
			rsStr                = gstr.NumberFormat(number, decimals, decPoint, thousandsSep)
		)
		fmt.Println(rsStr)
	}

	{
		var (
			number       = 1234.56
			decimals     = 1
			decPoint     = ","
			thousandsSep = " "
			rsStr        = gstr.NumberFormat(number, decimals, decPoint, thousandsSep)
		)
		fmt.Println(rsStr)
	}

	{
		var (
			number       = 1234.5678
			decimals     = 3
			decPoint     = "."
			thousandsSep = ","
			rsStr        = gstr.NumberFormat(number, decimals, decPoint, thousandsSep)
		)
		fmt.Println(rsStr)
	}

	// Output:
	// 123,456.00
	// 1 234,6
	// 1,234.568
}

func ExampleOctStr() {
	var (
		str   = `\346\200\241`
		rsStr = gstr.OctStr(str)
	)
	fmt.Println(rsStr)

	// Output:
	// 怡
}

// '103' is the 'g' in ASCII
func ExampleOrd() {
	var (
		str   = `goframe`
		rsStr = gstr.Ord(str)
	)

	fmt.Println(rsStr)

	// Output:
	// 103

}

func ExampleParse() {
	{
		var (
			str      = `v1=m&v2=n`
			rsStr, _ = gstr.Parse(str)
		)
		fmt.Println(rsStr)
	}

	{
		var (
			str      = `v[a][a]=m&v[a][b]=n`
			rsStr, _ = gstr.Parse(str)
		)
		fmt.Println(rsStr)
	}

	{
		// The form of nested Slice is not yet supported.
		var str = `v[][]=m&v[][]=n`
		rsStr, err := gstr.Parse(str)
		if err != nil {
			panic(err)
		}
		fmt.Println(rsStr)
	}

	{
		// This will produce an error.
		var str = `v=m&v[a]=n`
		rsStr, err := gstr.Parse(str)
		if err != nil {
			println(err)
		}
		fmt.Println(rsStr)
	}

	{
		var (
			str      = `a .[[b=c`
			rsStr, _ = gstr.Parse(str)
		)
		fmt.Println(rsStr)
	}

	// May Output:
	// map[v1:m v2:n]
	// map[v:map[a:map[a:m b:n]]]
	// map[v:map[]]
	// Error: expected type 'map[string]interface{}' for key 'v', but got 'string'
	// map[]
	// map[a___[b:c]
}

func ExamplePos() {
	{
		var (
			haystack = `goframe_is Very Nice.to-use`
			needle   = `Nice`
			rsStr    = gstr.Pos(haystack, needle)
		)
		fmt.Println(rsStr)
	}

	{
		var (
			haystack    = `goframe_is Very Nice.to-use`
			needle      = `Nice`
			startOffset = 16
			rsStr       = gstr.Pos(haystack, needle, startOffset)
		)
		fmt.Println(rsStr)
	}

	{
		var (
			haystack    = `goframe_is Very Nice.to-use`
			needle      = `Nice`
			startOffset = 17
			rsStr       = gstr.Pos(haystack, needle, startOffset)
		)
		fmt.Println(rsStr)
	}

	{
		var (
			haystack = `goframe_is Very Nice.to-use`
			needle   = `nice`
			rsStr    = gstr.Pos(haystack, needle)
		)
		fmt.Println(rsStr)
	}

	// Output:
	// 16
	// 16
	// -1
	// -1
}

func ExamplePosI() {
	{
		var (
			haystack = `goframe_is Very Nice.to-use`
			needle   = `Nice`
			rsStr    = gstr.PosI(haystack, needle)
		)

		fmt.Println(rsStr)
	}

	{
		var (
			haystack    = `goframe_is Very Nice.to-use`
			needle      = `Nice`
			startOffset = 16
			rsStr       = gstr.PosI(haystack, needle, startOffset)
		)
		fmt.Println(rsStr)
	}

	{
		var (
			haystack    = `goframe_is Very Nice.to-use`
			needle      = `Nice`
			startOffset = 17
			rsStr       = gstr.PosI(haystack, needle, startOffset)
		)
		fmt.Println(rsStr)
	}

	{
		var (
			haystack = `goframe_is Very Nice.to-use`
			needle   = `nice`
			rsStr    = gstr.PosI(haystack, needle)
		)
		fmt.Println(rsStr)
	}

	// Output:
	// 16
	// 16
	// -1
	// 16
}

func ExamplePosIRune() {
	{
		var (
			haystack = `goframe_is Very Nice.to-use`
			needle   = `Nice`
			rsStr    = gstr.PosIRune(haystack, needle)
		)

		fmt.Println(rsStr)
	}

	{
		var (
			haystack    = `Goframe是个非常好用的框架.`
			needle      = `Nice`
			startOffset = 16
			rsStr       = gstr.PosIRune(haystack, needle, startOffset)
		)
		fmt.Println(rsStr)
	}

	{
		var (
			haystack    = `Goframe是个非常好用的框架.`
			needle      = `Nice`
			startOffset = 17
			rsStr       = gstr.PosIRune(haystack, needle, startOffset)
		)
		fmt.Println(rsStr)
	}

	{
		var (
			haystack = `Goframe是个非常好用的框架.`
			needle   = `nice`
			rsStr    = gstr.PosIRune(haystack, needle)
		)
		fmt.Println(rsStr)
	}

	// Output:
	// 16
	// -1
	// -1
	// -1
}

func ExamplePosR() {
	{
		var (
			haystack = `goframe_is Very Nice.to-use`
			needle   = `Nice`
			rsStr    = gstr.PosR(haystack, needle)
		)
		fmt.Println(rsStr)
	}

	{
		var (
			haystack    = `goframe_is Very Nice.to-use`
			needle      = `Nice`
			startOffset = 16
			rsStr       = gstr.PosR(haystack, needle, startOffset)
		)
		fmt.Println(rsStr)
	}

	{
		var (
			haystack    = `goframe_is Very Nice.to-use`
			needle      = `Nice`
			startOffset = 17
			rsStr       = gstr.PosR(haystack, needle, startOffset)
		)
		fmt.Println(rsStr)
	}

	{
		var (
			haystack = `goframe_is Very Nice.to-use`
			needle   = `nice`
			rsStr    = gstr.PosR(haystack, needle)
		)
		fmt.Println(rsStr)
	}

	// Output:
	// 16
	// 16
	// -1
	// -1
}

func ExamplePosRI() {
	{
		var (
			haystack = `goframe_is Very Nice.to-use`
			needle   = `Nice`
			rsStr    = gstr.PosRI(haystack, needle)
		)
		fmt.Println(rsStr)
	}

	{
		var (
			haystack    = `goframe_is Very Nice.to-use`
			needle      = `Nice`
			startOffset = 16
			rsStr       = gstr.PosRI(haystack, needle, startOffset)
		)
		fmt.Println(rsStr)
	}

	{
		var (
			haystack    = `goframe_is Very Nice.to-use`
			needle      = `Nice`
			startOffset = 17
			rsStr       = gstr.PosRI(haystack, needle, startOffset)
		)
		fmt.Println(rsStr)
	}

	{
		var (
			haystack = `goframe_is Very Nice.to-use`
			needle   = `nice`
			rsStr    = gstr.PosRI(haystack, needle)
		)
		fmt.Println(rsStr)
	}

	// Output:
	// 16
	// 16
	// -1
	// 16

}

func ExamplePosRIRune() {
	{
		var (
			haystack = `Goframe是个非常好用的框架`
			needle   = `好用`
			rsStr    = gstr.PosRIRune(haystack, needle)
		)

		fmt.Println(rsStr)
	}

	{
		var (
			haystack    = `Goframe是个非常好用的框架`
			needle      = `框架`
			startOffset = 16
			rsStr       = gstr.PosRIRune(haystack, needle, startOffset)
		)
		fmt.Println(rsStr)
	}

	{
		var (
			haystack    = `Goframe是个非常好用的框架`
			needle      = `golang`
			startOffset = 17
			rsStr       = gstr.PosRIRune(haystack, needle, startOffset)
		)
		fmt.Println(rsStr)
	}

	{
		var (
			haystack = `Goframe是个非常好用的框架`
			needle   = `goframe`
			rsStr    = gstr.PosRIRune(haystack, needle)
		)
		fmt.Println(rsStr)
	}

	// Output:
	// 11
	// 14
	// -1
	// 0
}

func ExamplePosRRune() {
	{
		var (
			haystack = `goframe_is Very Nice.to-use`
			needle   = `Nice`
			rsStr    = gstr.PosRRune(haystack, needle)
		)
		fmt.Println(rsStr)
	}

	{
		var (
			haystack    = `goframe_is Very Nice.to-use`
			needle      = `Nice`
			startOffset = 16
			rsStr       = gstr.PosRRune(haystack, needle, startOffset)
		)
		fmt.Println(rsStr)
	}

	{
		var (
			haystack    = `Goframe是中国开发者的福利.`
			needle      = `Nice`
			startOffset = 17
			rsStr       = gstr.PosRRune(haystack, needle, startOffset)
		)
		fmt.Println(rsStr)
	}

	{
		var (
			haystack = `Goframe是中国开发者的福利.`
			needle   = `开发者`
			rsStr    = gstr.PosRRune(haystack, needle)
		)
		fmt.Println(rsStr)
	}

	// Output:
	// 16
	// 16
	// -1
	// 10
}

func ExamplePosRune() {
	{
		var (
			haystack = `goframe_is Very Nice.to-use`
			needle   = `Nice`
			rsStr    = gstr.PosRune(haystack, needle)
		)
		fmt.Println(rsStr)
	}

	{
		var (
			haystack    = `我喜欢Goframe框架`
			needle      = `框架`
			startOffset = 16
			rsStr       = gstr.PosRune(haystack, needle, startOffset)
		)
		fmt.Println(rsStr)
	}

	{
		var (
			haystack    = `我喜欢Goframe框架`
			needle      = `框架`
			startOffset = 17
			rsStr       = gstr.PosRune(haystack, needle, startOffset)
		)
		fmt.Println(rsStr)
	}

	{
		var (
			haystack = `我喜欢Goframe框架`
			needle   = `goframe`
			rsStr    = gstr.PosRune(haystack, needle)
		)
		fmt.Println(rsStr)
	}

	// Output:
	// 16
	// 10
	// -1
	// -1
}

func ExampleQuoteMeta() {
	{
		var (
			str   = `.\+?[^]()`
			rsStr = gstr.QuoteMeta(str)
		)
		fmt.Println(rsStr)
	}

	{
		var (
			str   = `https://goframe.org/pages/viewpage.action?pageId=1114327`
			rsStr = gstr.QuoteMeta(str)
		)
		fmt.Println(rsStr)
	}

	// Output:
	// \.\\\+\?\[\^\]\(\)
	// https://goframe\.org/pages/viewpage\.action\?pageId=1114327

}

func ExampleRepeat() {
	var (
		input      = `goframe `
		multiplier = 3
		rsStr      = gstr.Repeat(input, multiplier)
	)
	fmt.Println(rsStr)

	// Output:
	// goframe goframe goframe
}

func ExampleReplace() {
	{
		var (
			origin  = `goframe_is_Very_Nice_to_use!`
			search  = `_`
			replace = `+`
			rsStr   = gstr.Replace(origin, search, replace)
		)
		fmt.Println(rsStr)
	}

	{
		var (
			origin  = `goframe_is_Very_Nice_to_use!`
			search  = `_`
			replace = `+`
			count   = 2
			rsStr   = gstr.Replace(origin, search, replace, count)
		)
		fmt.Println(rsStr)
	}

	// Output:
	// goframe+is+Very+Nice+to+use!
	// goframe+is+Very_Nice_to_use!
}

func ExampleReplaceByArray() {
	{
		var (
			origin = `Golang is Very Good`
			array  = []string{"o", "O"}
			rsStr  = gstr.ReplaceByArray(origin, array)
		)
		fmt.Println(rsStr)
	}

	{
		var (
			origin = `Golang is Very Good`
			array  = []string{"Golang", "Goframe", "Good", "Nice"}
			rsStr  = gstr.ReplaceByArray(origin, array)
		)
		fmt.Println(rsStr)
	}

	// Output:
	// GOlang is Very GOOd
	// Goframe is Very Nice
}

func ExampleReplaceByMap() {
	{
		var (
			origin   = `Golang is Very Good`
			replaces = map[string]string{
				"Golang": "Goframe",
				"Good":   "Nice",
			}
			rsStr = gstr.ReplaceByMap(origin, replaces)
		)
		fmt.Println(rsStr)
	}

	{
		var (
			origin   = `Golang is Very Good`
			replaces = map[string]string{
				"golang": "Goframe",
				"good":   "Nice",
			}
			rsStr = gstr.ReplaceByMap(origin, replaces)
		)
		fmt.Println(rsStr)
	}

	// Output:
	// Goframe is Very Nice
	// Golang is Very Good
}

func ExampleReplaceI() {
	var (
		origin  = `goframe is Very Nice to use`
		search  = `Goframe`
		replace = `golang`
		count   = 3
		rsStr   = gstr.ReplaceI(origin, search, replace, count)
	)
	fmt.Println(rsStr)

	// Output:
	// golang is Very Nice to use
}

func ExampleReplaceIByArray() {
	var (
		origin = `golang is very nice`
		array  = []string{"Golang", "GoFrame", "Nice", "GOOD"}
		rsStr  = gstr.ReplaceIByArray(origin, array)
	)

	fmt.Println(rsStr)

	// Output:
	// GoFrame is very GOOD
}

func ExampleReplaceIByMap() {
	{
		var (
			origin   = `Golang is Very Good`
			replaces = map[string]string{
				"Golang": "Goframe",
				"Good":   "Nice",
			}
			rsStr = gstr.ReplaceIByMap(origin, replaces)
		)
		fmt.Println(rsStr)
	}

	{
		var (
			origin   = `Golang is Very Good`
			replaces = map[string]string{
				"golang": "Goframe",
				"good":   "Nice",
			}
			rsStr = gstr.ReplaceIByMap(origin, replaces)
		)
		fmt.Println(rsStr)
	}

	// Output:
	// Goframe is Very Nice
	// Goframe is Very Nice
}

func ExampleReverse() {
	var (
		str   = `123456`
		rsStr = gstr.Reverse(str)
	)
	fmt.Println(rsStr)

	// Output:
	// 654321
}

func ExampleSearchArray() {
	{
		var (
			a     = []string{"goframe", "is", "Very", "Nice"}
			s     = `goframe`
			rsStr = gstr.SearchArray(a, s)
		)
		fmt.Println(rsStr)
	}

	{
		var (
			a     = []string{"goframe", "is", "Very", "Nice"}
			s     = `Very`
			rsStr = gstr.SearchArray(a, s)
		)
		fmt.Println(rsStr)
	}

	{
		var (
			a     = []string{"goframe", "is", "Very", "Nice"}
			s     = `use`
			rsStr = gstr.SearchArray(a, s)
		)
		fmt.Println(rsStr)
	}

	// Output:
	// 0
	// 2
	// -1
}

func ExampleShuffle() {
	var (
		str   = `abcdefg`
		rsStr = gstr.Shuffle(str)
	)
	fmt.Println(rsStr)

	// May Output:
	// efcgbda
}

func ExampleSimilarText() {
	{
		var (
			percent *float64
			first   = `goframe_is`
			second  = `go-nice`
			rsStr   = gstr.SimilarText(first, second, percent)
		)
		fmt.Println(rsStr)
	}

	{
		var (
			first  = `AaBbCcDd`
			second = `ad`
		)
		p := 0.80
		rsStr := gstr.SimilarText(first, second, &p)
		fmt.Println(rsStr)
	}

	// Output:
	// 3
	// 2
}

func ExampleSoundex() {
	{
		var (
			str    = `Euler`
			str2   = `Ellery`
			rsStr  = gstr.Soundex(str)
			rsStr2 = gstr.Soundex(str2)
		)
		fmt.Println(rsStr, rsStr2)
	}

	{
		var (
			str    = `Gauss`
			str2   = `Ghosh`
			rsStr  = gstr.Soundex(str)
			rsStr2 = gstr.Soundex(str2)
		)
		fmt.Println(rsStr, rsStr2)
	}

	{
		var (
			str    = `Lloyd`
			str2   = `Ladd`
			rsStr  = gstr.Soundex(str)
			rsStr2 = gstr.Soundex(str2)
		)
		fmt.Println(rsStr, rsStr2)
	}

	// Output:
	// E406 E406
	// G020 G020
	// L030 L030
}

func ExampleSplit() {
	var (
		str       = `GoFrame_Nice`
		delimiter = `_`
		rsStr     = gstr.Split(str, delimiter)
	)
	fmt.Printf(`%#v`, rsStr)

	// Output:
	// []string{"GoFrame", "Nice"}
}

func ExampleSplitAndTrim() {
	var (
		str           = `    goframe, is,       Very ,Nice ,to,      use`
		delimiter     = `,`
		characterMask = ` `
		rsStr         = gstr.SplitAndTrim(str, delimiter, characterMask)
	)
	fmt.Printf("%#v", rsStr)

	// Output:
	// []string{"goframe", "is", "Very", "Nice", "to", "use"}
}

func ExampleStr() {
	var (
		haystack = `xxx.jpg`
		needle   = `.`
		rsStr    = gstr.Str(haystack, needle)
	)
	fmt.Println(rsStr)

	// Output:
	// .jpg
}

func ExampleStrEx() {
	var (
		haystack = `https://goframe.org/index.html?a=1&b=2`
		needle   = `?`
		rsStr    = gstr.StrEx(haystack, needle)
	)
	fmt.Println(rsStr)

	// Output:
	// a=1&b=2
}

func ExampleStrLimit() {
	{
		var (
			str    = `1234567890.png`
			length = 4
			suffix = `.jpg`
			rsStr  = gstr.StrLimit(str, length, suffix)
		)
		fmt.Println(rsStr)
	}

	{
		var (
			str    = `123456789`
			length = 3
			suffix = `...`
			rsStr  = gstr.StrLimit(str, length, suffix)
		)
		fmt.Println(rsStr)
	}

	// Output:
	// 1234.jpg
	// 123...
}

func ExampleStrLimitRune() {
	{
		var (
			str    = `123456789`
			length = 5
			suffix = "ABCD"
			rsStr  = gstr.StrLimitRune(str, length, suffix)
		)
		fmt.Println(rsStr)
	}

	{
		var (
			str    = `goframe是一个好用的框架她提供了非常丰富的工具给使用者.`
			length = 15
			suffix = "!"
			rsStr  = gstr.StrLimitRune(str, length, suffix)
		)
		fmt.Println(rsStr)
	}

	// Output:
	// 12345ABCD
	// goframe是一个好用的框架!
}

func ExampleStrTill() {
	{
		var (
			haystack = `goframe是一个好用的框架,她提供了非常丰富的工具给使用者.`
			needle   = `,`
			rsStr    = gstr.StrTill(haystack, needle)
		)
		fmt.Println(rsStr)
	}

	{
		var (
			haystack = `https://goframe.org/index.html?test=123456`
			needle   = `?`
			rsStr    = gstr.StrTill(haystack, needle)
		)
		fmt.Println(rsStr)
	}

	// Output:
	// goframe是一个好用的框架,
	// https://goframe.org/index.html?
}

func ExampleStrTillEx() {
	{
		var (
			haystack = `xxxx.txt`
			needle   = `.txt`
			rsStr    = gstr.StrTillEx(haystack, needle)
		)
		fmt.Println(rsStr)
	}

	{
		var (
			haystack = `文件名.zip`
			needle   = `.zip`
			rsStr    = gstr.StrTillEx(haystack, needle)
		)
		fmt.Println(rsStr)
	}

	// Output:
	// xxxx
	// 文件名
}

func ExampleStripSlashes() {
	var (
		str   = `C:\\windows\\GoFrame\\test`
		rsStr = gstr.StripSlashes(str)
	)
	fmt.Println(rsStr)

	// Output:
	// C:\windows\GoFrame\test
}

func ExampleSubStr() {
	{
		var (
			str    = `123456789`
			start  = 0
			length = 2
			subStr = gstr.SubStr(str, start, length)
		)
		fmt.Println(subStr)
	}

	{
		var (
			str    = `123456789`
			start  = 5
			length = 3
			subStr = gstr.SubStr(str, start, length)
		)
		fmt.Println(subStr)
	}

	// Output:
	// 12
	// 678
}

func ExampleSubStrRune() {
	{
		var (
			str    = `123456789`
			start  = 5
			length = 4
			subStr = gstr.SubStrRune(str, start, length)
		)
		fmt.Println(subStr)
	}
	{
		var (
			str    = `一二三四五六七八九零`
			start  = 5
			length = 4
			subStr = gstr.SubStrRune(str, start, length)
		)
		fmt.Println(subStr)
	}
	// Output:
	// 6789
	// 六七八九
}

func ExampleToLower() {
	var (
		s     = `GOFRAME IS VERY NICE TO USE`
		rsStr = gstr.ToLower(s)
	)
	fmt.Println(rsStr)

	// Output:
	// goframe is very nice to use
}

func ExampleToUpper() {
	var (
		s     = `goframe is very nice to use`
		rsStr = gstr.ToUpper(s)
	)
	fmt.Println(rsStr)

	// Output:
	// GOFRAME IS VERY NICE TO USE
}

func ExampleTrim() {
	var (
		str           = `.abc.def..`
		characterMask = "."
		rsStr         = gstr.Trim(str, characterMask)
	)
	fmt.Println(rsStr)

	// Output:
	// abc.def
}

func ExampleTrimAll() {
	var (
		str           = `.abc.def...`
		characterMask = "."
		rsStr         = gstr.TrimAll(str, characterMask)
	)
	fmt.Println(rsStr)

	// Output:
	// abcdef
}

func ExampleTrimLeft() {
	var (
		str           = `..abc.def.. `
		characterMask = "."
		rsStr         = gstr.TrimLeft(str, characterMask)
	)
	fmt.Println(rsStr)

	// Output:
	// abc.def..
}

func ExampleTrimLeftStr() {
	var (
		str   = `...abcd..efg...`
		cut   = "."
		count = 2
		rsStr = gstr.TrimLeftStr(str, cut, count)
	)
	fmt.Println(rsStr)

	// Output:
	// .abcd..efg...
}

func ExampleTrimRight() {
	var (
		str           = `abcdef    `
		characterMask = "def" // []byte{"d", "e", "f"}
		rsStr         = gstr.TrimRight(str, characterMask)
	)
	fmt.Println(rsStr)

	// Output:
	// abc
}

func ExampleTrimRightStr() {
	{
		var (
			str   = `aaa,bbb,ccc,ddd,eee,`
			cut   = ","
			count = 1
			rsStr = gstr.TrimRightStr(str, cut, count)
		)
		fmt.Println(rsStr)
	}

	{
		var (
			str   = `aaa,bbb,ccc,ddd,eee,,,,,`
			cut   = ","
			count = -1
			rsStr = gstr.TrimRightStr(str, cut, count)
		)
		fmt.Println(rsStr)
	}

	{
		var (
			str   = `aaa,bbb,ccc,ddd,eee,,,,,   `
			cut   = ","
			count = -1
			rsStr = gstr.TrimRightStr(str, cut, count)
		)
		fmt.Println(rsStr)
	}

	// Output:
	// aaa,bbb,ccc,ddd,eee
	// aaa,bbb,ccc,ddd,eee
	// aaa,bbb,ccc,ddd,eee,,,,,
}

func ExampleTrimStr() {
	{
		var (
			str   = `goframe is Very Nice to-use`
			cut   = "goframe"
			rsStr = gstr.TrimStr(str, cut)
		)
		fmt.Println(rsStr)
	}
	{
		var (
			str   = `goframe is Very Nice to-use`
			cut   = "use"
			count = -1
			rsStr = gstr.TrimStr(str, cut, count)
		)
		fmt.Println(rsStr)
	}

	// Output:
	// is Very Nice to-use
	// goframe is Very Nice to-
}

func ExampleUcFirst() {
	var (
		s     = `goframe is very nice to use.`
		rsStr = gstr.UcFirst(s)
	)
	fmt.Println(rsStr)

	// Output:
	// Goframe is very nice to use.
}

func ExampleUcWords() {
	var (
		str   = `goframe is very nice to use.`
		rsStr = gstr.UcWords(str)
	)
	fmt.Println(rsStr)

	// Output:
	// Goframe Is Very Nice To Use.
}

func ExampleWordWrap() {
	{
		var (
			str   = `A very long woooooooooooooooooord. and something`
			width = 8
			br    = "\n"
			rsStr = gstr.WordWrap(str, width, br)
		)
		fmt.Println(rsStr)
	}

	{
		var (
			str   = `The quick brown fox jumped over the lazy dog.`
			width = 20
			br    = "<br />\n"
			rsStr = gstr.WordWrap(str, width, br)
		)
		fmt.Printf("%v", rsStr)
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
