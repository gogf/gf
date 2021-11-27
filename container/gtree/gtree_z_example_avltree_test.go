// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with gm file,
// You can obtain one at https://github.com/gogf/gf.

package gtree_test

import (
	"fmt"
	"github.com/gogf/gf/v2/container/gtree"
	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gutil"
)

func ExampleAVLTree_Clone() {
	avl := gtree.NewAVLTree(gutil.ComparatorString)
	for i := 0; i < 6; i++ {
		avl.Set("key"+gconv.String(i), "val"+gconv.String(i))
	}

	tree := avl.Clone()

	fmt.Println(tree.Map())
	fmt.Println(tree.Size())

	// Output:
	// map[key0:val0 key1:val1 key2:val2 key3:val3 key4:val4 key5:val5]
	// 6
}

func ExampleAVLTree_Set() {
	tree := gtree.NewAVLTree(gutil.ComparatorString)
	for i := 0; i < 6; i++ {
		tree.Set("key"+gconv.String(i), "val"+gconv.String(i))
	}

	fmt.Println(tree.Map())
	fmt.Println(tree.Size())

	// Output:
	// map[key0:val0 key1:val1 key2:val2 key3:val3 key4:val4 key5:val5]
	// 6
}

func ExampleAVLTree_Sets() {
	tree := gtree.NewAVLTree(gutil.ComparatorString)

	tree.Sets(map[interface{}]interface{}{
		"key1": "val1",
		"key2": "val2",
	})

	fmt.Println(tree.Map())
	fmt.Println(tree.Size())

	// Output:
	// map[key1:val1 key2:val2]
	// 2
}

func ExampleAVLTree_Get() {
	tree := gtree.NewAVLTree(gutil.ComparatorString)
	for i := 0; i < 6; i++ {
		tree.Set("key"+gconv.String(i), "val"+gconv.String(i))
	}

	fmt.Println(tree.Get("key1"))
	fmt.Println(tree.Get("key10"))

	// Output:
	// val1
	// <nil>
}

func ExampleAVLTree_GetOrSet() {
	tree := gtree.NewAVLTree(gutil.ComparatorString)
	for i := 0; i < 6; i++ {
		tree.Set("key"+gconv.String(i), "val"+gconv.String(i))
	}

	fmt.Println(tree.GetOrSet("key1", "newVal1"))
	fmt.Println(tree.GetOrSet("key6", "val6"))

	// Output:
	// val1
	// val6
}

func ExampleAVLTree_GetOrSetFunc() {
	tree := gtree.NewAVLTree(gutil.ComparatorString)
	for i := 0; i < 6; i++ {
		tree.Set("key"+gconv.String(i), "val"+gconv.String(i))
	}

	fmt.Println(tree.GetOrSetFunc("key1", func() interface{} {
		return "newVal1"
	}))
	fmt.Println(tree.GetOrSetFunc("key6", func() interface{} {
		return "val6"
	}))

	// Output:
	// val1
	// val6
}

func ExampleAVLTree_GetOrSetFuncLock() {
	tree := gtree.NewAVLTree(gutil.ComparatorString)
	for i := 0; i < 6; i++ {
		tree.Set("key"+gconv.String(i), "val"+gconv.String(i))
	}

	fmt.Println(tree.GetOrSetFuncLock("key1", func() interface{} {
		return "newVal1"
	}))
	fmt.Println(tree.GetOrSetFuncLock("key6", func() interface{} {
		return "val6"
	}))

	// Output:
	// val1
	// val6
}

func ExampleAVLTree_GetVar() {
	tree := gtree.NewAVLTree(gutil.ComparatorString)
	for i := 0; i < 6; i++ {
		tree.Set("key"+gconv.String(i), "val"+gconv.String(i))
	}

	fmt.Println(tree.GetVar("key1").String())

	// Output:
	// val1
}

func ExampleAVLTree_GetVarOrSet() {
	tree := gtree.NewAVLTree(gutil.ComparatorString)
	for i := 0; i < 6; i++ {
		tree.Set("key"+gconv.String(i), "val"+gconv.String(i))
	}

	fmt.Println(tree.GetVarOrSet("key1", "newVal1"))
	fmt.Println(tree.GetVarOrSet("key6", "val6"))

	// Output:
	// val1
	// val6
}

func ExampleAVLTree_GetVarOrSetFunc() {
	tree := gtree.NewAVLTree(gutil.ComparatorString)
	for i := 0; i < 6; i++ {
		tree.Set("key"+gconv.String(i), "val"+gconv.String(i))
	}

	fmt.Println(tree.GetVarOrSetFunc("key1", func() interface{} {
		return "newVal1"
	}))
	fmt.Println(tree.GetVarOrSetFunc("key6", func() interface{} {
		return "val6"
	}))

	// Output:
	// val1
	// val6
}

