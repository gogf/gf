// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with gm file,
// You can obtain one at https://github.com/Agogf/gf.

package gtree_test

import (
	"fmt"
	"github.com/gogf/gf/v2/container/gtree"
	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gutil"
)

func ExampleBTree_Clone() {
	b := gtree.NewBTree(3, gutil.ComparatorString)
	for i := 0; i < 6; i++ {
		b.Set("key"+gconv.String(i), "val"+gconv.String(i))
	}

	newTree := b.Clone()

	fmt.Println(newTree.Map())
	fmt.Println(newTree.Size())

	// output:
	// map[key0:val0 key1:val1 key2:val2 key3:val3 key4:val4 key5:val5]
	// 6
}

func ExampleBTree_Set() {
	newTree := gtree.NewBTree(3, gutil.ComparatorString)
	for i := 0; i < 6; i++ {
		newTree.Set("key"+gconv.String(i), "val"+gconv.String(i))
	}

	fmt.Println(newTree.Map())
	fmt.Println(newTree.Size())

	// output:
	// map[key0:val0 key1:val1 key2:val2 key3:val3 key4:val4 key5:val5]
	// 6
}

func ExampleBTree_Sets() {
	newTree := gtree.NewBTree(3, gutil.ComparatorString)

	newTree.Sets(map[interface{}]interface{}{
		"key1": "val1",
		"key2": "val2",
	})

	fmt.Println(newTree.Map())
	fmt.Println(newTree.Size())

	// output:
	// map[key1:val1 key2:val2]
	// 2
}

func ExampleBTree_Get() {
	newTree := gtree.NewBTree(3, gutil.ComparatorString)
	for i := 0; i < 6; i++ {
		newTree.Set("key"+gconv.String(i), "val"+gconv.String(i))
	}

	fmt.Println(newTree.Get("key1"))
	fmt.Println(newTree.Get("key10"))

	// output:
	// val1
	// <nil>
}

func ExampleBTree_GetOrSet() {
	newTree := gtree.NewBTree(3, gutil.ComparatorString)
	for i := 0; i < 6; i++ {
		newTree.Set("key"+gconv.String(i), "val"+gconv.String(i))
	}

	fmt.Println(newTree.GetOrSet("key1", "newVal1"))
	fmt.Println(newTree.GetOrSet("key6", "val6"))

	// output:
	// val1
	// val6
}

func ExampleBTree_GetOrSetFunc() {
	newTree := gtree.NewBTree(3, gutil.ComparatorString)
	for i := 0; i < 6; i++ {
		newTree.Set("key"+gconv.String(i), "val"+gconv.String(i))
	}

	fmt.Println(newTree.GetOrSetFunc("key1", func() interface{} {
		return "newVal1"
	}))
	fmt.Println(newTree.GetOrSetFunc("key6", func() interface{} {
		return "val6"
	}))

	// output:
	// val1
	// val6
}

func ExampleBTree_GetOrSetFuncLock() {
	newTree := gtree.NewBTree(3, gutil.ComparatorString)
	for i := 0; i < 6; i++ {
		newTree.Set("key"+gconv.String(i), "val"+gconv.String(i))
	}

	fmt.Println(newTree.GetOrSetFuncLock("key1", func() interface{} {
		return "newVal1"
	}))
	fmt.Println(newTree.GetOrSetFuncLock("key6", func() interface{} {
		return "val6"
	}))

	// output:
	// val1
	// val6
}

func ExampleBTree_GetVar() {
	newTree := gtree.NewBTree(3, gutil.ComparatorString)
	for i := 0; i < 6; i++ {
		newTree.Set("key"+gconv.String(i), "val"+gconv.String(i))
	}

	fmt.Println(newTree.GetVar("key1"))

	// output:
	// val1
}

