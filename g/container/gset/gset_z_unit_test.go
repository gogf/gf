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

func TestSet_New(t *testing.T) {
	gtest.Case(t, func() {
		s := gset.New()
		s.Add(1).Add(1).Add(2)
		s.Add([]interface{}{3, 4}...)
		gtest.Assert(s.Size(), 4)
		gtest.AssertIN(1, s.Slice())
		gtest.AssertIN(2, s.Slice())
		gtest.AssertIN(3, s.Slice())
		gtest.AssertIN(4, s.Slice())
		gtest.AssertNI(0, s.Slice())
		gtest.Assert(s.Contains(4), true)
		gtest.Assert(s.Contains(5), false)
		s.Remove(1)
		gtest.Assert(s.Size(), 3)
		s.Clear()
		gtest.Assert(s.Size(), 0)
	})
}

func TestSet_Basic(t *testing.T) {
	gtest.Case(t, func() {
		s := gset.NewSet()
		s.Add(1).Add(1).Add(2)
		s.Add([]interface{}{3, 4}...)
		gtest.Assert(s.Size(), 4)
		gtest.AssertIN(1, s.Slice())
		gtest.AssertIN(2, s.Slice())
		gtest.AssertIN(3, s.Slice())
		gtest.AssertIN(4, s.Slice())
		gtest.AssertNI(0, s.Slice())
		gtest.Assert(s.Contains(4), true)
		gtest.Assert(s.Contains(5), false)
		s.Remove(1)
		gtest.Assert(s.Size(), 3)
		s.Clear()
		gtest.Assert(s.Size(), 0)
	})
}

func TestSet_Iterator(t *testing.T) {
	gtest.Case(t, func() {
		s := gset.NewSet()
		s.Add(1).Add(2).Add(3)
		gtest.Assert(s.Size(), 3)

		a1 := garray.New()
		a2 := garray.New()
		s.Iterator(func(v interface{}) bool {
			a1.Append(1)
			return false
		})
		s.Iterator(func(v interface{}) bool {
			a2.Append(1)
			return true
		})
		gtest.Assert(a1.Len(), 1)
		gtest.Assert(a2.Len(), 3)
	})
}

func TestSet_LockFunc(t *testing.T) {
	gtest.Case(t, func() {
		s := gset.NewSet()
		s.Add(1).Add(2).Add(3)
		gtest.Assert(s.Size(), 3)
		s.LockFunc(func(m map[interface{}]struct{}) {
			delete(m, 1)
		})
		gtest.Assert(s.Size(), 2)
		s.RLockFunc(func(m map[interface{}]struct{}) {
			gtest.Assert(m, map[interface{}]struct{}{
				3: struct{}{},
				2: struct{}{},
			})
		})
	})
}

func TestSet_Equal(t *testing.T) {
	gtest.Case(t, func() {
		s1 := gset.NewSet()
		s2 := gset.NewSet()
		s3 := gset.NewSet()
		s1.Add(1).Add(2).Add(3)
		s2.Add(1).Add(2).Add(3)
		s3.Add(1).Add(2).Add(3).Add(4)
		gtest.Assert(s1.Equal(s2), true)
		gtest.Assert(s1.Equal(s3), false)
	})
}

func TestSet_IsSubsetOf(t *testing.T) {
	gtest.Case(t, func() {
		s1 := gset.NewSet()
		s2 := gset.NewSet()
		s3 := gset.NewSet()
		s1.Add(1).Add(2)
		s2.Add(1).Add(2).Add(3)
		s3.Add(1).Add(2).Add(3).Add(4)
		gtest.Assert(s1.IsSubsetOf(s2), true)
		gtest.Assert(s2.IsSubsetOf(s3), true)
		gtest.Assert(s1.IsSubsetOf(s3), true)
		gtest.Assert(s2.IsSubsetOf(s1), false)
		gtest.Assert(s3.IsSubsetOf(s2), false)
	})
}

func TestSet_Union(t *testing.T) {
	gtest.Case(t, func() {
		s1 := gset.NewSet()
		s2 := gset.NewSet()
		s1.Add(1).Add(2)
		s2.Add(3).Add(4)
		s3 := s1.Union(s2)
		gtest.Assert(s3.Contains(1), true)
		gtest.Assert(s3.Contains(2), true)
		gtest.Assert(s3.Contains(3), true)
		gtest.Assert(s3.Contains(4), true)
	})
}

