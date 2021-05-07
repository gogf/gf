// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with gm file,
// You can obtain one at https://github.com/gogf/gf.

package gmap_test

import (
	"github.com/gogf/gf/container/garray"
	"github.com/gogf/gf/internal/json"
	"github.com/gogf/gf/util/gconv"
	"testing"

	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/test/gtest"
)

func Test_ListMap_Var(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var m gmap.ListMap
		m.Set("key1", "val1")
		t.Assert(m.Keys(), []interface{}{"key1"})

		t.Assert(m.Get("key1"), "val1")
		t.Assert(m.Size(), 1)
		t.Assert(m.IsEmpty(), false)

		t.Assert(m.GetOrSet("key2", "val2"), "val2")
		t.Assert(m.SetIfNotExist("key2", "val2"), false)

		t.Assert(m.SetIfNotExist("key3", "val3"), true)
		t.Assert(m.Remove("key2"), "val2")
		t.Assert(m.Contains("key2"), false)

		t.AssertIN("key3", m.Keys())
		t.AssertIN("key1", m.Keys())
		t.AssertIN("val3", m.Values())
		t.AssertIN("val1", m.Values())

		m.Flip()

		t.Assert(m.Map(), map[interface{}]interface{}{"val3": "key3", "val1": "key1"})

		m.Clear()
		t.Assert(m.Size(), 0)
		t.Assert(m.IsEmpty(), true)
	})
}

func Test_ListMap_Basic(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListMap()
		m.Set("key1", "val1")
		t.Assert(m.Keys(), []interface{}{"key1"})

		t.Assert(m.Get("key1"), "val1")
		t.Assert(m.Size(), 1)
		t.Assert(m.IsEmpty(), false)

		t.Assert(m.GetOrSet("key2", "val2"), "val2")
		t.Assert(m.SetIfNotExist("key2", "val2"), false)

		t.Assert(m.SetIfNotExist("key3", "val3"), true)
		t.Assert(m.Remove("key2"), "val2")
		t.Assert(m.Contains("key2"), false)

		t.AssertIN("key3", m.Keys())
		t.AssertIN("key1", m.Keys())
		t.AssertIN("val3", m.Values())
		t.AssertIN("val1", m.Values())

		m.Flip()

		t.Assert(m.Map(), map[interface{}]interface{}{"val3": "key3", "val1": "key1"})

		m.Clear()
		t.Assert(m.Size(), 0)
		t.Assert(m.IsEmpty(), true)

		m2 := gmap.NewListMapFrom(map[interface{}]interface{}{1: 1, "key1": "val1"})
		t.Assert(m2.Map(), map[interface{}]interface{}{1: 1, "key1": "val1"})
	})
}

func Test_ListMap_Set_Fun(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListMap()
		m.GetOrSetFunc("fun", getValue)
		m.GetOrSetFuncLock("funlock", getValue)
		t.Assert(m.Get("funlock"), 3)
		t.Assert(m.Get("fun"), 3)
		m.GetOrSetFunc("fun", getValue)
		t.Assert(m.SetIfNotExistFunc("fun", getValue), false)
		t.Assert(m.SetIfNotExistFuncLock("funlock", getValue), false)
	})
}

func Test_ListMap_Batch(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListMap()
		m.Sets(map[interface{}]interface{}{1: 1, "key1": "val1", "key2": "val2", "key3": "val3"})
		t.Assert(m.Map(), map[interface{}]interface{}{1: 1, "key1": "val1", "key2": "val2", "key3": "val3"})
		m.Removes([]interface{}{"key1", 1})
		t.Assert(m.Map(), map[interface{}]interface{}{"key2": "val2", "key3": "val3"})
	})
}
func Test_ListMap_Iterator(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		expect := map[interface{}]interface{}{1: 1, "key1": "val1"}

		m := gmap.NewListMapFrom(expect)
		m.Iterator(func(k interface{}, v interface{}) bool {
			t.Assert(expect[k], v)
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
		t.Assert(i, 2)
		t.Assert(j, 1)
	})
}

