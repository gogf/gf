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

func ExampleNewAVLTree() {

}

func ExampleNewAVLTreeFrom() {

}

func ExampleNewBTree() {
	bTree := gtree.NewBTree(3, gutil.ComparatorString)
	for i := 0; i < 6; i++ {
		bTree.Set("key"+gconv.String(i), "val"+gconv.String(i))
	}
	fmt.Println(bTree)

	// output:
	// key0
	// key1
	//     key2
	// key3
	//     key4
	//     key5
}

func ExampleNewBTreeFrom() {
	bTree := gtree.NewBTree(3, gutil.ComparatorString)
	for i := 0; i < 6; i++ {
		bTree.Set("key"+gconv.String(i), "val"+gconv.String(i))
	}

	otherBTree := gtree.NewBTreeFrom(3, gutil.ComparatorString, bTree.Map())
	fmt.Println(otherBTree)

	// output:
	// key0
	// key1
	//     key2
	// key3
	//     key4
	//     key5
}

func ExampleNewRedBlackTree() {

}

func ExampleNewRedBlackTreeFrom() {

}