func ExampleBTree_GetVarOrSet() {
	newTree := gtree.NewBTree(3, gutil.ComparatorString)
	for i := 0; i < 6; i++ {
		newTree.Set("key"+gconv.String(i), "val"+gconv.String(i))
	}

	fmt.Println(newTree.GetVarOrSet("key1", "newVal1"))
	fmt.Println(newTree.GetVarOrSet("key6", "val6"))

	// output:
	// val1
	// val6
}

func ExampleBTree_GetVarOrSetFunc() {
	newTree := gtree.NewBTree(3, gutil.ComparatorString)
	for i := 0; i < 6; i++ {
		newTree.Set("key"+gconv.String(i), "val"+gconv.String(i))
	}

	fmt.Println(newTree.GetVarOrSetFunc("key1", func() interface{} {
		return "newVal1"
	}))
	fmt.Println(newTree.GetVarOrSetFunc("key6", func() interface{} {
		return "val6"
	}))

	// output:
	// val1
	// val6
}

func ExampleBTree_GetVarOrSetFuncLock() {
	newTree := gtree.NewBTree(3, gutil.ComparatorString)
	for i := 0; i < 6; i++ {
		newTree.Set("key"+gconv.String(i), "val"+gconv.String(i))
	}

	fmt.Println(newTree.GetVarOrSetFuncLock("key1", func() interface{} {
		return "newVal1"
	}))
	fmt.Println(newTree.GetVarOrSetFuncLock("key6", func() interface{} {
		return "val6"
	}))

	// output:
	// val1
	// val6
}

func ExampleBTree_SetIfNotExist() {
	newTree := gtree.NewBTree(3, gutil.ComparatorString)
	for i := 0; i < 6; i++ {
		newTree.Set("key"+gconv.String(i), "val"+gconv.String(i))
	}

	fmt.Println(newTree.SetIfNotExist("key1", "newVal1"))
	fmt.Println(newTree.SetIfNotExist("key6", "val6"))

	// output:
	// false
	// true
}

func ExampleBTree_SetIfNotExistFunc() {
	newTree := gtree.NewBTree(3, gutil.ComparatorString)
	for i := 0; i < 6; i++ {
		newTree.Set("key"+gconv.String(i), "val"+gconv.String(i))
	}

	fmt.Println(newTree.SetIfNotExistFunc("key1", func() interface{} {
		return "newVal1"
	}))
	fmt.Println(newTree.SetIfNotExistFunc("key6", func() interface{} {
		return "val6"
	}))

	// output:
	// false
	// true
}

func ExampleBTree_SetIfNotExistFuncLock() {
	newTree := gtree.NewBTree(3, gutil.ComparatorString)
	for i := 0; i < 6; i++ {
		newTree.Set("key"+gconv.String(i), "val"+gconv.String(i))
	}

	fmt.Println(newTree.SetIfNotExistFuncLock("key1", func() interface{} {
		return "newVal1"
	}))
	fmt.Println(newTree.SetIfNotExistFuncLock("key6", func() interface{} {
		return "val6"
	}))

	// output:
	// false
	// true
}

func ExampleBTree_Contains() {
	newTree := gtree.NewBTree(3, gutil.ComparatorString)
	for i := 0; i < 6; i++ {
		newTree.Set("key"+gconv.String(i), "val"+gconv.String(i))
	}

	fmt.Println(newTree.Contains("key1"))
	fmt.Println(newTree.Contains("key6"))

	// output:
	// true
	// false
}

func ExampleBTree_Remove() {
	newTree := gtree.NewBTree(3, gutil.ComparatorString)
	for i := 0; i < 6; i++ {
		newTree.Set("key"+gconv.String(i), "val"+gconv.String(i))
	}

	fmt.Println(newTree.Remove("key1"))
	fmt.Println(newTree.Remove("key6"))
	fmt.Println(newTree.Map())

	// output:
	// val1
	// <nil>
	// map[key0:val0 key2:val2 key3:val3 key4:val4 key5:val5]
}

