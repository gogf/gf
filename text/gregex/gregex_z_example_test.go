// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
package gregex_test

import (
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/text/gregex"
)

func ExampleIsMatch() {
	var str = "hello 94 easy gf!"
	patternStr := `[1-9]\d*`
	fmt.Println(gregex.IsMatch(patternStr, []byte(str)))

	// output
	// true
}

func ExampleIsMatchString() {
	var str = "hello 94 easy gf!"
	patternStr := `[1-9]\d*`
	fmt.Println(gregex.IsMatchString(patternStr, str))

	// output
	// true
}

func ExampleMatch() {
	var str = "hello 94 easy gf!"
	patternStr := `[1-9]\d*`
	result, err := gregex.Match(patternStr, []byte(str))
	if err != nil {
		fmt.Println(err.Error())
	}
	g.Dump(result)
	g.Dump(result[0])

	// output
	// ["OTQ="]
	// 94
}

func ExampleMatchAll() {
	var str = "hello 94 easy gf!"
	patternStr := `[1-9]\d*`
	results, err := gregex.MatchAll(patternStr, []byte(str))
	if err != nil {
		fmt.Println(err.Error())
	}
	g.Dump(results)

	// output
	// [["OTQ="]]
	//
}

func ExampleMatchAllString() {
	var str = "hello 94 98 easy gf!"
	patternStr := `[1-9]\d*`
	results, err := gregex.MatchAllString(patternStr, str)
	if err != nil {
		fmt.Println(err.Error())
	}

	g.Dump(results)

	// output
	// [["94"],["98"]]
}

func ExampleMatchString() {
	var str = "hello 94 98 easy gf!"
	patternStr := `[1-9]\d*`

	// if you need a greed match, should use <..all> methods
	results, err := gregex.MatchString(patternStr, str)
	if err != nil {
		fmt.Println(err.Error())
	}

	g.Dump(results)

	// output
	// ["94"]
}

func ExampleQuote() {
	patternStr := `[1-9]\d*`
	result := gregex.Quote(patternStr)
	fmt.Println(result)

	// output
	// \[1-9\]\\d\*
}

func ExampleReplace() {
	patternStr := `[1-9]\d*`
	str := "hello gf 2020!"
	repStr := "2021"
	result, err := gregex.Replace(patternStr, []byte(repStr), []byte(str))

	g.Dump(err)
	g.Dump(result)

	// output
	// null
	// hello gf 2021!
}

func ExampleReplaceFunc() {
	patternStr := `[1-9]\d*`
	str := "hello gf 2018~2020!"

	// When the regular statement can match multiple results
	// func can be used to further control the value that needs to be modified
	result, err := gregex.ReplaceFunc(patternStr, []byte(str), func(b []byte) []byte {

		g.Dump(b)

		replaceStr := "2021"
		if string(b) == "2020" {
			return []byte(replaceStr)
		}
		return b
	})

	g.Dump(result)
	g.Dump(err)

	// output
	// 2018
	// 2020
	// hello gf 2018~2021!
	// null
}

func ExampleReplaceFuncMatch() {
	patternStr := `[1-9]\d*`
	str := "hello gf 2018~2020!"

	// In contrast to [ExampleReplaceFunc]
	// the result contains the `pattern' of all subpatterns that use the matching function
	result, err := gregex.ReplaceFuncMatch(patternStr, []byte(str), func(match [][]byte) []byte {

		g.Dump(match)

		replaceStr := "2021"
		for _, v := range match {
			if string(v) == "2020" {
				return []byte(replaceStr)
			}
		}
		return match[0]
	})

	g.Dump(result)
	g.Dump(err)

	// output
	// [
	//	"MjAxOA=="
	// ]
	//
	// [
	//	"MjAyMA=="
	// ]
	//
	// hello gf 2018~2021!
	// null
}

func ExampleReplaceString() {
	patternStr := `[1-9]\d*`
	str := "hello gf 2020!"
	replaceStr := "2021"
	result, err := gregex.ReplaceString(patternStr, replaceStr, str)

	g.Dump(result)
	g.Dump(err)

	// output
	// hello gf 2021!
	// null
}

func ExampleReplaceStringFunc() {
	patternStr := `[1-9]\d*`
	str := "hello gf 2018~2020!"

	// When the regular statement can match multiple results
	// func can be used to further control the value that needs to be modified
	result, err := gregex.ReplaceStringFunc(patternStr, str, func(b string) string {

		g.Dump(b)

		replaceStr := "2021"
		if b == "2020" {
			return replaceStr
		}
		return b
	})

	g.Dump(result)
	g.Dump(err)

	// output
	// 2018
	// 2020
	// hello gf 2018~2021!
	// null
}

func ExampleReplaceStringFuncMatch() {
	patternStr := `[1-9]\d*`
	str := "hello gf 2018~2020!"

	// When the regular statement can match multiple results
	// func can be used to further control the value that needs to be modified
	result, err := gregex.ReplaceStringFuncMatch(patternStr, str, func(b []string) string {

		g.Dump(b)

		replaceStr := "2021"
		for _, v := range b {
			if v == "2020" {
				return replaceStr
			}
		}
		return b[0]
	})

	g.Dump(result)
	g.Dump(err)

	// output
	// ["2018"]
	// ["2020"]
	// hello gf 2018~2021!
	// null
}

func ExampleSplit() {
	patternStr := `[1-9]\d*`
	str := "hello2020gf"
	result := gregex.Split(patternStr, str)

	g.Dump(result)

	// output
	// ["hello","gf"]
}

func ExampleValidate() {
	// Valid match statement
	g.Dump(gregex.Validate(`[1-9]\d*`))
	// Mismatched statement
	g.Dump(gregex.Validate(`[a-9]\d*`))

	// output
	// null
	// {
	//	"Code": "invalid character class range",
	//	"Expr": "a-9"
	// }
}
