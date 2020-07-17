// Copyright 2017-2019 gf Author(https://github.com/jin502437344/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with gm file,
// You can obtain one at https://github.com/jin502437344/gf.

package gmap_test

import (
	"github.com/jin502437344/gf/container/garray"
	"github.com/jin502437344/gf/frame/g"
	"github.com/jin502437344/gf/internal/json"
	"github.com/jin502437344/gf/util/gconv"
	"testing"

	"github.com/jin502437344/gf/container/gmap"
	"github.com/jin502437344/gf/test/gtest"
)

func Test_StrIntMap_Var(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var m gmap.StrIntMap
		m.Set("a", 1)

		t.Assert(m.Get("a"), 1)
		t.Assert(m.Size(), 1)
		t.Assert(m.IsEmpty(), false)

		t.Assert(m.GetOrSet("b", 2), 2)
		t.Assert(m.SetIfNotExist("b", 2), false)

		t.Assert(m.SetIfNotExist("c", 3), true)

		t.Assert(m.Remove("b"), 2)
		t.Assert(m.Contains("b"), false)

		t.AssertIN("c", m.Keys())
		t.AssertIN("a", m.Keys())
		t.AssertIN(3, m.Values())
		t.AssertIN(1, m.Values())

		m_f := gmap.NewStrIntMap()
		m_f.Set("1", 2)
		m_f.Flip()
		t.Assert(m_f.Map(), map[string]int{"2": 1})

		m.Clear()
		t.Assert(m.Size(), 0)
		t.Assert(m.IsEmpty(), true)
	})
}

func Test_StrIntMap_Basic(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewStrIntMap()
		m.Set("a", 1)

		t.Assert(m.Get("a"), 1)
		t.Assert(m.Size(), 1)
		t.Assert(m.IsEmpty(), false)

		t.Assert(m.GetOrSet("b", 2), 2)
		t.Assert(m.SetIfNotExist("b", 2), false)

		t.Assert(m.SetIfNotExist("c", 3), true)

		t.Assert(m.Remove("b"), 2)
		t.Assert(m.Contains("b"), false)

		t.AssertIN("c", m.Keys())
		t.AssertIN("a", m.Keys())
		t.AssertIN(3, m.Values())
		t.AssertIN(1, m.Values())

		m_f := gmap.NewStrIntMap()
		m_f.Set("1", 2)
		m_f.Flip()
		t.Assert(m_f.Map(), map[string]int{"2": 1})

		m.Clear()
		t.Assert(m.Size(), 0)
		t.Assert(m.IsEmpty(), true)

		m2 := gmap.NewStrIntMapFrom(map[string]int{"a": 1, "b": 2})
		t.Assert(m2.Map(), map[string]int{"a": 1, "b": 2})
	})
}

func Test_StrIntMap_Set_Fun(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewStrIntMap()

		m.GetOrSetFunc("a", getInt)
		m.GetOrSetFuncLock("b", getInt)
		t.Assert(m.Get("a"), 123)
		t.Assert(m.Get("b"), 123)
		t.Assert(m.SetIfNotExistFunc("a", getInt), false)
		t.Assert(m.SetIfNotExistFunc("c", getInt), true)

		t.Assert(m.SetIfNotExistFuncLock("b", getInt), false)
		t.Assert(m.SetIfNotExistFuncLock("d", getInt), true)
	})
}

func Test_StrIntMap_Batch(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewStrIntMap()

		m.Sets(map[string]int{"a": 1, "b": 2, "c": 3})
		t.Assert(m.Map(), map[string]int{"a": 1, "b": 2, "c": 3})
		m.Removes([]string{"a", "b"})
		t.Assert(m.Map(), map[string]int{"c": 3})
	})
}
func Test_StrIntMap_Iterator(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		expect := map[string]int{"a": 1, "b": 2}
		m := gmap.NewStrIntMapFrom(expect)
		m.Iterator(func(k string, v int) bool {
			t.Assert(expect[k], v)
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
		t.Assert(i, 2)
		t.Assert(j, 1)
	})
}

func Test_StrIntMap_Lock(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		expect := map[string]int{"a": 1, "b": 2}

		m := gmap.NewStrIntMapFrom(expect)
		m.LockFunc(func(m map[string]int) {
			t.Assert(m, expect)
		})
		m.RLockFunc(func(m map[string]int) {
			t.Assert(m, expect)
		})
	})
}

func Test_StrIntMap_Clone(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		//clone 方法是深克隆
		m := gmap.NewStrIntMapFrom(map[string]int{"a": 1, "b": 2, "c": 3})

		m_clone := m.Clone()
		m.Remove("a")
		//修改原 map,clone 后的 map 不影响
		t.AssertIN("a", m_clone.Keys())

		m_clone.Remove("b")
		//修改clone map,原 map 不影响
		t.AssertIN("b", m.Keys())
	})
}
func Test_StrIntMap_Merge(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m1 := gmap.NewStrIntMap()
		m2 := gmap.NewStrIntMap()
		m1.Set("a", 1)
		m2.Set("b", 2)
		m1.Merge(m2)
		t.Assert(m1.Map(), map[string]int{"a": 1, "b": 2})
	})
}