func ExampleBTree_Removes() {
	newTree := gtree.NewBTree(3, gutil.ComparatorString)
	for i := 0; i < 6; i++ {
		newTree.Set("key"+gconv.String(i), "val"+gconv.String(i))
	}

	removeKeys := make([]interface{}, 2)
	removeKeys = append(removeKeys, "key1")
	removeKeys = append(removeKeys, "key6")

	newTree.Removes(removeKeys)

	fmt.Println(newTree.Map())

	// output:
	// map[key0:val0 key2:val2 key3:val3 key4:val4 key5:val5]
}

func ExampleBTree_IsEmpty() {
	newTree := gtree.NewBTree(3, gutil.ComparatorString)

	fmt.Println(newTree.IsEmpty())

	for i := 0; i < 6; i++ {
		newTree.Set("key"+gconv.String(i), "val"+gconv.String(i))
	}

	fmt.Println(newTree.IsEmpty())

	// output:
	// true
	// false
}

func ExampleBTree_Size() {
	newTree := gtree.NewBTree(3, gutil.ComparatorString)

	fmt.Println(newTree.Size())

	for i := 0; i < 6; i++ {
		newTree.Set("key"+gconv.String(i), "val"+gconv.String(i))
	}

	fmt.Println(newTree.Size())

	// output:
	// 0
	// 6
}

func ExampleBTree_Keys() {
	newTree := gtree.NewBTree(3, gutil.ComparatorString)
	for i := 6; i > 0; i-- {
		newTree.Set("key"+gconv.String(i), "val"+gconv.String(i))
	}

	fmt.Println(newTree.Keys())

	// output:
	// [key1 key2 key3 key4 key5 key6]
}

func ExampleBTree_Values() {
	newTree := gtree.NewBTree(3, gutil.ComparatorString)
	for i := 6; i > 0; i-- {
		newTree.Set("key"+gconv.String(i), "val"+gconv.String(i))
	}

	fmt.Println(newTree.Values())

	// output:
	// [val1 val2 val3 val4 val5 val6]
}

func ExampleBTree_Map() {
	newTree := gtree.NewBTree(3, gutil.ComparatorString)
	for i := 0; i < 6; i++ {
		newTree.Set("key"+gconv.String(i), "val"+gconv.String(i))
	}

	fmt.Println(newTree.Map())

	// output:
	// map[key0:val0 key1:val1 key2:val2 key3:val3 key4:val4 key5:val5]
}

func ExampleBTree_MapStrAny() {
	newTree := gtree.NewBTree(3, gutil.ComparatorString)
	for i := 0; i < 6; i++ {
		newTree.Set(1000+i, "val"+gconv.String(i))
	}

	fmt.Println(newTree.MapStrAny())

	// output:
	// map[1000:val0 1001:val1 1002:val2 1003:val3 1004:val4 1005:val5]
}

func ExampleBTree_Clear() {
	newTree := gtree.NewBTree(3, gutil.ComparatorString)
	for i := 0; i < 6; i++ {
		newTree.Set(1000+i, "val"+gconv.String(i))
	}
	fmt.Println(newTree.Size())

	newTree.Clear()
	fmt.Println(newTree.Size())

	// output:
	// 6
	// 0
}

func ExampleBTree_Replace() {
	newTree := gtree.NewBTree(3, gutil.ComparatorString)
	for i := 0; i < 6; i++ {
		newTree.Set("key"+gconv.String(i), "val"+gconv.String(i))
	}

	fmt.Println(newTree.Map())

	data := map[interface{}]interface{}{
		"newKey0": "newVal0",
		"newKey1": "newVal1",
		"newKey2": "newVal2",
	}

	newTree.Replace(data)

	fmt.Println(newTree.Map())

	// output:
	// map[key0:val0 key1:val1 key2:val2 key3:val3 key4:val4 key5:val5]
	// map[newKey0:newVal0 newKey1:newVal1 newKey2:newVal2]
}

