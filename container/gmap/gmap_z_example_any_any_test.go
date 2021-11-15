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

func ExampleAnyAnyMap_Iterator() {

	// Output:
}

func ExampleAnyAnyMap_Clone() {

	// Output:
}

func ExampleAnyAnyMap_Map() {

	// Output:
}

func ExampleAnyAnyMap_MapCopy() {

	// Output:
}

func ExampleAnyAnyMap_MapStrAny() {

	// Output:
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

func ExampleAnyAnyMap_Set() {

	// Output:
}

func ExampleAnyAnyMap_Sets() {

	// Output:
}

func ExampleAnyAnyMap_Search() {

	// Output:
}

func ExampleAnyAnyMap_Get() {

	// Output:
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

func ExampleAnyAnyMap_GetOrSet() {

	// Output:
}

func ExampleAnyAnyMap_GetOrSetFunc() {

	// Output:
}

func ExampleAnyAnyMap_GetOrSetFuncLock() {

	// Output:
}

func ExampleAnyAnyMap_GetVar() {

	// Output:
}

func ExampleAnyAnyMap_GetVarOrSet() {

	// Output:
}

func ExampleAnyAnyMap_GetVarOrSetFunc() {

	// Output:
}

func ExampleAnyAnyMap_GetVarOrSetFuncLock() {

	// Output:
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

func ExampleAnyAnyMap_SetIfNotExistFunc() {

	// Output:
}

func ExampleAnyAnyMap_SetIfNotExistFuncLock() {

	// Output:
}

func ExampleAnyAnyMap_Remove() {

	// Output:
}

func ExampleAnyAnyMap_Removes() {

	// Output:
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

func ExampleAnyAnyMap_Contains() {

	// Output:
}

func ExampleAnyAnyMap_Size() {

	// Output:
}

func ExampleAnyAnyMap_IsEmpty() {

	// Output:
}

func ExampleAnyAnyMap_Clear() {

	// Output:
}

func ExampleAnyAnyMap_Replace() {

	// Output:
}

func ExampleAnyAnyMap_LockFunc() {

	// Output:
}

func ExampleAnyAnyMap_RLockFunc() {

	// Output:
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

func ExampleAnyAnyMap_Merge() {
	var m1, m2 gmap.Map
	m1.Set("key1", "val1")
	m2.Set("key2", "val2")
	m1.Merge(&m2)
	fmt.Println(m1.Map())

	// May Output:
	// map[key1:val1 key2:val2]
}

func ExampleAnyAnyMap_String() {

	// Output:
}

func ExampleAnyAnyMap_MarshalJSON() {

	// Output:
}

func ExampleAnyAnyMap_UnmarshalJSON() {

	// Output:
}

func ExampleAnyAnyMap_UnmarshalValue() {

	// Output:
}
