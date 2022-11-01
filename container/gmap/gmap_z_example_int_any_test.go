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

func ExampleIntAnyMap_Iterator() {
	m := gmap.NewIntAnyMap()
	for i := 0; i < 10; i++ {
		m.Set(i, i*2)
	}

	var totalKey, totalValue int
	m.Iterator(func(k int, v interface{}) bool {
		totalKey += k
		totalValue += v.(int)

		return totalKey < 10
	})

	fmt.Println("totalKey:", totalKey)
	fmt.Println("totalValue:", totalValue)

	// May Output:
	// totalKey: 11
	// totalValue: 22
}

func ExampleIntAnyMap_Clone() {
	m := gmap.NewIntAnyMap()

	m.Set(1, "val1")
	fmt.Println(m)

	n := m.Clone()
	fmt.Println(n)

	// Output:
	// {"1":"val1"}
	// {"1":"val1"}
}

func ExampleIntAnyMap_Map() {
	// non concurrent-safety, a pointer to the underlying data
	m1 := gmap.NewIntAnyMap()
	m1.Set(1, "val1")
	fmt.Println("m1:", m1)

	n1 := m1.Map()
	fmt.Println("before n1:", n1)
	m1.Set(1, "val2")
	fmt.Println("after n1:", n1)

	// concurrent-safety, copy of underlying data
	m2 := gmap.New(true)
	m2.Set(1, "val1")
	fmt.Println("m2:", m2)

	n2 := m2.Map()
	fmt.Println("before n2:", n2)
	m2.Set(1, "val2")
	fmt.Println("after n2:", n2)

	// Output:
	// m1: {"1":"val1"}
	// before n1: map[1:val1]
	// after n1: map[1:val2]
	// m2: {"1":"val1"}
	// before n2: map[1:val1]
	// after n2: map[1:val1]
}

func ExampleIntAnyMap_MapCopy() {
	m := gmap.NewIntAnyMap()

	m.Set(1, "val1")
	m.Set(2, "val2")
	fmt.Println(m)

	n := m.MapCopy()
	fmt.Println(n)

	// Output:
	// {"1":"val1","2":"val2"}
	// map[1:val1 2:val2]
}

func ExampleIntAnyMap_MapStrAny() {
	m := gmap.NewIntAnyMap()
	m.Set(1001, "val1")
	m.Set(1002, "val2")

	n := m.MapStrAny()
	fmt.Printf("%#v", n)

	// Output:
	// map[string]interface {}{"1001":"val1", "1002":"val2"}
}

func ExampleIntAnyMap_FilterEmpty() {
	m := gmap.NewIntAnyMapFrom(g.MapIntAny{
		1: "",
		2: nil,
		3: 0,
		4: 1,
	})
	m.FilterEmpty()
	fmt.Println(m.Map())

	// Output:
	// map[4:1]
}

func ExampleIntAnyMap_FilterNil() {
	m := gmap.NewIntAnyMapFrom(g.MapIntAny{
		1: "",
		2: nil,
		3: 0,
		4: 1,
	})
	m.FilterNil()
	fmt.Printf("%#v", m.Map())

	// Output:
	// map[int]interface {}{1:"", 3:0, 4:1}
}

func ExampleIntAnyMap_Set() {
	m := gmap.NewIntAnyMap()

	m.Set(1, "val1")
	fmt.Println(m)

	// Output:
	// {"1":"val1"}
}

func ExampleIntAnyMap_Sets() {
	m := gmap.NewIntAnyMap()

	addMap := make(map[int]interface{})
	addMap[1] = "val1"
	addMap[2] = "val2"
	addMap[3] = "val3"

	m.Sets(addMap)
	fmt.Println(m)

	// Output:
	// {"1":"val1","2":"val2","3":"val3"}
}

