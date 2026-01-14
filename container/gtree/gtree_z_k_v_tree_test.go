// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtree_test

import (
	"github.com/gogf/gf/v2/util/gutil"
	"testing"

	"github.com/gogf/gf/v2/container/gtree"
	"github.com/gogf/gf/v2/test/gtest"
)

func Test_KVAVLTree_TypedNil(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type Student struct {
			Name string
			Age  int
		}
		avlTree := gtree.NewAVLKVTree[int, *Student](gutil.ComparatorTStr[int], true)
		for i := 0; i < 10; i++ {
			if i%2 == 0 {
				avlTree.Set(i, &Student{})
			} else {
				var s *Student = nil
				avlTree.Set(i, s)
			}
		}
		t.Assert(avlTree.Size(), 10)
		avlTree2 := gtree.NewAVLKVTree[int, *Student](gutil.ComparatorTStr[int], true)
		avlTree2.RegisterNilChecker(func(student *Student) bool {
			return student == nil
		})
		for i := 0; i < 10; i++ {
			if i%2 == 0 {
				avlTree2.Set(i, &Student{})
			} else {
				var s *Student = nil
				avlTree2.Set(i, s)
			}
		}
		t.Assert(avlTree2.Size(), 5)

	})
}

func Test_KVBTree_TypedNil(t *testing.T) {
	type Student struct {
		Name string
		Age  int
	}
	gtest.C(t, func(t *gtest.T) {
		btree := gtree.NewBKVTree[int, *Student](100, gutil.ComparatorTStr[int], true)
		for i := 0; i < 10; i++ {
			if i%2 == 0 {
				btree.Set(i, &Student{})
			} else {
				var s *Student = nil
				btree.Set(i, s)
			}
		}
		t.Assert(btree.Size(), 10)
		btree2 := gtree.NewBKVTree[int, *Student](100, gutil.ComparatorTStr[int], true)
		btree2.RegisterNilChecker(func(student *Student) bool {
			return student == nil
		})
		for i := 0; i < 10; i++ {
			if i%2 == 0 {
				btree2.Set(i, &Student{})
			} else {
				var s *Student = nil
				btree2.Set(i, s)
			}
		}
		t.Assert(btree2.Size(), 5)
	})

}

func Test_KVRedBlackTree_TypedNil(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type Student struct {
			Name string
			Age  int
		}
		redBlackTree := gtree.NewRedBlackKVTree[int, *Student](gutil.ComparatorTStr[int], true)
		for i := 0; i < 10; i++ {
			if i%2 == 0 {
				redBlackTree.Set(i, &Student{})
			} else {
				var s *Student = nil
				redBlackTree.Set(i, s)
			}
		}
		t.Assert(redBlackTree.Size(), 10)
		redBlackTree2 := gtree.NewRedBlackKVTree[int, *Student](gutil.ComparatorTStr[int], true)

		redBlackTree2.RegisterNilChecker(func(student *Student) bool {
			return student == nil
		})
		for i := 0; i < 10; i++ {
			if i%2 == 0 {
				redBlackTree2.Set(i, &Student{})
			} else {
				var s *Student = nil
				redBlackTree2.Set(i, s)
			}
		}
		t.Assert(redBlackTree2.Size(), 5)
	})
}
