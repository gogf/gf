// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go

package gset_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/gogf/gf/container/garray"
	"github.com/gogf/gf/container/gset"
	"github.com/gogf/gf/test/gtest"
)

func TestStrSet_Basic(t *testing.T) {
	gtest.Case(t, func() {
		s := gset.NewStrSet()
		s.Add("1").Add("1").Add("2")
		s.Add([]string{"3", "4"}...)
		gtest.Assert(s.Size(), 4)
		gtest.AssertIN("1", s.Slice())
		gtest.AssertIN("2", s.Slice())
		gtest.AssertIN("3", s.Slice())
		gtest.AssertIN("4", s.Slice())
		gtest.AssertNI("0", s.Slice())
		gtest.Assert(s.Contains("4"), true)
		gtest.Assert(s.Contains("5"), false)
		s.Remove("1")
		gtest.Assert(s.Size(), 3)
		s.Clear()
		gtest.Assert(s.Size(), 0)
	})
}

func TestStrSet_Iterator(t *testing.T) {
	gtest.Case(t, func() {
		s := gset.NewStrSet()
		s.Add("1").Add("2").Add("3")
		gtest.Assert(s.Size(), 3)

		a1 := garray.New(true)
		a2 := garray.New(true)
		s.Iterator(func(v string) bool {
			a1.Append("1")
			return false
		})
		s.Iterator(func(v string) bool {
			a2.Append("1")
			return true
		})
		gtest.Assert(a1.Len(), 1)
		gtest.Assert(a2.Len(), 3)
	})
}

func TestStrSet_LockFunc(t *testing.T) {
	gtest.Case(t, func() {
		s := gset.NewStrSet()
		s.Add("1").Add("2").Add("3")
		gtest.Assert(s.Size(), 3)
		s.LockFunc(func(m map[string]struct{}) {
			delete(m, "1")
		})
		gtest.Assert(s.Size(), 2)
		s.RLockFunc(func(m map[string]struct{}) {
			gtest.Assert(m, map[string]struct{}{
				"3": struct{}{},
				"2": struct{}{},
			})
		})
	})
}

func TestStrSet_Equal(t *testing.T) {
	gtest.Case(t, func() {
		s1 := gset.NewStrSet()
		s2 := gset.NewStrSet()
		s3 := gset.NewStrSet()
		s1.Add("1").Add("2").Add("3")
		s2.Add("1").Add("2").Add("3")
		s3.Add("1").Add("2").Add("3").Add("4")
		gtest.Assert(s1.Equal(s2), true)
		gtest.Assert(s1.Equal(s3), false)
	})
}

func TestStrSet_IsSubsetOf(t *testing.T) {
	gtest.Case(t, func() {
		s1 := gset.NewStrSet()
		s2 := gset.NewStrSet()
		s3 := gset.NewStrSet()
		s1.Add("1").Add("2")
		s2.Add("1").Add("2").Add("3")
		s3.Add("1").Add("2").Add("3").Add("4")
		gtest.Assert(s1.IsSubsetOf(s2), true)
		gtest.Assert(s2.IsSubsetOf(s3), true)
		gtest.Assert(s1.IsSubsetOf(s3), true)
		gtest.Assert(s2.IsSubsetOf(s1), false)
		gtest.Assert(s3.IsSubsetOf(s2), false)
	})
}

func TestStrSet_Union(t *testing.T) {
	gtest.Case(t, func() {
		s1 := gset.NewStrSet()
		s2 := gset.NewStrSet()
		s1.Add("1").Add("2")
		s2.Add("3").Add("4")
		s3 := s1.Union(s2)
		gtest.Assert(s3.Contains("1"), true)
		gtest.Assert(s3.Contains("2"), true)
		gtest.Assert(s3.Contains("3"), true)
		gtest.Assert(s3.Contains("4"), true)
	})
}

func TestStrSet_Diff(t *testing.T) {
	gtest.Case(t, func() {
		s1 := gset.NewStrSet()
		s2 := gset.NewStrSet()
		s1.Add("1").Add("2").Add("3")
		s2.Add("3").Add("4").Add("5")
		s3 := s1.Diff(s2)
		gtest.Assert(s3.Contains("1"), true)
		gtest.Assert(s3.Contains("2"), true)
		gtest.Assert(s3.Contains("3"), false)
		gtest.Assert(s3.Contains("4"), false)
	})
}

func TestStrSet_Intersect(t *testing.T) {
	gtest.Case(t, func() {
		s1 := gset.NewStrSet()
		s2 := gset.NewStrSet()
		s1.Add("1").Add("2").Add("3")
		s2.Add("3").Add("4").Add("5")
		s3 := s1.Intersect(s2)
		gtest.Assert(s3.Contains("1"), false)
		gtest.Assert(s3.Contains("2"), false)
		gtest.Assert(s3.Contains("3"), true)
		gtest.Assert(s3.Contains("4"), false)
	})
}

func TestStrSet_Complement(t *testing.T) {
	gtest.Case(t, func() {
		s1 := gset.NewStrSet()
		s2 := gset.NewStrSet()
		s1.Add("1").Add("2").Add("3")
		s2.Add("3").Add("4").Add("5")
		s3 := s1.Complement(s2)
		gtest.Assert(s3.Contains("1"), false)
		gtest.Assert(s3.Contains("2"), false)
		gtest.Assert(s3.Contains("4"), true)
		gtest.Assert(s3.Contains("5"), true)
	})
}