func ExampleIntAnyMap_Search() {
	m := gmap.NewIntAnyMap()

	m.Set(1, "val1")

	value, found := m.Search(1)
	if found {
		fmt.Println("find key1 value:", value)
	}

	value, found = m.Search(2)
	if !found {
		fmt.Println("key2 not find")
	}

	// Output:
	// find key1 value: val1
	// key2 not find
}

func ExampleIntAnyMap_Get() {
	m := gmap.NewIntAnyMap()

	m.Set(1, "val1")

	fmt.Println("key1 value:", m.Get(1))
	fmt.Println("key2 value:", m.Get(2))

	// Output:
	// key1 value: val1
	// key2 value: <nil>
}

func ExampleIntAnyMap_Pop() {
	var m gmap.IntAnyMap
	m.Sets(g.MapIntAny{
		1: "v1",
		2: "v2",
		3: "v3",
		4: "v4",
	})

	fmt.Println(m.Pop())

	// May Output:
	// 1 v1
}

func ExampleIntAnyMap_Pops() {
	var m gmap.IntAnyMap
	m.Sets(g.MapIntAny{
		1: "v1",
		2: "v2",
		3: "v3",
		4: "v4",
	})
	fmt.Println(m.Pops(-1))
	fmt.Println("size:", m.Size())

	m.Sets(g.MapIntAny{
		1: "v1",
		2: "v2",
		3: "v3",
		4: "v4",
	})
	fmt.Println(m.Pops(2))
	fmt.Println("size:", m.Size())

	// May Output:
	// map[1:v1 2:v2 3:v3 4:v4]
	// size: 0
	// map[1:v1 2:v2]
	// size: 2
}

func ExampleIntAnyMap_GetOrSet() {
	m := gmap.NewIntAnyMap()
	m.Set(1, "val1")

	fmt.Println(m.GetOrSet(1, "NotExistValue"))
	fmt.Println(m.GetOrSet(2, "val2"))

	// Output:
	// val1
	// val2
}

func ExampleIntAnyMap_GetOrSetFunc() {
	m := gmap.NewIntAnyMap()
	m.Set(1, "val1")

	fmt.Println(m.GetOrSetFunc(1, func() interface{} {
		return "NotExistValue"
	}))
	fmt.Println(m.GetOrSetFunc(2, func() interface{} {
		return "NotExistValue"
	}))

	// Output:
	// val1
	// NotExistValue
}

func ExampleIntAnyMap_GetOrSetFuncLock() {
	m := gmap.NewIntAnyMap()
	m.Set(1, "val1")

	fmt.Println(m.GetOrSetFuncLock(1, func() interface{} {
		return "NotExistValue"
	}))
	fmt.Println(m.GetOrSetFuncLock(2, func() interface{} {
		return "NotExistValue"
	}))

	// Output:
	// val1
	// NotExistValue
}

func ExampleIntAnyMap_GetVar() {
	m := gmap.NewIntAnyMap()
	m.Set(1, "val1")

	fmt.Println(m.GetVar(1))
	fmt.Println(m.GetVar(2).IsNil())

	// Output:
	// val1
	// true
}

func ExampleIntAnyMap_GetVarOrSet() {
	m := gmap.NewIntAnyMap()
	m.Set(1, "val1")

	fmt.Println(m.GetVarOrSet(1, "NotExistValue"))
	fmt.Println(m.GetVarOrSet(2, "val2"))

	// Output:
	// val1
	// val2
}

func ExampleIntAnyMap_GetVarOrSetFunc() {
	m := gmap.NewIntAnyMap()
	m.Set(1, "val1")

	fmt.Println(m.GetVarOrSetFunc(1, func() interface{} {
		return "NotExistValue"
	}))
	fmt.Println(m.GetVarOrSetFunc(2, func() interface{} {
		return "NotExistValue"
	}))

	// Output:
	// val1
	// NotExistValue
}

