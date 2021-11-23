// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with gm file,
// You can obtain one at https://github.com/gogf/gf.

package gmap_test

import (
	"fmt"
	"github.com/gogf/gf/v2/util/gutil"

	"github.com/gogf/gf/v2/container/gmap"
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

func ExampleNewFrom() {
	m := gmap.New()

	m.Set("key1", "val1")
	fmt.Println(m)

	n := gmap.NewFrom(m.MapCopy(), true)
	fmt.Println(n)

	// Output:
	// {"key1":"val1"}
	// {"key1":"val1"}
}

func ExampleNewHashMap() {
	m := gmap.New()

	m.Set("key1", "val1")
	fmt.Println(m)

	// Output:
	// {"key1":"val1"}
}

func ExampleNewHashMapFrom() {
	m := gmap.New()

	m.Set("key1", "val1")
	fmt.Println(m)

	n := gmap.NewFrom(m.MapCopy(), true)
	fmt.Println(n)

	// Output:
	// {"key1":"val1"}
	// {"key1":"val1"}
}

func ExampleNewAnyAnyMap() {
	m := gmap.NewAnyAnyMap()

	m.Set("key1", "val1")
	fmt.Println(m)

	// Output:
	// {"key1":"val1"}
}

func ExampleNewAnyAnyMapFrom() {
	m := gmap.NewAnyAnyMap()

	m.Set("key1", "val1")
	fmt.Println(m)

	n := gmap.NewAnyAnyMapFrom(m.MapCopy(), true)
	fmt.Println(n)

	// Output:
	// {"key1":"val1"}
	// {"key1":"val1"}
}

func ExampleNewIntAnyMap() {
	m := gmap.NewIntAnyMap()

	m.Set(1, "val1")
	fmt.Println(m)

	// Output:
	// {"1":"val1"}
}

func ExampleNewIntAnyMapFrom() {
	m := gmap.NewIntAnyMap()

	m.Set(1, "val1")
	fmt.Println(m)

	n := gmap.NewIntAnyMapFrom(m.MapCopy(), true)
	fmt.Println(n)

	// Output:
	// {"1":"val1"}
	// {"1":"val1"}
}

func ExampleNewIntIntMap() {
	m := gmap.NewIntIntMap()

	m.Set(1, 1)
	fmt.Println(m)

	// Output:
	// {"1":1}
}

func ExampleNewIntIntMapFrom() {
	m := gmap.NewIntIntMap()

	m.Set(1, 1)
	fmt.Println(m)

	n := gmap.NewIntIntMapFrom(m.MapCopy(), true)
	fmt.Println(n)

	// Output:
	// {"1":1}
	// {"1":1}
}

func ExampleNewStrAnyMap() {
	m := gmap.NewStrAnyMap()

	m.Set("key1", "var1")
	fmt.Println(m)

	// Output:
	// {"key1":"var1"}
}

func ExampleNewStrAnyMapFrom() {
	m := gmap.NewStrAnyMap()

	m.Set("key1", "var1")
	fmt.Println(m)

	n := gmap.NewStrAnyMapFrom(m.MapCopy(), true)
	fmt.Println(n)

	// Output:
	// {"key1":"var1"}
	// {"key1":"var1"}
}

func ExampleNewStrIntMap() {
	m := gmap.NewStrIntMap()

	m.Set("key1", 1)
	fmt.Println(m)

	// Output:
	// {"key1":1}
}

func ExampleNewStrIntMapFrom() {
	m := gmap.NewStrIntMap()

	m.Set("key1", 1)
	fmt.Println(m)

	n := gmap.NewStrIntMapFrom(m.MapCopy(), true)
	fmt.Println(n)

	// Output:
	// {"key1":1}
	// {"key1":1}
}

func ExampleNewStrStrMap() {
	m := gmap.NewStrStrMap()

	m.Set("key1", "var1")
	fmt.Println(m)

	// Output:
	// {"key1":"var1"}
}

func ExampleNewStrStrMapFrom() {
	m := gmap.NewStrStrMap()

	m.Set("key1", "var1")
	fmt.Println(m)

	n := gmap.NewStrStrMapFrom(m.MapCopy(), true)
	fmt.Println(n)

	// Output:
	// {"key1":"var1"}
	// {"key1":"var1"}
}

func ExampleNewListMap() {
	m := gmap.NewListMap()

	m.Set("key1", "var1")
	m.Set("key2", "var2")
	fmt.Println(m)

	// Output:
	// {"key1":"var1","key2":"var2"}
}

func ExampleNewListMapFrom() {
	m := gmap.NewListMap()

	m.Set("key1", "var1")
	m.Set("key2", "var2")
	fmt.Println(m)

	n := gmap.NewListMapFrom(m.Map(), true)
	fmt.Println(n)

	// Output:
	// {"key1":"var1","key2":"var2"}
	// {"key1":"var1","key2":"var2"}
}

func ExampleNewTreeMap() {
	m := gmap.NewTreeMap(gutil.ComparatorString)

	m.Set("key2", "var2")
	m.Set("key1", "var1")

	fmt.Println(m.Map())

	// Output:
	// map[key1:var1 key2:var2]
}

func ExampleNewTreeMapFrom() {
	m := gmap.NewTreeMap(gutil.ComparatorString)

	m.Set("key2", "var2")
	m.Set("key1", "var1")

	fmt.Println(m.Map())

	n := gmap.NewListMapFrom(m.Map(), true)
	fmt.Println(n.Map())

	// Output:
	// map[key1:var1 key2:var2]
	// map[key1:var1 key2:var2]
}
