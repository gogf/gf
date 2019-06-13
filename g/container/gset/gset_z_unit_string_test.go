// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go

package gset_test

import (
	"github.com/gogf/gf/g/container/garray"
	"github.com/gogf/gf/g/container/gset"
	"github.com/gogf/gf/g/test/gtest"
	"strings"
	"testing"
)

func TestStringSet_Basic(t *testing.T) {
	gtest.Case(t, func() {
		s := gset.NewStringSet()
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

func TestStringSet_Iterator(t *testing.T) {
	gtest.Case(t, func() {
		s := gset.NewStringSet()
		s.Add("1").Add("2").Add("3")
		gtest.Assert(s.Size(), 3)

		a1 := garray.New()
		a2 := garray.New()
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

func TestStringSet_LockFunc(t *testing.T) {
	gtest.Case(t, func() {
		s := gset.NewStringSet()
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

func TestStringSet_Equal(t *testing.T) {
	gtest.Case(t, func() {
		s1 := gset.NewStringSet()
		s2 := gset.NewStringSet()
		s3 := gset.NewStringSet()
		s1.Add("1").Add("2").Add("3")
		s2.Add("1").Add("2").Add("3")
		s3.Add("1").Add("2").Add("3").Add("4")
		gtest.Assert(s1.Equal(s2), true)
		gtest.Assert(s1.Equal(s3), false)
	})
}

func TestStringSet_IsSubsetOf(t *testing.T) {
	gtest.Case(t, func() {
		s1 := gset.NewStringSet()
		s2 := gset.NewStringSet()
		s3 := gset.NewStringSet()
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

func TestStringSet_Union(t *testing.T) {
	gtest.Case(t, func() {
		s1 := gset.NewStringSet()
		s2 := gset.NewStringSet()
		s1.Add("1").Add("2")
		s2.Add("3").Add("4")
		s3 := s1.Union(s2)
		gtest.Assert(s3.Contains("1"), true)
		gtest.Assert(s3.Contains("2"), true)
		gtest.Assert(s3.Contains("3"), true)
		gtest.Assert(s3.Contains("4"), true)
	})
}

func TestStringSet_Diff(t *testing.T) {
	gtest.Case(t, func() {
		s1 := gset.NewStringSet()
		s2 := gset.NewStringSet()
		s1.Add("1").Add("2").Add("3")
		s2.Add("3").Add("4").Add("5")
		s3 := s1.Diff(s2)
		gtest.Assert(s3.Contains("1"), true)
		gtest.Assert(s3.Contains("2"), true)
		gtest.Assert(s3.Contains("3"), false)
		gtest.Assert(s3.Contains("4"), false)
	})
}

func TestStringSet_Intersect(t *testing.T) {
	gtest.Case(t, func() {
		s1 := gset.NewStringSet()
		s2 := gset.NewStringSet()
		s1.Add("1").Add("2").Add("3")
		s2.Add("3").Add("4").Add("5")
		s3 := s1.Intersect(s2)
		gtest.Assert(s3.Contains("1"), false)
		gtest.Assert(s3.Contains("2"), false)
		gtest.Assert(s3.Contains("3"), true)
		gtest.Assert(s3.Contains("4"), false)
	})
}

func TestStringSet_Complement(t *testing.T) {
	gtest.Case(t, func() {
		s1 := gset.NewStringSet()
		s2 := gset.NewStringSet()
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

func TestStringSet_Merge(t *testing.T) {
	gtest.Case(t, func() {
		s1 := gset.NewStringSet()
		s2 := gset.NewStringSet()
		s1.Add("1").Add("2").Add("3")
		s2.Add("3").Add("4").Add("5")
		s3 := s1.Merge(s2)
		gtest.Assert(s3.Contains("1"), true)
		gtest.Assert(s3.Contains("6"), false)
		gtest.Assert(s3.Contains("4"), true)
		gtest.Assert(s3.Contains("5"), true)
	})
}

func TestNewStringSetFrom(t *testing.T) {
	gtest.Case(t, func() {
		s1 := gset.NewStringSetFrom([]string{"a", "b", "c"}, true)
		gtest.Assert(s1.Contains("b"), true)
		gtest.Assert(s1.Contains("d"), false)
	})
}

func TestStringSet_Join(t *testing.T) {
	s1 := gset.NewStringSetFrom([]string{"a", "b", "c"}, true)
	str1 := s1.Join(",")
	gtest.Assert(strings.Contains(str1, "b"), true)
	gtest.Assert(strings.Contains(str1, "d"), false)
}

func TestStringSet_String(t *testing.T) {
	gtest.Case(t, func() {
		s1 := gset.NewStringSetFrom([]string{"a", "b", "c"}, true)
		str1 := s1.String()
		gtest.Assert(strings.Contains(str1, "b"), true)
		gtest.Assert(strings.Contains(str1, "d"), false)
	})

}

func TestStringSet_Sum(t *testing.T) {
	gtest.Case(t, func() {
		s1 := gset.NewStringSetFrom([]string{"a", "b", "c"}, true)
		s2 := gset.NewIntSetFrom([]int{2, 3, 4}, true)
		gtest.Assert(s1.Sum(), 0)
		gtest.Assert(s2.Sum(), 9)
	})
}

func TestStringSet_Size(t *testing.T) {
	gtest.Case(t, func() {
		s1 := gset.NewStringSetFrom([]string{"a", "b", "c"}, true)
		gtest.Assert(s1.Size(), 3)

	})
}

func TestStringSet_Remove(t *testing.T) {
	gtest.Case(t, func() {
		s1 := gset.NewStringSetFrom([]string{"a", "b", "c"}, true)
		s1 = s1.Remove("b")
		gtest.Assert(s1.Contains("b"), false)
		gtest.Assert(s1.Contains("c"), true)
	})
}

func TestStringSet_Pop(t *testing.T) {
	gtest.Case(t, func() {
		s1 := gset.NewStringSetFrom([]string{"a", "b", "c"}, true)
		str1 := s1.Pop(1)
		gtest.Assert(strings.Contains("a,b,c", str1), true)
	})
}

func TestStringSet_Pops(t *testing.T) {
	gtest.Case(t, func() {
		s1 := gset.NewStringSetFrom([]string{"a", "b", "c"}, true)
		strs1 := s1.Pops(2)
		gtest.AssertIN(strs1, []string{"a", "b", "c"})
		gtest.Assert(len(strs1), 2)
		str2 := s1.Pops(7)
		gtest.AssertIN(str2, []string{"a", "b", "c"})
	})
}
