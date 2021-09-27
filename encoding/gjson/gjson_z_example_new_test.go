// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gjson_test

import (
	"fmt"
	"github.com/gogf/gf/encoding/gjson"
)

func Example_newFromJson() {
	jsonContent := `{"name":"john", "score":"100"}`
	j := gjson.New(jsonContent)
	fmt.Println(j.Get("name"))
	fmt.Println(j.Get("score"))
	// Output:
	// john
	// 100
}

func Example_newFromXml() {
	jsonContent := `<?xml version="1.0" encoding="UTF-8"?><doc><name>john</name><score>100</score></doc>`
	j := gjson.New(jsonContent)
	// Note that there's root node in the XML content.
	fmt.Println(j.Get("doc.name"))
	fmt.Println(j.Get("doc.score"))
	// Output:
	// john
	// 100
}

func Example_newFromStruct() {
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

func Example_newFromStructWithTag() {
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
	// The parameter `tags` specifies custom priority tags for struct conversion to map,
	// multiple tags joined with char ','.
	j := gjson.NewWithTag(me, "tag")
	fmt.Println(j.Get("name"))
	fmt.Println(j.Get("score"))
	fmt.Println(j.Get("Title"))
	// Output:
	// john
	// 100
	// engineer
}