func ExampleIntAnyMap_GetVarOrSetFuncLock() {
	m := gmap.NewIntAnyMap()
	m.Set(1, "val1")

	fmt.Println(m.GetVarOrSetFuncLock(1, func() interface{} {
		return "NotExistValue"
	}))
	fmt.Println(m.GetVarOrSetFuncLock(2, func() interface{} {
		return "NotExistValue"
	}))

	// Output:
	// val1
	// NotExistValue
}

func ExampleIntAnyMap_SetIfNotExist() {
	var m gmap.IntAnyMap
	fmt.Println(m.SetIfNotExist(1, "v1"))
	fmt.Println(m.SetIfNotExist(1, "v2"))
	fmt.Println(m.Map())

	// Output:
	// true
	// false
	// map[1:v1]
}

func ExampleIntAnyMap_SetIfNotExistFunc() {
	var m gmap.IntAnyMap
	fmt.Println(m.SetIfNotExistFunc(1, func() interface{} {
		return "v1"
	}))
	fmt.Println(m.SetIfNotExistFunc(1, func() interface{} {
		return "v2"
	}))
	fmt.Println(m.Map())

	// Output:
	// true
	// false
	// map[1:v1]
}

func ExampleIntAnyMap_SetIfNotExistFuncLock() {
	var m gmap.IntAnyMap
	fmt.Println(m.SetIfNotExistFuncLock(1, func() interface{} {
		return "v1"
	}))
	fmt.Println(m.SetIfNotExistFuncLock(1, func() interface{} {
		return "v2"
	}))
	fmt.Println(m.Map())

	// Output:
	// true
	// false
	// map[1:v1]
}

func ExampleIntAnyMap_Remove() {
	var m gmap.IntAnyMap
	m.Set(1, "v1")

	fmt.Println(m.Remove(1))
	fmt.Println(m.Remove(2))
	fmt.Println(m.Size())

	// Output:
	// v1
	// <nil>
	// 0
}

func ExampleIntAnyMap_Removes() {
	var m gmap.IntAnyMap
	m.Sets(g.MapIntAny{
		1: "v1",
		2: "v2",
		3: "v3",
		4: "v4",
	})

	removeList := make([]int, 2)
	removeList = append(removeList, 1)
	removeList = append(removeList, 2)

	m.Removes(removeList)

	fmt.Println(m.Map())

	// Output:
	// map[3:v3 4:v4]
}

func ExampleIntAnyMap_Keys() {
	var m gmap.IntAnyMap
	m.Sets(g.MapIntAny{
		1: "v1",
		2: "v2",
		3: "v3",
		4: "v4",
	})
	fmt.Println(m.Keys())

	// May Output:
	// [1 2 3 4]
}

func ExampleIntAnyMap_Values() {
	var m gmap.IntAnyMap
	m.Sets(g.MapIntAny{
		1: "v1",
		2: "v2",
		3: "v3",
		4: "v4",
	})
	fmt.Println(m.Values())

	// May Output:
	// [v1 v2 v3 v4]
}

func ExampleIntAnyMap_Contains() {
	var m gmap.IntAnyMap
	m.Sets(g.MapIntAny{
		1: "v1",
		2: "v2",
		3: "v3",
		4: "v4",
	})

	fmt.Println(m.Contains(1))
	fmt.Println(m.Contains(5))

	// Output:
	// true
	// false
}

func ExampleIntAnyMap_Size() {
	var m gmap.IntAnyMap
	m.Sets(g.MapIntAny{
		1: "v1",
		2: "v2",
		3: "v3",
		4: "v4",
	})

	fmt.Println(m.Size())

	// Output:
	// 4
}

func ExampleIntAnyMap_IsEmpty() {
	var m gmap.IntAnyMap
	fmt.Println(m.IsEmpty())

	m.Set(1, "v1")
	fmt.Println(m.IsEmpty())

	// Output:
	// true
	// false
}

