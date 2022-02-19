// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gjson_test

import (
	"fmt"

	"github.com/gogf/gf/v2/debug/gdebug"
	"github.com/gogf/gf/v2/encoding/gjson"
)

func ExampleLoad() {
	jsonFilePath := gdebug.TestDataPath("json", "data1.json")
	j, _ := gjson.Load(jsonFilePath)
	fmt.Println(j.Get("name"))
	fmt.Println(j.Get("score"))

	// Output:
	// john
	// 100
}

func ExampleLoadJson() {
	jsonContent := `{"name":"john", "score":"100"}`
	j, _ := gjson.LoadJson(jsonContent)
	fmt.Println(j.Get("name"))
	fmt.Println(j.Get("score"))

	// Output:
	// john
	// 100
}

func ExampleLoadXml() {
	xmlContent := `<?xml version="1.0" encoding="UTF-8"?>
	<base>
		<name>john</name>
		<score>100</score>
	</base>`
	j, _ := gjson.LoadXml(xmlContent)
	fmt.Println(j.Get("base.name"))
	fmt.Println(j.Get("base.score"))

	// Output:
	// john
	// 100
}

func ExampleLoadIni() {
	iniContent := `
	[base]
	name = john
	score = 100
	`
	j, _ := gjson.LoadIni(iniContent)
	fmt.Println(j.Get("base.name"))
	fmt.Println(j.Get("base.score"))

	// Output:
	// john
	// 100
}

func ExampleLoadYaml() {
	yamlContent :=
		`base:
  name: john
  score: 100`

	j, _ := gjson.LoadYaml(yamlContent)
	fmt.Println(j.Get("base.name"))
	fmt.Println(j.Get("base.score"))

	// Output:
	// john
	// 100
}

func ExampleLoadToml() {
	tomlContent :=
		`[base]
  name = "john"
  score = 100`

	j, _ := gjson.LoadToml(tomlContent)
	fmt.Println(j.Get("base.name"))
	fmt.Println(j.Get("base.score"))

	// Output:
	// john
	// 100
}

func ExampleLoadContent() {
	jsonContent := `{"name":"john", "score":"100"}`
	xmlContent := `<?xml version="1.0" encoding="UTF-8"?>
	<base>
		<name>john</name>
		<score>100</score>
	</base>`

	j, _ := gjson.LoadContent(jsonContent)
	x, _ := gjson.LoadContent(xmlContent)

	fmt.Println(j.Get("name"))
	fmt.Println(j.Get("score"))
	fmt.Println(x.Get("base.name"))
	fmt.Println(x.Get("base.score"))

	// Output:
	// john
	// 100
	// john
	// 100
}

func ExampleLoadContentType() {
	jsonContent := `{"name":"john", "score":"100"}`
	xmlContent := `<?xml version="1.0" encoding="UTF-8"?>
	<base>
		<name>john</name>
		<score>100</score>
	</base>`

	j, _ := gjson.LoadContentType("json", jsonContent)
	x, _ := gjson.LoadContentType("xml", xmlContent)

	fmt.Println(j.Get("name"))
	fmt.Println(j.Get("score"))
	fmt.Println(x.Get("base.name"))
	fmt.Println(x.Get("base.score"))

	// Output:
	// john
	// 100
	// john
	// 100
}

func ExampleIsValidDataType() {
	fmt.Println(gjson.IsValidDataType("json"))
	fmt.Println(gjson.IsValidDataType("yml"))
	fmt.Println(gjson.IsValidDataType("js"))
	fmt.Println(gjson.IsValidDataType("mp4"))
	fmt.Println(gjson.IsValidDataType("xsl"))
	fmt.Println(gjson.IsValidDataType("txt"))

	// Output:
	// true
	// true
	// true
	// false
	// false
	// false
}

func ExampleLoad_Xml() {
	jsonFilePath := gdebug.TestDataPath("xml", "data1.xml")
	j, _ := gjson.Load(jsonFilePath)
	fmt.Println(j.Get("doc.name"))
	fmt.Println(j.Get("doc.score"))
}
