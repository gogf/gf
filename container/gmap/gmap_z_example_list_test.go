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
	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/util/gconv"
)

func ExampleListMap_Iterator() {
	m := gmap.NewListMap()
	for i := 0; i < 10; i++ {
		m.Set(i, i*2)
	}

	var totalKey, totalValue int
	m.Iterator(func(k interface{}, v interface{}) bool {
		totalKey += k.(int)
		totalValue += v.(int)

		return totalKey < 10
	})

	fmt.Println("totalKey:", totalKey)
	fmt.Println("totalValue:", totalValue)

	// Output:
	// totalKey: 10
	// totalValue: 20
}

func ExampleListMap_IteratorAsc() {
	m := gmap.NewListMap()
	for i := 0; i < 10; i++ {
		m.Set(i, i*2)
	}

	var totalKey, totalValue int
	m.IteratorAsc(func(k interface{}, v interface{}) bool {
		totalKey += k.(int)
		totalValue += v.(int)

		return totalKey < 10
	})

	fmt.Println("totalKey:", totalKey)
	fmt.Println("totalValue:", totalValue)

	// Output:
	// totalKey: 10
	// totalValue: 20
}

func ExampleListMap_IteratorDesc() {
	m := gmap.NewListMap()
	for i := 0; i < 10; i++ {
		m.Set(i, i*2)
	}

	var totalKey, totalValue int
	m.IteratorDesc(func(k interface{}, v interface{}) bool {
		totalKey += k.(int)
		totalValue += v.(int)

		return totalKey < 10
	})

	fmt.Println("totalKey:", totalKey)
	fmt.Println("totalValue:", totalValue)

	// Output:
	// totalKey: 17
	// totalValue: 34
}

func ExampleListMap_Clone() {
	m := gmap.NewListMap()

	m.Set("key1", "val1")
	fmt.Println(m)

	n := m.Clone()
	fmt.Println(n)

	// Output:
	// {"key1":"val1"}
	// {"key1":"val1"}
}

func ExampleListMap_Clear() {
	var m gmap.ListMap
	m.Sets(g.MapAnyAny{
		"k1": "v1",
		"k2": "v2",
		"k3": "v3",
		"k4": "v4",
	})

	m.Clear()

	fmt.Println(m.Map())

	// Output:
	// map[]
}

func ExampleListMap_Replace() {
	var m gmap.ListMap
	m.Sets(g.MapAnyAny{
		"k1": "v1",
	})

	var n gmap.ListMap
	n.Sets(g.MapAnyAny{
		"k2": "v2",
	})

	fmt.Println(m.Map())

	m.Replace(n.Map())
	fmt.Println(m.Map())

	// Output:
	// map[k1:v1]
	// map[k2:v2]
}

func ExampleListMap_Map() {
	m1 := gmap.NewListMap()
	m1.Set("key1", "val1")
	fmt.Println("m1:", m1)

	n1 := m1.Map()
	fmt.Println("before n1:", n1)
	m1.Set("key1", "val2")
	fmt.Println("after n1:", n1)

	// Output:
	// m1: {"key1":"val1"}
	// before n1: map[key1:val1]
	// after n1: map[key1:val1]
}

func ExampleListMap_MapStrAny() {
	m := gmap.NewListMap()
	m.Set("key1", "val1")
	m.Set("key2", "val2")

	n := m.MapStrAny()
	fmt.Printf("%#v", n)

	// Output:
	// map[string]interface {}{"key1":"val1", "key2":"val2"}
}

func ExampleListMap_FilterEmpty() {
	m := gmap.NewListMapFrom(g.MapAnyAny{
		"k1": "",
		"k2": nil,
		"k3": 0,
		"k4": 1,
	})
	m.FilterEmpty()
	fmt.Println(m.Map())

	// Output:
	// map[k4:1]
}

func ExampleListMap_Set() {
	m := gmap.NewListMap()

	m.Set("key1", "val1")
	fmt.Println(m)

	// Output:
	// {"key1":"val1"}
}

func ExampleListMap_Sets() {
	m := gmap.NewListMap()

	addMap := make(map[interface{}]interface{})
	addMap["key1"] = "val1"
	addMap["key2"] = "val2"
	addMap["key3"] = "val3"

	m.Sets(addMap)
	fmt.Println(m)

	// Output:
	// {"key1":"val1","key2":"val2","key3":"val3"}
}

func ExampleListMap_Search() {
	m := gmap.NewListMap()

	m.Set("key1", "val1")

	value, found := m.Search("key1")
	if found {
		fmt.Println("find key1 value:", value)
	}

	value, found = m.Search("key2")
	if !found {
		fmt.Println("key2 not find")
	}

	// Output:
	// find key1 value: val1
	// key2 not find
}

