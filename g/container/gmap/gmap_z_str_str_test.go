// Copyright 2017-2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with gm file,
// You can obtain one at https://github.com/gogf/gf.

package gmap_test

import (
	"github.com/gogf/gf/g/container/gmap"
	"github.com/gogf/gf/g/test/gtest"
	"testing"
)

func stringStrCallBack(string, string) bool {
	return true
}
func Test_StrStrMap_Basic(t *testing.T) {
	gtest.Case(t, func() {
		m := gmap.NewStrStrMap()
		m.Set("a", "a")

		gtest.Assert(m.Get("a"), "a")
		gtest.Assert(m.Size(), 1)
		gtest.Assert(m.IsEmpty(), false)

		gtest.Assert(m.GetOrSet("b", "b"), "b")
		gtest.Assert(m.SetIfNotExist("b", "b"), false)

		gtest.Assert(m.SetIfNotExist("c", "c"), true)

		gtest.Assert(m.Remove("b"), "b")
		gtest.Assert(m.Contains("b"), false)

		gtest.AssertIN("c", m.Keys())
		gtest.AssertIN("a", m.Keys())
		gtest.AssertIN("a", m.Values())
		gtest.AssertIN("c", m.Values())

		m.Flip()

		gtest.Assert(m.Map(), map[string]string{"a": "a", "c": "c"})

		m.Clear()
		gtest.Assert(m.Size(), 0)
		gtest.Assert(m.IsEmpty(), true)

		m2 := gmap.NewStrStrMapFrom(map[string]string{"a": "a", "b": "b"})
		gtest.Assert(m2.Map(), map[string]string{"a": "a", "b": "b"})
	})
}
func Test_StrStrMap_Set_Fun(t *testing.T) {
	m := gmap.NewStrStrMap()

	m.GetOrSetFunc("a", getStr)
	m.GetOrSetFuncLock("b", getStr)
	gtest.Assert(m.Get("a"), "z")
	gtest.Assert(m.Get("b"), "z")
	gtest.Assert(m.SetIfNotExistFunc("a", getStr), false)
	gtest.Assert(m.SetIfNotExistFunc("c", getStr), true)

	gtest.Assert(m.SetIfNotExistFuncLock("b", getStr), false)
	gtest.Assert(m.SetIfNotExistFuncLock("d", getStr), true)

}

func Test_StrStrMap_Batch(t *testing.T) {
	m := gmap.NewStrStrMap()

	m.Sets(map[string]string{"a": "a", "b": "b", "c": "c"})
	gtest.Assert(m.Map(), map[string]string{"a": "a", "b": "b", "c": "c"})
	m.Removes([]string{"a", "b"})
	gtest.Assert(m.Map(), map[string]string{"c": "c"})
}
func Test_StrStrMap_Iterator(t *testing.T) {
	expect := map[string]string{"a": "a", "b": "b"}
	m := gmap.NewStrStrMapFrom(expect)
	m.Iterator(func(k string, v string) bool {
		gtest.Assert(expect[k], v)
		return true
	})
	// 断言返回值对遍历控制
	i := 0
	j := 0
	m.Iterator(func(k string, v string) bool {
		i++
		return true
	})
	m.Iterator(func(k string, v string) bool {
		j++
		return false
	})
	gtest.Assert(i, 2)
	gtest.Assert(j, 1)
}

func Test_StrStrMap_Lock(t *testing.T) {
	expect := map[string]string{"a": "a", "b": "b"}

	m      := gmap.NewStrStrMapFrom(expect)
	m.LockFunc(func(m map[string]string) {
		gtest.Assert(m, expect)
	})
	m.RLockFunc(func(m map[string]string) {
		gtest.Assert(m, expect)
	})
}
func Test_StrStrMap_Clone(t *testing.T) {
	//clone 方法是深克隆
	m := gmap.NewStrStrMapFrom(map[string]string{"a": "a", "b": "b", "c": "c"})

	m_clone := m.Clone()
	m.Remove("a")
	//修改原 map,clone 后的 map 不影响
	gtest.AssertIN("a", m_clone.Keys())

	m_clone.Remove("b")
	//修改clone map,原 map 不影响
	gtest.AssertIN("b", m.Keys())
}
func Test_StrStrMap_Merge(t *testing.T) {
	m1 := gmap.NewStrStrMap()
	m2 := gmap.NewStrStrMap()
	m1.Set("a", "a")
	m2.Set("b", "b")
	m1.Merge(m2)
	gtest.Assert(m1.Map(), map[string]string{"a": "a", "b": "b"})
}
