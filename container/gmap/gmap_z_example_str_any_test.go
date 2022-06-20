// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with gm file,
// You can obtain one at https://github.com/gogf/gf.

package gmap_test

import (
	"fmt"
	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/util/gconv"

	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/frame/g"
)

func ExampleStrAnyMap_Iterator() {
	m := gmap.NewStrAnyMap()
	for i := 1; i <= 10; i++ {
		m.Set(gconv.String(i), i*2)
	}

	var totalValue int
	m.Iterator(func(k string, v interface{}) bool {
		totalValue += v.(int)

		return totalValue < 50
	})

	fmt.Println("totalValue:", totalValue)

	// May Output:
	// totalValue: 52
}

func ExampleStrAnyMap_Clone() {
	m := gmap.NewStrAnyMap()

	m.Set("key1", "val1")
	fmt.Println(m)

	n := m.Clone()
	fmt.Println(n)

	// Output:
	// {"key1":"val1"}
	// {"key1":"val1"}
}

func ExampleStrAnyMap_Map() {
	// non concurrent-safety, a pointer to the underlying data
	m1 := gmap.NewStrAnyMap()
	m1.Set("key1", "val1")
	fmt.Println("m1:", m1)

	n1 := m1.Map()
	fmt.Println("before n1:", n1)
	m1.Set("key1", "val2")
	fmt.Println("after n1:", n1)

	// concurrent-safety, copy of underlying data
	m2 := gmap.NewStrAnyMap(true)
	m2.Set("key1", "val1")
	fmt.Println("m2:", m2)

	n2 := m2.Map()
	fmt.Println("before n2:", n2)
	m2.Set("key1", "val2")
	fmt.Println("after n2:", n2)

	// Output:
	// m1: {"key1":"val1"}
	// before n1: map[key1:val1]
	// after n1: map[key1:val2]
	// m2: {"key1":"val1"}
	// before n2: map[key1:val1]
	// after n2: map[key1:val1]
}

func ExampleStrAnyMap_MapCopy() {
	m := gmap.NewStrAnyMap()

	m.Set("key1", "val1")
	m.Set("key2", "val2")
	fmt.Println(m)

	n := m.MapCopy()
	fmt.Println(n)

	// Output:
	// {"key1":"val1","key2":"val2"}
	// map[key1:val1 key2:val2]
}

func ExampleStrAnyMap_MapStrAny() {
	m := gmap.NewStrAnyMap()
	m.Set("key1", "val1")
	m.Set("key2", "val2")

	n := m.MapStrAny()
	fmt.Printf("%#v", n)

	// Output:
	// map[string]interface {}{"key1":"val1", "key2":"val2"}
}

