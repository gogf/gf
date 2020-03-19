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
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewStrAnyMap()
		m.Set("a", 1)

		t.Assert(m.Get("a"), 1)
		t.Assert(m.Size(), 1)
		t.Assert(m.IsEmpty(), false)

		t.Assert(m.GetOrSet("b", "2"), "2")
		t.Assert(m.SetIfNotExist("b", "2"), false)

		t.Assert(m.SetIfNotExist("c", 3), true)

		t.Assert(m.Remove("b"), "2")
		t.Assert(m.Contains("b"), false)

		t.AssertIN("c", m.Keys())
		t.AssertIN("a", m.Keys())
		t.AssertIN(3, m.Values())
		t.AssertIN(1, m.Values())

		m.Flip()
		t.Assert(m.Map(), map[string]interface{}{"1": "a", "3": "c"})

		m.Clear()
		t.Assert(m.Size(), 0)
		t.Assert(m.IsEmpty(), true)

		m2 := gmap.NewStrAnyMapFrom(map[string]interface{}{"a": 1, "b": "2"})
		t.Assert(m2.Map(), map[string]interface{}{"a": 1, "b": "2"})
	})
}
func Test_StrAnyMap_Set_Fun(t *testing.T) {
	m := gmap.NewStrAnyMap()

	m.GetOrSetFunc("a", getAny)
	m.GetOrSetFuncLock("b", getAny)
	t.Assert(m.Get("a"), 123)
	t.Assert(m.Get("b"), 123)
	t.Assert(m.SetIfNotExistFunc("a", getAny), false)
	t.Assert(m.SetIfNotExistFunc("c", getAny), true)

	t.Assert(m.SetIfNotExistFuncLock("b", getAny), false)
	t.Assert(m.SetIfNotExistFuncLock("d", getAny), true)

}

func Test_StrAnyMap_Batch(t *testing.T) {
	m := gmap.NewStrAnyMap()

	m.Sets(map[string]interface{}{"a": 1, "b": "2", "c": 3})
	t.Assert(m.Map(), map[string]interface{}{"a": 1, "b": "2", "c": 3})
	m.Removes([]string{"a", "b"})
	t.Assert(m.Map(), map[string]interface{}{"c": 3})
}