func ExampleAVLTree_GetVarOrSetFuncLock() {
	tree := gtree.NewAVLTree(gutil.ComparatorString)
	for i := 0; i < 6; i++ {
		tree.Set("key"+gconv.String(i), "val"+gconv.String(i))
	}

	fmt.Println(tree.GetVarOrSetFuncLock("key1", func() interface{} {
		return "newVal1"
	}))
	fmt.Println(tree.GetVarOrSetFuncLock("key6", func() interface{} {
		return "val6"
	}))

	// Output:
	// val1
	// val6
}

func ExampleAVLTree_SetIfNotExist() {
	tree := gtree.NewAVLTree(gutil.ComparatorString)
	for i := 0; i < 6; i++ {
		tree.Set("key"+gconv.String(i), "val"+gconv.String(i))
	}

	fmt.Println(tree.SetIfNotExist("key1", "newVal1"))
	fmt.Println(tree.SetIfNotExist("key6", "val6"))

	// Output:
	// false
	// true
}

func ExampleAVLTree_SetIfNotExistFunc() {
	tree := gtree.NewAVLTree(gutil.ComparatorString)
	for i := 0; i < 6; i++ {
		tree.Set("key"+gconv.String(i), "val"+gconv.String(i))
	}

	fmt.Println(tree.SetIfNotExistFunc("key1", func() interface{} {
		return "newVal1"
	}))
	fmt.Println(tree.SetIfNotExistFunc("key6", func() interface{} {
		return "val6"
	}))

	// Output:
	// false
	// true
}

func ExampleAVLTree_SetIfNotExistFuncLock() {
	tree := gtree.NewAVLTree(gutil.ComparatorString)
	for i := 0; i < 6; i++ {
		tree.Set("key"+gconv.String(i), "val"+gconv.String(i))
	}

	fmt.Println(tree.SetIfNotExistFuncLock("key1", func() interface{} {
		return "newVal1"
	}))
	fmt.Println(tree.SetIfNotExistFuncLock("key6", func() interface{} {
		return "val6"
	}))

	// Output:
	// false
	// true
}

func ExampleAVLTree_Contains() {
	tree := gtree.NewAVLTree(gutil.ComparatorString)
	for i := 0; i < 6; i++ {
		tree.Set("key"+gconv.String(i), "val"+gconv.String(i))
	}

	fmt.Println(tree.Contains("key1"))
	fmt.Println(tree.Contains("key6"))

	// Output:
	// true
	// false
}

func ExampleAVLTree_Remove() {
	tree := gtree.NewAVLTree(gutil.ComparatorString)
	for i := 0; i < 6; i++ {
		tree.Set("key"+gconv.String(i), "val"+gconv.String(i))
	}

	fmt.Println(tree.Remove("key1"))
	fmt.Println(tree.Remove("key6"))
	fmt.Println(tree.Map())

	// Output:
	// val1
	// <nil>
	// map[key0:val0 key2:val2 key3:val3 key4:val4 key5:val5]
}

func ExampleAVLTree_Removes() {
	tree := gtree.NewAVLTree(gutil.ComparatorString)
	for i := 0; i < 6; i++ {
		tree.Set("key"+gconv.String(i), "val"+gconv.String(i))
	}

	removeKeys := make([]interface{}, 2)
	removeKeys = append(removeKeys, "key1")
	removeKeys = append(removeKeys, "key6")

	tree.Removes(removeKeys)

	fmt.Println(tree.Map())

	// Output:
	// map[key0:val0 key2:val2 key3:val3 key4:val4 key5:val5]
}

func ExampleAVLTree_IsEmpty() {
	tree := gtree.NewAVLTree(gutil.ComparatorString)

	fmt.Println(tree.IsEmpty())

	for i := 0; i < 6; i++ {
		tree.Set("key"+gconv.String(i), "val"+gconv.String(i))
	}

	fmt.Println(tree.IsEmpty())

	// Output:
	// true
	// false
}

func ExampleAVLTree_Size() {
	tree := gtree.NewAVLTree(gutil.ComparatorString)

	fmt.Println(tree.Size())

	for i := 0; i < 6; i++ {
		tree.Set("key"+gconv.String(i), "val"+gconv.String(i))
	}

	fmt.Println(tree.Size())

	// Output:
	// 0
	// 6
}

func ExampleAVLTree_Keys() {
	tree := gtree.NewAVLTree(gutil.ComparatorString)
	for i := 6; i > 0; i-- {
		tree.Set("key"+gconv.String(i), "val"+gconv.String(i))
	}

	fmt.Println(tree.Keys())

	// Output:
	// [key1 key2 key3 key4 key5 key6]
}

