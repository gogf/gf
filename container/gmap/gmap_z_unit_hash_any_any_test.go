// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with gm file,
// You can obtain one at https://github.com/gogf/gf.

package gmap_test

import (
	"testing"
	"time"

	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gconv"
)

func Test_AnyAnyMap_Var(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var m gmap.AnyAnyMap
		m.Set(1, 1)

		t.Assert(m.Get(1), 1)
		t.Assert(m.Size(), 1)
		t.Assert(m.IsEmpty(), false)

		t.Assert(m.GetOrSet(2, "2"), "2")
		t.Assert(m.SetIfNotExist(2, "2"), false)

		t.Assert(m.SetIfNotExist(3, 3), true)

		t.Assert(m.Remove(2), "2")
		t.Assert(m.Contains(2), false)

		t.AssertIN(3, m.Keys())
		t.AssertIN(1, m.Keys())
		t.AssertIN(3, m.Values())
		t.AssertIN(1, m.Values())
		m.Flip()
		t.Assert(m.Map(), map[interface{}]int{1: 1, 3: 3})

		m.Clear()
		t.Assert(m.Size(), 0)
		t.Assert(m.IsEmpty(), true)
	})
}

func Test_AnyAnyMap_Basic(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewAnyAnyMap()
		m.Set(1, 1)

		t.Assert(m.Get(1), 1)
		t.Assert(m.Size(), 1)
		t.Assert(m.IsEmpty(), false)

		t.Assert(m.GetOrSet(2, "2"), "2")
		t.Assert(m.SetIfNotExist(2, "2"), false)

		t.Assert(m.SetIfNotExist(3, 3), true)

		t.Assert(m.Remove(2), "2")
		t.Assert(m.Contains(2), false)

		t.AssertIN(3, m.Keys())
		t.AssertIN(1, m.Keys())
		t.AssertIN(3, m.Values())
		t.AssertIN(1, m.Values())
		m.Flip()
		t.Assert(m.Map(), map[interface{}]int{1: 1, 3: 3})

		m.Clear()
		t.Assert(m.Size(), 0)
		t.Assert(m.IsEmpty(), true)

		m2 := gmap.NewAnyAnyMapFrom(map[interface{}]interface{}{1: 1, 2: "2"})
		t.Assert(m2.Map(), map[interface{}]interface{}{1: 1, 2: "2"})
	})
}

func Test_AnyAnyMap_Set_Fun(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewAnyAnyMap()

		m.GetOrSetFunc(1, getAny)
		m.GetOrSetFuncLock(2, getAny)
		t.Assert(m.Get(1), 123)
		t.Assert(m.Get(2), 123)

		t.Assert(m.SetIfNotExistFunc(1, getAny), false)
		t.Assert(m.SetIfNotExistFunc(3, getAny), true)

		t.Assert(m.SetIfNotExistFuncLock(2, getAny), false)
		t.Assert(m.SetIfNotExistFuncLock(4, getAny), true)
	})

}

func Test_AnyAnyMap_Batch(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewAnyAnyMap()

		m.Sets(map[interface{}]interface{}{1: 1, 2: "2", 3: 3})
		t.Assert(m.Map(), map[interface{}]interface{}{1: 1, 2: "2", 3: 3})
		m.Removes([]interface{}{1, 2})
		t.Assert(m.Map(), map[interface{}]interface{}{3: 3})
	})
}

func Test_AnyAnyMap_Iterator_Deadlock(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewAnyAnyMapFrom(map[interface{}]interface{}{1: 1, 2: "2", "3": "3", "4": 4}, true)
		m.Iterator(func(k interface{}, _ interface{}) bool {
			if gconv.Int(k)%2 == 0 {
				m.Remove(k)
			}
			return true
		})
		t.Assert(m.Map(), map[interface{}]interface{}{
			1:   1,
			"3": "3",
		})
	})
}

func Test_AnyAnyMap_Iterator(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		expect := map[interface{}]interface{}{1: 1, 2: "2"}
		m := gmap.NewAnyAnyMapFrom(expect)
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
		t.Assert(i, "2")
		t.Assert(j, 1)
	})
}

func Test_AnyAnyMap_Lock(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		expect := map[interface{}]interface{}{1: 1, 2: "2"}
		m := gmap.NewAnyAnyMapFrom(expect)
		m.LockFunc(func(m map[interface{}]interface{}) {
			t.Assert(m, expect)
		})
		m.RLockFunc(func(m map[interface{}]interface{}) {
			t.Assert(m, expect)
		})
	})
}

func Test_AnyAnyMap_Clone(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// clone 方法是深克隆
		m := gmap.NewAnyAnyMapFrom(map[interface{}]interface{}{1: 1, 2: "2"})

		m_clone := m.Clone()
		m.Remove(1)
		// 修改原 map,clone 后的 map 不影响
		t.AssertIN(1, m_clone.Keys())

		m_clone.Remove(2)
		// 修改clone map,原 map 不影响
		t.AssertIN(2, m.Keys())
	})
}

