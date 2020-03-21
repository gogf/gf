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

func Example_Pattern_Get() {
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
		fmt.Println("John Score:", j.GetFloat32("users.list.1.score"))
	}
	// Output:
	// John Score: 99.5
}

func Example_Pattern_CustomSplitChar() {
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
		fmt.Println("John Score:", j.GetFloat32("users#list#1#score"))
	}
	// Output:
	// John Score: 99.5
}

func Example_Pattern_ViolenceCheck() {
	data :=
		`{
        "users" : {
            "count" : 100
        },
        "users.count" : 101
    }`
	if j, err := gjson.DecodeToJson(data); err != nil {
		panic(err)
	} else {
		j.SetViolenceCheck(true)
		fmt.Println("Users Count:", j.GetInt("users.count"))
	}
	// Output:
	// Users Count: 101
}