func TestSet_Diff(t *testing.T) {
	gtest.Case(t, func() {
		s1 := gset.NewSet()
		s2 := gset.NewSet()
		s1.Add(1).Add(2).Add(3)
		s2.Add(3).Add(4).Add(5)
		s3 := s1.Diff(s2)
		gtest.Assert(s3.Contains(1), true)
		gtest.Assert(s3.Contains(2), true)
		gtest.Assert(s3.Contains(3), false)
		gtest.Assert(s3.Contains(4), false)
	})
}

func TestSet_Intersect(t *testing.T) {
	gtest.Case(t, func() {
		s1 := gset.NewSet()
		s2 := gset.NewSet()
		s1.Add(1).Add(2).Add(3)
		s2.Add(3).Add(4).Add(5)
		s3 := s1.Intersect(s2)
		gtest.Assert(s3.Contains(1), false)
		gtest.Assert(s3.Contains(2), false)
		gtest.Assert(s3.Contains(3), true)
		gtest.Assert(s3.Contains(4), false)
	})
}

func TestSet_Complement(t *testing.T) {
	gtest.Case(t, func() {
		s1 := gset.NewSet()
		s2 := gset.NewSet()
		s1.Add(1).Add(2).Add(3)
		s2.Add(3).Add(4).Add(5)
		s3 := s1.Complement(s2)
		gtest.Assert(s3.Contains(1), false)
		gtest.Assert(s3.Contains(2), false)
		gtest.Assert(s3.Contains(4), true)
		gtest.Assert(s3.Contains(5), true)
	})
}

func TestNewFrom(t *testing.T) {
	gtest.Case(t, func() {
		s1 := gset.NewFrom("a")
		s2 := gset.NewFrom("b", false)
		s3 := gset.NewFrom(3, true)
		s4 := gset.NewFrom([]string{"s1", "s2"}, true)
		gtest.Assert(s1.Contains("a"), true)
		gtest.Assert(s2.Contains("b"), true)
		gtest.Assert(s3.Contains(3), true)
		gtest.Assert(s4.Contains("s1"), true)
		gtest.Assert(s4.Contains("s3"), false)

	})
}

func TestNew(t *testing.T) {
	gtest.Case(t, func() {
		s1 := gset.New()
		s1.Add("a").Add(2)
		s2 := gset.New(true)
		s2.Add("b").Add(3)
		gtest.Assert(s1.Contains("a"), true)

	})
}

func TestSet_Join(t *testing.T) {
	gtest.Case(t, func() {
		s1 := gset.New(true)
		s1.Add("a").Add("a1").Add("b").Add("c")
		str1 := s1.Join(",")
		gtest.Assert(strings.Contains(str1, "a1"), true)

	})
}

func TestSet_String(t *testing.T) {
	gtest.Case(t, func() {
		s1 := gset.New(true)
		s1.Add("a").Add("a2").Add("b").Add("c")
		str1 := s1.String()
		gtest.Assert(strings.Contains(str1, "a2"), true)

	})
}

func TestSet_Merge(t *testing.T) {
	gtest.Case(t, func() {
		s1 := gset.New(true)
		s2 := gset.New(true)
		s1.Add("a").Add("a2").Add("b").Add("c")
		s2.Add("b").Add("b1").Add("e").Add("f")
		ss := s1.Merge(s2)
		gtest.Assert(ss.Contains("a2"), true)
		gtest.Assert(ss.Contains("b1"), true)

	})
}

func TestSet_Sum(t *testing.T) {
	gtest.Case(t, func() {
		s1 := gset.New(true)
		s1.Add(1).Add(2).Add(3).Add(4)
		gtest.Assert(s1.Sum(), int(10))

	})
}

func TestSet_Pop(t *testing.T) {
	gtest.Case(t, func() {
		s1 := gset.New(true)
		s1.Add(1).Add(2).Add(3).Add(4)
		gtest.AssertIN(s1.Pop(1), []int{1, 2, 3, 4})
	})
}

func TestSet_Pops(t *testing.T) {
	gtest.Case(t, func() {
		s1 := gset.New(true)
		s1.Add(1).Add(2).Add(3).Add(4)
		gtest.AssertIN(s1.Pops(1), []int{1, 2, 3, 4})
		gtest.AssertIN(s1.Pops(6), []int{1, 2, 3, 4})
		gtest.Assert(len(s1.Pops(2)), 2)
	})
}