func Test_ListMap_Clone(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		//clone 方法是深克隆
		m := gmap.NewListMapFrom(map[interface{}]interface{}{1: 1, "key1": "val1"})
		m_clone := m.Clone()
		m.Remove(1)
		//修改原 map,clone 后的 map 不影响
		t.AssertIN(1, m_clone.Keys())

		m_clone.Remove("key1")
		//修改clone map,原 map 不影响
		t.AssertIN("key1", m.Keys())
	})
}

func Test_ListMap_Basic_Merge(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m1 := gmap.NewListMap()
		m2 := gmap.NewListMap()
		m1.Set("key1", "val1")
		m2.Set("key2", "val2")
		m1.Merge(m2)
		t.Assert(m1.Map(), map[interface{}]interface{}{"key1": "val1", "key2": "val2"})
	})
}

func Test_ListMap_Order(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListMap()
		m.Set("k1", "v1")
		m.Set("k2", "v2")
		m.Set("k3", "v3")
		t.Assert(m.Keys(), g.Slice{"k1", "k2", "k3"})
		t.Assert(m.Values(), g.Slice{"v1", "v2", "v3"})
	})
}

func Test_ListMap_FilterEmpty(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListMap()
		m.Set(1, "")
		m.Set(2, "2")
		t.Assert(m.Size(), 2)
		t.Assert(m.Get(2), "2")
		m.FilterEmpty()
		t.Assert(m.Size(), 1)
		t.Assert(m.Get(2), "2")
	})
}

func Test_ListMap_Json(t *testing.T) {
	// Marshal
	gtest.C(t, func(t *gtest.T) {
		data := g.MapAnyAny{
			"k1": "v1",
			"k2": "v2",
		}
		m1 := gmap.NewListMapFrom(data)
		b1, err1 := json.Marshal(m1)
		b2, err2 := json.Marshal(gconv.Map(data))
		t.Assert(err1, err2)
		t.Assert(b1, b2)
	})
	// Unmarshal
	gtest.C(t, func(t *gtest.T) {
		data := g.MapAnyAny{
			"k1": "v1",
			"k2": "v2",
		}
		b, err := json.Marshal(gconv.Map(data))
		t.Assert(err, nil)

		m := gmap.NewListMap()
		err = json.Unmarshal(b, m)
		t.Assert(err, nil)
		t.Assert(m.Get("k1"), data["k1"])
		t.Assert(m.Get("k2"), data["k2"])
	})

	gtest.C(t, func(t *gtest.T) {
		data := g.MapAnyAny{
			"k1": "v1",
			"k2": "v2",
		}
		b, err := json.Marshal(gconv.Map(data))
		t.Assert(err, nil)

		var m gmap.ListMap
		err = json.Unmarshal(b, &m)
		t.Assert(err, nil)
		t.Assert(m.Get("k1"), data["k1"])
		t.Assert(m.Get("k2"), data["k2"])
	})
}

func Test_ListMap_Pop(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListMapFrom(g.MapAnyAny{
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

func Test_ListMap_Pops(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewListMapFrom(g.MapAnyAny{
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

func TestListMap_UnmarshalValue(t *testing.T) {
	type V struct {
		Name string
		Map  *gmap.ListMap
	}
	// JSON
	gtest.C(t, func(t *gtest.T) {
		var v *V
		err := gconv.Struct(map[string]interface{}{
			"name": "john",
			"map":  []byte(`{"1":"v1","2":"v2"}`),
		}, &v)
		t.Assert(err, nil)
		t.Assert(v.Name, "john")
		t.Assert(v.Map.Size(), 2)
		t.Assert(v.Map.Get("1"), "v1")
		t.Assert(v.Map.Get("2"), "v2")
	})
	// Map
	gtest.C(t, func(t *gtest.T) {
		var v *V
		err := gconv.Struct(map[string]interface{}{
			"name": "john",
			"map": g.MapIntAny{
				1: "v1",
				2: "v2",
			},
		}, &v)
		t.Assert(err, nil)
		t.Assert(v.Name, "john")
		t.Assert(v.Map.Size(), 2)
		t.Assert(v.Map.Get("1"), "v1")
		t.Assert(v.Map.Get("2"), "v2")
	})
}
