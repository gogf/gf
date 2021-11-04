// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
package gregex_test

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/text/gregex"
)

func ExampleIsMatch() {
	patternStr := `[1-9]\d+`
	g.Dump(gregex.IsMatch(patternStr, []byte("hello 2022! hello gf!")))
	g.Dump(gregex.IsMatch(patternStr, nil))
	g.Dump(gregex.IsMatch(patternStr, []byte("hello gf!")))

	// Output
	// true
	// false
	// false
}

func ExampleIsMatchString() {
	patternStr := `[1-9]\d+`
	g.Dump(gregex.IsMatchString(patternStr, "hello 2022! hello gf!"))
	g.Dump(gregex.IsMatchString(patternStr, "hello gf!"))
	g.Dump(gregex.IsMatchString(patternStr, ""))

	// Output
	// true
	// false
	// false
}

func ExampleMatch() {
	patternStr := `[1-9]\d+`
	result, err := gregex.Match(patternStr, []byte("hello 2022! hello gf!"))
	g.Dump(result)
	g.Dump(err)

	result, err = gregex.Match(patternStr, nil)
	g.Dump(result)
	g.Dump(err)

	result, err = gregex.Match(patternStr, []byte("hello easy gf!"))
	g.Dump(result)
	g.Dump(err)

	// Output
	// ["2022"]
	// <nil>
	// []
	// <nil>
	// []
	// <nil>
}

func ExampleMatchAll() {
	patternStr := `[1-9]\d+`
	results, err := gregex.MatchAll(patternStr, []byte("goodBye 2021! hello 2022! hello gf!"))
	g.Dump(results)
	g.Dump(err)

	results, err = gregex.MatchAll(patternStr, []byte("hello gf!"))
	g.Dump(results)
	g.Dump(err)

	results, err = gregex.MatchAll(patternStr, nil)
	g.Dump(results)
	g.Dump(err)

	// Output
	// [["2021"],["2022"]]
	// <nil>
	// []
	// <nil>
	// []
	// <nil>
}

func ExampleMatchAllString() {
	patternStr := `[1-9]\d+`
	results, err := gregex.MatchAllString(patternStr, "goodBye 2021! hello 2022! hello gf!")
	g.Dump(results)
	g.Dump(err)

	results, err = gregex.MatchAllString(patternStr, "hello gf!")
	g.Dump(results)
	g.Dump(err)

	results, err = gregex.MatchAllString(patternStr, "")
	g.Dump(results)
	g.Dump(err)

	// Output
	// [["2021"],["2022"]]
	// <nil>
	// []
	// <nil>
	// []
	// <nil>
}

func ExampleMatchString() {
	var str = "goodBye 2021! hello 2022! hello gf!"
	patternStr := `[1-9]\d+`
	// if you need a greed match, should use <..all> methods
	results, err := gregex.MatchString(patternStr, str)

	g.Dump(results)
	g.Dump(err)

	// Output
	// ["2021"]
	// <nil>
}

func ExampleQuote() {
	patternStr := `[1-9]\d+`
	result := gregex.Quote(patternStr)

	g.Dump(result)

	// Output
	// "\[1-9\]\\d\+"
}

func ExampleReplace() {
	var (
		patternStr  = `[1-9]\d+`
		str         = "hello gf 2020!"
		repStr      = "2021"
		result, err = gregex.Replace(patternStr, []byte(repStr), []byte(str))
	)

	g.Dump(err)
	g.Dump(result)

	// Output
	// <nil>
	// "hello gf 2021!"
}

func ExampleReplaceFunc() {
	var (
		patternStr    = `[1-9]\d+`
		str           = "hello gf 2018~2020!"
		replaceStrMap = map[string][]byte{
			"2020": []byte("2021"),
		}
	)

	// When the regular statement can match multiple results
	// func can be used to further control the value that needs to be modified
	result, err := gregex.ReplaceFunc(patternStr, []byte(str), func(b []byte) []byte {
		g.Dump(b)
		if replaceStr, ok := replaceStrMap[string(b)]; ok {
			return replaceStr
		}
		return b
	})

	g.Dump(result)
	g.Dump(err)

	// Output
	// "2018"
	// "2020"
	// "hello gf 2018~2021!"
	// <nil>
}

func ExampleReplaceFuncMatch() {
	var (
		patternStr = `[1-9]\d+`
		str        = "hello gf 2018~2020!"
		replaceMap = map[string][]byte{
			"2020": []byte("2021"),
		}
	)
	// In contrast to [ExampleReplaceFunc]
	// the result contains the `pattern' of all subpatterns that use the matching function
	result, err := gregex.ReplaceFuncMatch(patternStr, []byte(str), func(match [][]byte) []byte {
		g.Dump(match)
		for _, v := range match {
			replaceStr, ok := replaceMap[string(v)]
			if ok {
				return replaceStr
			}
		}
		return match[0]
	})
	g.Dump(result)
	g.Dump(err)

	// Output
	// ["2018"]
	// ["2020"]
	// "hello gf 2018~2021!"
	// <nil>
}

func ExampleReplaceString() {
	patternStr := `[1-9]\d+`
	str := "hello gf 2020!"
	replaceStr := "2021"
	result, err := gregex.ReplaceString(patternStr, replaceStr, str)

	g.Dump(result)
	g.Dump(err)

	// Output
	// "hello gf 2021!"
	// <nil>
}

func ExampleReplaceStringFunc() {
	var (
		patternStr    = `[1-9]\d+`
		str           = "hello gf 2018~2020!"
		replaceStrMap = map[string]string{
			"2020": "2021",
		}
	)
	// When the regular statement can match multiple results
	// func can be used to further control the value that needs to be modified
	result, err := gregex.ReplaceStringFunc(patternStr, str, func(b string) string {
		g.Dump(b)
		if replaceStr, ok := replaceStrMap[b]; ok {
			return replaceStr
		}
		return b
	})
	g.Dump(result)
	g.Dump(err)

	// Output
	// "2018"
	// "2020"
	// "hello gf 2018~2021!"
	// <nil>
}

func ExampleReplaceStringFuncMatch() {
	var (
		patternStr    = `[1-9]\d+`
		str           = "hello gf 2018~2020!"
		replaceStrMap = map[string]string{
			"2020": "2021",
		}
	)
	// When the regular statement can match multiple results
	// func can be used to further control the value that needs to be modified
	result, err := gregex.ReplaceStringFuncMatch(patternStr, str, func(b []string) string {
		g.Dump(b)
		for _, v := range b {
			if replaceStr, ok := replaceStrMap[v]; ok {
				return replaceStr
			}
		}
		return b[0]
	})
	g.Dump(result)
	g.Dump(err)

	// Output
	// ["2018"]
	// ["2020"]
	// "hello gf 2018~2021!"
	// <nil>
}

func ExampleSplit() {
	patternStr := `[1-9]\d+`
	str := "hello2020gf"
	result := gregex.Split(patternStr, str)
	g.Dump(result)

	// Output
	// ["hello","gf"]
}

func ExampleValidate() {
	// Valid match statement
	g.Dump(gregex.Validate(`[1-9]\d+`))
	// Mismatched statement
	g.Dump(gregex.Validate(`[a-9]\d+`))

	// Output
	// <nil>
	// {
	//	"Code": "invalid character class range",
	//	"Expr": "a-9"
	// }
}