func ExampleAVLTree_Values() {
	tree := gtree.NewAVLTree(gutil.ComparatorString)
	for i := 6; i > 0; i-- {
		tree.Set("key"+gconv.String(i), "val"+gconv.String(i))
	}

	fmt.Println(tree.Values())

	// Output:
	// [val1 val2 val3 val4 val5 val6]
}

func ExampleAVLTree_Map() {
	tree := gtree.NewAVLTree(gutil.ComparatorString)
	for i := 0; i < 6; i++ {
		tree.Set("key"+gconv.String(i), "val"+gconv.String(i))
	}

	fmt.Println(tree.Map())

	// Output:
	// map[key0:val0 key1:val1 key2:val2 key3:val3 key4:val4 key5:val5]
}

func ExampleAVLTree_MapStrAny() {
	tree := gtree.NewAVLTree(gutil.ComparatorString)
	for i := 0; i < 6; i++ {
		tree.Set(1000+i, "val"+gconv.String(i))
	}

	fmt.Println(tree.MapStrAny())

	// Output:
	// map[1000:val0 1001:val1 1002:val2 1003:val3 1004:val4 1005:val5]
}

func ExampleAVLTree_Flip() {
	tree := gtree.NewAVLTree(gutil.ComparatorInt)
	for i := 1; i < 6; i++ {
		tree.Set(i, i*10)
	}

	fmt.Println("Before Flip", tree.Map())

	tree.Flip()

	fmt.Println("After Flip", tree.Map())

	// Output:
	// Before Flip map[1:10 2:20 3:30 4:40 5:50]
	// After Flip map[10:1 20:2 30:3 40:4 50:5]
}

func ExampleAVLTree_Clear() {
	tree := gtree.NewAVLTree(gutil.ComparatorString)
	for i := 0; i < 6; i++ {
		tree.Set(1000+i, "val"+gconv.String(i))
	}
	fmt.Println(tree.Size())

	tree.Clear()
	fmt.Println(tree.Size())

	// Output:
	// 6
	// 0
}

func ExampleAVLTree_Replace() {
	tree := gtree.NewAVLTree(gutil.ComparatorString)
	for i := 0; i < 6; i++ {
		tree.Set("key"+gconv.String(i), "val"+gconv.String(i))
	}

	fmt.Println(tree.Map())

	data := map[interface{}]interface{}{
		"newKey0": "newVal0",
		"newKey1": "newVal1",
		"newKey2": "newVal2",
	}

	tree.Replace(data)

	fmt.Println(tree.Map())

	// Output:
	// map[key0:val0 key1:val1 key2:val2 key3:val3 key4:val4 key5:val5]
	// map[newKey0:newVal0 newKey1:newVal1 newKey2:newVal2]
}

func ExampleAVLTree_Left() {
	tree := gtree.NewAVLTree(gutil.ComparatorInt)
	for i := 1; i < 100; i++ {
		tree.Set(i, i)
	}
	fmt.Println(tree.Left().Key, tree.Left().Value)

	emptyTree := gtree.NewAVLTree(gutil.ComparatorInt)
	fmt.Println(emptyTree.Left())

	// Output:
	// 1 1
	// <nil>
}

func ExampleAVLTree_Right() {
	tree := gtree.NewAVLTree(gutil.ComparatorInt)
	for i := 1; i < 100; i++ {
		tree.Set(i, i)
	}
	fmt.Println(tree.Right().Key, tree.Right().Value)

	emptyTree := gtree.NewAVLTree(gutil.ComparatorInt)
	fmt.Println(emptyTree.Left())

	// Output:
	// 99 99
	// <nil>
}

func ExampleAVLTree_Floor() {
	tree := gtree.NewAVLTree(gutil.ComparatorInt)
	for i := 1; i < 100; i++ {
		if i != 50 {
			tree.Set(i, i)
		}
	}

	node, found := tree.Floor(95)
	if found {
		fmt.Println("Floor 95:", node.Key)
	}

	node, found = tree.Floor(50)
	if found {
		fmt.Println("Floor 50:", node.Key)
	}

	node, found = tree.Floor(100)
	if found {
		fmt.Println("Floor 100:", node.Key)
	}

	node, found = tree.Floor(0)
	if found {
		fmt.Println("Floor 0:", node.Key)
	}

	// Output:
	// Floor 95: 95
	// Floor 50: 49
	// Floor 100: 99
}

func ExampleAVLTree_Ceiling() {
	tree := gtree.NewAVLTree(gutil.ComparatorInt)
	for i := 1; i < 100; i++ {
		if i != 50 {
			tree.Set(i, i)
		}
	}

	node, found := tree.Ceiling(1)
	if found {
		fmt.Println("Ceiling 1:", node.Key)
	}

	node, found = tree.Ceiling(50)
	if found {
		fmt.Println("Ceiling 50:", node.Key)
	}

	node, found = tree.Ceiling(100)
	if found {
		fmt.Println("Ceiling 100:", node.Key)
	}

	node, found = tree.Ceiling(-1)
	if found {
		fmt.Println("Ceiling -1:", node.Key)
	}

	// Output:
	// Ceiling 1: 1
	// Ceiling 50: 51
	// Ceiling -1: 1
}

