// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with gm file,
// You can obtain one at https://github.com/gogf/gf.

package gmap_test

import (
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/util/gconv"

	"github.com/gogf/gf/v2/container/gmap"
)

func ExampleStrIntMap_Iterator() {
	m := gmap.NewStrIntMap()
	for i := 0; i < 10; i++ {
		m.Set(gconv.String(i), i*2)
	}

	var totalValue int
	m.Iterator(func(k string, v int) bool {
		totalValue += v

		return totalValue < 50
	})

	fmt.Println("totalValue:", totalValue)

	// May Output:
	// totalValue: 52
}

func ExampleStrIntMap_Clone() {
	m := gmap.NewStrIntMap()

	m.Set("key1", 1)
	fmt.Println(m)

	n := m.Clone()
	fmt.Println(n)

	// Output:
	// {"key1":1}
	// {"key1":1}
}

func ExampleStrIntMap_Map() {
	// non concurrent-safety, a pointer to the underlying data
	m1 := gmap.NewStrIntMap()
	m1.Set("key1", 1)
	fmt.Println("m1:", m1)

	n1 := m1.Map()
	fmt.Println("before n1:", n1)
	m1.Set("key1", 2)
	fmt.Println("after n1:", n1)

	// concurrent-safety, copy of underlying data
	m2 := gmap.NewStrIntMap(true)
	m2.Set("key1", 1)
	fmt.Println("m2:", m2)

	n2 := m2.Map()
	fmt.Println("before n2:", n2)
	m2.Set("key1", 2)
	fmt.Println("after n2:", n2)

	// Output:
	// m1: {"key1":1}
	// before n1: map[key1:1]
	// after n1: map[key1:2]
	// m2: {"key1":1}
	// before n2: map[key1:1]
	// after n2: map[key1:1]
}

func ExampleStrIntMap_MapCopy() {
	m := gmap.NewStrIntMap()

	m.Set("key1", 1)
	m.Set("key2", 2)
	fmt.Println(m)

	n := m.MapCopy()
	fmt.Println(n)

	// Output:
	// {"key1":1,"key2":2}
	// map[key1:1 key2:2]
}

func ExampleStrIntMap_MapStrAny() {
	m := gmap.NewStrIntMap()
	m.Set("key1", 1)
	m.Set("key2", 2)

	n := m.MapStrAny()
	fmt.Printf("%#v", n)

	// Output:
	// map[string]interface {}{"key1":1, "key2":2}
}

func ExampleStrIntMap_FilterEmpty() {
	m := gmap.NewStrIntMapFrom(g.MapStrInt{
		"k1": 0,
		"k2": 1,
	})
	m.FilterEmpty()
	fmt.Println(m.Map())

	// Output:
	// map[k2:1]
}

func ExampleStrIntMap_Set() {
	m := gmap.NewStrIntMap()

	m.Set("key1", 1)
	fmt.Println(m)

	// Output:
	// {"key1":1}
}

func ExampleStrIntMap_Sets() {
	m := gmap.NewStrIntMap()

	addMap := make(map[string]int)
	addMap["key1"] = 1
	addMap["key2"] = 2
	addMap["key3"] = 3

	m.Sets(addMap)
	fmt.Println(m)

	// Output:
	// {"key1":1,"key2":2,"key3":3}
}

func ExampleStrIntMap_Search() {
	m := gmap.NewStrIntMap()

	m.Set("key1", 1)

	value, found := m.Search("key1")
	if found {
		fmt.Println("find key1 value:", value)
	}

	value, found = m.Search("key2")
	if !found {
		fmt.Println("key2 not find")
	}

	// Output:
	// find key1 value: 1
	// key2 not find
}

func ExampleStrIntMap_Get() {
	m := gmap.NewStrIntMap()

	m.Set("key1", 1)

	fmt.Println("key1 value:", m.Get("key1"))
	fmt.Println("key2 value:", m.Get("key2"))

	// Output:
	// key1 value: 1
	// key2 value: 0
}

func ExampleStrIntMap_Pop() {
	var m gmap.StrIntMap
	m.Sets(g.MapStrInt{
		"k1": 1,
		"k2": 2,
		"k3": 3,
		"k4": 4,
	})

	fmt.Println(m.Pop())

	// May Output:
	// k1 1
}

