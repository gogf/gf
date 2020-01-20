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

	m := gmap.NewStrStrMapFrom(expect)
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

func Test_StrStrMap_Map(t *testing.T) {
	m := gmap.NewStrStrMap()
	m.Set("1", "1")
	m.Set("2", "2")
	gtest.Assert(m.Get("1"), "1")
	gtest.Assert(m.Get("2"), "2")
	data := m.Map()
	gtest.Assert(data["1"], "1")
	gtest.Assert(data["2"], "2")
	data["3"] = "3"
	gtest.Assert(m.Get("3"), "3")
	m.Set("4", "4")
	gtest.Assert(data["4"], "4")
}

func Test_StrStrMap_MapCopy(t *testing.T) {
	m := gmap.NewStrStrMap()
	m.Set("1", "1")
	m.Set("2", "2")
	gtest.Assert(m.Get("1"), "1")
	gtest.Assert(m.Get("2"), "2")
	data := m.MapCopy()
	gtest.Assert(data["1"], "1")
	gtest.Assert(data["2"], "2")
	data["3"] = "3"
	gtest.Assert(m.Get("3"), "")
	m.Set("4", "4")
	gtest.Assert(data["4"], "")
}

func Test_StrStrMap_FilterEmpty(t *testing.T) {
	m := gmap.NewStrStrMap()
	m.Set("1", "")
	m.Set("2", "2")
	gtest.Assert(m.Size(), 2)
	gtest.Assert(m.Get("1"), "")
	gtest.Assert(m.Get("2"), "2")
	m.FilterEmpty()
	gtest.Assert(m.Size(), 1)
	gtest.Assert(m.Get("2"), "2")
}

func Test_StrStrMap_Json(t *testing.T) {
	// Marshal
	gtest.Case(t, func() {
		data := g.MapStrStr{
			"k1": "v1",
			"k2": "v2",
		}
		m1 := gmap.NewStrStrMapFrom(data)
		b1, err1 := json.Marshal(m1)
		b2, err2 := json.Marshal(data)
		gtest.Assert(err1, err2)
		gtest.Assert(b1, b2)
	})
	// Unmarshal
	gtest.Case(t, func() {
		data := g.MapStrStr{
			"k1": "v1",
			"k2": "v2",
		}
		b, err := json.Marshal(data)
		gtest.Assert(err, nil)

		m := gmap.NewStrStrMap()
		err = json.Unmarshal(b, m)
		gtest.Assert(err, nil)
		gtest.Assert(m.Get("k1"), data["k1"])
		gtest.Assert(m.Get("k2"), data["k2"])
	})
	gtest.Case(t, func() {
		data := g.MapStrStr{
			"k1": "v1",
			"k2": "v2",
		}
		b, err := json.Marshal(data)
		gtest.Assert(err, nil)

		var m gmap.StrStrMap
		err = json.Unmarshal(b, &m)
		gtest.Assert(err, nil)
		gtest.Assert(m.Get("k1"), data["k1"])
		gtest.Assert(m.Get("k2"), data["k2"])
	})
}

func Test_StrStrMap_Pop(t *testing.T) {
	gtest.Case(t, func() {
		m := gmap.NewStrStrMapFrom(g.MapStrStr{
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

func Test_StrStrMap_Pops(t *testing.T) {
	gtest.Case(t, func() {
		m := gmap.NewStrStrMapFrom(g.MapStrStr{
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

func TestStrStrMap_UnmarshalValue(t *testing.T) {
	type T struct {
		Name string
		Map  *gmap.StrStrMap
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
