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

func stringIntCallBack(string, int) bool {
	return true
}
func Test_StrIntMap_Basic(t *testing.T) {
	gtest.Case(t, func() {
		m := gmap.NewStrIntMap()
		m.Set("a", 1)

		gtest.Assert(m.Get("a"), 1)
		gtest.Assert(m.Size(), 1)
		gtest.Assert(m.IsEmpty(), false)

		gtest.Assert(m.GetOrSet("b", 2), 2)
		gtest.Assert(m.SetIfNotExist("b", 2), false)

		gtest.Assert(m.SetIfNotExist("c", 3), true)

		gtest.Assert(m.Remove("b"), 2)
		gtest.Assert(m.Contains("b"), false)

		gtest.AssertIN("c", m.Keys())
		gtest.AssertIN("a", m.Keys())
		gtest.AssertIN(3, m.Values())
		gtest.AssertIN(1, m.Values())

		m_f := gmap.NewStrIntMap()
		m_f.Set("1", 2)
		m_f.Flip()
		gtest.Assert(m_f.Map(), map[string]int{"2": 1})

		m.Clear()
		gtest.Assert(m.Size(), 0)
		gtest.Assert(m.IsEmpty(), true)

		m2 := gmap.NewStrIntMapFrom(map[string]int{"a": 1, "b": 2})
		gtest.Assert(m2.Map(), map[string]int{"a": 1, "b": 2})
	})
}
func Test_StrIntMap_Set_Fun(t *testing.T) {
	m := gmap.NewStrIntMap()

	m.GetOrSetFunc("a", getInt)
	m.GetOrSetFuncLock("b", getInt)
	gtest.Assert(m.Get("a"), 123)
	gtest.Assert(m.Get("b"), 123)
	gtest.Assert(m.SetIfNotExistFunc("a", getInt), false)
	gtest.Assert(m.SetIfNotExistFunc("c", getInt), true)

	gtest.Assert(m.SetIfNotExistFuncLock("b", getInt), false)
	gtest.Assert(m.SetIfNotExistFuncLock("d", getInt), true)

}

func Test_StrIntMap_Batch(t *testing.T) {
	m := gmap.NewStrIntMap()

	m.Sets(map[string]int{"a": 1, "b": 2, "c": 3})
	gtest.Assert(m.Map(), map[string]int{"a": 1, "b": 2, "c": 3})
	m.Removes([]string{"a", "b"})
	gtest.Assert(m.Map(), map[string]int{"c": 3})
}
func Test_StrIntMap_Iterator(t *testing.T) {
	expect := map[string]int{"a": 1, "b": 2}
	m := gmap.NewStrIntMapFrom(expect)
	m.Iterator(func(k string, v int) bool {
		gtest.Assert(expect[k], v)
		return true
	})
	// 断言返回值对遍历控制
	i := 0
	j := 0
	m.Iterator(func(k string, v int) bool {
		i++
		return true
	})
	m.Iterator(func(k string, v int) bool {
		j++
		return false
	})
	gtest.Assert(i, 2)
	gtest.Assert(j, 1)

}

func Test_StrIntMap_Lock(t *testing.T) {
	expect := map[string]int{"a": 1, "b": 2}

	m := gmap.NewStrIntMapFrom(expect)
	m.LockFunc(func(m map[string]int) {
		gtest.Assert(m, expect)
	})
	m.RLockFunc(func(m map[string]int) {
		gtest.Assert(m, expect)
	})
}

func Test_StrIntMap_Clone(t *testing.T) {
	//clone 方法是深克隆
	m := gmap.NewStrIntMapFrom(map[string]int{"a": 1, "b": 2, "c": 3})

	m_clone := m.Clone()
	m.Remove("a")
	//修改原 map,clone 后的 map 不影响
	gtest.AssertIN("a", m_clone.Keys())

	m_clone.Remove("b")
	//修改clone map,原 map 不影响
	gtest.AssertIN("b", m.Keys())
}
func Test_StrIntMap_Merge(t *testing.T) {
	m1 := gmap.NewStrIntMap()
	m2 := gmap.NewStrIntMap()
	m1.Set("a", 1)
	m2.Set("b", 2)
	m1.Merge(m2)
	gtest.Assert(m1.Map(), map[string]int{"a": 1, "b": 2})
}

