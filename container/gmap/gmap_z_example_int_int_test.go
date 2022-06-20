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

func ExampleIntIntMap_Iterator() {
	m := gmap.NewIntIntMap()
	for i := 0; i < 10; i++ {
		m.Set(i, i*2)
	}

	var totalKey, totalValue int
	m.Iterator(func(k int, v int) bool {
		totalKey += k
		totalValue += v

		return totalKey < 10
	})

	fmt.Println("totalKey:", totalKey)
	fmt.Println("totalValue:", totalValue)

	// May Output:
	// totalKey: 11
	// totalValue: 22
}

func ExampleIntIntMap_Clone() {
	m := gmap.NewIntIntMap()

	m.Set(1, 1)
	fmt.Println(m)

	n := m.Clone()
	fmt.Println(n)

	// Output:
	// {"1":1}
	// {"1":1}
}

func ExampleIntIntMap_Map() {
	// non concurrent-safety, a pointer to the underlying data
	m1 := gmap.NewIntIntMap()
	m1.Set(1, 1)
	fmt.Println("m1:", m1)

	n1 := m1.Map()
	fmt.Println("before n1:", n1)
	m1.Set(1, 2)
	fmt.Println("after n1:", n1)

	// concurrent-safety, copy of underlying data
	m2 := gmap.New(true)
	m2.Set(1, "1")
	fmt.Println("m2:", m2)

	n2 := m2.Map()
	fmt.Println("before n2:", n2)
	m2.Set(1, "2")
	fmt.Println("after n2:", n2)

	// Output:
	// m1: {"1":1}
	// before n1: map[1:1]
	// after n1: map[1:2]
	// m2: {"1":"1"}
	// before n2: map[1:1]
	// after n2: map[1:1]
}

func ExampleIntIntMap_MapCopy() {
	m := gmap.NewIntIntMap()

	m.Set(1, 1)
	m.Set(2, 2)
	fmt.Println(m)

	n := m.MapCopy()
	fmt.Println(n)

	// Output:
	// {"1":1,"2":2}
	// map[1:1 2:2]
}

func ExampleIntIntMap_MapStrAny() {
	m := gmap.NewIntIntMap()
	m.Set(1001, 1)
	m.Set(1002, 2)

	n := m.MapStrAny()
	fmt.Printf("%#v", n)

	// Output:
	// map[string]interface {}{"1001":1, "1002":2}
}

func ExampleIntIntMap_FilterEmpty() {
	m := gmap.NewIntIntMapFrom(g.MapIntInt{
		1: 0,
		2: 1,
	})
	m.FilterEmpty()
	fmt.Println(m.Map())

	// Output:
	// map[2:1]
}

func ExampleIntIntMap_Set() {
	m := gmap.NewIntIntMap()

	m.Set(1, 1)
	fmt.Println(m)

	// Output:
	// {"1":1}
}

func ExampleIntIntMap_Sets() {
	m := gmap.NewIntIntMap()

	addMap := make(map[int]int)
	addMap[1] = 1
	addMap[2] = 12
	addMap[3] = 123

	m.Sets(addMap)
	fmt.Println(m)

	// Output:
	// {"1":1,"2":12,"3":123}
}

func ExampleIntIntMap_Search() {
	m := gmap.NewIntIntMap()

	m.Set(1, 1)

	value, found := m.Search(1)
	if found {
		fmt.Println("find key1 value:", value)
	}

	value, found = m.Search(2)
	if !found {
		fmt.Println("key2 not find")
	}

	// Output:
	// find key1 value: 1
	// key2 not find
}

func ExampleIntIntMap_Get() {
	m := gmap.NewIntIntMap()

	m.Set(1, 1)

	fmt.Println("key1 value:", m.Get(1))
	fmt.Println("key2 value:", m.Get(2))

	// Output:
	// key1 value: 1
	// key2 value: 0
}

func ExampleIntIntMap_Pop() {
	var m gmap.IntIntMap
	m.Sets(g.MapIntInt{
		1: 1,
		2: 2,
		3: 3,
		4: 4,
	})

	fmt.Println(m.Pop())

	// May Output:
	// 1 1
}