func ExampleBTree_Height() {
	newTree := gtree.NewBTree(3, gutil.ComparatorInt)
	for i := 0; i < 100; i++ {
		newTree.Set(i, i)
	}
	fmt.Println(newTree.Height())

	// output:
	// 6
}

func ExampleBTree_Left() {
	newTree := gtree.NewBTree(3, gutil.ComparatorInt)
	for i := 1; i < 100; i++ {
		newTree.Set(i, i)
	}
	fmt.Println(newTree.Left().Key, newTree.Left().Value)

	emptyTree := gtree.NewBTree(3, gutil.ComparatorInt)
	fmt.Println(emptyTree.Left())

	// output:
	// 1 1
	// <nil>
}

func ExampleBTree_Right() {
	newTree := gtree.NewBTree(3, gutil.ComparatorInt)
	for i := 1; i < 100; i++ {
		newTree.Set(i, i)
	}
	fmt.Println(newTree.Right().Key, newTree.Right().Value)

	emptyTree := gtree.NewBTree(3, gutil.ComparatorInt)
	fmt.Println(emptyTree.Left())

	// output:
	// 99 99
	// <nil>
}

func ExampleBTree_String() {
	newTree := gtree.NewBTree(3, gutil.ComparatorString)
	for i := 0; i < 6; i++ {
		newTree.Set("key"+gconv.String(i), "val"+gconv.String(i))
	}

	fmt.Println(newTree.String())

	// output:
	// key0
	// key1
	//     key2
	// key3
	//     key4
	//     key5
}

func ExampleBTree_Search() {
	newTree := gtree.NewBTree(3, gutil.ComparatorString)
	for i := 0; i < 6; i++ {
		newTree.Set("key"+gconv.String(i), "val"+gconv.String(i))
	}

	fmt.Println(newTree.Search("key0"))
	fmt.Println(newTree.Search("key6"))

	// output:
	// val0 true
	// <nil> false
}

func ExampleBTree_Print() {
	newTree := gtree.NewBTree(3, gutil.ComparatorString)
	for i := 0; i < 6; i++ {
		newTree.Set("key"+gconv.String(i), "val"+gconv.String(i))
	}

	newTree.Print()

	// output:
	// key0
	// key1
	//     key2
	// key3
	//     key4
	//     key5
}

func ExampleBTree_Iterator() {
	newTree := gtree.NewBTree(3, gutil.ComparatorString)
	for i := 0; i < 10; i++ {
		newTree.Set(i, 10-i)
	}

	var totalKey, totalValue int
	newTree.Iterator(func(key, value interface{}) bool {
		totalKey += key.(int)
		totalValue += value.(int)

		return totalValue < 20
	})

	fmt.Println("totalKey:", totalKey)
	fmt.Println("totalValue:", totalValue)

	// Output:
	// totalKey: 3
	// totalValue: 27
}

func ExampleBTree_IteratorFrom() {
	m := make(map[interface{}]interface{})
	for i := 1; i <= 5; i++ {
		m[i] = i * 10
	}
	tree := gtree.NewBTreeFrom(3, gutil.ComparatorInt, m)

	tree.IteratorFrom(1, true, func(key, value interface{}) bool {
		fmt.Println("key:", key)
		fmt.Println("value:", value)
		return true
	})

	// Output:
	// key: 1
	// value: 10
	// key: 2
	// value: 20
	// key: 3
	// value: 30
	// key: 4
	// value: 40
	// key: 5
	// value: 50
}

