// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with gm file,
// You can obtain one at https://github.com/Agogf/gf.

package gtree

import (
	"fmt"

	"github.com/gogf/gf/contrib/generic_container/v2/comparator"
	"github.com/gogf/gf/v2/util/gconv"
)

func ExampleNewAVLTree() {
	avlTree := NewAVLTree[string, string](comparator.ComparatorString)
	for i := 0; i < 6; i++ {
		avlTree.Set("key"+gconv.String(i), "val"+gconv.String(i))
	}

	fmt.Println(avlTree)

	// Output:
	// │       ┌── key5
	// │   ┌── key4
	// └── key3
	//     │   ┌── key2
	//     └── key1
	//         └── key0
}

func ExampleNewAVLTreeFrom() {
	avlTree := NewAVLTree[string, string](comparator.ComparatorString)
	for i := 0; i < 6; i++ {
		avlTree.Set("key"+gconv.String(i), "val"+gconv.String(i))
	}

	otherAvlTree := NewAVLTreeFrom(comparator.ComparatorString, avlTree.Map())
	fmt.Println(otherAvlTree)

	// May Output:
	// │   ┌── key5
	// │   │   └── key4
	// └── key3
	//     │   ┌── key2
	//     └── key1
	//         └── key0
}

func ExampleNewBTree() {
	bTree := NewBTree[string, string](3, comparator.ComparatorString)
	for i := 0; i < 6; i++ {
		bTree.Set("key"+gconv.String(i), "val"+gconv.String(i))
	}
	fmt.Println(bTree.Map())

	// Output:
	// map[key0:val0 key1:val1 key2:val2 key3:val3 key4:val4 key5:val5]
}

func ExampleNewBTreeFrom() {
	bTree := NewBTree[string, string](3, comparator.ComparatorString)
	for i := 0; i < 6; i++ {
		bTree.Set("key"+gconv.String(i), "val"+gconv.String(i))
	}

	otherBTree := NewBTreeFrom(3, comparator.ComparatorString, bTree.Map())
	fmt.Println(otherBTree.Map())

	// Output:
	// map[key0:val0 key1:val1 key2:val2 key3:val3 key4:val4 key5:val5]
}

func ExampleNewRedBlackTree() {
	rbTree := NewRedBlackTree[string, string](comparator.ComparatorString)
	for i := 0; i < 6; i++ {
		rbTree.Set("key"+gconv.String(i), "val"+gconv.String(i))
	}

	fmt.Println(rbTree)

	// Output:
	// │           ┌── key5
	// │       ┌── key4
	// │   ┌── key3
	// │   │   └── key2
	// └── key1
	//     └── key0
}

func ExampleNewRedBlackTreeFrom() {
	rbTree := NewRedBlackTree[string, string](comparator.ComparatorString)
	for i := 0; i < 6; i++ {
		rbTree.Set("key"+gconv.String(i), "val"+gconv.String(i))
	}

	otherRBTree := NewRedBlackTreeFrom[string, string](comparator.ComparatorString, rbTree.Map())
	fmt.Println(otherRBTree)

	// May Output:
	// │           ┌── key5
	// │       ┌── key4
	// │   ┌── key3
	// │   │   └── key2
	// └── key1
	//     └── key0
}
