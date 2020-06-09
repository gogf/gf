// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with gm file,
// You can obtain one at https://github.com/gogf/gf.

package gmap_test

import (
	"fmt"
	"github.com/gogf/gf/frame/g"

	"github.com/gogf/gf/container/gmap"
)

func Example_normalBasic() {
	m := gmap.New()

	//Add data
	m.Set("key1", "val1")

	//Print size
	fmt.Println(m.Size())
	//output 1

	add_map := make(map[interface{}]interface{})
	add_map["key2"] = "val2"
	add_map["key3"] = "val3"
	add_map[1] = 1

	fmt.Println(m.Values())

	//Batch add data
	m.Sets(add_map)

	//Gets the value of the corresponding key
	key3_val := m.Get("key3")
	fmt.Println(key3_val)

	//Get the value by key, or set it with given key-value if not exist.
	get_or_set_val := m.GetOrSet("key4", "val4")
	fmt.Println(get_or_set_val)

	// Set key-value if the key does not exist, then return true; or else return false.
	is_set := m.SetIfNotExist("key3", "val3")
	fmt.Println(is_set)

	//Remove key
	m.Remove("key2")
	fmt.Println(m.Keys())

	//Batch remove keys
	remove_keys := []interface{}{"key1", 1}
	m.Removes(remove_keys)
	fmt.Println(m.Keys())

	//Contains checks whether a key exists.
	is_contain := m.Contains("key3")
	fmt.Println(is_contain)

	//Flip exchanges key-value of the map, it will change key-value to value-key.
	m.Flip()
	fmt.Println(m.Map())

	// Clear deletes all data of the map,
	m.Clear()

	fmt.Println(m.Size())
}

func Example_keysValues() {
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

func Example_flip() {
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

func Example_pop() {
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

func Example_filter() {
	m1 := gmap.NewFrom(g.MapAnyAny{
		"k1": "",
		"k2": nil,
		"k3": 0,
		"k4": 1,
	})
	m2 := gmap.NewFrom(g.MapAnyAny{
		"k1": "",
		"k2": nil,
		"k3": 0,
		"k4": 1,
	})
	m1.FilterEmpty()
	m2.FilterNil()
	fmt.Println(m1.Map())
	fmt.Println(m2.Map())

	// May Output:
	// map[k4:1]
	// map[k1: k3:0 k4:1]
}

func Example_setIfNotExist() {
	var m gmap.Map
	fmt.Println(m.SetIfNotExist("k1", "v1"))
	fmt.Println(m.SetIfNotExist("k1", "v1"))
	fmt.Println(m.Map())

	// Output:
	// true
	// false
	// map[k1:v1]
}

func Example_normalMerge() {
	var m1, m2 gmap.Map
	m1.Set("key1", "val1")
	m2.Set("key2", "val2")
	m1.Merge(&m2)
	fmt.Println(m1.Map())

	// May Output:
	// map[key1:val1 key2:val2]
}
