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

func Example_conversionNormalFormats() {
	data :=
		`{
        "users" : {
            "count" : 1,
            "array" : ["John", "Ming"]
        }
    }`

	if j, err := gjson.DecodeToJson(data); err != nil {
		panic(err)
	} else {
		fmt.Println("JSON:")
		fmt.Println(j.MustToJsonString())
		fmt.Println("======================")

		fmt.Println("XML:")
		fmt.Println(j.MustToXmlString())
		fmt.Println("======================")

		fmt.Println("YAML:")
		fmt.Println(j.MustToYamlString())
		fmt.Println("======================")

		fmt.Println("TOML:")
		fmt.Println(j.MustToTomlString())
	}

	// Output:
	// JSON:
	// {"users":{"array":["John","Ming"],"count":1}}
	// ======================
	// XML:
	// <users><array>John</array><array>Ming</array><count>1</count></users>
	// ======================
	// YAML:
	// users:
	//     array:
	//         - John
	//         - Ming
	//     count: 1
	//
	// ======================
	// TOML:
	// [users]
	//   array = ["John", "Ming"]
	//   count = 1.0
}

func ExampleJson_ConversionGetStruct() {
	data :=
		`{
        "users" : {
            "count" : 1,
            "array" : ["John", "Ming"]
        }
    }`
	if j, err := gjson.DecodeToJson(data); err != nil {
		panic(err)
	} else {
		type Users struct {
			Count int
			Array []string
		}
		users := new(Users)
		if err := j.Get("users").Scan(users); err != nil {
			panic(err)
		}
		fmt.Printf(`%+v`, users)
	}

	// Output:
	// &{Count:1 Array:[John Ming]}
}

func ExampleJson_ConversionToStruct() {
	data :=
		`
	{
        "count" : 1,
        "array" : ["John", "Ming"]
    }`
	if j, err := gjson.DecodeToJson(data); err != nil {
		panic(err)
	} else {
		type Users struct {
			Count int
			Array []string
		}
		users := new(Users)
		if err := j.Var().Scan(users); err != nil {
			panic(err)
		}
		fmt.Printf(`%+v`, users)
	}

	// Output:
	// &{Count:1 Array:[John Ming]}
}

func ExampleValid() {
	data1 := []byte(`{"n":123456789, "m":{"k":"v"}, "a":[1,2,3]}`)
	data2 := []byte(`{"n":123456789, "m":{"k":"v"}, "a":[1,2,3]`)
	fmt.Println(gjson.Valid(data1))
	fmt.Println(gjson.Valid(data2))

	// Output:
	// true
	// false
}

func ExampleMarshal() {
	data := map[string]interface{}{
		"name":  "john",
		"score": 100,
	}

	jsonData, _ := gjson.Marshal(data)
	fmt.Println(string(jsonData))

	type BaseInfo struct {
		Name string
		Age  int
	}

	info := BaseInfo{
		Name: "Guo Qiang",
		Age:  18,
	}

	infoData, _ := gjson.Marshal(info)
	fmt.Println(string(infoData))

	// Output:
	// {"name":"john","score":100}
	// {"Name":"Guo Qiang","Age":18}
}

func ExampleMarshalIndent() {
	type BaseInfo struct {
		Name string
		Age  int
	}

	info := BaseInfo{
		Name: "John",
		Age:  18,
	}

	infoData, _ := gjson.MarshalIndent(info, "", "\t")
	fmt.Println(string(infoData))

	// Output:
	// {
	//	"Name": "John",
	//	"Age": 18
	// }
}

func ExampleUnmarshal() {
	type BaseInfo struct {
		Name  string
		Score int
	}

	var info BaseInfo

	jsonContent := "{\"name\":\"john\",\"score\":100}"
	gjson.Unmarshal([]byte(jsonContent), &info)
	fmt.Printf("%+v", info)

	// Output:
	// {Name:john Score:100}
}

func ExampleEncode() {
	type BaseInfo struct {
		Name string
		Age  int
	}

	info := BaseInfo{
		Name: "John",
		Age:  18,
	}

	infoData, _ := gjson.Encode(info)
	fmt.Println(string(infoData))

	// Output:
	// {"Name":"John","Age":18}
}

func ExampleMustEncode() {
	type BaseInfo struct {
		Name string
		Age  int
	}

	info := BaseInfo{
		Name: "John",
		Age:  18,
	}

	infoData := gjson.MustEncode(info)
	fmt.Println(string(infoData))

	// Output:
	// {"Name":"John","Age":18}
}

func ExampleEncodeString() {
	type BaseInfo struct {
		Name string
		Age  int
	}

	info := BaseInfo{
		Name: "John",
		Age:  18,
	}

	infoData, _ := gjson.EncodeString(info)
	fmt.Println(infoData)

	// Output:
	// {"Name":"John","Age":18}
}

func ExampleMustEncodeString() {
	type BaseInfo struct {
		Name string
		Age  int
	}

	info := BaseInfo{
		Name: "John",
		Age:  18,
	}

	infoData := gjson.MustEncodeString(info)
	fmt.Println(infoData)

	// Output:
	// {"Name":"John","Age":18}
}

func ExampleDecode() {
	jsonContent := `{"name":"john","score":100}`
	info, _ := gjson.Decode([]byte(jsonContent))
	fmt.Println(info)

	// Output:
	// map[name:john score:100]
}

func ExampleDecodeTo() {
	type BaseInfo struct {
		Name  string
		Score int
	}

	var info BaseInfo

	jsonContent := "{\"name\":\"john\",\"score\":100}"
	gjson.DecodeTo([]byte(jsonContent), &info)
	fmt.Printf("%+v", info)

	// Output:
	// {Name:john Score:100}
}

func ExampleDecodeToJson() {
	jsonContent := `{"name":"john","score":100}"`
	j, _ := gjson.DecodeToJson([]byte(jsonContent))
	fmt.Println(j.Map())

	// May Output:
	// map[name:john score:100]
}
