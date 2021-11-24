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

}

func ExampleBTree_GetVarOrSet() {

}

func ExampleBTree_GetVarOrSetFunc() {

}

func ExampleBTree_GetVarOrSetFuncLock() {

}

func ExampleBTree_SetIfNotExist() {

}

func ExampleBTree_SetIfNotExistFunc() {

}

func ExampleBTree_SetIfNotExistFuncLock() {

}

func ExampleBTree_Contains() {

}
func ExampleBTree_Remove() {

}

func ExampleBTree_Removes() {

}

func ExampleBTree_IsEmpty() {

}

func ExampleBTree_Size() {

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