func ExampleIntIntMap_Pops() {
	var m gmap.IntIntMap
	m.Sets(g.MapIntInt{
		1: 1,
		2: 2,
		3: 3,
		4: 4,
	})
	fmt.Println(m.Pops(-1))
	fmt.Println("size:", m.Size())

	m.Sets(g.MapIntInt{
		1: 1,
		2: 2,
		3: 3,
		4: 4,
	})
	fmt.Println(m.Pops(2))
	fmt.Println("size:", m.Size())

	// May Output:
	// map[1:1 2:2 3:3 4:4]
	// size: 0
	// map[1:1 2:2]
	// size: 2
}

func ExampleIntIntMap_GetOrSet() {
	m := gmap.NewIntIntMap()
	m.Set(1, 1)

	fmt.Println(m.GetOrSet(1, 0))
	fmt.Println(m.GetOrSet(2, 2))

	// Output:
	// 1
	// 2
}

func ExampleIntIntMap_GetOrSetFunc() {
	m := gmap.NewIntIntMap()
	m.Set(1, 1)

	fmt.Println(m.GetOrSetFunc(1, func() int {
		return 0
	}))
	fmt.Println(m.GetOrSetFunc(2, func() int {
		return 0
	}))

	// Output:
	// 1
	// 0
}

func ExampleIntIntMap_GetOrSetFuncLock() {
	m := gmap.NewIntIntMap()
	m.Set(1, 1)

	fmt.Println(m.GetOrSetFuncLock(1, func() int {
		return 0
	}))
	fmt.Println(m.GetOrSetFuncLock(2, func() int {
		return 0
	}))

	// Output:
	// 1
	// 0
}

func ExampleIntIntMap_SetIfNotExist() {
	var m gmap.IntIntMap
	fmt.Println(m.SetIfNotExist(1, 1))
	fmt.Println(m.SetIfNotExist(1, 2))
	fmt.Println(m.Map())

	// Output:
	// true
	// false
	// map[1:1]
}

func ExampleIntIntMap_SetIfNotExistFunc() {
	var m gmap.IntIntMap
	fmt.Println(m.SetIfNotExistFunc(1, func() int {
		return 1
	}))
	fmt.Println(m.SetIfNotExistFunc(1, func() int {
		return 2
	}))
	fmt.Println(m.Map())

	// Output:
	// true
	// false
	// map[1:1]
}

func ExampleIntIntMap_SetIfNotExistFuncLock() {
	var m gmap.IntIntMap
	fmt.Println(m.SetIfNotExistFuncLock(1, func() int {
		return 1
	}))
	fmt.Println(m.SetIfNotExistFuncLock(1, func() int {
		return 2
	}))
	fmt.Println(m.Map())

	// Output:
	// true
	// false
	// map[1:1]
}

func ExampleIntIntMap_Remove() {
	var m gmap.IntIntMap
	m.Set(1, 1)

	fmt.Println(m.Remove(1))
	fmt.Println(m.Remove(2))
	fmt.Println(m.Size())

	// Output:
	// 1
	// 0
	// 0
}

func ExampleIntIntMap_Removes() {
	var m gmap.IntIntMap
	m.Sets(g.MapIntInt{
		1: 1,
		2: 2,
		3: 3,
		4: 4,
	})

	removeList := make([]int, 2)
	removeList = append(removeList, 1)
	removeList = append(removeList, 2)

	m.Removes(removeList)

	fmt.Println(m.Map())

	// Output:
	// map[3:3 4:4]
}

func ExampleIntIntMap_Keys() {
	var m gmap.IntIntMap
	m.Sets(g.MapIntInt{
		1: 1,
		2: 2,
		3: 3,
		4: 4,
	})
	fmt.Println(m.Keys())

	// May Output:
	// [1 2 3 4]
}

func ExampleIntIntMap_Values() {
	var m gmap.IntIntMap
	m.Sets(g.MapIntInt{
		1: 1,
		2: 2,
		3: 3,
		4: 4,
	})
	fmt.Println(m.Values())

	// May Output:
	// [1 v2 v3 4]
}

func ExampleIntIntMap_Contains() {
	var m gmap.IntIntMap
	m.Sets(g.MapIntInt{
		1: 1,
		2: 2,
		3: 3,
		4: 4,
	})

	fmt.Println(m.Contains(1))
	fmt.Println(m.Contains(5))

	// Output:
	// true
	// false
}

