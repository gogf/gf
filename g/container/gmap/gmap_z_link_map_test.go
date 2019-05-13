// Copyright 2017-2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with gm file,
// You can obtain one at https://github.com/gogf/gf.

package gmap_test

import (
	"github.com/gogf/gf/g"
	"github.com/gogf/gf/g/container/gmap"
	"github.com/gogf/gf/g/test/gtest"
	"testing"
)

func Test_List_Map_Basic(t *testing.T) {
	gtest.Case(t, func() {
		m := gmap.NewListMap()
		m.Set("key1", "val1")
		gtest.Assert(m.Keys(), []interface{}{"key1"})

		gtest.Assert(m.Get("key1"), "val1")
		gtest.Assert(m.Size(), 1)
		gtest.Assert(m.IsEmpty(), false)

		gtest.Assert(m.GetOrSet("key2", "val2"), "val2")
		gtest.Assert(m.SetIfNotExist("key2", "val2"), false)

		gtest.Assert(m.SetIfNotExist("key3", "val3"), true)
		gtest.Assert(m.Remove("key2"), "val2")
		gtest.Assert(m.Contains("key2"), false)

		gtest.AssertIN("key3", m.Keys())
		gtest.AssertIN("key1", m.Keys())
		gtest.AssertIN("val3", m.Values())
		gtest.AssertIN("val1", m.Values())

		m.Flip()

		gtest.Assert(m.Map(), map[interface{}]interface{}{"val3": "key3", "val1": "key1"})

		m.Clear()
		gtest.Assert(m.Size(), 0)
		gtest.Assert(m.IsEmpty(), true)

		m2 := gmap.NewListMapFrom(map[interface{}]interface{}{1: 1, "key1": "val1"})
		gtest.Assert(m2.Map(), map[interface{}]interface{}{1: 1, "key1": "val1"})
	})
}
func Test_List_Map_Set_Fun(t *testing.T) {
	m := gmap.NewListMap()
	m.GetOrSetFunc("fun", getValue)
	m.GetOrSetFuncLock("funlock", getValue)
	gtest.Assert(m.Get("funlock"), 3)
	gtest.Assert(m.Get("fun"), 3)
	m.GetOrSetFunc("fun", getValue)
	gtest.Assert(m.SetIfNotExistFunc("fun", getValue), false)
	gtest.Assert(m.SetIfNotExistFuncLock("funlock", getValue), false)
}

func Test_List_Map_Batch(t *testing.T) {
	m := gmap.NewListMap()
	m.Sets(map[interface{}]interface{}{1: 1, "key1": "val1", "key2": "val2", "key3": "val3"})
	gtest.Assert(m.Map(), map[interface{}]interface{}{1: 1, "key1": "val1", "key2": "val2", "key3": "val3"})
	m.Removes([]interface{}{"key1", 1})
	gtest.Assert(m.Map(), map[interface{}]interface{}{"key2": "val2", "key3": "val3"})
}
func Test_List_Map_Iterator(t *testing.T){
	expect :=map[interface{}]interface{}{1: 1, "key1": "val1"}

	m := gmap.NewListMapFrom(expect)
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
	gtest.Assert(i, 2)
	gtest.Assert(j, 1)
}

func Test_List_Map_Clone(t *testing.T) {
	//clone 方法是深克隆
	m := gmap.NewListMapFrom(map[interface{}]interface{}{1: 1, "key1": "val1"})
	m_clone := m.Clone()
	m.Remove(1)
	//修改原 map,clone 后的 map 不影响
	gtest.AssertIN(1, m_clone.Keys())

	m_clone.Remove("key1")
	//修改clone map,原 map 不影响
	gtest.AssertIN("key1", m.Keys())
}

func Test_List_Map_Basic_Merge(t *testing.T) {
	m1 := gmap.NewListMap()
	m2 := gmap.NewListMap()
	m1.Set("key1", "val1")
	m2.Set("key2", "val2")
	m1.Merge(m2)
	gtest.Assert(m1.Map(), map[interface{}]interface{}{"key1": "val1", "key2": "val2"})
}

func Test_List_Map_Order(t *testing.T) {
	m := gmap.NewListMap()
	m.Set("k1", "v1")
	m.Set("k2", "v2")
	m.Set("k3", "v3")
	gtest.Assert(m.Keys(),   g.Slice{"k1", "k2", "k3"})
	gtest.Assert(m.Values(), g.Slice{"v1", "v2", "v3"})
}
