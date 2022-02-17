package gjson_test

import (
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
)

func ExampleJson_SetSplitChar() {
	data :=
		`{
        "users" : {
            "count" : 2,
            "list"  : [
                {"name" : "Ming",  "score" : 60},
                {"name" : "John", "score" : 99.5}
            ]
        }
    }`
	if j, err := gjson.DecodeToJson(data); err != nil {
		panic(err)
	} else {
		j.SetSplitChar('#')
		fmt.Println("John Score:", j.Get("users#list#1#score").Float32())
	}
	// Output:
	// John Score: 99.5
}

func ExampleJson_SetViolenceCheck() {
	data :=
		`{
        "users" : {
            "count" : 100
        },
        "users.count" : 101
    }`
	if j, err := gjson.DecodeToJson(data); err != nil {
		fmt.Println(err)
	} else {
		j.SetViolenceCheck(true)
		fmt.Println("Users Count:", j.Get("users.count"))
	}
	// Output:
	// Users Count: 101
}

// ========================================================================
// JSON
// ========================================================================
func ExampleJson_ToJson() {
	type BaseInfo struct {
		Name string
		Age  int
	}

	info := BaseInfo{
		Name: "John",
		Age:  18,
	}

	j := gjson.New(info)
	jsonBytes, _ := j.ToJson()
	fmt.Println(string(jsonBytes))

	// Output:
	// {"Age":18,"Name":"John"}
}

func ExampleJson_ToJsonString() {
	type BaseInfo struct {
		Name string
		Age  int
	}

	info := BaseInfo{
		Name: "John",
		Age:  18,
	}

	j := gjson.New(info)
	jsonStr, _ := j.ToJsonString()
	fmt.Println(jsonStr)

	// Output:
	// {"Age":18,"Name":"John"}
}

func ExampleJson_ToJsonIndent() {
	type BaseInfo struct {
		Name string
		Age  int
	}

	info := BaseInfo{
		Name: "John",
		Age:  18,
	}

	j := gjson.New(info)
	jsonBytes, _ := j.ToJsonIndent()
	fmt.Println(string(jsonBytes))

	// Output:
	//{
	//	"Age": 18,
	//	"Name": "John"
	//}
}

func ExampleJson_ToJsonIndentString() {
	type BaseInfo struct {
		Name string
		Age  int
	}

	info := BaseInfo{
		Name: "John",
		Age:  18,
	}

	j := gjson.New(info)
	jsonStr, _ := j.ToJsonIndentString()
	fmt.Println(jsonStr)

	// Output:
	//{
	//	"Age": 18,
	//	"Name": "John"
	//}
}

func ExampleJson_MustToJson() {
	type BaseInfo struct {
		Name string
		Age  int
	}

	info := BaseInfo{
		Name: "John",
		Age:  18,
	}

	j := gjson.New(info)
	jsonBytes := j.MustToJson()
	fmt.Println(string(jsonBytes))

	// Output:
	// {"Age":18,"Name":"John"}
}

func ExampleJson_MustToJsonString() {
	type BaseInfo struct {
		Name string
		Age  int
	}

	info := BaseInfo{
		Name: "John",
		Age:  18,
	}

	j := gjson.New(info)
	jsonStr := j.MustToJsonString()
	fmt.Println(jsonStr)

	// Output:
	// {"Age":18,"Name":"John"}
}

func ExampleJson_MustToJsonIndent() {
	type BaseInfo struct {
		Name string
		Age  int
	}

	info := BaseInfo{
		Name: "John",
		Age:  18,
	}

	j := gjson.New(info)
	jsonBytes := j.MustToJsonIndent()
	fmt.Println(string(jsonBytes))

	// Output:
	//{
	//	"Age": 18,
	//	"Name": "John"
	//}
}

func ExampleJson_MustToJsonIndentString() {
	type BaseInfo struct {
		Name string
		Age  int
	}

	info := BaseInfo{
		Name: "John",
		Age:  18,
	}

	j := gjson.New(info)
	jsonStr := j.MustToJsonIndentString()
	fmt.Println(jsonStr)

	// Output:
	//{
	//	"Age": 18,
	//	"Name": "John"
	//}
}

// ========================================================================
// XML
// ========================================================================
func ExampleJson_ToXml() {
	type BaseInfo struct {
		Name string
		Age  int
	}

	info := BaseInfo{
		Name: "John",
		Age:  18,
	}

	j := gjson.New(info)
	xmlBytes, _ := j.ToXml()
	fmt.Println(string(xmlBytes))

	// Output:
	// <doc><Age>18</Age><Name>John</Name></doc>
}

