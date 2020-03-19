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

	"github.com/gogf/gf/container/garray"
	"github.com/gogf/gf/container/gset"
	"github.com/gogf/gf/test/gtest"

	"testing"
)

func TestSet_New(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := gset.New()
		s.Add(1).Add(1).Add(2)
		s.Add([]interface{}{3, 4}...)
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

func TestSet_Basic(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := gset.NewSet()
		s.Add(1).Add(1).Add(2)
		s.Add([]interface{}{3, 4}...)
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

func TestSet_Iterator(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := gset.NewSet()
		s.Add(1).Add(2).Add(3)
		t.Assert(s.Size(), 3)

		a1 := garray.New(true)
		a2 := garray.New(true)
		s.Iterator(func(v interface{}) bool {
			a1.Append(1)
			return false
		})
		s.Iterator(func(v interface{}) bool {
			a2.Append(1)
			return true
		})
		t.Assert(a1.Len(), 1)
		t.Assert(a2.Len(), 3)
	})
}

func TestSet_LockFunc(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := gset.NewSet()
		s.Add(1).Add(2).Add(3)
		t.Assert(s.Size(), 3)
		s.LockFunc(func(m map[interface{}]struct{}) {
			delete(m, 1)
		})
		t.Assert(s.Size(), 2)
		s.RLockFunc(func(m map[interface{}]struct{}) {
			t.Assert(m, map[interface{}]struct{}{
				3: struct{}{},
				2: struct{}{},
			})
		})
	})
}

func TestSet_Equal(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := gset.NewSet()
		s2 := gset.NewSet()
		s3 := gset.NewSet()
		s1.Add(1).Add(2).Add(3)
		s2.Add(1).Add(2).Add(3)
		s3.Add(1).Add(2).Add(3).Add(4)
		t.Assert(s1.Equal(s2), true)
		t.Assert(s1.Equal(s3), false)
	})
}

func TestSet_IsSubsetOf(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := gset.NewSet()
		s2 := gset.NewSet()
		s3 := gset.NewSet()
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

func TestSet_Union(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := gset.NewSet()
		s2 := gset.NewSet()
		s1.Add(1).Add(2)
		s2.Add(3).Add(4)
		s3 := s1.Union(s2)
		t.Assert(s3.Contains(1), true)
		t.Assert(s3.Contains(2), true)
		t.Assert(s3.Contains(3), true)
		t.Assert(s3.Contains(4), true)
	})
}

func TestSet_Diff(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := gset.NewSet()
		s2 := gset.NewSet()
		s1.Add(1).Add(2).Add(3)
		s2.Add(3).Add(4).Add(5)
		s3 := s1.Diff(s2)
		t.Assert(s3.Contains(1), true)
		t.Assert(s3.Contains(2), true)
		t.Assert(s3.Contains(3), false)
		t.Assert(s3.Contains(4), false)
	})
}

func TestSet_Intersect(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := gset.NewSet()
		s2 := gset.NewSet()
		s1.Add(1).Add(2).Add(3)
		s2.Add(3).Add(4).Add(5)
		s3 := s1.Intersect(s2)
		t.Assert(s3.Contains(1), false)
		t.Assert(s3.Contains(2), false)
		t.Assert(s3.Contains(3), true)
		t.Assert(s3.Contains(4), false)
	})
}

func TestSet_Complement(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := gset.NewSet()
		s2 := gset.NewSet()
		s1.Add(1).Add(2).Add(3)
		s2.Add(3).Add(4).Add(5)
		s3 := s1.Complement(s2)
		t.Assert(s3.Contains(1), false)
		t.Assert(s3.Contains(2), false)
		t.Assert(s3.Contains(4), true)
		t.Assert(s3.Contains(5), true)
	})
}

func TestNewFrom(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := gset.NewFrom("a")
		s2 := gset.NewFrom("b", false)
		s3 := gset.NewFrom(3, true)
		s4 := gset.NewFrom([]string{"s1", "s2"}, true)
		t.Assert(s1.Contains("a"), true)
		t.Assert(s2.Contains("b"), true)
		t.Assert(s3.Contains(3), true)
		t.Assert(s4.Contains("s1"), true)
		t.Assert(s4.Contains("s3"), false)

	})
}

func TestNew(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := gset.New()
		s1.Add("a").Add(2)
		s2 := gset.New(true)
		s2.Add("b").Add(3)
		t.Assert(s1.Contains("a"), true)

	})
}

func TestSet_Join(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := gset.New(true)
		s1.Add("a").Add("a1").Add("b").Add("c")
		str1 := s1.Join(",")
		t.Assert(strings.Contains(str1, "a1"), true)
	})
	gtest.C(t, func(t *gtest.T) {
		s1 := gset.New(true)
		s1.Add("a").Add(`"b"`).Add(`\c`)
		str1 := s1.Join(",")
		t.Assert(strings.Contains(str1, `"b"`), true)
		t.Assert(strings.Contains(str1, `\c`), true)
		t.Assert(strings.Contains(str1, `a`), true)
	})
}

func TestSet_String(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := gset.New(true)
		s1.Add("a").Add("a2").Add("b").Add("c")
		str1 := s1.String()
		t.Assert(strings.Contains(str1, "["), true)
		t.Assert(strings.Contains(str1, "]"), true)
		t.Assert(strings.Contains(str1, "a2"), true)
	})
}