func ExampleListMap_Get() {
	m := gmap.NewListMap()

	m.Set("key1", "val1")

	fmt.Println("key1 value:", m.Get("key1"))
	fmt.Println("key2 value:", m.Get("key2"))

	// Output:
	// key1 value: val1
	// key2 value: <nil>
}

func ExampleListMap_Pop() {
	var m gmap.ListMap
	m.Sets(g.MapAnyAny{
		"k1": "v1",
		"k2": "v2",
		"k3": "v3",
		"k4": "v4",
	})

	fmt.Println(m.Pop())

	// May Output:
	// k1 v1
}

func ExampleListMap_Pops() {
	var m gmap.ListMap
	m.Sets(g.MapAnyAny{
		"k1": "v1",
		"k2": "v2",
		"k3": "v3",
		"k4": "v4",
	})
	fmt.Println(m.Pops(-1))
	fmt.Println("size:", m.Size())

	m.Sets(g.MapAnyAny{
		"k1": "v1",
		"k2": "v2",
		"k3": "v3",
		"k4": "v4",
	})
	fmt.Println(m.Pops(2))
	fmt.Println("size:", m.Size())

	// May Output:
	// map[k1:v1 k2:v2 k3:v3 k4:v4]
	// size: 0
	// map[k1:v1 k2:v2]
	// size: 2
}

func ExampleListMap_GetOrSet() {
	m := gmap.NewListMap()
	m.Set("key1", "val1")

	fmt.Println(m.GetOrSet("key1", "NotExistValue"))
	fmt.Println(m.GetOrSet("key2", "val2"))

	// Output:
	// val1
	// val2
}

func ExampleListMap_GetOrSetFunc() {
	m := gmap.NewListMap()
	m.Set("key1", "val1")

	fmt.Println(m.GetOrSetFunc("key1", func() interface{} {
		return "NotExistValue"
	}))
	fmt.Println(m.GetOrSetFunc("key2", func() interface{} {
		return "NotExistValue"
	}))

	// Output:
	// val1
	// NotExistValue
}

func ExampleListMap_GetOrSetFuncLock() {
	m := gmap.NewListMap()
	m.Set("key1", "val1")

	fmt.Println(m.GetOrSetFuncLock("key1", func() interface{} {
		return "NotExistValue"
	}))
	fmt.Println(m.GetOrSetFuncLock("key2", func() interface{} {
		return "NotExistValue"
	}))

	// Output:
	// val1
	// NotExistValue
}

func ExampleListMap_GetVar() {
	m := gmap.NewListMap()
	m.Set("key1", "val1")

	fmt.Println(m.GetVar("key1"))
	fmt.Println(m.GetVar("key2").IsNil())

	// Output:
	// val1
	// true
}

func ExampleListMap_GetVarOrSet() {
	m := gmap.NewListMap()
	m.Set("key1", "val1")

	fmt.Println(m.GetVarOrSet("key1", "NotExistValue"))
	fmt.Println(m.GetVarOrSet("key2", "val2"))

	// Output:
	// val1
	// val2
}

func ExampleListMap_GetVarOrSetFunc() {
	m := gmap.NewListMap()
	m.Set("key1", "val1")

	fmt.Println(m.GetVarOrSetFunc("key1", func() interface{} {
		return "NotExistValue"
	}))
	fmt.Println(m.GetVarOrSetFunc("key2", func() interface{} {
		return "NotExistValue"
	}))

	// Output:
	// val1
	// NotExistValue
}

func ExampleListMap_GetVarOrSetFuncLock() {
	m := gmap.NewListMap()
	m.Set("key1", "val1")

	fmt.Println(m.GetVarOrSetFuncLock("key1", func() interface{} {
		return "NotExistValue"
	}))
	fmt.Println(m.GetVarOrSetFuncLock("key2", func() interface{} {
		return "NotExistValue"
	}))

	// Output:
	// val1
	// NotExistValue
}

func ExampleListMap_SetIfNotExist() {
	var m gmap.ListMap
	fmt.Println(m.SetIfNotExist("k1", "v1"))
	fmt.Println(m.SetIfNotExist("k1", "v2"))
	fmt.Println(m.Map())

	// Output:
	// true
	// false
	// map[k1:v1]
}

func ExampleListMap_SetIfNotExistFunc() {
	var m gmap.ListMap
	fmt.Println(m.SetIfNotExistFunc("k1", func() interface{} {
		return "v1"
	}))
	fmt.Println(m.SetIfNotExistFunc("k1", func() interface{} {
		return "v2"
	}))
	fmt.Println(m.Map())

	// Output:
	// true
	// false
	// map[k1:v1]
}

func ExampleListMap_SetIfNotExistFuncLock() {
	var m gmap.ListMap
	fmt.Println(m.SetIfNotExistFuncLock("k1", func() interface{} {
		return "v1"
	}))
	fmt.Println(m.SetIfNotExistFuncLock("k1", func() interface{} {
		return "v2"
	}))
	fmt.Println(m.Map())

	// Output:
	// true
	// false
	// map[k1:v1]
}

