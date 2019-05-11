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

func stringAnyCallBack(string, interface{}) bool {
	return true
}
func Test_StrAnyMap_Basic(t *testing.T) {
	gtest.Case(t, func() {
		m := gmap.NewStrAnyMap()
		m.Set("a", 1)

		gtest.Assert(m.Get("a"), 1)
		gtest.Assert(m.Size(), 1)
		gtest.Assert(m.IsEmpty(), false)

		gtest.Assert(m.GetOrSet("b", "2"), "2")
		gtest.Assert(m.SetIfNotExist("b", "2"), false)

		gtest.Assert(m.SetIfNotExist("c", 3), true)

		gtest.Assert(m.Remove("b"), "2")
		gtest.Assert(m.Contains("b"), false)

		gtest.AssertIN("c", m.Keys())
		gtest.AssertIN("a", m.Keys())
		gtest.AssertIN(3, m.Values())
		gtest.AssertIN(1, m.Values())

		m.Flip()
		gtest.Assert(m.Map(), map[string]interface{}{"1": "a", "3": "c"})

		m.Clear()
		gtest.Assert(m.Size(), 0)
		gtest.Assert(m.IsEmpty(), true)

		m2 := gmap.NewStrAnyMapFrom(map[string]interface{}{"a": 1, "b": "2"})
		gtest.Assert(m2.Map(), map[string]interface{}{"a": 1, "b": "2"})
	})
}
func Test_StrAnyMap_Set_Fun(t *testing.T) {
	m := gmap.NewStrAnyMap()

	m.GetOrSetFunc("a", getAny)
	m.GetOrSetFuncLock("b", getAny)
	gtest.Assert(m.Get("a"), 123)
	gtest.Assert(m.Get("b"), 123)
	gtest.Assert(m.SetIfNotExistFunc("a", getAny), false)
	gtest.Assert(m.SetIfNotExistFunc("c", getAny), true)

	gtest.Assert(m.SetIfNotExistFuncLock("b", getAny), false)
	gtest.Assert(m.SetIfNotExistFuncLock("d", getAny), true)

}

func Test_StrAnyMap_Batch(t *testing.T) {
	m := gmap.NewStrAnyMap()

	m.Sets(map[string]interface{}{"a": 1, "b": "2", "c": 3})
	gtest.Assert(m.Map(), map[string]interface{}{"a": 1, "b": "2", "c": 3})
	m.Removes([]string{"a", "b"})
	gtest.Assert(m.Map(), map[string]interface{}{"c": 3})
}

func Test_StrAnyMap_Iterator(t *testing.T) {
	expect := map[string]interface{}{"a": true, "b": false}
	m := gmap.NewStrAnyMapFrom(expect)
	m.Iterator(func(k string, v interface{}) bool {
		gtest.Assert(expect[k], v)
		return true
	})
	// 断言返回值对遍历控制
	i := 0
	j := 0
	m.Iterator(func(k string, v interface{}) bool {
		i++
		return true
	})
	m.Iterator(func(k string, v interface{}) bool {
		j++
		return false
	})
	gtest.Assert(i, 2)
	gtest.Assert(j, 1)
}

func Test_StrAnyMap_Lock(t *testing.T) {
	expect := map[string]interface{}{"a": true, "b": false}

	m      := gmap.NewStrAnyMapFrom(expect)
	m.LockFunc(func(m map[string]interface{}) {
		gtest.Assert(m, expect)
	})
	m.RLockFunc(func(m map[string]interface{}) {
		gtest.Assert(m, expect)
	})
}
func Test_StrAnyMap_Clone(t *testing.T) {
	//clone 方法是深克隆
	m := gmap.NewStrAnyMapFrom(map[string]interface{}{"a": 1, "b": "2"})

	m_clone := m.Clone()
	m.Remove("a")
	//修改原 map,clone 后的 map 不影响
	gtest.AssertIN("a", m_clone.Keys())

	m_clone.Remove("b")
	//修改clone map,原 map 不影响
	gtest.AssertIN("b", m.Keys())
}
func Test_StrAnyMap_Merge(t *testing.T) {
	m1 := gmap.NewStrAnyMap()
	m2 := gmap.NewStrAnyMap()
	m1.Set("a", 1)
	m2.Set("b", "2")
	m1.Merge(m2)
	gtest.Assert(m1.Map(), map[string]interface{}{"a": 1, "b": "2"})
}
