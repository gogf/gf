// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gjson_test

import (
	"fmt"

	"github.com/gogf/gf/v2/encoding/gjson"
)

func ExampleNew() {
	jsonContent := `{"name":"john", "score":"100"}`
	j := gjson.New(jsonContent)
	fmt.Println(j.Get("name"))
	fmt.Println(j.Get("score"))

	// Output:
	// john
	// 100
}

func ExampleNewWithTag() {
	type Me struct {
		Name  string `tag:"name"`
		Score int    `tag:"score"`
		Title string
	}
	me := Me{
		Name:  "john",
		Score: 100,
		Title: "engineer",
	}
	j := gjson.NewWithTag(me, "tag", true)
	fmt.Println(j.Get("name"))
	fmt.Println(j.Get("score"))
	fmt.Println(j.Get("Title"))

	// Output:
	// john
	// 100
	// engineer
}

func ExampleNewWithOptions() {
	type Me struct {
		Name  string `tag:"name"`
		Score int    `tag:"score"`
		Title string
	}
	me := Me{
		Name:  "john",
		Score: 100,
		Title: "engineer",
	}

	j := gjson.NewWithOptions(me, gjson.Options{
		Tags: "tag",
	})
	fmt.Println(j.Get("name"))
	fmt.Println(j.Get("score"))
	fmt.Println(j.Get("Title"))

	// Output:
	// john
	// 100
	// engineer
}

func ExampleNewWithOptions_UTF8BOM() {
	jsonContent := `{"name":"john", "score":"100"}`

	content := make([]byte, 3, len(jsonContent)+3)
	content[0] = 0xEF
	content[1] = 0xBB
	content[2] = 0xBF
	content = append(content, jsonContent...)

	j := gjson.NewWithOptions(content, gjson.Options{
		Tags: "tag",
	})
	fmt.Println(j.Get("name"))
	fmt.Println(j.Get("score"))

	// Output:
	// john
	// 100
}

func ExampleNew_Xml() {
	jsonContent := `<?xml version="1.0" encoding="UTF-8"?><doc><name>john</name><score>100</score></doc>`
	j := gjson.New(jsonContent)
	// Note that there's root node in the XML content.
	fmt.Println(j.Get("doc.name"))
	fmt.Println(j.Get("doc.score"))
	// Output:
	// john
	// 100
}

func ExampleNew_Struct() {
	type Me struct {
		Name  string `json:"name"`
		Score int    `json:"score"`
	}
	me := Me{
		Name:  "john",
		Score: 100,
	}
	j := gjson.New(me)
	fmt.Println(j.Get("name"))
	fmt.Println(j.Get("score"))
	// Output:
	// john
	// 100
}