func ExampleListMap_Remove() {
	var m gmap.ListMap
	m.Set("k1", "v1")

	fmt.Println(m.Remove("k1"))
	fmt.Println(m.Remove("k2"))
	fmt.Println(m.Size())

	// Output:
	// v1
	// <nil>
	// 0
}

func ExampleListMap_Removes() {
	var m gmap.ListMap
	m.Sets(g.MapAnyAny{
		"k1": "v1",
		"k2": "v2",
		"k3": "v3",
		"k4": "v4",
	})

	removeList := make([]interface{}, 2)
	removeList = append(removeList, "k1")
	removeList = append(removeList, "k2")

	m.Removes(removeList)

	fmt.Println(m.Map())

	// Output:
	// map[k3:v3 k4:v4]
}

func ExampleListMap_Keys() {
	var m gmap.ListMap
	m.Sets(g.MapAnyAny{
		"k1": "v1",
		"k2": "v2",
		"k3": "v3",
		"k4": "v4",
	})
	fmt.Println(m.Keys())

	// May Output:
	// [k1 k2 k3 k4]
}

func ExampleListMap_Values() {
	var m gmap.ListMap
	m.Sets(g.MapAnyAny{
		"k1": "v1",
		"k2": "v2",
		"k3": "v3",
		"k4": "v4",
	})
	fmt.Println(m.Values())

	// May Output:
	// [v1 v2 v3 v4]
}

func ExampleListMap_Contains() {
	var m gmap.ListMap
	m.Sets(g.MapAnyAny{
		"k1": "v1",
		"k2": "v2",
		"k3": "v3",
		"k4": "v4",
	})

	fmt.Println(m.Contains("k1"))
	fmt.Println(m.Contains("k5"))

	// Output:
	// true
	// false
}

func ExampleListMap_Size() {
	var m gmap.ListMap
	m.Sets(g.MapAnyAny{
		"k1": "v1",
		"k2": "v2",
		"k3": "v3",
		"k4": "v4",
	})

	fmt.Println(m.Size())

	// Output:
	// 4
}

func ExampleListMap_IsEmpty() {
	var m gmap.ListMap
	fmt.Println(m.IsEmpty())

	m.Set("k1", "v1")
	fmt.Println(m.IsEmpty())

	// Output:
	// true
	// false
}

func ExampleListMap_Flip() {
	var m gmap.ListMap
	m.Sets(g.MapAnyAny{
		"k1": "v1",
	})
	m.Flip()
	fmt.Println(m.Map())

	// Output:
	// map[v1:k1]
}

func ExampleListMap_Merge() {
	var m1, m2 gmap.ListMap
	m1.Set("key1", "val1")
	m2.Set("key2", "val2")
	m1.Merge(&m2)
	fmt.Println(m1.Map())

	// May Output:
	// map[key1:val1 key2:val2]
}

func ExampleListMap_String() {
	var m gmap.ListMap
	m.Sets(g.MapAnyAny{
		"k1": "v1",
	})

	fmt.Println(m.String())

	// Output:
	// {"k1":"v1"}
}

func ExampleListMap_MarshalJSON() {
	var m gmap.ListMap
	m.Set("k1", "v1")
	m.Set("k2", "v2")
	m.Set("k3", "v3")
	m.Set("k4", "v4")

	bytes, err := json.Marshal(&m)
	if err == nil {
		fmt.Println(gconv.String(bytes))
	}

	// Output:
	// {"k1":"v1","k2":"v2","k3":"v3","k4":"v4"}
}

func ExampleListMap_UnmarshalJSON() {
	var m gmap.ListMap
	m.Sets(g.MapAnyAny{
		"k1": "v1",
		"k2": "v2",
		"k3": "v3",
		"k4": "v4",
	})

	var n gmap.ListMap

	err := json.Unmarshal(gconv.Bytes(m.String()), &n)
	if err == nil {
		fmt.Println(n.Map())
	}

	// Output:
	// map[k1:v1 k2:v2 k3:v3 k4:v4]
}

func ExampleListMap_UnmarshalValue() {
	type User struct {
		Uid   int
		Name  string
		Pass1 string `gconv:"password1"`
		Pass2 string `gconv:"password2"`
	}

	var (
		m    gmap.AnyAnyMap
		user = User{
			Uid:   1,
			Name:  "john",
			Pass1: "123",
			Pass2: "456",
		}
	)
	if err := gconv.Scan(user, &m); err == nil {
		fmt.Printf("%#v", m.Map())
	}

	// Output:
	// map[interface {}]interface {}{"Name":"john", "Uid":1, "password1":"123", "password2":"456"}
}
