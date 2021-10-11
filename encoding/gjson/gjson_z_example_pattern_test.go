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

func Example_patternGet() {
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
		fmt.Println("John Score:", j.Get("users.list.1.score").Float32())
	}
	// Output:
	// John Score: 99.5
}

func Example_patternCustomSplitChar() {
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

func Example_patternViolenceCheck() {
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
		fmt.Println("Users Count:", j.Get("users.count").Int())
	}
	// Output:
	// Users Count: 101
}

func Example_mapSliceChange() {
	jsonContent := `{"map":{"key":"value"}, "slice":[59,90]}`
	j, _ := gjson.LoadJson(jsonContent)
	m := j.Get("map").Map()
	fmt.Println(m)

	// Change the key-value pair.
	m["key"] = "john"

	// It changes the underlying key-value pair.
	fmt.Println(j.Get("map").Map())

	s := j.Get("slice").Array()
	fmt.Println(s)

	// Change the value of specified index.
	s[0] = 100

	// It changes the underlying slice.
	fmt.Println(j.Get("slice").Array())

	// output:
	// map[key:value]
	// map[key:john]
	// [59 90]
	// [100 90]
}