func ExampleIntIntMap_Size() {
	var m gmap.IntIntMap
	m.Sets(g.MapIntInt{
		1: 1,
		2: 2,
		3: 3,
		4: 4,
	})

	fmt.Println(m.Size())

	// Output:
	// 4
}

func ExampleIntIntMap_IsEmpty() {
	var m gmap.IntIntMap
	fmt.Println(m.IsEmpty())

	m.Set(1, 1)
	fmt.Println(m.IsEmpty())

	// Output:
	// true
	// false
}

func ExampleIntIntMap_Clear() {
	var m gmap.IntIntMap
	m.Sets(g.MapIntInt{
		1: 1,
		2: 2,
		3: 3,
		4: 4,
	})

	m.Clear()

	fmt.Println(m.Map())

	// Output:
	// map[]
}

func ExampleIntIntMap_Replace() {
	var m gmap.IntIntMap
	m.Sets(g.MapIntInt{
		1: 1,
	})

	var n gmap.IntIntMap
	n.Sets(g.MapIntInt{
		2: 2,
	})

	fmt.Println(m.Map())

	m.Replace(n.Map())
	fmt.Println(m.Map())

	n.Set(2, 1)
	fmt.Println(m.Map())

	// Output:
	// map[1:1]
	// map[2:2]
	// map[2:1]
}

func ExampleIntIntMap_LockFunc() {
	var m gmap.IntIntMap
	m.Sets(g.MapIntInt{
		1: 1,
		2: 2,
		3: 3,
		4: 4,
	})

	m.LockFunc(func(m map[int]int) {
		totalValue := 0
		for _, v := range m {
			totalValue += v
		}
		fmt.Println("totalValue:", totalValue)
	})

	// Output:
	// totalValue: 10
}

func ExampleIntIntMap_RLockFunc() {
	var m gmap.IntIntMap
	m.Sets(g.MapIntInt{
		1: 1,
		2: 2,
		3: 3,
		4: 4,
	})

	m.RLockFunc(func(m map[int]int) {
		totalValue := 0
		for _, v := range m {
			totalValue += v
		}
		fmt.Println("totalValue:", totalValue)
	})

	// Output:
	// totalValue: 10
}

func ExampleIntIntMap_Flip() {
	var m gmap.IntIntMap
	m.Sets(g.MapIntInt{
		1: 10,
	})
	m.Flip()
	fmt.Println(m.Map())

	// Output:
	// map[10:1]
}

func ExampleIntIntMap_Merge() {
	var m1, m2 gmap.Map
	m1.Set(1, "1")
	m2.Set(2, "2")
	m1.Merge(&m2)
	fmt.Println(m1.Map())

	// May Output:
	// map[key1:1 key2:2]
}

func ExampleIntIntMap_String() {
	var m gmap.IntIntMap
	m.Sets(g.MapIntInt{
		1: 1,
	})

	fmt.Println(m.String())

	var m1 *gmap.IntIntMap = nil
	fmt.Println(len(m1.String()))

	// Output:
	// {"1":1}
	// 0
}

func ExampleIntIntMap_MarshalJSON() {
	var m gmap.IntIntMap
	m.Sets(g.MapIntInt{
		1: 1,
		2: 2,
		3: 3,
		4: 4,
	})

	bytes, err := json.Marshal(&m)
	if err == nil {
		fmt.Println(gconv.String(bytes))
	}

	// Output:
	// {"1":1,"2":2,"3":3,"4":4}
}

func ExampleIntIntMap_UnmarshalJSON() {
	var m gmap.IntIntMap
	m.Sets(g.MapIntInt{
		1: 1,
		2: 2,
		3: 3,
		4: 4,
	})

	var n gmap.Map

	err := json.Unmarshal(gconv.Bytes(m.String()), &n)
	if err == nil {
		fmt.Println(n.Map())
	}

	// Output:
	// map[1:1 2:2 3:3 4:4]
}

func ExampleIntIntMap_UnmarshalValue() {
	var m gmap.IntIntMap

	n := map[int]int{
		1: 1001,
		2: 1002,
		3: 1003,
	}

	if err := gconv.Scan(n, &m); err == nil {
		fmt.Printf("%#v", m.Map())
	}
	// Output:
	// map[int]int{1:1001, 2:1002, 3:1003}
}