func Test_StrIntMap_Map(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewStrIntMap()
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

func Test_StrIntMap_MapCopy(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewStrIntMap()
		m.Set("1", 1)
		m.Set("2", 2)
		t.Assert(m.Get("1"), 1)
		t.Assert(m.Get("2"), 2)
		data := m.MapCopy()
		t.Assert(data["1"], 1)
		t.Assert(data["2"], 2)
		data["3"] = 3
		t.Assert(m.Get("3"), 0)
		m.Set("4", 4)
		t.Assert(data["4"], 0)
	})
}

func Test_StrIntMap_FilterEmpty(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewStrIntMap()
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

func Test_StrIntMap_Json(t *testing.T) {
	// Marshal
	gtest.C(t, func(t *gtest.T) {
		data := g.MapStrInt{
			"k1": 1,
			"k2": 2,
		}
		m1 := gmap.NewStrIntMapFrom(data)
		b1, err1 := json.Marshal(m1)
		b2, err2 := json.Marshal(data)
		t.Assert(err1, err2)
		t.Assert(b1, b2)
	})
	// Unmarshal
	gtest.C(t, func(t *gtest.T) {
		data := g.MapStrInt{
			"k1": 1,
			"k2": 2,
		}
		b, err := json.Marshal(data)
		t.Assert(err, nil)

		m := gmap.NewStrIntMap()
		err = json.Unmarshal(b, m)
		t.Assert(err, nil)
		t.Assert(m.Get("k1"), data["k1"])
		t.Assert(m.Get("k2"), data["k2"])
	})
	gtest.C(t, func(t *gtest.T) {
		data := g.MapStrInt{
			"k1": 1,
			"k2": 2,
		}
		b, err := json.Marshal(data)
		t.Assert(err, nil)

		var m gmap.StrIntMap
		err = json.Unmarshal(b, &m)
		t.Assert(err, nil)
		t.Assert(m.Get("k1"), data["k1"])
		t.Assert(m.Get("k2"), data["k2"])
	})
}

func Test_StrIntMap_Pop(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewStrIntMapFrom(g.MapStrInt{
			"k1": 11,
			"k2": 22,
		})
		t.Assert(m.Size(), 2)

		k1, v1 := m.Pop()
		t.AssertIN(k1, g.Slice{"k1", "k2"})
		t.AssertIN(v1, g.Slice{11, 22})
		t.Assert(m.Size(), 1)
		k2, v2 := m.Pop()
		t.AssertIN(k2, g.Slice{"k1", "k2"})
		t.AssertIN(v2, g.Slice{11, 22})
		t.Assert(m.Size(), 0)

		t.AssertNE(k1, k2)
		t.AssertNE(v1, v2)
	})
}

func Test_StrIntMap_Pops(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewStrIntMapFrom(g.MapStrInt{
			"k1": 11,
			"k2": 22,
			"k3": 33,
		})
		t.Assert(m.Size(), 3)

		kArray := garray.New()
		vArray := garray.New()
		for k, v := range m.Pops(1) {
			t.AssertIN(k, g.Slice{"k1", "k2", "k3"})
			t.AssertIN(v, g.Slice{11, 22, 33})
			kArray.Append(k)
			vArray.Append(v)
		}
		t.Assert(m.Size(), 2)
		for k, v := range m.Pops(2) {
			t.AssertIN(k, g.Slice{"k1", "k2", "k3"})
			t.AssertIN(v, g.Slice{11, 22, 33})
			kArray.Append(k)
			vArray.Append(v)
		}
		t.Assert(m.Size(), 0)

		t.Assert(kArray.Unique().Len(), 3)
		t.Assert(vArray.Unique().Len(), 3)
	})
}

func TestStrIntMap_UnmarshalValue(t *testing.T) {
	type V struct {
		Name string
		Map  *gmap.StrIntMap
	}
	// JSON
	gtest.C(t, func(t *gtest.T) {
		var v *V
		err := gconv.Struct(map[string]interface{}{
			"name": "john",
			"map":  []byte(`{"k1":1,"k2":2}`),
		}, &v)
		t.Assert(err, nil)
		t.Assert(v.Name, "john")
		t.Assert(v.Map.Size(), 2)
		t.Assert(v.Map.Get("k1"), 1)
		t.Assert(v.Map.Get("k2"), 2)
	})
	// Map
	gtest.C(t, func(t *gtest.T) {
		var v *V
		err := gconv.Struct(map[string]interface{}{
			"name": "john",
			"map": g.Map{
				"k1": 1,
				"k2": 2,
			},
		}, &v)
		t.Assert(err, nil)
		t.Assert(v.Name, "john")
		t.Assert(v.Map.Size(), 2)
		t.Assert(v.Map.Get("k1"), 1)
		t.Assert(v.Map.Get("k2"), 2)
	})
}
