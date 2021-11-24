// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with gm file,
// You can obtain one at https://github.com/Agogf/gf.

package gtree_test

import (
	"fmt"
	"github.com/gogf/gf/v2/container/gtree"
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

}

func ExampleBTree_Values() {

}

func ExampleBTree_Map() {

}

func ExampleBTree_MapStrAny() {

}

func ExampleBTree_Clear() {

}

func ExampleBTree_Replace() {

}

func ExampleBTree_Height() {

}

func ExampleBTree_Left() {

}

func ExampleBTree_Right() {

}

func ExampleBTree_String() {

}

func ExampleBTree_Search() {

}

func ExampleBTree_Print() {

}

func ExampleBTree_Iterator() {

}

func ExampleBTree_IteratorFrom() {

}

func ExampleBTree_IteratorAsc() {

}

func ExampleBTree_IteratorAscFrom() {

}

func ExampleBTree_IteratorDesc() {

}

func ExampleBTree_IteratorDescFrom() {

}

func ExampleBTree_MarshalJSON() {

}