func ExampleStrAnyMap_FilterEmpty() {
	m := gmap.NewStrAnyMapFrom(g.MapStrAny{
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

func ExampleStrAnyMap_FilterNil() {
	m := gmap.NewStrAnyMapFrom(g.MapStrAny{
		"k1": "",
		"k2": nil,
		"k3": 0,
		"k4": 1,
	})
	m.FilterNil()
	fmt.Printf("%#v", m.Map())

	// Output:
	// map[string]interface {}{"k1":"", "k3":0, "k4":1}
}

func ExampleStrAnyMap_Set() {
	m := gmap.NewStrAnyMap()

	m.Set("key1", "val1")
	fmt.Println(m)

	// Output:
	// {"key1":"val1"}
}

func ExampleStrAnyMap_Sets() {
	m := gmap.NewStrAnyMap()

	addMap := make(map[string]interface{})
	addMap["key1"] = "val1"
	addMap["key2"] = "val2"
	addMap["key3"] = "val3"

	m.Sets(addMap)
	fmt.Println(m)

	// Output:
	// {"key1":"val1","key2":"val2","key3":"val3"}
}

func ExampleStrAnyMap_Search() {
	m := gmap.NewStrAnyMap()

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

func ExampleStrAnyMap_Get() {
	m := gmap.NewStrAnyMap()

	m.Set("key1", "val1")

	fmt.Println("key1 value:", m.Get("key1"))
	fmt.Println("key2 value:", m.Get("key2"))

	// Output:
	// key1 value: val1
	// key2 value: <nil>
}

func ExampleStrAnyMap_Pop() {
	var m gmap.StrAnyMap
	m.Sets(g.MapStrAny{
		"k1": "v1",
		"k2": "v2",
		"k3": "v3",
		"k4": "v4",
	})

	fmt.Println(m.Pop())

	// May Output:
	// k1 v1
}

func ExampleStrAnyMap_Pops() {
	var m gmap.StrAnyMap
	m.Sets(g.MapStrAny{
		"k1": "v1",
		"k2": "v2",
		"k3": "v3",
		"k4": "v4",
	})
	fmt.Println(m.Pops(-1))
	fmt.Println("size:", m.Size())

	m.Sets(g.MapStrAny{
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

func ExampleStrAnyMap_GetOrSet() {
	m := gmap.NewStrAnyMap()
	m.Set("key1", "val1")

	fmt.Println(m.GetOrSet("key1", "NotExistValue"))
	fmt.Println(m.GetOrSet("key2", "val2"))

	// Output:
	// val1
	// val2
}

func ExampleStrAnyMap_GetOrSetFunc() {
	m := gmap.NewStrAnyMap()
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

func ExampleStrAnyMap_GetOrSetFuncLock() {
	m := gmap.NewStrAnyMap()
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

func ExampleStrAnyMap_GetVar() {
	m := gmap.NewStrAnyMap()
	m.Set("key1", "val1")

	fmt.Println(m.GetVar("key1"))
	fmt.Println(m.GetVar("key2").IsNil())

	// Output:
	// val1
	// true
}

func ExampleStrAnyMap_GetVarOrSet() {
	m := gmap.NewStrAnyMap()
	m.Set("key1", "val1")

	fmt.Println(m.GetVarOrSet("key1", "NotExistValue"))
	fmt.Println(m.GetVarOrSet("key2", "val2"))

	// Output:
	// val1
	// val2
}

func ExampleStrAnyMap_GetVarOrSetFunc() {
	m := gmap.NewStrAnyMap()
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

func ExampleStrAnyMap_GetVarOrSetFuncLock() {
	m := gmap.NewStrAnyMap()
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

func ExampleStrAnyMap_SetIfNotExist() {
	var m gmap.StrAnyMap
	fmt.Println(m.SetIfNotExist("k1", "v1"))
	fmt.Println(m.SetIfNotExist("k1", "v2"))
	fmt.Println(m.Map())

	// Output:
	// true
	// false
	// map[k1:v1]
}

func ExampleStrAnyMap_SetIfNotExistFunc() {
	var m gmap.StrAnyMap
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

func ExampleStrAnyMap_SetIfNotExistFuncLock() {
	var m gmap.StrAnyMap
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

func ExampleStrAnyMap_Remove() {
	var m gmap.StrAnyMap
	m.Set("k1", "v1")

	fmt.Println(m.Remove("k1"))
	fmt.Println(m.Remove("k2"))
	fmt.Println(m.Size())

	// Output:
	// v1
	// <nil>
	// 0
}

func ExampleStrAnyMap_Removes() {
	var m gmap.StrAnyMap
	m.Sets(g.MapStrAny{
		"k1": "v1",
		"k2": "v2",
		"k3": "v3",
		"k4": "v4",
	})

	removeList := make([]string, 2)
	removeList = append(removeList, "k1")
	removeList = append(removeList, "k2")

	m.Removes(removeList)

	fmt.Println(m.Map())

	// Output:
	// map[k3:v3 k4:v4]
}

func ExampleStrAnyMap_Keys() {
	var m gmap.StrAnyMap
	m.Sets(g.MapStrAny{
		"k1": "v1",
		"k2": "v2",
		"k3": "v3",
		"k4": "v4",
	})
	fmt.Println(m.Keys())

	// May Output:
	// [k1 k2 k3 k4]
}

func ExampleStrAnyMap_Values() {
	var m gmap.StrAnyMap
	m.Sets(g.MapStrAny{
		"k1": "v1",
		"k2": "v2",
		"k3": "v3",
		"k4": "v4",
	})
	fmt.Println(m.Values())

	// May Output:
	// [v1 v2 v3 v4]
}

func ExampleStrAnyMap_Contains() {
	var m gmap.StrAnyMap
	m.Sets(g.MapStrAny{
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

func ExampleStrAnyMap_Size() {
	var m gmap.StrAnyMap
	m.Sets(g.MapStrAny{
		"k1": "v1",
		"k2": "v2",
		"k3": "v3",
		"k4": "v4",
	})

	fmt.Println(m.Size())

	// Output:
	// 4
}

func ExampleStrAnyMap_IsEmpty() {
	var m gmap.StrAnyMap
	fmt.Println(m.IsEmpty())

	m.Set("k1", "v1")
	fmt.Println(m.IsEmpty())

	// Output:
	// true
	// false
}

func ExampleStrAnyMap_Clear() {
	var m gmap.StrAnyMap
	m.Sets(g.MapStrAny{
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

func ExampleStrAnyMap_Replace() {
	var m gmap.StrAnyMap
	m.Sets(g.MapStrAny{
		"k1": "v1",
	})

	var n gmap.StrAnyMap
	n.Sets(g.MapStrAny{
		"k2": "v2",
	})

	fmt.Println(m.Map())

	m.Replace(n.Map())
	fmt.Println(m.Map())

	n.Set("k2", "v1")
	fmt.Println(m.Map())

	// Output:
	// map[k1:v1]
	// map[k2:v2]
	// map[k2:v1]
}

func ExampleStrAnyMap_LockFunc() {
	var m gmap.StrAnyMap
	m.Sets(g.MapStrAny{
		"k1": 1,
		"k2": 2,
		"k3": 3,
		"k4": 4,
	})

	m.LockFunc(func(m map[string]interface{}) {
		totalValue := 0
		for _, v := range m {
			totalValue += v.(int)
		}
		fmt.Println("totalValue:", totalValue)
	})

	// Output:
	// totalValue: 10
}

func ExampleStrAnyMap_RLockFunc() {
	var m gmap.StrAnyMap
	m.Sets(g.MapStrAny{
		"k1": 1,
		"k2": 2,
		"k3": 3,
		"k4": 4,
	})

	m.RLockFunc(func(m map[string]interface{}) {
		totalValue := 0
		for _, v := range m {
			totalValue += v.(int)
		}
		fmt.Println("totalValue:", totalValue)
	})

	// Output:
	// totalValue: 10
}

func ExampleStrAnyMap_Flip() {
	var m gmap.StrAnyMap
	m.Sets(g.MapStrAny{
		"k1": "v1",
	})
	m.Flip()
	fmt.Println(m.Map())

	// Output:
	// map[v1:k1]
}

func ExampleStrAnyMap_Merge() {
	var m1, m2 gmap.StrAnyMap
	m1.Set("key1", "val1")
	m2.Set("key2", "val2")
	m1.Merge(&m2)
	fmt.Println(m1.Map())

	// May Output:
	// map[key1:val1 key2:val2]
}

func ExampleStrAnyMap_String() {
	var m gmap.StrAnyMap
	m.Sets(g.MapStrAny{
		"k1": "v1",
	})

	fmt.Println(m.String())

	var m1 *gmap.StrAnyMap = nil
	fmt.Println(len(m1.String()))

	// Output:
	// {"k1":"v1"}
	// 0
}

func ExampleStrAnyMap_MarshalJSON() {
	var m gmap.StrAnyMap
	m.Sets(g.MapStrAny{
		"k1": "v1",
		"k2": "v2",
		"k3": "v3",
		"k4": "v4",
	})

	bytes, err := json.Marshal(&m)
	if err == nil {
		fmt.Println(gconv.String(bytes))
	}

	// Output:
	// {"k1":"v1","k2":"v2","k3":"v3","k4":"v4"}
}

func ExampleStrAnyMap_UnmarshalJSON() {
	var m gmap.StrAnyMap
	m.Sets(g.MapStrAny{
		"k1": "v1",
		"k2": "v2",
		"k3": "v3",
		"k4": "v4",
	})

	var n gmap.StrAnyMap

	err := json.Unmarshal(gconv.Bytes(m.String()), &n)
	if err == nil {
		fmt.Println(n.Map())
	}

	// Output:
	// map[k1:v1 k2:v2 k3:v3 k4:v4]
}

func ExampleStrAnyMap_UnmarshalValue() {
	var m gmap.StrAnyMap

	goWeb := map[string]interface{}{
		"goframe": "https://goframe.org",
		"gin":     "https://gin-gonic.com/",
		"echo":    "https://echo.labstack.com/",
	}

	if err := gconv.Scan(goWeb, &m); err == nil {
		fmt.Printf("%#v", m.Map())
	}
	// Output:
	// map[string]interface {}{"echo":"https://echo.labstack.com/", "gin":"https://gin-gonic.com/", "goframe":"https://goframe.org"}
}
