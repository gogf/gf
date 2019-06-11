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
    "testing"
)

func TestStringSet_Basic(t *testing.T) {
    gtest.Case(t, func() {
        s := gset.NewStringSet()
        s.Add("1").Add("1").Add("2")
        s.Add([]string{"3","4"}...)
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
                "3" : struct{}{},
                "2" : struct{}{},
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