func Test_StrAnyMap_Iterator(t *testing.T) {
	expect := map[string]interface{}{"a": true, "b": false}
	m := gmap.NewStrAnyMapFrom(expect)
	m.Iterator(func(k string, v interface{}) bool {
		t.Assert(expect[k], v)
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
	t.Assert(i, 2)
	t.Assert(j, 1)
}

func Test_StrAnyMap_Lock(t *testing.T) {
	expect := map[string]interface{}{"a": true, "b": false}

	m := gmap.NewStrAnyMapFrom(expect)
	m.LockFunc(func(m map[string]interface{}) {
		t.Assert(m, expect)
	})
	m.RLockFunc(func(m map[string]interface{}) {
		t.Assert(m, expect)
	})
}
func Test_StrAnyMap_Clone(t *testing.T) {
	//clone 方法是深克隆
	m := gmap.NewStrAnyMapFrom(map[string]interface{}{"a": 1, "b": "2"})

	m_clone := m.Clone()
	m.Remove("a")
	//修改原 map,clone 后的 map 不影响
	t.AssertIN("a", m_clone.Keys())

	m_clone.Remove("b")
	//修改clone map,原 map 不影响
	t.AssertIN("b", m.Keys())
}
func Test_StrAnyMap_Merge(t *testing.T) {
	m1 := gmap.NewStrAnyMap()
	m2 := gmap.NewStrAnyMap()
	m1.Set("a", 1)
	m2.Set("b", "2")
	m1.Merge(m2)
	t.Assert(m1.Map(), map[string]interface{}{"a": 1, "b": "2"})
}

func Test_StrAnyMap_Map(t *testing.T) {
	m := gmap.NewStrAnyMap()
	m.Set("1", 1)
	m.Set("2", 2)
	t.Assert(m.Get("1"), 1)
	t.Assert(m.Get("2"), 2)
	data := m.Map()
	t.Assert(data["1"], 1)
	t.Assert(data["2"], 2)
	data["3"] = 3
	t.Assert(m.Get("3"), 3)
	m.Set("4", 4)
	t.Assert(data["4"], 4)
}

func Test_StrAnyMap_MapCopy(t *testing.T) {
	m := gmap.NewStrAnyMap()
	m.Set("1", 1)
	m.Set("2", 2)
	t.Assert(m.Get("1"), 1)
	t.Assert(m.Get("2"), 2)
	data := m.MapCopy()
	t.Assert(data["1"], 1)
	t.Assert(data["2"], 2)
	data["3"] = 3
	t.Assert(m.Get("3"), nil)
	m.Set("4", 4)
	t.Assert(data["4"], nil)
}

func Test_StrAnyMap_FilterEmpty(t *testing.T) {
	m := gmap.NewStrAnyMap()
	m.Set("1", 0)
	m.Set("2", 2)
	t.Assert(m.Size(), 2)
	t.Assert(m.Get("1"), 0)
	t.Assert(m.Get("2"), 2)
	m.FilterEmpty()
	t.Assert(m.Size(), 1)
	t.Assert(m.Get("2"), 2)
}

func Test_StrAnyMap_Json(t *testing.T) {
	// Marshal
	gtest.C(t, func(t *gtest.T) {
		data := g.MapStrAny{
			"k1": "v1",
			"k2": "v2",
		}
		m1 := gmap.NewStrAnyMapFrom(data)
		b1, err1 := json.Marshal(m1)
		b2, err2 := json.Marshal(data)
		t.Assert(err1, err2)
		t.Assert(b1, b2)
	})
	// Unmarshal
	gtest.C(t, func(t *gtest.T) {
		data := g.MapStrAny{
			"k1": "v1",
			"k2": "v2",
		}
		b, err := json.Marshal(data)
		t.Assert(err, nil)

		m := gmap.NewStrAnyMap()
		err = json.Unmarshal(b, m)
		t.Assert(err, nil)
		t.Assert(m.Get("k1"), data["k1"])
		t.Assert(m.Get("k2"), data["k2"])
	})
	gtest.C(t, func(t *gtest.T) {
		data := g.MapStrAny{
			"k1": "v1",
			"k2": "v2",
		}
		b, err := json.Marshal(data)
		t.Assert(err, nil)

		var m gmap.StrAnyMap
		err = json.Unmarshal(b, &m)
		t.Assert(err, nil)
		t.Assert(m.Get("k1"), data["k1"])
		t.Assert(m.Get("k2"), data["k2"])
	})
}

func Test_StrAnyMap_Pop(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewStrAnyMapFrom(g.MapStrAny{
			"k1": "v1",
			"k2": "v2",
		})
		t.Assert(m.Size(), 2)

		k1, v1 := m.Pop()
		t.AssertIN(k1, g.Slice{"k1", "k2"})
		t.AssertIN(v1, g.Slice{"v1", "v2"})
		t.Assert(m.Size(), 1)
		k2, v2 := m.Pop()
		t.AssertIN(k2, g.Slice{"k1", "k2"})
		t.AssertIN(v2, g.Slice{"v1", "v2"})
		t.Assert(m.Size(), 0)

		t.AssertNE(k1, k2)
		t.AssertNE(v1, v2)
	})
}

func Test_StrAnyMap_Pops(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewStrAnyMapFrom(g.MapStrAny{
			"k1": "v1",
			"k2": "v2",
			"k3": "v3",
		})
		t.Assert(m.Size(), 3)

		kArray := garray.New()
		vArray := garray.New()
		for k, v := range m.Pops(1) {
			t.AssertIN(k, g.Slice{"k1", "k2", "k3"})
			t.AssertIN(v, g.Slice{"v1", "v2", "v3"})
			kArray.Append(k)
			vArray.Append(v)
		}
		t.Assert(m.Size(), 2)
		for k, v := range m.Pops(2) {
			t.AssertIN(k, g.Slice{"k1", "k2", "k3"})
			t.AssertIN(v, g.Slice{"v1", "v2", "v3"})
			kArray.Append(k)
			vArray.Append(v)
		}
		t.Assert(m.Size(), 0)

		t.Assert(kArray.Unique().Len(), 3)
		t.Assert(vArray.Unique().Len(), 3)
	})
}

func TestStrAnyMap_UnmarshalValue(t *testing.T) {
	type T struct {
		Name string
		Map  *gmap.StrAnyMap
	}
	// JSON
	gtest.C(t, func(t *gtest.T) {
		var t *T
		err := gconv.Struct(map[string]interface{}{
			"name": "john",
			"map":  []byte(`{"k1":"v1","k2":"v2"}`),
		}, &t)
		t.Assert(err, nil)
		t.Assert(t.Name, "john")
		t.Assert(t.Map.Size(), 2)
		t.Assert(t.Map.Get("k1"), "v1")
		t.Assert(t.Map.Get("k2"), "v2")
	})
	// Map
	gtest.C(t, func(t *gtest.T) {
		var t *T
		err := gconv.Struct(map[string]interface{}{
			"name": "john",
			"map": g.Map{
				"k1": "v1",
				"k2": "v2",
			},
		}, &t)
		t.Assert(err, nil)
		t.Assert(t.Name, "john")
		t.Assert(t.Map.Size(), 2)
		t.Assert(t.Map.Get("k1"), "v1")
		t.Assert(t.Map.Get("k2"), "v2")
	})
}
