// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with gm file,
// You can obtain one at https://github.com/gogf/gf.

package gmap_test

import (
	"fmt"

	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/frame/g"
)

func ExampleNew() {
	m := gmap.New()

	// Add data.
	m.Set("key1", "val1")

	// Print size.
	fmt.Println(m.Size())

	addMap := make(map[interface{}]interface{})
	addMap["key2"] = "val2"
	addMap["key3"] = "val3"
	addMap[1] = 1

	fmt.Println(m.Values())

	// Batch add data.
	m.Sets(addMap)

	// Gets the value of the corresponding key.
	fmt.Println(m.Get("key3"))

	// Get the value by key, or set it with given key-value if not exist.
	fmt.Println(m.GetOrSet("key4", "val4"))

	// Set key-value if the key does not exist, then return true; or else return false.
	fmt.Println(m.SetIfNotExist("key3", "val3"))

	// Remove key
	m.Remove("key2")
	fmt.Println(m.Keys())

	// Batch remove keys.
	m.Removes([]interface{}{"key1", 1})
	fmt.Println(m.Keys())

	// Contains checks whether a key exists.
	fmt.Println(m.Contains("key3"))

	// Flip exchanges key-value of the map, it will change key-value to value-key.
	m.Flip()
	fmt.Println(m.Map())

	// Clear deletes all data of the map.
	m.Clear()

	fmt.Println(m.Size())

	// May Output:
	// 1
	// [val1]
	// val3
	// val4
	// false
	// [key4 key1 key3 1]
	// [key4 key3]
	// true
	// map[val3:key3 val4:key4]
	// 0
}

func ExampleAnyAnyMap_Keys() {
	var m gmap.Map
	m.Sets(g.MapAnyAny{
		"k1": "v1",
		"k2": "v2",
		"k3": "v3",
		"k4": "v4",
	})
	fmt.Println(m.Keys())
	fmt.Println(m.Values())

	// May Output:
	// [k1 k2 k3 k4]
	// [v2 v3 v4 v1]
}

func ExampleAnyAnyMap_Values() {
	var m gmap.Map
	m.Sets(g.MapAnyAny{
		"k1": "v1",
		"k2": "v2",
		"k3": "v3",
		"k4": "v4",
	})
	fmt.Println(m.Keys())
	fmt.Println(m.Values())

	// May Output:
	// [k1 k2 k3 k4]
	// [v2 v3 v4 v1]
}

func ExampleAnyAnyMap_Flip() {
	var m gmap.Map
	m.Sets(g.MapAnyAny{
		"k1": "v1",
		"k2": "v2",
	})
	m.Flip()
	fmt.Println(m.Map())

	// May Output:
	// map[v1:k1 v2:k2]
}

func ExampleAnyAnyMap_Pop() {
	var m gmap.Map
	m.Sets(g.MapAnyAny{
		"k1": "v1",
		"k2": "v2",
		"k3": "v3",
		"k4": "v4",
	})
	fmt.Println(m.Pop())
	fmt.Println(m.Pops(2))
	fmt.Println(m.Size())

	// May Output:
	// k1 v1
	// map[k2:v2 k4:v4]
	// 1
}

func ExampleAnyAnyMap_Pops() {
	var m gmap.Map
	m.Sets(g.MapAnyAny{
		"k1": "v1",
		"k2": "v2",
		"k3": "v3",
		"k4": "v4",
	})
	fmt.Println(m.Pop())
	fmt.Println(m.Pops(2))
	fmt.Println(m.Size())

	// May Output:
	// k1 v1
	// map[k2:v2 k4:v4]
	// 1
}

func ExampleAnyAnyMap_FilterEmpty() {
	m := gmap.NewFrom(g.MapAnyAny{
		"k1": "",
		"k2": nil,
		"k3": 0,
		"k4": 1,
	})
	m.FilterEmpty()
	fmt.Println(m.Map())

	// May Output:
	// map[k4:1]
}

func ExampleAnyAnyMap_FilterNil() {
	m := gmap.NewFrom(g.MapAnyAny{
		"k1": "",
		"k2": nil,
		"k3": 0,
		"k4": 1,
	})
	m.FilterNil()
	fmt.Println(m.Map())

	// May Output:
	// map[k1: k3:0 k4:1]
}

func ExampleAnyAnyMap_SetIfNotExist() {
	var m gmap.Map
	fmt.Println(m.SetIfNotExist("k1", "v1"))
	fmt.Println(m.SetIfNotExist("k1", "v1"))
	fmt.Println(m.Map())

	// Output:
	// true
	// false
	// map[k1:v1]
}

func ExampleAnyAnyMap_Merge() {
	var m1, m2 gmap.Map
	m1.Set("key1", "val1")
	m2.Set("key2", "val2")
	m1.Merge(&m2)
	fmt.Println(m1.Map())

	// May Output:
	// map[key1:val1 key2:val2]
}