func ExampleJson_ToXmlString() {
	type BaseInfo struct {
		Name string
		Age  int
	}

	info := BaseInfo{
		Name: "John",
		Age:  18,
	}

	j := gjson.New(info)
	xmlStr, _ := j.ToXmlString()
	fmt.Println(string(xmlStr))

	// Output:
	// <doc><Age>18</Age><Name>John</Name></doc>
}

func ExampleJson_ToXmlIndent() {
	type BaseInfo struct {
		Name string
		Age  int
	}

	info := BaseInfo{
		Name: "John",
		Age:  18,
	}

	j := gjson.New(info)
	xmlBytes, _ := j.ToXmlIndent()
	fmt.Println(string(xmlBytes))

	// Output:
	//<doc>
	//	<Age>18</Age>
	//	<Name>John</Name>
	//</doc>
}

func ExampleJson_ToXmlIndentString() {
	type BaseInfo struct {
		Name string
		Age  int
	}

	info := BaseInfo{
		Name: "John",
		Age:  18,
	}

	j := gjson.New(info)
	xmlStr, _ := j.ToXmlIndentString()
	fmt.Println(string(xmlStr))

	// Output:
	//<doc>
	//	<Age>18</Age>
	//	<Name>John</Name>
	//</doc>
}

func ExampleJson_MustToXml() {
	type BaseInfo struct {
		Name string
		Age  int
	}

	info := BaseInfo{
		Name: "John",
		Age:  18,
	}

	j := gjson.New(info)
	xmlBytes := j.MustToXml()
	fmt.Println(string(xmlBytes))

	// Output:
	// <doc><Age>18</Age><Name>John</Name></doc>
}

func ExampleJson_MustToXmlString() {
	type BaseInfo struct {
		Name string
		Age  int
	}

	info := BaseInfo{
		Name: "John",
		Age:  18,
	}

	j := gjson.New(info)
	xmlStr := j.MustToXmlString()
	fmt.Println(string(xmlStr))

	// Output:
	// <doc><Age>18</Age><Name>John</Name></doc>
}

func ExampleJson_MustToXmlIndent() {
	type BaseInfo struct {
		Name string
		Age  int
	}

	info := BaseInfo{
		Name: "John",
		Age:  18,
	}

	j := gjson.New(info)
	xmlBytes := j.MustToXmlIndent()
	fmt.Println(string(xmlBytes))

	// Output:
	//<doc>
	//	<Age>18</Age>
	//	<Name>John</Name>
	//</doc>
}

func ExampleJson_MustToXmlIndentString() {
	type BaseInfo struct {
		Name string
		Age  int
	}

	info := BaseInfo{
		Name: "John",
		Age:  18,
	}

	j := gjson.New(info)
	xmlStr := j.MustToXmlIndentString()
	fmt.Println(string(xmlStr))

	// Output:
	//<doc>
	//	<Age>18</Age>
	//	<Name>John</Name>
	//</doc>
}

// ========================================================================
// YAML
// ========================================================================
func ExampleJson_ToYaml() {
	type BaseInfo struct {
		Name string
		Age  int
	}

	info := BaseInfo{
		Name: "John",
		Age:  18,
	}

	j := gjson.New(info)
	YamlBytes, _ := j.ToYaml()
	fmt.Println(string(YamlBytes))

	// Output:
	//Age: 18
	//Name: John
}

func ExampleJson_ToYamlString() {
	type BaseInfo struct {
		Name string
		Age  int
	}

	info := BaseInfo{
		Name: "John",
		Age:  18,
	}

	j := gjson.New(info)
	YamlStr, _ := j.ToYamlString()
	fmt.Println(string(YamlStr))

	// Output:
	//Age: 18
	//Name: John
}

func ExampleJson_ToYamlIndent() {
	type BaseInfo struct {
		Name string
		Age  int
	}

	info := BaseInfo{
		Name: "John",
		Age:  18,
	}

	j := gjson.New(info)
	YamlBytes, _ := j.ToYamlIndent("")
	fmt.Println(string(YamlBytes))

	// Output:
	//Age: 18
	//Name: John
}

func ExampleJson_MustToYaml() {
	type BaseInfo struct {
		Name string
		Age  int
	}

	info := BaseInfo{
		Name: "John",
		Age:  18,
	}

	j := gjson.New(info)
	YamlBytes := j.MustToYaml()
	fmt.Println(string(YamlBytes))

	// Output:
	//Age: 18
	//Name: John
}