func ExampleStrIntMap_Pops() {
	var m gmap.StrIntMap
	m.Sets(g.MapStrInt{
		"k1": 1,
		"k2": 2,
		"k3": 3,
		"k4": 4,
	})
	fmt.Println(m.Pops(-1))
	fmt.Println("size:", m.Size())

	m.Sets(g.MapStrInt{
		"k1": 1,
		"k2": 2,
		"k3": 3,
		"k4": 4,
	})
	fmt.Println(m.Pops(2))
	fmt.Println("size:", m.Size())

	// May Output:
	// map[k1:1 k2:2 k3:3 k4:4]
	// size: 0
	// map[k1:1 k2:2]
	// size: 2
}

func ExampleStrIntMap_GetOrSet() {
	m := gmap.NewStrIntMap()
	m.Set("key1", 1)

	fmt.Println(m.GetOrSet("key1", 0))
	fmt.Println(m.GetOrSet("key2", 2))

	// Output:
	// 1
	// 2
}

func ExampleStrIntMap_GetOrSetFunc() {
	m := gmap.NewStrIntMap()
	m.Set("key1", 1)

	fmt.Println(m.GetOrSetFunc("key1", func() int {
		return 0
	}))
	fmt.Println(m.GetOrSetFunc("key2", func() int {
		return 0
	}))

	// Output:
	// 1
	// 0
}

func ExampleStrIntMap_GetOrSetFuncLock() {
	m := gmap.NewStrIntMap()
	m.Set("key1", 1)

	fmt.Println(m.GetOrSetFuncLock("key1", func() int {
		return 0
	}))
	fmt.Println(m.GetOrSetFuncLock("key2", func() int {
		return 0
	}))

	// Output:
	// 1
	// 0
}

func ExampleStrIntMap_SetIfNotExist() {
	var m gmap.StrIntMap
	fmt.Println(m.SetIfNotExist("k1", 1))
	fmt.Println(m.SetIfNotExist("k1", 2))
	fmt.Println(m.Map())

	// Output:
	// true
	// false
	// map[k1:1]
}

func ExampleStrIntMap_SetIfNotExistFunc() {
	var m gmap.StrIntMap
	fmt.Println(m.SetIfNotExistFunc("k1", func() int {
		return 1
	}))
	fmt.Println(m.SetIfNotExistFunc("k1", func() int {
		return 2
	}))
	fmt.Println(m.Map())

	// Output:
	// true
	// false
	// map[k1:1]
}

func ExampleStrIntMap_SetIfNotExistFuncLock() {
	var m gmap.StrIntMap
	fmt.Println(m.SetIfNotExistFuncLock("k1", func() int {
		return 1
	}))
	fmt.Println(m.SetIfNotExistFuncLock("k1", func() int {
		return 2
	}))
	fmt.Println(m.Map())

	// Output:
	// true
	// false
	// map[k1:1]
}

func ExampleStrIntMap_Remove() {
	var m gmap.StrIntMap
	m.Set("k1", 1)

	fmt.Println(m.Remove("k1"))
	fmt.Println(m.Remove("k2"))
	fmt.Println(m.Size())

	// Output:
	// 1
	// 0
	// 0
}

func ExampleStrIntMap_Removes() {
	var m gmap.StrIntMap
	m.Sets(g.MapStrInt{
		"k1": 1,
		"k2": 2,
		"k3": 3,
		"k4": 4,
	})

	removeList := make([]string, 2)
	removeList = append(removeList, "k1")
	removeList = append(removeList, "k2")

	m.Removes(removeList)

	fmt.Println(m.Map())

	// Output:
	// map[k3:3 k4:4]
}

func ExampleStrIntMap_Keys() {
	var m gmap.StrIntMap
	m.Sets(g.MapStrInt{
		"k1": 1,
		"k2": 2,
		"k3": 3,
		"k4": 4,
	})
	fmt.Println(m.Keys())

	// May Output:
	// [k1 k2 k3 k4]
}

func ExampleStrIntMap_Values() {
	var m gmap.StrIntMap
	m.Sets(g.MapStrInt{
		"k1": 1,
		"k2": 2,
		"k3": 3,
		"k4": 4,
	})
	fmt.Println(m.Values())

	// May Output:
	// [1 2 3 4]
}

func ExampleStrIntMap_Contains() {
	var m gmap.StrIntMap
	m.Sets(g.MapStrInt{
		"k1": 1,
		"k2": 2,
		"k3": 3,
		"k4": 4,
	})

	fmt.Println(m.Contains("k1"))
	fmt.Println(m.Contains("k5"))

	// Output:
	// true
	// false
}