func ExampleIntAnyMap_Clear() {
	var m gmap.IntAnyMap
	m.Sets(g.MapIntAny{
		1: "v1",
		2: "v2",
		3: "v3",
		4: "v4",
	})

	m.Clear()

	fmt.Println(m.Map())

	// Output:
	// map[]
}

func ExampleIntAnyMap_Replace() {
	var m gmap.IntAnyMap
	m.Sets(g.MapIntAny{
		1: "v1",
	})

	var n gmap.IntAnyMap
	n.Sets(g.MapIntAny{
		2: "v2",
	})

	fmt.Println(m.Map())

	m.Replace(n.Map())
	fmt.Println(m.Map())

	n.Set(2, "v1")
	fmt.Println(m.Map())

	// Output:
	// map[1:v1]
	// map[2:v2]
	// map[2:v1]
}

func ExampleIntAnyMap_LockFunc() {
	var m gmap.IntAnyMap
	m.Sets(g.MapIntAny{
		1: 1,
		2: 2,
		3: 3,
		4: 4,
	})

	m.LockFunc(func(m map[int]interface{}) {
		totalValue := 0
		for _, v := range m {
			totalValue += v.(int)
		}
		fmt.Println("totalValue:", totalValue)
	})

	// Output:
	// totalValue: 10
}

func ExampleIntAnyMap_RLockFunc() {
	var m gmap.IntAnyMap
	m.Sets(g.MapIntAny{
		1: 1,
		2: 2,
		3: 3,
		4: 4,
	})

	m.RLockFunc(func(m map[int]interface{}) {
		totalValue := 0
		for _, v := range m {
			totalValue += v.(int)
		}
		fmt.Println("totalValue:", totalValue)
	})

	// Output:
	// totalValue: 10
}

func ExampleIntAnyMap_Flip() {
	var m gmap.IntAnyMap
	m.Sets(g.MapIntAny{
		1: 10,
	})
	m.Flip()
	fmt.Println(m.Map())

	// Output:
	// map[10:1]
}

func ExampleIntAnyMap_Merge() {
	var m1, m2 gmap.Map
	m1.Set(1, "val1")
	m2.Set(2, "val2")
	m1.Merge(&m2)
	fmt.Println(m1.Map())

	// May Output:
	// map[key1:val1 key2:val2]
}

func ExampleIntAnyMap_String() {
	var m gmap.IntAnyMap
	m.Sets(g.MapIntAny{
		1: "v1",
	})

	fmt.Println(m.String())

	var m1 *gmap.IntAnyMap = nil
	fmt.Println(len(m1.String()))

	// Output:
	// {"1":"v1"}
	// 0
}

func ExampleIntAnyMap_MarshalJSON() {
	var m gmap.IntAnyMap
	m.Sets(g.MapIntAny{
		1: "v1",
		2: "v2",
		3: "v3",
		4: "v4",
	})

	bytes, err := json.Marshal(&m)
	if err == nil {
		fmt.Println(gconv.String(bytes))
	}

	// Output:
	// {"1":"v1","2":"v2","3":"v3","4":"v4"}
}

func ExampleIntAnyMap_UnmarshalJSON() {
	var m gmap.IntAnyMap
	m.Sets(g.MapIntAny{
		1: "v1",
		2: "v2",
		3: "v3",
		4: "v4",
	})

	var n gmap.Map

	err := json.Unmarshal(gconv.Bytes(m.String()), &n)
	if err == nil {
		fmt.Println(n.Map())
	}

	// Output:
	// map[1:v1 2:v2 3:v3 4:v4]
}

func ExampleIntAnyMap_UnmarshalValue() {
	var m gmap.IntAnyMap

	goWeb := map[int]interface{}{
		1: "goframe",
		2: "gin",
		3: "echo",
	}

	if err := gconv.Scan(goWeb, &m); err == nil {
		fmt.Printf("%#v", m.Map())
	}

	// Output:
	// map[int]interface {}{1:"goframe", 2:"gin", 3:"echo"}
}
