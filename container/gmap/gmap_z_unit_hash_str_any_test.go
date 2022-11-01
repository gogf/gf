// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with gm file,
// You can obtain one at https://github.com/gogf/gf.

package gmap_test

import (
	"testing"

	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gconv"
)

func Test_StrAnyMap_Var(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var m gmap.StrAnyMap
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
	})
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
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewStrAnyMap()

		m.GetOrSetFunc("a", getAny)
		m.GetOrSetFuncLock("b", getAny)
		t.Assert(m.Get("a"), 123)
		t.Assert(m.Get("b"), 123)
		t.Assert(m.SetIfNotExistFunc("a", getAny), false)
		t.Assert(m.SetIfNotExistFunc("c", getAny), true)

		t.Assert(m.SetIfNotExistFuncLock("b", getAny), false)
		t.Assert(m.SetIfNotExistFuncLock("d", getAny), true)
	})
}

func Test_StrAnyMap_Batch(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewStrAnyMap()

		m.Sets(map[string]interface{}{"a": 1, "b": "2", "c": 3})
		t.Assert(m.Map(), map[string]interface{}{"a": 1, "b": "2", "c": 3})
		m.Removes([]string{"a", "b"})
		t.Assert(m.Map(), map[string]interface{}{"c": 3})
	})
}

func Test_StrAnyMap_Iterator(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
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
	})
}

func Test_StrAnyMap_Lock(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		expect := map[string]interface{}{"a": true, "b": false}

		m := gmap.NewStrAnyMapFrom(expect)
		m.LockFunc(func(m map[string]interface{}) {
			t.Assert(m, expect)
		})
		m.RLockFunc(func(m map[string]interface{}) {
			t.Assert(m, expect)
		})
	})
}

func Test_StrAnyMap_Clone(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// clone 方法是深克隆
		m := gmap.NewStrAnyMapFrom(map[string]interface{}{"a": 1, "b": "2"})

		m_clone := m.Clone()
		m.Remove("a")
		// 修改原 map,clone 后的 map 不影响
		t.AssertIN("a", m_clone.Keys())

		m_clone.Remove("b")
		// 修改clone map,原 map 不影响
		t.AssertIN("b", m.Keys())
	})
}

func Test_StrAnyMap_Merge(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m1 := gmap.NewStrAnyMap()
		m2 := gmap.NewStrAnyMap()
		m1.Set("a", 1)
		m2.Set("b", "2")
		m1.Merge(m2)
		t.Assert(m1.Map(), map[string]interface{}{"a": 1, "b": "2"})

		m3 := gmap.NewStrAnyMapFrom(nil)
		m3.Merge(m2)
		t.Assert(m3.Map(), m2.Map())
	})
}

func Test_StrAnyMap_Map(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
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
	})
}

func Test_StrAnyMap_MapCopy(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
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
	})
}

func Test_StrAnyMap_FilterEmpty(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewStrAnyMap()
		m.Set("1", 0)
		m.Set("2", 2)
		t.Assert(m.Size(), 2)
		t.Assert(m.Get("1"), 0)
		t.Assert(m.Get("2"), 2)
		m.FilterEmpty()
		t.Assert(m.Size(), 1)
		t.Assert(m.Get("2"), 2)
	})
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
		t.AssertNil(err)

		m := gmap.NewStrAnyMap()
		err = json.UnmarshalUseNumber(b, m)
		t.AssertNil(err)
		t.Assert(m.Get("k1"), data["k1"])
		t.Assert(m.Get("k2"), data["k2"])
	})
	gtest.C(t, func(t *gtest.T) {
		data := g.MapStrAny{
			"k1": "v1",
			"k2": "v2",
		}
		b, err := json.Marshal(data)
		t.AssertNil(err)

		var m gmap.StrAnyMap
		err = json.UnmarshalUseNumber(b, &m)
		t.AssertNil(err)
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

		k3, v3 := m.Pop()
		t.Assert(k3, "")
		t.Assert(v3, "")
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

		v := m.Pops(1)
		t.AssertNil(v)
		v = m.Pops(-1)
		t.AssertNil(v)
	})
}

func TestStrAnyMap_UnmarshalValue(t *testing.T) {
	type V struct {
		Name string
		Map  *gmap.StrAnyMap
	}
	// JSON
	gtest.C(t, func(t *gtest.T) {
		var v *V
		err := gconv.Struct(map[string]interface{}{
			"name": "john",
			"map":  []byte(`{"k1":"v1","k2":"v2"}`),
		}, &v)
		t.AssertNil(err)
		t.Assert(v.Name, "john")
		t.Assert(v.Map.Size(), 2)
		t.Assert(v.Map.Get("k1"), "v1")
		t.Assert(v.Map.Get("k2"), "v2")
	})
	// Map
	gtest.C(t, func(t *gtest.T) {
		var v *V
		err := gconv.Struct(map[string]interface{}{
			"name": "john",
			"map": g.Map{
				"k1": "v1",
				"k2": "v2",
			},
		}, &v)
		t.AssertNil(err)
		t.Assert(v.Name, "john")
		t.Assert(v.Map.Size(), 2)
		t.Assert(v.Map.Get("k1"), "v1")
		t.Assert(v.Map.Get("k2"), "v2")
	})
}

func Test_StrAnyMap_DeepCopy(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewStrAnyMapFrom(g.MapStrAny{
			"key1": "val1",
			"key2": "val2",
		})
		t.Assert(m.Size(), 2)

		n := m.DeepCopy().(*gmap.StrAnyMap)
		n.Set("key1", "v1")
		t.AssertNE(m.Get("key1"), n.Get("key1"))
	})
}
