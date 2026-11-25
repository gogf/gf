// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtree_test

import (
	"testing"
	"time"

	"github.com/gogf/gf/v2/container/gtree"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gutil"
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
		t.Assert(avlTree.Size(), 5)

		avlTree2 := gtree.NewAVLKVTree[int, *Student](gutil.ComparatorTStr[int], true)
		avlTree2.SetNilChecker(func(student *Student) bool {
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
		t.Assert(btree.Size(), 5)

		btree2 := gtree.NewBKVTree[int, *Student](100, gutil.ComparatorTStr[int], true)
		btree2.SetNilChecker(func(student *Student) bool {
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
		t.Assert(redBlackTree.Size(), 5)

		redBlackTree2 := gtree.NewRedBlackKVTree[int, *Student](gutil.ComparatorTStr[int], true)
		redBlackTree2.SetNilChecker(func(student *Student) bool {
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

func Test_NewKVAVLTreeWithChecker_TypedNil(t *testing.T) {
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
		t.Assert(avlTree.Size(), 5)

		avlTree2 := gtree.NewAVLKVTreeWithChecker[int, *Student](gutil.ComparatorTStr[int], func(student *Student) bool {
			return student == nil
		}, true)
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

func Test_NewKVBTreeWithChecker_TypedNil(t *testing.T) {
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
		t.Assert(btree.Size(), 5)

		btree2 := gtree.NewBKVTreeWithChecker[int, *Student](100, gutil.ComparatorTStr[int], func(student *Student) bool {
			return student == nil
		}, true)
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

func Test_NewRedBlackKVTreeWithChecker_TypedNil(t *testing.T) {
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
		t.Assert(redBlackTree.Size(), 5)

		redBlackTree2 := gtree.NewRedBlackKVTreeWithChecker[int, *Student](gutil.ComparatorTStr[int], func(student *Student) bool {
			return student == nil
		}, true)

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

type kvIteratorFromTree interface {
	Set(key int, value int)
	IteratorAscFrom(key int, match bool, f func(key int, value int) bool)
	IteratorDescFrom(key int, match bool, f func(key int, value int) bool)
}

func newKVIteratorFromTrees() map[string]kvIteratorFromTree {
	return map[string]kvIteratorFromTree{
		"AVLKVTree":      gtree.NewAVLKVTree[int, int](gutil.ComparatorTStr[int], true),
		"BKVTree":        gtree.NewBKVTree[int, int](3, gutil.ComparatorTStr[int], true),
		"RedBlackKVTree": gtree.NewRedBlackKVTree[int, int](gutil.ComparatorTStr[int], true),
	}
}

// runWithin fails the test instead of hanging when `f` doesn't return in time.
func runWithin(t *gtest.T, name string, f func()) {
	done := make(chan struct{})
	go func() {
		defer close(done)
		f()
	}()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatalf("%s: iteration did not return", name)
	}
}

func Test_KVTree_IteratorFrom_Stop(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		for name, tree := range newKVIteratorFromTrees() {
			for i := 1; i <= 10; i++ {
				tree.Set(i, i*10)
			}
			var asc, desc []int
			runWithin(t, name, func() {
				tree.IteratorAscFrom(3, true, func(key, value int) bool {
					asc = append(asc, value)
					return len(asc) < 3
				})
				tree.IteratorDescFrom(8, true, func(key, value int) bool {
					desc = append(desc, value)
					return len(desc) < 3
				})
			})
			t.AssertEQ(asc, []int{30, 40, 50})
			t.AssertEQ(desc, []int{80, 70, 60})
		}
	})
}

func Test_KVTree_IteratorFrom_ConcurrentWrite(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		for name, tree := range newKVIteratorFromTrees() {
			for i := 1; i <= 5; i++ {
				tree.Set(i, i)
			}
			var (
				visited int
				written = make(chan struct{})
			)
			runWithin(t, name, func() {
				tree.IteratorAscFrom(1, true, func(key, value int) bool {
					if visited == 0 {
						// A writer waiting for the lock must not block the
						// rest of the iteration.
						go func() {
							tree.Set(100, 100)
							close(written)
						}()
						time.Sleep(50 * time.Millisecond)
					}
					visited++
					return true
				})
			})
			<-written
			t.Assert(visited, 5)
		}
	})
}