func Test_AnyAnyMap_Merge(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m1 := gmap.NewAnyAnyMap()
		m2 := gmap.NewAnyAnyMap()
		m1.Set(1, 1)
		m2.Set(2, "2")
		m1.Merge(m2)
		t.Assert(m1.Map(), map[interface{}]interface{}{1: 1, 2: "2"})
		m3 := gmap.NewAnyAnyMapFrom(nil)
		m3.Merge(m2)
		t.Assert(m3.Map(), m2.Map())
	})
}

func Test_AnyAnyMap_Map(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewAnyAnyMap()
		m.Set(1, 0)
		m.Set(2, 2)
		t.Assert(m.Get(1), 0)
		t.Assert(m.Get(2), 2)
		data := m.Map()
		t.Assert(data[1], 0)
		t.Assert(data[2], 2)
		data[3] = 3
		t.Assert(m.Get(3), 3)
		m.Set(4, 4)
		t.Assert(data[4], 4)
	})
}

func Test_AnyAnyMap_MapCopy(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewAnyAnyMap()
		m.Set(1, 0)
		m.Set(2, 2)
		t.Assert(m.Get(1), 0)
		t.Assert(m.Get(2), 2)
		data := m.MapCopy()
		t.Assert(data[1], 0)
		t.Assert(data[2], 2)
		data[3] = 3
		t.Assert(m.Get(3), nil)
		m.Set(4, 4)
		t.Assert(data[4], nil)
	})
}

func Test_AnyAnyMap_FilterEmpty(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewAnyAnyMap()
		m.Set(1, 0)
		m.Set(2, 2)
		t.Assert(m.Get(1), 0)
		t.Assert(m.Get(2), 2)
		m.FilterEmpty()
		t.Assert(m.Get(1), nil)
		t.Assert(m.Get(2), 2)
	})
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewAnyAnyMap()
		m.Set(1, 0)
		m.Set("time1", time.Time{})
		m.Set("time2", time.Now())
		t.Assert(m.Get(1), 0)
		t.Assert(m.Get("time1"), time.Time{})
		m.FilterEmpty()
		t.Assert(m.Get(1), nil)
		t.Assert(m.Get("time1"), nil)
		t.AssertNE(m.Get("time2"), nil)
	})
}

func Test_AnyAnyMap_Json(t *testing.T) {
	// Marshal
	gtest.C(t, func(t *gtest.T) {
		data := g.MapAnyAny{
			"k1": "v1",
			"k2": "v2",
		}
		m1 := gmap.NewAnyAnyMapFrom(data)
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
		t.AssertNil(err)

		m := gmap.New()
		err = json.UnmarshalUseNumber(b, m)
		t.AssertNil(err)
		t.Assert(m.Get("k1"), data["k1"])
		t.Assert(m.Get("k2"), data["k2"])
	})
	gtest.C(t, func(t *gtest.T) {
		data := g.MapAnyAny{
			"k1": "v1",
			"k2": "v2",
		}
		b, err := json.Marshal(gconv.Map(data))
		t.AssertNil(err)

		var m gmap.Map
		err = json.UnmarshalUseNumber(b, &m)
		t.AssertNil(err)
		t.Assert(m.Get("k1"), data["k1"])
		t.Assert(m.Get("k2"), data["k2"])
	})
}

func Test_AnyAnyMap_Pop(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewAnyAnyMapFrom(g.MapAnyAny{
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
		t.AssertNil(k3)
		t.AssertNil(v3)
	})
}

func Test_AnyAnyMap_Pops(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewAnyAnyMapFrom(g.MapAnyAny{
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

func TestAnyAnyMap_UnmarshalValue(t *testing.T) {
	type V struct {
		Name string
		Map  *gmap.Map
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

func Test_AnyAnyMap_DeepCopy(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewAnyAnyMapFrom(g.MapAnyAny{
			"k1": "v1",
			"k2": "v2",
		})
		t.Assert(m.Size(), 2)

		n := m.DeepCopy().(*gmap.AnyAnyMap)
		n.Set("k1", "val1")
		t.AssertNE(m.Get("k1"), n.Get("k1"))
	})
}

func Test_AnyAnyMap_IsSubOf(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m1 := gmap.NewAnyAnyMapFrom(g.MapAnyAny{
			"k1": "v1",
			"k2": "v2",
		})
		m2 := gmap.NewAnyAnyMapFrom(g.MapAnyAny{
			"k2": "v2",
		})
		t.Assert(m1.IsSubOf(m2), false)
		t.Assert(m2.IsSubOf(m1), true)
		t.Assert(m2.IsSubOf(m2), true)
	})
}

func Test_AnyAnyMap_Diff(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m1 := gmap.NewAnyAnyMapFrom(g.MapAnyAny{
			"0": "v0",
			"1": "v1",
			2:   "v2",
			3:   3,
		})
		m2 := gmap.NewAnyAnyMapFrom(g.MapAnyAny{
			"0": "v0",
			2:   "v2",
			3:   "v3",
			4:   "v4",
		})
		addedKeys, removedKeys, updatedKeys := m1.Diff(m2)
		t.Assert(addedKeys, []interface{}{4})
		t.Assert(removedKeys, []interface{}{"1"})
		t.Assert(updatedKeys, []interface{}{3})
	})
}
