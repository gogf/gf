// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go

package gset_test

import (
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/container/gset"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gconv"
)

func TestIntSet_Var(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var s gset.IntSet
		s.Add(1, 1, 2)
		s.Add([]int{3, 4}...)
		t.Assert(s.Size(), 4)
		t.AssertIN(1, s.Slice())
		t.AssertIN(2, s.Slice())
		t.AssertIN(3, s.Slice())
		t.AssertIN(4, s.Slice())
		t.AssertNI(0, s.Slice())
		t.Assert(s.Contains(4), true)
		t.Assert(s.Contains(5), false)
		s.Remove(1)
		t.Assert(s.Size(), 3)
		s.Clear()
		t.Assert(s.Size(), 0)
	})
}

func TestIntSet_Basic(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := gset.NewIntSet()
		s.Add(1, 1, 2)
		s.Add([]int{3, 4}...)
		t.Assert(s.Size(), 4)
		t.AssertIN(1, s.Slice())
		t.AssertIN(2, s.Slice())
		t.AssertIN(3, s.Slice())
		t.AssertIN(4, s.Slice())
		t.AssertNI(0, s.Slice())
		t.Assert(s.Contains(4), true)
		t.Assert(s.Contains(5), false)
		s.Remove(1)
		t.Assert(s.Size(), 3)
		s.Clear()
		t.Assert(s.Size(), 0)
	})
}

func TestIntSet_Iterator_Deadlock(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		set := gset.NewIntSetFrom([]int{1, 2, 3, 4, 5}, true)
		set.Iterator(func(k int) bool {
			if k%2 == 0 {
				set.Remove(k)
			}
			return true
		})
		t.Assert(set.Contains(1), true)
		t.Assert(set.Contains(2), false)
		t.Assert(set.Contains(3), true)
		t.Assert(set.Contains(4), false)
		t.Assert(set.Contains(5), true)
	})
}

func TestIntSet_Iterator(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := gset.NewIntSet()
		s.Add(1, 2, 3)
		t.Assert(s.Size(), 3)

		a1 := garray.New(true)
		a2 := garray.New(true)
		s.Iterator(func(v int) bool {
			a1.Append(1)
			return false
		})
		s.Iterator(func(v int) bool {
			a2.Append(1)
			return true
		})
		t.Assert(a1.Len(), 1)
		t.Assert(a2.Len(), 3)
	})
}

func TestIntSet_LockFunc(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := gset.NewIntSet()
		s.Add(1, 2, 3)
		t.Assert(s.Size(), 3)
		s.LockFunc(func(m map[int]struct{}) {
			delete(m, 1)
		})
		t.Assert(s.Size(), 2)
		s.RLockFunc(func(m map[int]struct{}) {
			t.Assert(m, map[int]struct{}{
				3: {},
				2: {},
			})
		})
	})
}

func TestIntSet_Equal(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := gset.NewIntSet()
		s2 := gset.NewIntSet()
		s3 := gset.NewIntSet()
		s4 := gset.NewIntSet()
		s1.Add(1, 2, 3)
		s2.Add(1, 2, 3)
		s3.Add(1, 2, 3, 4)
		s4.Add(4, 5, 6)
		t.Assert(s1.Equal(s2), true)
		t.Assert(s1.Equal(s3), false)
		t.Assert(s1.Equal(s4), false)
		s5 := s1
		t.Assert(s1.Equal(s5), true)
	})
}

func TestIntSet_IsSubsetOf(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := gset.NewIntSet()
		s2 := gset.NewIntSet()
		s3 := gset.NewIntSet()
		s1.Add(1, 2)
		s2.Add(1, 2, 3)
		s3.Add(1, 2, 3, 4)
		t.Assert(s1.IsSubsetOf(s2), true)
		t.Assert(s2.IsSubsetOf(s3), true)
		t.Assert(s1.IsSubsetOf(s3), true)
		t.Assert(s2.IsSubsetOf(s1), false)
		t.Assert(s3.IsSubsetOf(s2), false)

		s4 := s1
		t.Assert(s1.IsSubsetOf(s4), true)
	})
}

func TestIntSet_Union(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := gset.NewIntSet()
		s2 := gset.NewIntSet()
		s1.Add(1, 2)
		s2.Add(3, 4)
		s3 := s1.Union(s2)
		t.Assert(s3.Contains(1), true)
		t.Assert(s3.Contains(2), true)
		t.Assert(s3.Contains(3), true)
		t.Assert(s3.Contains(4), true)
	})
}

