// Copyright 2017-2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with gm file,
// You can obtain one at https://github.com/gogf/gf.

package gmap_test

import (
	"encoding/json"
	"github.com/gogf/gf/container/garray"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/util/gconv"
	"testing"

	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/test/gtest"
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

	m := gmap.NewStrAnyMapFrom(expect)
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

func Test_StrAnyMap_Map(t *testing.T) {
	m := gmap.NewStrAnyMap()
	m.Set("1", 1)
	m.Set("2", 2)
	gtest.Assert(m.Get("1"), 1)
	gtest.Assert(m.Get("2"), 2)
	data := m.Map()
	gtest.Assert(data["1"], 1)
	gtest.Assert(data["2"], 2)
	data["3"] = 3
	gtest.Assert(m.Get("3"), 3)
	m.Set("4", 4)
	gtest.Assert(data["4"], 4)
}

func Test_StrAnyMap_MapCopy(t *testing.T) {
	m := gmap.NewStrAnyMap()
	m.Set("1", 1)
	m.Set("2", 2)
	gtest.Assert(m.Get("1"), 1)
	gtest.Assert(m.Get("2"), 2)
	data := m.MapCopy()
	gtest.Assert(data["1"], 1)
	gtest.Assert(data["2"], 2)
	data["3"] = 3
	gtest.Assert(m.Get("3"), nil)
	m.Set("4", 4)
	gtest.Assert(data["4"], nil)
}

func Test_StrAnyMap_FilterEmpty(t *testing.T) {
	m := gmap.NewStrAnyMap()
	m.Set("1", 0)
	m.Set("2", 2)
	gtest.Assert(m.Size(), 2)
	gtest.Assert(m.Get("1"), 0)
	gtest.Assert(m.Get("2"), 2)
	m.FilterEmpty()
	gtest.Assert(m.Size(), 1)
	gtest.Assert(m.Get("2"), 2)
}

func Test_StrAnyMap_Json(t *testing.T) {
	// Marshal
	gtest.Case(t, func() {
		data := g.MapStrAny{
			"k1": "v1",
			"k2": "v2",
		}
		m1 := gmap.NewStrAnyMapFrom(data)
		b1, err1 := json.Marshal(m1)
		b2, err2 := json.Marshal(data)
		gtest.Assert(err1, err2)
		gtest.Assert(b1, b2)
	})
	// Unmarshal
	gtest.Case(t, func() {
		data := g.MapStrAny{
			"k1": "v1",
			"k2": "v2",
		}
		b, err := json.Marshal(data)
		gtest.Assert(err, nil)

		m := gmap.NewStrAnyMap()
		err = json.Unmarshal(b, m)
		gtest.Assert(err, nil)
		gtest.Assert(m.Get("k1"), data["k1"])
		gtest.Assert(m.Get("k2"), data["k2"])
	})
	gtest.Case(t, func() {
		data := g.MapStrAny{
			"k1": "v1",
			"k2": "v2",
		}
		b, err := json.Marshal(data)
		gtest.Assert(err, nil)

		var m gmap.StrAnyMap
		err = json.Unmarshal(b, &m)
		gtest.Assert(err, nil)
		gtest.Assert(m.Get("k1"), data["k1"])
		gtest.Assert(m.Get("k2"), data["k2"])
	})
}

func Test_StrAnyMap_Pop(t *testing.T) {
	gtest.Case(t, func() {
		m := gmap.NewStrAnyMapFrom(g.MapStrAny{
			"k1": "v1",
			"k2": "v2",
		})
		gtest.Assert(m.Size(), 2)

		k1, v1 := m.Pop()
		gtest.AssertIN(k1, g.Slice{"k1", "k2"})
		gtest.AssertIN(v1, g.Slice{"v1", "v2"})
		gtest.Assert(m.Size(), 1)
		k2, v2 := m.Pop()
		gtest.AssertIN(k2, g.Slice{"k1", "k2"})
		gtest.AssertIN(v2, g.Slice{"v1", "v2"})
		gtest.Assert(m.Size(), 0)

		gtest.AssertNE(k1, k2)
		gtest.AssertNE(v1, v2)
	})
}

func Test_StrAnyMap_Pops(t *testing.T) {
	gtest.Case(t, func() {
		m := gmap.NewStrAnyMapFrom(g.MapStrAny{
			"k1": "v1",
			"k2": "v2",
			"k3": "v3",
		})
		gtest.Assert(m.Size(), 3)

		kArray := garray.New()
		vArray := garray.New()
		for k, v := range m.Pops(1) {
			gtest.AssertIN(k, g.Slice{"k1", "k2", "k3"})
			gtest.AssertIN(v, g.Slice{"v1", "v2", "v3"})
			kArray.Append(k)
			vArray.Append(v)
		}
		gtest.Assert(m.Size(), 2)
		for k, v := range m.Pops(2) {
			gtest.AssertIN(k, g.Slice{"k1", "k2", "k3"})
			gtest.AssertIN(v, g.Slice{"v1", "v2", "v3"})
			kArray.Append(k)
			vArray.Append(v)
		}
		gtest.Assert(m.Size(), 0)

		gtest.Assert(kArray.Unique().Len(), 3)
		gtest.Assert(vArray.Unique().Len(), 3)
	})
}

func TestStrAnyMap_UnmarshalValue(t *testing.T) {
	type T struct {
		Name string
		Map  *gmap.StrAnyMap
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
