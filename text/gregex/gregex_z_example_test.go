// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
package gregex_test

import (
	"fmt"
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
	for _, v := range result {
		fmt.Println(string(v))
	}

	// output
	// 94
}

func ExampleMatchAll() {
	var str = "hello 94 easy gf!"
	patternStr := `[1-9]\d*`
	results, err := gregex.MatchAll(patternStr, []byte(str))
	if err != nil {
		fmt.Println(err.Error())
	}
	for _, result := range results {
		for _, v := range result {
			fmt.Println(string(v))
		}
	}

	// output
	// 94
}

func ExampleMatchAllString() {
	var str = "hello 94 easy gf!"
	patternStr := `[1-9]\d*`
	results, err := gregex.MatchAllString(patternStr, str)
	if err != nil {
		fmt.Println(err.Error())
	}
	for _, result := range results {
		for _, v := range result {
			fmt.Println(v)
		}
	}

	// output
	// 94
}

func ExampleMatchString() {
	var str = "hello 94 easy gf!"
	patternStr := `[1-9]\d*`
	results, err := gregex.MatchString(patternStr, str)
	if err != nil {
		fmt.Println(err.Error())
	}
	for _, result := range results {
		fmt.Println(result)
	}

	// output
	// 94
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
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(string(result))

	// output
	// hello gf 2021!
}

func ExampleReplaceFunc() {
	patternStr := `[1-9]\d*`
	str := "hello gf 2020!"
	result, err := gregex.ReplaceFunc(patternStr, []byte(str), func(b []byte) []byte {
		replaceStr := "2021"
		if string(b) == "2020" {
			return []byte(replaceStr)
		}
		return nil
	})
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(string(result))

	// output
	// hello gf 2021!
}

func ExampleReplaceFuncMatch() {
	patternStr := `[1-9]\d*`
	str := "hello gf 2020!"
	result, err := gregex.ReplaceFuncMatch(patternStr, []byte(str), func(match [][]byte) []byte {
		replaceStr := "2021"
		for _, v := range match {
			if string(v) == "2020" {
				return []byte(replaceStr)
			}
		}
		return nil
	})
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(string(result))

	// output
	// hello gf 2021!
}

func ExampleReplaceString() {
	patternStr := `[1-9]\d*`
	str := "hello gf 2020!"
	replaceStr := "2021"
	result, err := gregex.ReplaceString(patternStr, replaceStr, str)
	if err != nil {
		fmt.Println(result)
	}
	fmt.Println(result)

	// output
	// hello gf 2021!
}

func ExampleReplaceStringFunc() {
	patternStr := `[1-9]\d*`
	str := "hello gf 2020!"
	result, err := gregex.ReplaceStringFunc(patternStr, str, func(s string) string {
		replaceStr := "2021"
		if s == "2020" {
			return replaceStr
		}
		return ""
	})
	if err != nil {
		fmt.Println(result)
	}
	fmt.Println(result)

	// output
	// hello gf 2021!
}

func ExampleReplaceStringFuncMatch() {
	patternStr := `[1-9]\d*`
	str := "hello gf 2020!"
	result, err := gregex.ReplaceStringFuncMatch(patternStr, str, func(match []string) string {
		replaceStr := "2021"
		for _, v := range match {
			if v == "2020" {
				return replaceStr
			}
		}
		return ""
	})
	if err != nil {
		fmt.Println(result)
	}
	fmt.Println(result)

	// output
	// hello gf 2021!
}

func ExampleSplit() {
	patternStr := `[1-9]\d*`
	str := "hello2020gf"
	result := gregex.Split(patternStr, str)
	for _, v := range result {
		fmt.Println(v)
	}

	// output
	// hello
	// gf
}

func ExampleValidate() {
	patternStr := `[1-9]\d*`
	err := gregex.Validate(patternStr)
	if err != nil {
		fmt.Println(err)
	}
	if err == nil {
		fmt.Println("ok")
	}

	// output
	// ok
}
