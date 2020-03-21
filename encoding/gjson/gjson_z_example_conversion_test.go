// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gjson_test

import (
	"fmt"
	"github.com/gogf/gf/encoding/gjson"
)

func Example_Conversion_NormalFormats() {
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
	//       - John
	//       - Ming
	//     count: 1
	//
	// ======================
	// TOML:
	// [users]
	//   array = ["John", "Ming"]
	//   count = 1.0
}

func Example_Conversion_GetStruct() {
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
		if err := j.GetStruct("users", users); err != nil {
			panic(err)
		}
		fmt.Printf(`%+v`, users)
	}

	// Output:
	// &{Count:1 Array:[John Ming]}
}

func Example_Conversion_ToStruct() {
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
		if err := j.ToStruct(users); err != nil {
			panic(err)
		}
		fmt.Printf(`%+v`, users)
	}

	// Output:
	// &{Count:1 Array:[John Ming]}
}