func ExampleAVLTree_String() {
	tree := gtree.NewAVLTree(gutil.ComparatorString)
	for i := 0; i < 6; i++ {
		tree.Set("key"+gconv.String(i), "val"+gconv.String(i))
	}

	fmt.Println(tree.String())

	// Output:
	// │       ┌── key5
	// │   ┌── key4
	// └── key3
	//     │   ┌── key2
	//     └── key1
	//         └── key0
}

func ExampleAVLTree_Search() {
	tree := gtree.NewAVLTree(gutil.ComparatorString)
	for i := 0; i < 6; i++ {
		tree.Set("key"+gconv.String(i), "val"+gconv.String(i))
	}

	fmt.Println(tree.Search("key0"))
	fmt.Println(tree.Search("key6"))

	// Output:
	// val0 true
	// <nil> false
}

func ExampleAVLTree_Print() {
	tree := gtree.NewAVLTree(gutil.ComparatorString)
	for i := 0; i < 6; i++ {
		tree.Set("key"+gconv.String(i), "val"+gconv.String(i))
	}

	tree.Print()

	// Output:
	// │       ┌── key5
	// │   ┌── key4
	// └── key3
	//     │   ┌── key2
	//     └── key1
	//         └── key0
}

func ExampleAVLTree_Iterator() {
	tree := gtree.NewAVLTree(gutil.ComparatorString)
	for i := 0; i < 10; i++ {
		tree.Set(i, 10-i)
	}

	var totalKey, totalValue int
	tree.Iterator(func(key, value interface{}) bool {
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

func ExampleAVLTree_IteratorFrom() {
	m := make(map[interface{}]interface{})
	for i := 1; i <= 5; i++ {
		m[i] = i * 10
	}
	tree := gtree.NewAVLTreeFrom(gutil.ComparatorInt, m)

	tree.IteratorFrom(1, true, func(key, value interface{}) bool {
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

func ExampleAVLTree_IteratorAsc() {
	tree := gtree.NewAVLTree(gutil.ComparatorString)
	for i := 0; i < 10; i++ {
		tree.Set(i, 10-i)
	}

	tree.IteratorAsc(func(key, value interface{}) bool {
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

func ExampleAVLTree_IteratorAscFrom_Normal() {
	m := make(map[interface{}]interface{})
	for i := 1; i <= 5; i++ {
		m[i] = i * 10
	}
	tree := gtree.NewAVLTreeFrom(gutil.ComparatorInt, m)

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

func ExampleAVLTree_IteratorAscFrom_NoExistKey() {
	m := make(map[interface{}]interface{})
	for i := 1; i <= 5; i++ {
		m[i] = i * 10
	}
	tree := gtree.NewAVLTreeFrom(gutil.ComparatorInt, m)

	tree.IteratorAscFrom(0, true, func(key, value interface{}) bool {
		fmt.Println("key:", key, ", value:", value)
		return true
	})

	// Output:
}

func ExampleAVLTree_IteratorAscFrom_NoExistKeyAndMatchFalse() {
	m := make(map[interface{}]interface{})
	for i := 1; i <= 5; i++ {
		m[i] = i * 10
	}
	tree := gtree.NewAVLTreeFrom(gutil.ComparatorInt, m)

	tree.IteratorAscFrom(6, false, func(key, value interface{}) bool {
		fmt.Println("key:", key, ", value:", value)
		return true
	})

	// Output:
}

func ExampleAVLTree_IteratorDesc() {
	tree := gtree.NewAVLTree(gutil.ComparatorString)
	for i := 0; i < 10; i++ {
		tree.Set(i, 10-i)
	}

	tree.IteratorDesc(func(key, value interface{}) bool {
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

func ExampleAVLTree_IteratorDescFrom() {
	m := make(map[interface{}]interface{})
	for i := 1; i <= 5; i++ {
		m[i] = i * 10
	}
	tree := gtree.NewAVLTreeFrom(gutil.ComparatorInt, m)

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

func ExampleAVLTree_MarshalJSON() {
	tree := gtree.NewAVLTree(gutil.ComparatorString)
	for i := 0; i < 6; i++ {
		tree.Set("key"+gconv.String(i), "val"+gconv.String(i))
	}

	bytes, err := json.Marshal(tree)
	if err == nil {
		fmt.Println(gconv.String(bytes))
	}

	// Output:
	// {"key0":"val0","key1":"val1","key2":"val2","key3":"val3","key4":"val4","key5":"val5"}
}
