// Copyright 2017-2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with gm file,
// You can obtain one at https://github.com/gogf/gf.

package gmap_test

import (
	"encoding/json"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/util/gconv"
	"testing"

	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/test/gtest"
	"github.com/gogf/gf/util/gutil"
)

func Test_TreeMap_Basic(t *testing.T) {
	gtest.Case(t, func() {
		m := gmap.NewTreeMap(gutil.ComparatorString)
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

		m2 := gmap.NewTreeMapFrom(gutil.ComparatorString, map[interface{}]interface{}{1: 1, "key1": "val1"})
		gtest.Assert(m2.Map(), map[interface{}]interface{}{1: 1, "key1": "val1"})
	})
}
func Test_TreeMap_Set_Fun(t *testing.T) {
	m := gmap.NewTreeMap(gutil.ComparatorString)
	m.GetOrSetFunc("fun", getValue)
	m.GetOrSetFuncLock("funlock", getValue)
	gtest.Assert(m.Get("funlock"), 3)
	gtest.Assert(m.Get("fun"), 3)
	m.GetOrSetFunc("fun", getValue)
	gtest.Assert(m.SetIfNotExistFunc("fun", getValue), false)
	gtest.Assert(m.SetIfNotExistFuncLock("funlock", getValue), false)
}

func Test_TreeMap_Batch(t *testing.T) {
	m := gmap.NewTreeMap(gutil.ComparatorString)
	m.Sets(map[interface{}]interface{}{1: 1, "key1": "val1", "key2": "val2", "key3": "val3"})
	gtest.Assert(m.Map(), map[interface{}]interface{}{1: 1, "key1": "val1", "key2": "val2", "key3": "val3"})
	m.Removes([]interface{}{"key1", 1})
	gtest.Assert(m.Map(), map[interface{}]interface{}{"key2": "val2", "key3": "val3"})
}
func Test_TreeMap_Iterator(t *testing.T) {
	expect := map[interface{}]interface{}{1: 1, "key1": "val1"}

	m := gmap.NewTreeMapFrom(gutil.ComparatorString, expect)
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

func Test_TreeMap_Clone(t *testing.T) {
	//clone 方法是深克隆
	m := gmap.NewTreeMapFrom(gutil.ComparatorString, map[interface{}]interface{}{1: 1, "key1": "val1"})
	m_clone := m.Clone()
	m.Remove(1)
	//修改原 map,clone 后的 map 不影响
	gtest.AssertIN(1, m_clone.Keys())

	m_clone.Remove("key1")
	//修改clone map,原 map 不影响
	gtest.AssertIN("key1", m.Keys())
}

func Test_TreeMap_Json(t *testing.T) {
	// Marshal
	gtest.Case(t, func() {
		data := g.MapAnyAny{
			"k1": "v1",
			"k2": "v2",
		}
		m1 := gmap.NewTreeMapFrom(gutil.ComparatorString, data)
		b1, err1 := json.Marshal(m1)
		b2, err2 := json.Marshal(gconv.Map(data))
		gtest.Assert(err1, err2)
		gtest.Assert(b1, b2)
	})
	// Unmarshal
	gtest.Case(t, func() {
		data := g.MapAnyAny{
			"k1": "v1",
			"k2": "v2",
		}
		b, err := json.Marshal(gconv.Map(data))
		gtest.Assert(err, nil)

		m := gmap.NewTreeMap(gutil.ComparatorString)
		err = json.Unmarshal(b, m)
		gtest.Assert(err, nil)
		gtest.Assert(m.Get("k1"), data["k1"])
		gtest.Assert(m.Get("k2"), data["k2"])
	})
	gtest.Case(t, func() {
		data := g.MapAnyAny{
			"k1": "v1",
			"k2": "v2",
		}
		b, err := json.Marshal(gconv.Map(data))
		gtest.Assert(err, nil)

		var m gmap.TreeMap
		err = json.Unmarshal(b, &m)
		gtest.Assert(err, nil)
		gtest.Assert(m.Get("k1"), data["k1"])
		gtest.Assert(m.Get("k2"), data["k2"])
	})
}

func TestTreeMap_UnmarshalValue(t *testing.T) {
	type T struct {
		Name string
		Map  *gmap.TreeMap
	}
	// JSON
	gtest.Case(t, func() {
		var t *T
		err := gconv.Struct(map[string]interface{}{
			"name": "john",
			"map":  []byte(`{"k1":"v1","k2":"v2"}`),
		}, &t)
		gtest.Assert(err, nil)
		gtest.Assert(t.Name, "john")
		gtest.Assert(t.Map.Size(), 2)
		gtest.Assert(t.Map.Get("k1"), "v1")
		gtest.Assert(t.Map.Get("k2"), "v2")
	})
	// Map
	gtest.Case(t, func() {
		var t *T
		err := gconv.Struct(map[string]interface{}{
			"name": "john",
			"map": g.Map{
				"k1": "v1",
				"k2": "v2",
			},
		}, &t)
		gtest.Assert(err, nil)
		gtest.Assert(t.Name, "john")
		gtest.Assert(t.Map.Size(), 2)
		gtest.Assert(t.Map.Get("k1"), "v1")
		gtest.Assert(t.Map.Get("k2"), "v2")
	})
}