func ExampleJson_MustToYamlString() {
	type BaseInfo struct {
		Name string
		Age  int
	}

	info := BaseInfo{
		Name: "John",
		Age:  18,
	}

	j := gjson.New(info)
	YamlStr := j.MustToYamlString()
	fmt.Println(string(YamlStr))

	// Output:
	//Age: 18
	//Name: John
}

// ========================================================================
// TOML
// ========================================================================
func ExampleJson_ToToml() {
	type BaseInfo struct {
		Name string
		Age  int
	}

	info := BaseInfo{
		Name: "John",
		Age:  18,
	}

	j := gjson.New(info)
	TomlBytes, _ := j.ToToml()
	fmt.Println(string(TomlBytes))

	// Output:
	//Age = 18
	//Name = "John"
}

func ExampleJson_ToTomlString() {
	type BaseInfo struct {
		Name string
		Age  int
	}

	info := BaseInfo{
		Name: "John",
		Age:  18,
	}

	j := gjson.New(info)
	TomlStr, _ := j.ToTomlString()
	fmt.Println(string(TomlStr))

	// Output:
	//Age = 18
	//Name = "John"
}

func ExampleJson_MustToToml() {
	type BaseInfo struct {
		Name string
		Age  int
	}

	info := BaseInfo{
		Name: "John",
		Age:  18,
	}

	j := gjson.New(info)
	TomlBytes := j.MustToToml()
	fmt.Println(string(TomlBytes))

	// Output:
	//Age = 18
	//Name = "John"
}

func ExampleJson_MustToTomlString() {
	type BaseInfo struct {
		Name string
		Age  int
	}

	info := BaseInfo{
		Name: "John",
		Age:  18,
	}

	j := gjson.New(info)
	TomlStr := j.MustToTomlString()
	fmt.Println(string(TomlStr))

	// Output:
	//Age = 18
	//Name = "John"
}

// ========================================================================
// INI
// ========================================================================
func ExampleJson_ToIni() {
	type BaseInfo struct {
		Name string
		Age  int
	}

	info := BaseInfo{
		Name: "John",
		Age:  18,
	}

	j := gjson.New(info)
	IniBytes, _ := j.ToIni()
	fmt.Println(string(IniBytes))

	// May Output:
	//Name=John
	//Age=18
}

func ExampleJson_ToIniString() {
	type BaseInfo struct {
		Name string
		Age  int
	}

	info := BaseInfo{
		Name: "John",
		Age:  18,
	}

	j := gjson.New(info)
	IniStr, _ := j.ToIniString()
	fmt.Println(string(IniStr))

	// May Output:
	//Name=John
	//Age=18
}

func ExampleJson_MustToIni() {
	type BaseInfo struct {
		Name string
		Age  int
	}

	info := BaseInfo{
		Name: "John",
		Age:  18,
	}

	j := gjson.New(info)
	IniBytes := j.MustToIni()
	fmt.Println(string(IniBytes))

	// May Output:
	//Name=John
	//Age=18
}

func ExampleJson_MustToIniString() {
	type BaseInfo struct {
		Name string
		Age  int
	}

	info := BaseInfo{
		Name: "John",
		Age:  18,
	}

	j := gjson.New(info)
	IniStr := j.MustToIniString()
	fmt.Println(string(IniStr))

	// May Output:
	//Name=John
	//Age=18
}

func ExampleJson_MarshalJSON() {
	type BaseInfo struct {
		Name string
		Age  int
	}

	info := BaseInfo{
		Name: "John",
		Age:  18,
	}

	j := gjson.New(info)
	jsonBytes, _ := j.MarshalJSON()
	fmt.Println(string(jsonBytes))

	// Output:
	// {"Age":18,"Name":"John"}
}

func ExampleJson_UnmarshalJSON() {
	jsonStr := `{"Age":18,"Name":"John"}`

	j := gjson.New("")
	j.UnmarshalJSON([]byte(jsonStr))
	fmt.Println(j.Map())

	// Output:
	// map[Age:18 Name:John]
}

func ExampleJson_UnmarshalValue_Xml() {
	xmlStr := `<doc><Age>18</Age><Name>John</Name></doc>`

	j := gjson.New("")
	j.UnmarshalValue([]byte(xmlStr))
	fmt.Println(j.Var().String())

	// Output:
	// {"doc":{"Age":"18","Name":"John"}}
}