func Test_StrIntMap_Map(t *testing.T) {
	m := gmap.NewStrIntMap()
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

func Test_StrIntMap_MapCopy(t *testing.T) {
	m := gmap.NewStrIntMap()
	m.Set("1", 1)
	m.Set("2", 2)
	gtest.Assert(m.Get("1"), 1)
	gtest.Assert(m.Get("2"), 2)
	data := m.MapCopy()
	gtest.Assert(data["1"], 1)
	gtest.Assert(data["2"], 2)
	data["3"] = 3
	gtest.Assert(m.Get("3"), 0)
	m.Set("4", 4)
	gtest.Assert(data["4"], 0)
}

func Test_StrIntMap_FilterEmpty(t *testing.T) {
	m := gmap.NewStrIntMap()
	m.Set("1", 0)
	m.Set("2", 2)
	gtest.Assert(m.Size(), 2)
	gtest.Assert(m.Get("1"), 0)
	gtest.Assert(m.Get("2"), 2)
	m.FilterEmpty()
	gtest.Assert(m.Size(), 1)
	gtest.Assert(m.Get("2"), 2)
}

func Test_StrIntMap_Json(t *testing.T) {
	// Marshal
	gtest.Case(t, func() {
		data := g.MapStrInt{
			"k1": 1,
			"k2": 2,
		}
		m1 := gmap.NewStrIntMapFrom(data)
		b1, err1 := json.Marshal(m1)
		b2, err2 := json.Marshal(data)
		gtest.Assert(err1, err2)
		gtest.Assert(b1, b2)
	})
	// Unmarshal
	gtest.Case(t, func() {
		data := g.MapStrInt{
			"k1": 1,
			"k2": 2,
		}
		b, err := json.Marshal(data)
		gtest.Assert(err, nil)

		m := gmap.NewStrIntMap()
		err = json.Unmarshal(b, m)
		gtest.Assert(err, nil)
		gtest.Assert(m.Get("k1"), data["k1"])
		gtest.Assert(m.Get("k2"), data["k2"])
	})
	gtest.Case(t, func() {
		data := g.MapStrInt{
			"k1": 1,
			"k2": 2,
		}
		b, err := json.Marshal(data)
		gtest.Assert(err, nil)

		var m gmap.StrIntMap
		err = json.Unmarshal(b, &m)
		gtest.Assert(err, nil)
		gtest.Assert(m.Get("k1"), data["k1"])
		gtest.Assert(m.Get("k2"), data["k2"])
	})
}

func Test_StrIntMap_Pop(t *testing.T) {
	gtest.Case(t, func() {
		m := gmap.NewStrIntMapFrom(g.MapStrInt{
			"k1": 11,
			"k2": 22,
		})
		gtest.Assert(m.Size(), 2)

		k1, v1 := m.Pop()
		gtest.AssertIN(k1, g.Slice{"k1", "k2"})
		gtest.AssertIN(v1, g.Slice{11, 22})
		gtest.Assert(m.Size(), 1)
		k2, v2 := m.Pop()
		gtest.AssertIN(k2, g.Slice{"k1", "k2"})
		gtest.AssertIN(v2, g.Slice{11, 22})
		gtest.Assert(m.Size(), 0)

		gtest.AssertNE(k1, k2)
		gtest.AssertNE(v1, v2)
	})
}

func Test_StrIntMap_Pops(t *testing.T) {
	gtest.Case(t, func() {
		m := gmap.NewStrIntMapFrom(g.MapStrInt{
			"k1": 11,
			"k2": 22,
			"k3": 33,
		})
		gtest.Assert(m.Size(), 3)

		kArray := garray.New()
		vArray := garray.New()
		for k, v := range m.Pops(1) {
			gtest.AssertIN(k, g.Slice{"k1", "k2", "k3"})
			gtest.AssertIN(v, g.Slice{11, 22, 33})
			kArray.Append(k)
			vArray.Append(v)
		}
		gtest.Assert(m.Size(), 2)
		for k, v := range m.Pops(2) {
			gtest.AssertIN(k, g.Slice{"k1", "k2", "k3"})
			gtest.AssertIN(v, g.Slice{11, 22, 33})
			kArray.Append(k)
			vArray.Append(v)
		}
		gtest.Assert(m.Size(), 0)

		gtest.Assert(kArray.Unique().Len(), 3)
		gtest.Assert(vArray.Unique().Len(), 3)
	})
}

func TestStrIntMap_UnmarshalValue(t *testing.T) {
	type T struct {
		Name string
		Map  *gmap.StrIntMap
	}
	// JSON
	gtest.Case(t, func() {
		var t *T
		err := gconv.Struct(map[string]interface{}{
			"name": "john",
			"map":  []byte(`{"k1":1,"k2":2}`),
		}, &t)
		gtest.Assert(err, nil)
		gtest.Assert(t.Name, "john")
		gtest.Assert(t.Map.Size(), 2)
		gtest.Assert(t.Map.Get("k1"), 1)
		gtest.Assert(t.Map.Get("k2"), 2)
	})
	// Map
	gtest.Case(t, func() {
		var t *T
		err := gconv.Struct(map[string]interface{}{
			"name": "john",
			"map": g.Map{
				"k1": 1,
				"k2": 2,
			},
		}, &t)
		gtest.Assert(err, nil)
		gtest.Assert(t.Name, "john")
		gtest.Assert(t.Map.Size(), 2)
		gtest.Assert(t.Map.Get("k1"), 1)
		gtest.Assert(t.Map.Get("k2"), 2)
	})
}