func ExampleStrIntMap_Size() {
	var m gmap.StrIntMap
	m.Sets(g.MapStrInt{
		"k1": 1,
		"k2": 2,
		"k3": 3,
		"k4": 4,
	})

	fmt.Println(m.Size())

	// Output:
	// 4
}

func ExampleStrIntMap_IsEmpty() {
	var m gmap.StrIntMap
	fmt.Println(m.IsEmpty())

	m.Set("k1", 1)
	fmt.Println(m.IsEmpty())

	// Output:
	// true
	// false
}

func ExampleStrIntMap_Clear() {
	var m gmap.StrIntMap
	m.Sets(g.MapStrInt{
		"k1": 1,
		"k2": 2,
		"k3": 3,
		"k4": 4,
	})

	m.Clear()

	fmt.Println(m.Map())

	// Output:
	// map[]
}

func ExampleStrIntMap_Replace() {
	var m gmap.StrIntMap
	m.Sets(g.MapStrInt{
		"k1": 1,
	})

	var n gmap.StrIntMap
	n.Sets(g.MapStrInt{
		"k2": 2,
	})

	fmt.Println(m.Map())

	m.Replace(n.Map())
	fmt.Println(m.Map())

	n.Set("k2", 1)
	fmt.Println(m.Map())

	// Output:
	// map[k1:1]
	// map[k2:2]
	// map[k2:1]
}

func ExampleStrIntMap_LockFunc() {
	var m gmap.StrIntMap
	m.Sets(g.MapStrInt{
		"k1": 1,
		"k2": 2,
		"k3": 3,
		"k4": 4,
	})

	m.LockFunc(func(m map[string]int) {
		totalValue := 0
		for _, v := range m {
			totalValue += v
		}
		fmt.Println("totalValue:", totalValue)
	})

	// Output:
	// totalValue: 10
}

func ExampleStrIntMap_RLockFunc() {
	var m gmap.StrIntMap
	m.Sets(g.MapStrInt{
		"k1": 1,
		"k2": 2,
		"k3": 3,
		"k4": 4,
	})

	m.RLockFunc(func(m map[string]int) {
		totalValue := 0
		for _, v := range m {
			totalValue += v
		}
		fmt.Println("totalValue:", totalValue)
	})

	// Output:
	// totalValue: 10
}

func ExampleStrIntMap_Flip() {
	var m gmap.StrIntMap
	m.Sets(g.MapStrInt{
		"k1": 1,
	})
	m.Flip()
	fmt.Println(m.Map())

	var n gmap.StrIntMap
	n.Sets(g.MapStrInt{
		"11": 1,
	})
	n.Flip()
	fmt.Println(n.Map())

	// Output:
	// map[1:0]
	// map[1:11]
}

func ExampleStrIntMap_Merge() {
	var m1, m2 gmap.StrIntMap
	m1.Set("key1", 1)
	m2.Set("key2", 2)
	m1.Merge(&m2)
	fmt.Println(m1.Map())

	// May Output:
	// map[key1:1 key2:2]
}

func ExampleStrIntMap_String() {
	var m gmap.StrIntMap
	m.Sets(g.MapStrInt{
		"k1": 1,
	})

	fmt.Println(m.String())

	// Output:
	// {"k1":1}
}

func ExampleStrIntMap_MarshalJSON() {
	var m gmap.StrIntMap
	m.Sets(g.MapStrInt{
		"k1": 1,
		"k2": 2,
		"k3": 3,
		"k4": 4,
	})

	bytes, err := json.Marshal(&m)
	if err == nil {
		fmt.Println(gconv.String(bytes))
	}

	// Output:
	// {"k1":1,"k2":2,"k3":3,"k4":4}
}

func ExampleStrIntMap_UnmarshalJSON() {
	var m gmap.StrIntMap
	m.Sets(g.MapStrInt{
		"k1": 1,
		"k2": 2,
		"k3": 3,
		"k4": 4,
	})

	var n gmap.StrIntMap

	err := json.Unmarshal(gconv.Bytes(m.String()), &n)
	if err == nil {
		fmt.Println(n.Map())
	}

	// Output:
	// map[k1:1 k2:2 k3:3 k4:4]
}

func ExampleStrIntMap_UnmarshalValue() {
	var m gmap.StrIntMap
	m.Sets(g.MapStrInt{
		"goframe": 1,
		"gin":     2,
		"echo":    3,
	})

	var goweb map[string]int

	err := gconv.Scan(m.String(), &goweb)
	if err == nil {
		fmt.Printf("%#v", goweb)
	}
	// Output:
	// map[string]int{"echo":3, "gin":2, "goframe":1}
}