func TestIntSet_Diff(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := gset.NewIntSet()
		s2 := gset.NewIntSet()
		s1.Add(1, 2, 3)
		s2.Add(3, 4, 5)
		s3 := s1.Diff(s2)
		t.Assert(s3.Contains(1), true)
		t.Assert(s3.Contains(2), true)
		t.Assert(s3.Contains(3), false)
		t.Assert(s3.Contains(4), false)

		s4 := s1
		s5 := s1.Diff(s2, s4)
		t.Assert(s5.Contains(1), true)
		t.Assert(s5.Contains(2), true)
		t.Assert(s5.Contains(3), false)
		t.Assert(s5.Contains(4), false)
	})
}

func TestIntSet_Intersect(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := gset.NewIntSet()
		s2 := gset.NewIntSet()
		s1.Add(1, 2, 3)
		s2.Add(3, 4, 5)
		s3 := s1.Intersect(s2)
		t.Assert(s3.Contains(1), false)
		t.Assert(s3.Contains(2), false)
		t.Assert(s3.Contains(3), true)
		t.Assert(s3.Contains(4), false)
	})
}

func TestIntSet_Complement(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := gset.NewIntSet()
		s2 := gset.NewIntSet()
		s1.Add(1, 2, 3)
		s2.Add(3, 4, 5)
		s3 := s1.Complement(s2)
		t.Assert(s3.Contains(1), false)
		t.Assert(s3.Contains(2), false)
		t.Assert(s3.Contains(4), true)
		t.Assert(s3.Contains(5), true)
	})
}

func TestIntSet_Size(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := gset.NewIntSet(true)
		s1.Add(1, 2, 3)
		t.Assert(s1.Size(), 3)

	})

}

func TestIntSet_Merge(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := gset.NewIntSet()
		s2 := gset.NewIntSet()
		s1.Add(1, 2, 3)
		s2.Add(3, 4, 5)
		s3 := s1.Merge(s2)
		t.Assert(s3.Contains(1), true)
		t.Assert(s3.Contains(5), true)
		t.Assert(s3.Contains(6), false)
	})
}

func TestIntSet_Join(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := gset.NewIntSet()
		t.Assert(s1.Join(","), "")
		s1.Add(1, 2, 3)
		s3 := s1.Join(",")
		t.Assert(strings.Contains(s3, "1"), true)
		t.Assert(strings.Contains(s3, "2"), true)
		t.Assert(strings.Contains(s3, "3"), true)
	})
}

func TestIntSet_String(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := gset.NewIntSet()
		s1.Add(1, 2, 3)
		s3 := s1.String()
		t.Assert(strings.Contains(s3, "["), true)
		t.Assert(strings.Contains(s3, "]"), true)
		t.Assert(strings.Contains(s3, "1"), true)
		t.Assert(strings.Contains(s3, "2"), true)
		t.Assert(strings.Contains(s3, "3"), true)
		s1 = nil
		t.Assert(s1.String(), "")
	})
}

func TestIntSet_Sum(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := gset.NewIntSet()
		s1.Add(1, 2, 3)
		s2 := gset.NewIntSet()
		s2.Add(5, 6, 7)
		t.Assert(s2.Sum(), 18)

	})

}

func TestIntSet_Pop(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := gset.NewIntSet()
		t.Assert(s.Pop(), 0)
		s.Add(4, 2, 3)
		t.Assert(s.Size(), 3)
		t.AssertIN(s.Pop(), []int{4, 2, 3})
		t.AssertIN(s.Pop(), []int{4, 2, 3})
		t.Assert(s.Size(), 1)
	})
}

func TestIntSet_Pops(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := gset.NewIntSet()
		s.Add(1, 4, 2, 3)
		t.Assert(s.Size(), 4)
		t.Assert(s.Pops(0), nil)
		t.AssertIN(s.Pops(1), []int{1, 4, 2, 3})
		t.Assert(s.Size(), 3)
		a := s.Pops(2)
		t.Assert(len(a), 2)
		t.AssertIN(a, []int{1, 4, 2, 3})
		t.Assert(s.Size(), 1)
	})

	gtest.C(t, func(t *gtest.T) {
		s := gset.NewIntSet(true)
		a := []int{1, 2, 3, 4}
		s.Add(a...)
		t.Assert(s.Size(), 4)
		t.Assert(s.Pops(-2), nil)
		t.AssertIN(s.Pops(-1), a)
	})
}

func TestIntSet_AddIfNotExist(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := gset.NewIntSet(true)
		s.Add(1)
		t.Assert(s.Contains(1), true)
		t.Assert(s.AddIfNotExist(1), false)
		t.Assert(s.AddIfNotExist(2), true)
		t.Assert(s.Contains(2), true)
		t.Assert(s.AddIfNotExist(2), false)
		t.Assert(s.Contains(2), true)
	})
	gtest.C(t, func(t *gtest.T) {
		s := gset.IntSet{}
		t.Assert(s.AddIfNotExist(1), true)
	})
}