func ExampleBTree_IteratorAsc() {
	newTree := gtree.NewBTree(3, gutil.ComparatorString)
	for i := 0; i < 10; i++ {
		newTree.Set(i, 10-i)
	}

	newTree.IteratorAsc(func(key, value interface{}) bool {
		fmt.Println("key:", key, ", value:", value)
		return true
	})

	// Output:
	// key: 0 , value: 10
	// key: 1 , value: 9
	// key: 2 , value: 8
	// key: 3 , value: 7
	// key: 4 , value: 6
	// key: 5 , value: 5
	// key: 6 , value: 4
	// key: 7 , value: 3
	// key: 8 , value: 2
	// key: 9 , value: 1
}

func ExampleBTree_IteratorAscFrom_Normal() {
	m := make(map[interface{}]interface{})
	for i := 1; i <= 5; i++ {
		m[i] = i * 10
	}
	tree := gtree.NewBTreeFrom(3, gutil.ComparatorInt, m)

	tree.IteratorAscFrom(1, true, func(key, value interface{}) bool {
		fmt.Println("key:", key, ", value:", value)
		return true
	})

	// Output:
	// key: 1 , value: 10
	// key: 2 , value: 20
	// key: 3 , value: 30
	// key: 4 , value: 40
	// key: 5 , value: 50
}

func ExampleBTree_IteratorAscFrom_NoExistKey() {
	m := make(map[interface{}]interface{})
	for i := 1; i <= 5; i++ {
		m[i] = i * 10
	}
	tree := gtree.NewBTreeFrom(3, gutil.ComparatorInt, m)

	tree.IteratorAscFrom(0, true, func(key, value interface{}) bool {
		fmt.Println("key:", key, ", value:", value)
		return true
	})

	// Output:
}

func ExampleBTree_IteratorAscFrom_NoExistKeyAndMatchFalse() {
	m := make(map[interface{}]interface{})
	for i := 1; i <= 5; i++ {
		m[i] = i * 10
	}
	tree := gtree.NewBTreeFrom(3, gutil.ComparatorInt, m)

	tree.IteratorAscFrom(0, false, func(key, value interface{}) bool {
		fmt.Println("key:", key, ", value:", value)
		return true
	})

	// Output:
	// key: 1 , value: 10
	// key: 2 , value: 20
	// key: 3 , value: 30
	// key: 4 , value: 40
	// key: 5 , value: 50
}

func ExampleBTree_IteratorDesc() {
	newTree := gtree.NewBTree(3, gutil.ComparatorString)
	for i := 0; i < 10; i++ {
		newTree.Set(i, 10-i)
	}

	newTree.IteratorDesc(func(key, value interface{}) bool {
		fmt.Println("key:", key, ", value:", value)
		return true
	})

	// Output:
	// key: 9 , value: 1
	// key: 8 , value: 2
	// key: 7 , value: 3
	// key: 6 , value: 4
	// key: 5 , value: 5
	// key: 4 , value: 6
	// key: 3 , value: 7
	// key: 2 , value: 8
	// key: 1 , value: 9
	// key: 0 , value: 10
}

func ExampleBTree_IteratorDescFrom() {
	m := make(map[interface{}]interface{})
	for i := 1; i <= 5; i++ {
		m[i] = i * 10
	}
	tree := gtree.NewBTreeFrom(3, gutil.ComparatorInt, m)

	tree.IteratorDescFrom(5, true, func(key, value interface{}) bool {
		fmt.Println("key:", key, ", value:", value)
		return true
	})

	// Output:
	// key: 5 , value: 50
	// key: 4 , value: 40
	// key: 3 , value: 30
	// key: 2 , value: 20
	// key: 1 , value: 10
}

func ExampleBTree_MarshalJSON() {
	newTree := gtree.NewBTree(3, gutil.ComparatorString)
	for i := 0; i < 6; i++ {
		newTree.Set("key"+gconv.String(i), "val"+gconv.String(i))
	}

	bytes, err := json.Marshal(newTree)
	if err == nil {
		fmt.Println(gconv.String(bytes))
	}

	// output:
	// {"key0":"val0","key1":"val1","key2":"val2","key3":"val3","key4":"val4","key5":"val5"}
}
