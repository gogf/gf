// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go

package gset_test

import (
	"encoding/json"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/util/gconv"
	"strings"
	"testing"

	"github.com/gogf/gf/container/garray"
	"github.com/gogf/gf/container/gset"
	"github.com/gogf/gf/test/gtest"
)

func TestIntSet_Var(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var s gset.IntSet
		s.Add(1).Add(1).Add(2)
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
		s.Add(1).Add(1).Add(2)
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

func TestIntSet_Iterator(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := gset.NewIntSet()
		s.Add(1).Add(2).Add(3)
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
		s.Add(1).Add(2).Add(3)
		t.Assert(s.Size(), 3)
		s.LockFunc(func(m map[int]struct{}) {
			delete(m, 1)
		})
		t.Assert(s.Size(), 2)
		s.RLockFunc(func(m map[int]struct{}) {
			t.Assert(m, map[int]struct{}{
				3: struct{}{},
				2: struct{}{},
			})
		})
	})
}

func TestIntSet_Equal(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := gset.NewIntSet()
		s2 := gset.NewIntSet()
		s3 := gset.NewIntSet()
		s1.Add(1).Add(2).Add(3)
		s2.Add(1).Add(2).Add(3)
		s3.Add(1).Add(2).Add(3).Add(4)
		t.Assert(s1.Equal(s2), true)
		t.Assert(s1.Equal(s3), false)
	})
}

func TestIntSet_IsSubsetOf(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := gset.NewIntSet()
		s2 := gset.NewIntSet()
		s3 := gset.NewIntSet()
		s1.Add(1).Add(2)
		s2.Add(1).Add(2).Add(3)
		s3.Add(1).Add(2).Add(3).Add(4)
		t.Assert(s1.IsSubsetOf(s2), true)
		t.Assert(s2.IsSubsetOf(s3), true)
		t.Assert(s1.IsSubsetOf(s3), true)
		t.Assert(s2.IsSubsetOf(s1), false)
		t.Assert(s3.IsSubsetOf(s2), false)
	})
}

func TestIntSet_Union(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := gset.NewIntSet()
		s2 := gset.NewIntSet()
		s1.Add(1).Add(2)
		s2.Add(3).Add(4)
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
		s1.Add(1).Add(2).Add(3)
		s2.Add(3).Add(4).Add(5)
		s3 := s1.Diff(s2)
		t.Assert(s3.Contains(1), true)
		t.Assert(s3.Contains(2), true)
		t.Assert(s3.Contains(3), false)
		t.Assert(s3.Contains(4), false)
	})
}

func TestIntSet_Intersect(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := gset.NewIntSet()
		s2 := gset.NewIntSet()
		s1.Add(1).Add(2).Add(3)
		s2.Add(3).Add(4).Add(5)
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
		s1.Add(1).Add(2).Add(3)
		s2.Add(3).Add(4).Add(5)
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
		s1.Add(1).Add(2).Add(3)
		t.Assert(s1.Size(), 3)

	})

}

func TestIntSet_Merge(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := gset.NewIntSet()
		s2 := gset.NewIntSet()
		s1.Add(1).Add(2).Add(3)
		s2.Add(3).Add(4).Add(5)
		s3 := s1.Merge(s2)
		t.Assert(s3.Contains(1), true)
		t.Assert(s3.Contains(5), true)
		t.Assert(s3.Contains(6), false)
	})
}

func TestIntSet_Join(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := gset.NewIntSet()
		s1.Add(1).Add(2).Add(3)
		s3 := s1.Join(",")
		t.Assert(strings.Contains(s3, "1"), true)
		t.Assert(strings.Contains(s3, "2"), true)
		t.Assert(strings.Contains(s3, "3"), true)
	})
}

func TestIntSet_String(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := gset.NewIntSet()
		s1.Add(1).Add(2).Add(3)
		s3 := s1.String()
		t.Assert(strings.Contains(s3, "["), true)
		t.Assert(strings.Contains(s3, "]"), true)
		t.Assert(strings.Contains(s3, "1"), true)
		t.Assert(strings.Contains(s3, "2"), true)
		t.Assert(strings.Contains(s3, "3"), true)
	})
}

func TestIntSet_Sum(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := gset.NewIntSet()
		s1.Add(1).Add(2).Add(3)
		s2 := gset.NewIntSet()
		s2.Add(5).Add(6).Add(7)
		t.Assert(s2.Sum(), 18)

	})

}

func TestIntSet_Pop(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := gset.NewIntSet()
		s.Add(4).Add(2).Add(3)
		t.Assert(s.Size(), 3)
		t.AssertIN(s.Pop(), []int{4, 2, 3})
		t.AssertIN(s.Pop(), []int{4, 2, 3})
		t.Assert(s.Size(), 1)
	})
}

func TestIntSet_Pops(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := gset.NewIntSet()
		s.Add(1).Add(4).Add(2).Add(3)
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

func TestIntSet_Json(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := []int{1, 3, 2, 4}
		a1 := gset.NewIntSetFrom(s1)
		b1, err1 := json.Marshal(a1)
		b2, err2 := json.Marshal(s1)
		t.Assert(len(b1), len(b2))
		t.Assert(err1, err2)

		a2 := gset.NewIntSet()
		err2 = json.Unmarshal(b2, &a2)
		t.Assert(err2, nil)
		t.Assert(a2.Contains(1), true)
		t.Assert(a2.Contains(2), true)
		t.Assert(a2.Contains(3), true)
		t.Assert(a2.Contains(4), true)
		t.Assert(a2.Contains(5), false)

		var a3 gset.IntSet
		err := json.Unmarshal(b2, &a3)
		t.Assert(err, nil)
		t.Assert(a2.Contains(1), true)
		t.Assert(a2.Contains(2), true)
		t.Assert(a2.Contains(3), true)
		t.Assert(a2.Contains(4), true)
		t.Assert(a2.Contains(5), false)
	})
}

func TestIntSet_AddIfNotExistFunc(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := gset.NewIntSet(true)
		s.Add(1)
		t.Assert(s.Contains(1), true)
		t.Assert(s.Contains(2), false)

		s.AddIfNotExistFunc(2, func() int {
			return 3
		})
		t.Assert(s.Contains(2), false)
		t.Assert(s.Contains(3), true)

		s.AddIfNotExistFunc(3, func() int {
			return 4
		})
		t.Assert(s.Contains(3), true)
		t.Assert(s.Contains(4), false)
	})

	gtest.C(t, func(t *gtest.T) {
		s := gset.NewIntSet(true)
		s.Add(1)
		t.Assert(s.Contains(1), true)
		t.Assert(s.Contains(2), false)

		s.AddIfNotExistFuncLock(2, func() int {
			return 3
		})
		t.Assert(s.Contains(2), false)
		t.Assert(s.Contains(3), true)

		s.AddIfNotExistFuncLock(3, func() int {
			return 4
		})
		t.Assert(s.Contains(3), true)
		t.Assert(s.Contains(4), false)
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
		t.Assert(err, nil)
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
		t.Assert(err, nil)
		t.Assert(v.Name, "john")
		t.Assert(v.Set.Size(), 3)
		t.Assert(v.Set.Contains(1), true)
		t.Assert(v.Set.Contains(2), true)
		t.Assert(v.Set.Contains(3), true)
		t.Assert(v.Set.Contains(4), false)
	})
}