func TestIntSet_AddIfNotExistFunc(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := gset.NewIntSet(true)
		s.Add(1)
		t.Assert(s.Contains(1), true)
		t.Assert(s.Contains(2), false)
		t.Assert(s.AddIfNotExistFunc(2, func() bool { return false }), false)
		t.Assert(s.Contains(2), false)
		t.Assert(s.AddIfNotExistFunc(2, func() bool { return true }), true)
		t.Assert(s.Contains(2), true)
		t.Assert(s.AddIfNotExistFunc(2, func() bool { return true }), false)
		t.Assert(s.Contains(2), true)
	})
	gtest.C(t, func(t *gtest.T) {
		s := gset.NewIntSet(true)
		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			r := s.AddIfNotExistFunc(1, func() bool {
				time.Sleep(100 * time.Millisecond)
				return true
			})
			t.Assert(r, false)
		}()
		s.Add(1)
		wg.Wait()
	})
	gtest.C(t, func(t *gtest.T) {
		s := gset.IntSet{}
		t.Assert(s.AddIfNotExistFunc(1, func() bool { return true }), true)
	})
}

func TestIntSet_AddIfNotExistFuncLock(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := gset.NewIntSet(true)
		wg := sync.WaitGroup{}
		wg.Add(2)
		go func() {
			defer wg.Done()
			r := s.AddIfNotExistFuncLock(1, func() bool {
				time.Sleep(500 * time.Millisecond)
				return true
			})
			t.Assert(r, true)
		}()
		time.Sleep(100 * time.Millisecond)
		go func() {
			defer wg.Done()
			r := s.AddIfNotExistFuncLock(1, func() bool {
				return true
			})
			t.Assert(r, false)
		}()
		wg.Wait()
	})
	gtest.C(t, func(t *gtest.T) {
		s := gset.IntSet{}
		t.Assert(s.AddIfNotExistFuncLock(1, func() bool { return true }), true)
	})
}

func TestIntSet_Json(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := []int{1, 3, 2, 4}
		a1 := gset.NewIntSetFrom(s1)
		b1, err1 := json.Marshal(a1)
		b2, err2 := json.Marshal(s1)
		t.Assert(len(b1), len(b2))
		t.Assert(err1, err2)

		a2 := gset.NewIntSet()
		err2 = json.UnmarshalUseNumber(b2, &a2)
		t.Assert(err2, nil)
		t.Assert(a2.Contains(1), true)
		t.Assert(a2.Contains(2), true)
		t.Assert(a2.Contains(3), true)
		t.Assert(a2.Contains(4), true)
		t.Assert(a2.Contains(5), false)

		var a3 gset.IntSet
		err := json.UnmarshalUseNumber(b2, &a3)
		t.AssertNil(err)
		t.Assert(a2.Contains(1), true)
		t.Assert(a2.Contains(2), true)
		t.Assert(a2.Contains(3), true)
		t.Assert(a2.Contains(4), true)
		t.Assert(a2.Contains(5), false)
	})
}

func TestIntSet_Walk(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var set gset.IntSet
		set.Add(g.SliceInt{1, 2}...)
		set.Walk(func(item int) int {
			return item + 10
		})
		t.Assert(set.Size(), 2)
		t.Assert(set.Contains(11), true)
		t.Assert(set.Contains(12), true)
	})
}

func TestIntSet_UnmarshalValue(t *testing.T) {
	type V struct {
		Name string
		Set  *gset.IntSet
	}
	// JSON
	gtest.C(t, func(t *gtest.T) {
		var v *V
		err := gconv.Struct(g.Map{
			"name": "john",
			"set":  []byte(`[1,2,3]`),
		}, &v)
		t.AssertNil(err)
		t.Assert(v.Name, "john")
		t.Assert(v.Set.Size(), 3)
		t.Assert(v.Set.Contains(1), true)
		t.Assert(v.Set.Contains(2), true)
		t.Assert(v.Set.Contains(3), true)
		t.Assert(v.Set.Contains(4), false)
	})
	// Map
	gtest.C(t, func(t *gtest.T) {
		var v *V
		err := gconv.Struct(g.Map{
			"name": "john",
			"set":  g.Slice{1, 2, 3},
		}, &v)
		t.AssertNil(err)
		t.Assert(v.Name, "john")
		t.Assert(v.Set.Size(), 3)
		t.Assert(v.Set.Contains(1), true)
		t.Assert(v.Set.Contains(2), true)
		t.Assert(v.Set.Contains(3), true)
		t.Assert(v.Set.Contains(4), false)
	})
}

func TestIntSet_DeepCopy(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		set := gset.NewIntSet()
		set.Add(1, 2, 3)

		copySet := set.DeepCopy().(*gset.IntSet)
		copySet.Add(4)
		t.AssertNE(set.Size(), copySet.Size())
		t.AssertNE(set.String(), copySet.String())

		set = nil
		t.AssertNil(set.DeepCopy())
	})
}
