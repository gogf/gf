// Copyright 2017-2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with gm file,
// You can obtain one at https://github.com/gogf/gf.

package gmap_test

import (
	"testing"

	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/test/gtest"
)

func anyAnyCallBack(int, interface{}) bool {
	return true
}

func Test_AnyAnyMap_Basic(t *testing.T) {
	gtest.Case(t, func() {
		m := gmap.NewAnyAnyMap()
		m.Set(1, 1)

		gtest.Assert(m.Get(1), 1)
		gtest.Assert(m.Size(), 1)
		gtest.Assert(m.IsEmpty(), false)

		gtest.Assert(m.GetOrSet(2, "2"), "2")
		gtest.Assert(m.SetIfNotExist(2, "2"), false)

		gtest.Assert(m.SetIfNotExist(3, 3), true)

		gtest.Assert(m.Remove(2), "2")
		gtest.Assert(m.Contains(2), false)

		gtest.AssertIN(3, m.Keys())
		gtest.AssertIN(1, m.Keys())
		gtest.AssertIN(3, m.Values())
		gtest.AssertIN(1, m.Values())
		m.Flip()
		gtest.Assert(m.Map(), map[interface{}]int{1: 1, 3: 3})

		m.Clear()
		gtest.Assert(m.Size(), 0)
		gtest.Assert(m.IsEmpty(), true)

		m2 := gmap.NewAnyAnyMapFrom(map[interface{}]interface{}{1: 1, 2: "2"})
		gtest.Assert(m2.Map(), map[interface{}]interface{}{1: 1, 2: "2"})
	})
}

func Test_AnyAnyMap_Set_Fun(t *testing.T) {
	m := gmap.NewAnyAnyMap()

	m.GetOrSetFunc(1, getAny)
	m.GetOrSetFuncLock(2, getAny)
	gtest.Assert(m.Get(1), 123)
	gtest.Assert(m.Get(2), 123)

	gtest.Assert(m.SetIfNotExistFunc(1, getAny), false)
	gtest.Assert(m.SetIfNotExistFunc(3, getAny), true)

	gtest.Assert(m.SetIfNotExistFuncLock(2, getAny), false)
	gtest.Assert(m.SetIfNotExistFuncLock(4, getAny), true)

}

func Test_AnyAnyMap_Batch(t *testing.T) {
	m := gmap.NewAnyAnyMap()

	m.Sets(map[interface{}]interface{}{1: 1, 2: "2", 3: 3})
	gtest.Assert(m.Map(), map[interface{}]interface{}{1: 1, 2: "2", 3: 3})
	m.Removes([]interface{}{1, 2})
	gtest.Assert(m.Map(), map[interface{}]interface{}{3: 3})
}

func Test_AnyAnyMap_Iterator(t *testing.T) {
	expect := map[interface{}]interface{}{1: 1, 2: "2"}
	m := gmap.NewAnyAnyMapFrom(expect)
	m.Iterator(func(k interface{}, v interface{}) bool {
		gtest.Assert(expect[k], v)
		return true
	})
	// 断言返回值对遍历控制
	i := 0
	j := 0
	m.Iterator(func(k interface{}, v interface{}) bool {
		i++
		return true
	})
	m.Iterator(func(k interface{}, v interface{}) bool {
		j++
		return false
	})
	gtest.Assert(i, "2")
	gtest.Assert(j, 1)

}

func Test_AnyAnyMap_Lock(t *testing.T) {
	expect := map[interface{}]interface{}{1: 1, 2: "2"}
	m := gmap.NewAnyAnyMapFrom(expect)
	m.LockFunc(func(m map[interface{}]interface{}) {
		gtest.Assert(m, expect)
	})
	m.RLockFunc(func(m map[interface{}]interface{}) {
		gtest.Assert(m, expect)
	})
}

func Test_AnyAnyMap_Clone(t *testing.T) {
	gtest.Case(t, func() {
		//clone 方法是深克隆
		m := gmap.NewAnyAnyMapFrom(map[interface{}]interface{}{1: 1, 2: "2"})

		m_clone := m.Clone()
		m.Remove(1)
		//修改原 map,clone 后的 map 不影响
		gtest.AssertIN(1, m_clone.Keys())

		m_clone.Remove(2)
		//修改clone map,原 map 不影响
		gtest.AssertIN(2, m.Keys())
	})
}

func Test_AnyAnyMap_Merge(t *testing.T) {
	m1 := gmap.NewAnyAnyMap()
	m2 := gmap.NewAnyAnyMap()
	m1.Set(1, 1)
	m2.Set(2, "2")
	m1.Merge(m2)
	gtest.Assert(m1.Map(), map[interface{}]interface{}{1: 1, 2: "2"})
}

func Test_AnyAnyMap_Map(t *testing.T) {
	m := gmap.NewAnyAnyMap()
	m.Set(1, 0)
	m.Set(2, 2)
	gtest.Assert(m.Get(1), 0)
	gtest.Assert(m.Get(2), 2)
	data := m.Map()
	gtest.Assert(data[1], 0)
	gtest.Assert(data[2], 2)
	data[3] = 3
	gtest.Assert(m.Get(3), 3)
	m.Set(4, 4)
	gtest.Assert(data[4], 4)
}

func Test_AnyAnyMap_MapCopy(t *testing.T) {
	m := gmap.NewAnyAnyMap()
	m.Set(1, 0)
	m.Set(2, 2)
	gtest.Assert(m.Get(1), 0)
	gtest.Assert(m.Get(2), 2)
	data := m.MapCopy()
	gtest.Assert(data[1], 0)
	gtest.Assert(data[2], 2)
	data[3] = 3
	gtest.Assert(m.Get(3), nil)
	m.Set(4, 4)
	gtest.Assert(data[4], nil)
}

func Test_AnyAnyMap_FilterEmpty(t *testing.T) {
	m := gmap.NewAnyAnyMap()
	m.Set(1, 0)
	m.Set(2, 2)
	gtest.Assert(m.Get(1), 0)
	gtest.Assert(m.Get(2), 2)
	m.FilterEmpty()
	gtest.Assert(m.Get(1), nil)
	gtest.Assert(m.Get(2), 2)
}