func TestNewIntSetFrom(t *testing.T) {
	gtest.Case(t, func() {
		s1 := gset.NewIntSetFrom([]int{1, 2, 3, 4})
		s2 := gset.NewIntSetFrom([]int{5, 6, 7, 8})
		gtest.Assert(s1.Contains(3), true)
		gtest.Assert(s1.Contains(5), false)
		gtest.Assert(s2.Contains(3), false)
		gtest.Assert(s2.Contains(5), true)
	})
}

func TestStrSet_Merge(t *testing.T) {
	gtest.Case(t, func() {
		s1 := gset.NewStrSet()
		s2 := gset.NewStrSet()
		s1.Add("1").Add("2").Add("3")
		s2.Add("3").Add("4").Add("5")
		s3 := s1.Merge(s2)
		gtest.Assert(s3.Contains("1"), true)
		gtest.Assert(s3.Contains("6"), false)
		gtest.Assert(s3.Contains("4"), true)
		gtest.Assert(s3.Contains("5"), true)
	})
}

func TestNewStrSetFrom(t *testing.T) {
	gtest.Case(t, func() {
		s1 := gset.NewStrSetFrom([]string{"a", "b", "c"}, true)
		gtest.Assert(s1.Contains("b"), true)
		gtest.Assert(s1.Contains("d"), false)
	})
}

func TestStrSet_Join(t *testing.T) {
	gtest.Case(t, func() {
		s1 := gset.NewStrSetFrom([]string{"a", "b", "c"}, true)
		str1 := s1.Join(",")
		gtest.Assert(strings.Contains(str1, "b"), true)
		gtest.Assert(strings.Contains(str1, "d"), false)
	})

	gtest.Case(t, func() {
		s1 := gset.NewStrSet()
		s1.Add("a").Add(`"b"`).Add(`\c`)
		str1 := s1.Join(",")
		gtest.Assert(strings.Contains(str1, `\"b\"`), true)
		gtest.Assert(strings.Contains(str1, `\\c`), true)
		gtest.Assert(strings.Contains(str1, `a`), true)
	})
}

func TestStrSet_String(t *testing.T) {
	gtest.Case(t, func() {
		s1 := gset.NewStrSetFrom([]string{"a", "b", "c"}, true)
		str1 := s1.String()
		gtest.Assert(strings.Contains(str1, "b"), true)
		gtest.Assert(strings.Contains(str1, "d"), false)
	})

	gtest.Case(t, func() {
		s1 := gset.New(true)
		s1.Add("a").Add("a2").Add("b").Add("c")
		str1 := s1.String()
		gtest.Assert(strings.Contains(str1, "["), true)
		gtest.Assert(strings.Contains(str1, "]"), true)
		gtest.Assert(strings.Contains(str1, "a2"), true)
	})
}

func TestStrSet_Sum(t *testing.T) {
	gtest.Case(t, func() {
		s1 := gset.NewStrSetFrom([]string{"a", "b", "c"}, true)
		s2 := gset.NewIntSetFrom([]int{2, 3, 4}, true)
		gtest.Assert(s1.Sum(), 0)
		gtest.Assert(s2.Sum(), 9)
	})
}

func TestStrSet_Size(t *testing.T) {
	gtest.Case(t, func() {
		s1 := gset.NewStrSetFrom([]string{"a", "b", "c"}, true)
		gtest.Assert(s1.Size(), 3)

	})
}

func TestStrSet_Remove(t *testing.T) {
	gtest.Case(t, func() {
		s1 := gset.NewStrSetFrom([]string{"a", "b", "c"}, true)
		s1 = s1.Remove("b")
		gtest.Assert(s1.Contains("b"), false)
		gtest.Assert(s1.Contains("c"), true)
	})
}

func TestStrSet_Pop(t *testing.T) {
	gtest.Case(t, func() {
		s1 := gset.NewStrSetFrom([]string{"a", "b", "c"}, true)
		str1 := s1.Pop()
		gtest.Assert(strings.Contains("a,b,c", str1), true)
	})
}

func TestStrSet_Pops(t *testing.T) {
	gtest.Case(t, func() {
		s1 := gset.NewStrSetFrom([]string{"a", "b", "c"}, true)
		strs1 := s1.Pops(2)
		gtest.AssertIN(strs1, []string{"a", "b", "c"})
		gtest.Assert(len(strs1), 2)
		str2 := s1.Pops(7)
		gtest.AssertIN(str2, []string{"a", "b", "c"})
	})
}

func TestStrSet_Json(t *testing.T) {
	gtest.Case(t, func() {
		s1 := []string{"a", "b", "d", "c"}
		a1 := gset.NewStrSetFrom(s1)
		b1, err1 := json.Marshal(a1)
		b2, err2 := json.Marshal(s1)
		gtest.Assert(len(b1), len(b2))
		gtest.Assert(err1, err2)

		a2 := gset.NewStrSet()
		err2 = json.Unmarshal(b2, &a2)
		gtest.Assert(err2, nil)
		gtest.Assert(a2.Contains("a"), true)
		gtest.Assert(a2.Contains("b"), true)
		gtest.Assert(a2.Contains("c"), true)
		gtest.Assert(a2.Contains("d"), true)
		gtest.Assert(a2.Contains("e"), false)

		var a3 gset.StrSet
		err := json.Unmarshal(b2, &a3)
		gtest.Assert(err, nil)
		gtest.Assert(a3.Contains("a"), true)
		gtest.Assert(a3.Contains("b"), true)
		gtest.Assert(a3.Contains("c"), true)
		gtest.Assert(a3.Contains("d"), true)
		gtest.Assert(a3.Contains("e"), false)
	})
}