func TestSet_Merge(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := gset.New(true)
		s2 := gset.New(true)
		s1.Add("a").Add("a2").Add("b").Add("c")
		s2.Add("b").Add("b1").Add("e").Add("f")
		ss := s1.Merge(s2)
		t.Assert(ss.Contains("a2"), true)
		t.Assert(ss.Contains("b1"), true)

	})
}

func TestSet_Sum(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := gset.New(true)
		s1.Add(1).Add(2).Add(3).Add(4)
		t.Assert(s1.Sum(), int(10))

	})
}

func TestSet_Pop(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := gset.New(true)
		s.Add(1).Add(2).Add(3).Add(4)
		t.Assert(s.Size(), 4)
		t.AssertIN(s.Pop(), []int{1, 2, 3, 4})
		t.Assert(s.Size(), 3)
	})
}

func TestSet_Pops(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := gset.New(true)
		s.Add(1).Add(2).Add(3).Add(4)
		t.Assert(s.Size(), 4)
		t.Assert(s.Pops(0), nil)
		t.AssertIN(s.Pops(1), []int{1, 2, 3, 4})
		t.Assert(s.Size(), 3)
		a := s.Pops(6)
		t.Assert(len(a), 3)
		t.AssertIN(a, []int{1, 2, 3, 4})
		t.Assert(s.Size(), 0)
	})

	gtest.C(t, func(t *gtest.T) {
		s := gset.New(true)
		a := []interface{}{1, 2, 3, 4}
		s.Add(a...)
		t.Assert(s.Size(), 4)
		t.Assert(s.Pops(-2), nil)
		t.AssertIN(s.Pops(-1), a)
	})
}

func TestSet_Json(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := []interface{}{"a", "b", "d", "c"}
		a1 := gset.NewFrom(s1)
		b1, err1 := json.Marshal(a1)
		b2, err2 := json.Marshal(s1)
		t.Assert(len(b1), len(b2))
		t.Assert(err1, err2)

		a2 := gset.New()
		err2 = json.Unmarshal(b2, &a2)
		t.Assert(err2, nil)
		t.Assert(a2.Contains("a"), true)
		t.Assert(a2.Contains("b"), true)
		t.Assert(a2.Contains("c"), true)
		t.Assert(a2.Contains("d"), true)
		t.Assert(a2.Contains("e"), false)

		var a3 gset.Set
		err := json.Unmarshal(b2, &a3)
		t.Assert(err, nil)
		t.Assert(a3.Contains("a"), true)
		t.Assert(a3.Contains("b"), true)
		t.Assert(a3.Contains("c"), true)
		t.Assert(a3.Contains("d"), true)
		t.Assert(a3.Contains("e"), false)
	})
}

func TestSet_AddIfNotExistFunc(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := gset.New(true)
		s.Add(1)
		t.Assert(s.Contains(1), true)
		t.Assert(s.Contains(2), false)

		s.AddIfNotExistFunc(2, func() interface{} {
			return 3
		})
		t.Assert(s.Contains(2), false)
		t.Assert(s.Contains(3), true)

		s.AddIfNotExistFunc(3, func() interface{} {
			return 4
		})
		t.Assert(s.Contains(3), true)
		t.Assert(s.Contains(4), false)
	})

	gtest.C(t, func(t *gtest.T) {
		s := gset.New(true)
		s.Add(1)
		t.Assert(s.Contains(1), true)
		t.Assert(s.Contains(2), false)

		s.AddIfNotExistFuncLock(2, func() interface{} {
			return 3
		})
		t.Assert(s.Contains(2), false)
		t.Assert(s.Contains(3), true)

		s.AddIfNotExistFuncLock(3, func() interface{} {
			return 4
		})
		t.Assert(s.Contains(3), true)
		t.Assert(s.Contains(4), false)
	})
}

func TestSet_UnmarshalValue(t *testing.T) {
	type T struct {
		Name string
		Set  *gset.Set
	}
	// JSON
	gtest.C(t, func(t *gtest.T) {
		var t *T
		err := gconv.Struct(g.Map{
			"name": "john",
			"set":  []byte(`["k1","k2","k3"]`),
		}, &t)
		t.Assert(err, nil)
		t.Assert(t.Name, "john")
		t.Assert(t.Set.Size(), 3)
		t.Assert(t.Set.Contains("k1"), true)
		t.Assert(t.Set.Contains("k2"), true)
		t.Assert(t.Set.Contains("k3"), true)
		t.Assert(t.Set.Contains("k4"), false)
	})
	// Map
	gtest.C(t, func(t *gtest.T) {
		var t *T
		err := gconv.Struct(g.Map{
			"name": "john",
			"set":  g.Slice{"k1", "k2", "k3"},
		}, &t)
		t.Assert(err, nil)
		t.Assert(t.Name, "john")
		t.Assert(t.Set.Size(), 3)
		t.Assert(t.Set.Contains("k1"), true)
		t.Assert(t.Set.Contains("k2"), true)
		t.Assert(t.Set.Contains("k3"), true)
		t.Assert(t.Set.Contains("k4"), false)
	})
}
