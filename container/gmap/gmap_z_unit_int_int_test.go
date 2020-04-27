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

func getInt() int {
	return 123
}

func intIntCallBack(int, int) bool {
	return true
}

func Test_IntIntMap_Var(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var m gmap.IntIntMap
		m.Set(1, 1)

		t.Assert(m.Get(1), 1)
		t.Assert(m.Size(), 1)
		t.Assert(m.IsEmpty(), false)

		t.Assert(m.GetOrSet(2, 2), 2)
		t.Assert(m.SetIfNotExist(2, 2), false)

		t.Assert(m.SetIfNotExist(3, 3), true)

		t.Assert(m.Remove(2), 2)
		t.Assert(m.Contains(2), false)

		t.AssertIN(3, m.Keys())
		t.AssertIN(1, m.Keys())
		t.AssertIN(3, m.Values())
		t.AssertIN(1, m.Values())
		m.Flip()
		t.Assert(m.Map(), map[int]int{1: 1, 3: 3})

		m.Clear()
		t.Assert(m.Size(), 0)
		t.Assert(m.IsEmpty(), true)
	})
}

func Test_IntIntMap_Basic(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewIntIntMap()
		m.Set(1, 1)

		t.Assert(m.Get(1), 1)
		t.Assert(m.Size(), 1)
		t.Assert(m.IsEmpty(), false)

		t.Assert(m.GetOrSet(2, 2), 2)
		t.Assert(m.SetIfNotExist(2, 2), false)

		t.Assert(m.SetIfNotExist(3, 3), true)

		t.Assert(m.Remove(2), 2)
		t.Assert(m.Contains(2), false)

		t.AssertIN(3, m.Keys())
		t.AssertIN(1, m.Keys())
		t.AssertIN(3, m.Values())
		t.AssertIN(1, m.Values())
		m.Flip()
		t.Assert(m.Map(), map[int]int{1: 1, 3: 3})

		m.Clear()
		t.Assert(m.Size(), 0)
		t.Assert(m.IsEmpty(), true)

		m2 := gmap.NewIntIntMapFrom(map[int]int{1: 1, 2: 2})
		t.Assert(m2.Map(), map[int]int{1: 1, 2: 2})
	})
}

func Test_IntIntMap_Set_Fun(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewIntIntMap()

		m.GetOrSetFunc(1, getInt)
		m.GetOrSetFuncLock(2, getInt)
		t.Assert(m.Get(1), 123)
		t.Assert(m.Get(2), 123)
		t.Assert(m.SetIfNotExistFunc(1, getInt), false)
		t.Assert(m.SetIfNotExistFunc(3, getInt), true)

		t.Assert(m.SetIfNotExistFuncLock(2, getInt), false)
		t.Assert(m.SetIfNotExistFuncLock(4, getInt), true)
	})
}

func Test_IntIntMap_Batch(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewIntIntMap()

		m.Sets(map[int]int{1: 1, 2: 2, 3: 3})
		m.Iterator(intIntCallBack)
		t.Assert(m.Map(), map[int]int{1: 1, 2: 2, 3: 3})
		m.Removes([]int{1, 2})
		t.Assert(m.Map(), map[int]int{3: 3})
	})
}

func Test_IntIntMap_Iterator(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		expect := map[int]int{1: 1, 2: 2}
		m := gmap.NewIntIntMapFrom(expect)
		m.Iterator(func(k int, v int) bool {
			t.Assert(expect[k], v)
			return true
		})
		// 断言返回值对遍历控制
		i := 0
		j := 0
		m.Iterator(func(k int, v int) bool {
			i++
			return true
		})
		m.Iterator(func(k int, v int) bool {
			j++
			return false
		})
		t.Assert(i, 2)
		t.Assert(j, 1)
	})
}

func Test_IntIntMap_Lock(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		expect := map[int]int{1: 1, 2: 2}
		m := gmap.NewIntIntMapFrom(expect)
		m.LockFunc(func(m map[int]int) {
			t.Assert(m, expect)
		})
		m.RLockFunc(func(m map[int]int) {
			t.Assert(m, expect)
		})
	})
}

func Test_IntIntMap_Clone(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		//clone 方法是深克隆
		m := gmap.NewIntIntMapFrom(map[int]int{1: 1, 2: 2})

		m_clone := m.Clone()
		m.Remove(1)
		//修改原 map,clone 后的 map 不影响
		t.AssertIN(1, m_clone.Keys())

		m_clone.Remove(2)
		//修改clone map,原 map 不影响
		t.AssertIN(2, m.Keys())
	})
}

func Test_IntIntMap_Merge(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m1 := gmap.NewIntIntMap()
		m2 := gmap.NewIntIntMap()
		m1.Set(1, 1)
		m2.Set(2, 2)
		m1.Merge(m2)
		t.Assert(m1.Map(), map[int]int{1: 1, 2: 2})
	})
}

func Test_IntIntMap_Map(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewIntIntMap()
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

func Test_IntIntMap_MapCopy(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewIntIntMap()
		m.Set(1, 0)
		m.Set(2, 2)
		t.Assert(m.Get(1), 0)
		t.Assert(m.Get(2), 2)
		data := m.MapCopy()
		t.Assert(data[1], 0)
		t.Assert(data[2], 2)
		data[3] = 3
		t.Assert(m.Get(3), 0)
		m.Set(4, 4)
		t.Assert(data[4], 0)
	})
}

func Test_IntIntMap_FilterEmpty(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewIntIntMap()
		m.Set(1, 0)
		m.Set(2, 2)
		t.Assert(m.Size(), 2)
		t.Assert(m.Get(1), 0)
		t.Assert(m.Get(2), 2)
		m.FilterEmpty()
		t.Assert(m.Size(), 1)
		t.Assert(m.Get(2), 2)
	})
}

func Test_IntIntMap_Json(t *testing.T) {
	// Marshal
	gtest.C(t, func(t *gtest.T) {
		data := g.MapIntInt{
			1: 10,
			2: 20,
		}
		m1 := gmap.NewIntIntMapFrom(data)
		b1, err1 := json.Marshal(m1)
		b2, err2 := json.Marshal(data)
		t.Assert(err1, err2)
		t.Assert(b1, b2)
	})
	// Unmarshal
	gtest.C(t, func(t *gtest.T) {
		data := g.MapIntInt{
			1: 10,
			2: 20,
		}
		b, err := json.Marshal(data)
		t.Assert(err, nil)

		m := gmap.NewIntIntMap()
		err = json.Unmarshal(b, m)
		t.Assert(err, nil)
		t.Assert(m.Get(1), data[1])
		t.Assert(m.Get(2), data[2])
	})
}

func Test_IntIntMap_Pop(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewIntIntMapFrom(g.MapIntInt{
			1: 11,
			2: 22,
		})
		t.Assert(m.Size(), 2)

		k1, v1 := m.Pop()
		t.AssertIN(k1, g.Slice{1, 2})
		t.AssertIN(v1, g.Slice{11, 22})
		t.Assert(m.Size(), 1)
		k2, v2 := m.Pop()
		t.AssertIN(k2, g.Slice{1, 2})
		t.AssertIN(v2, g.Slice{11, 22})
		t.Assert(m.Size(), 0)

		t.AssertNE(k1, k2)
		t.AssertNE(v1, v2)
	})
}

func Test_IntIntMap_Pops(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewIntIntMapFrom(g.MapIntInt{
			1: 11,
			2: 22,
			3: 33,
		})
		t.Assert(m.Size(), 3)

		kArray := garray.New()
		vArray := garray.New()
		for k, v := range m.Pops(1) {
			t.AssertIN(k, g.Slice{1, 2, 3})
			t.AssertIN(v, g.Slice{11, 22, 33})
			kArray.Append(k)
			vArray.Append(v)
		}
		t.Assert(m.Size(), 2)
		for k, v := range m.Pops(2) {
			t.AssertIN(k, g.Slice{1, 2, 3})
			t.AssertIN(v, g.Slice{11, 22, 33})
			kArray.Append(k)
			vArray.Append(v)
		}
		t.Assert(m.Size(), 0)

		t.Assert(kArray.Unique().Len(), 3)
		t.Assert(vArray.Unique().Len(), 3)
	})
}

func TestIntIntMap_UnmarshalValue(t *testing.T) {
	type V struct {
		Name string
		Map  *gmap.IntIntMap
	}
	// JSON
	gtest.C(t, func(t *gtest.T) {
		var v *V
		err := gconv.Struct(map[string]interface{}{
			"name": "john",
			"map":  []byte(`{"1":1,"2":2}`),
		}, &v)
		t.Assert(err, nil)
		t.Assert(v.Name, "john")
		t.Assert(v.Map.Size(), 2)
		t.Assert(v.Map.Get(1), "1")
		t.Assert(v.Map.Get(2), "2")
	})
	// Map
	gtest.C(t, func(t *gtest.T) {
		var v *V
		err := gconv.Struct(map[string]interface{}{
			"name": "john",
			"map": g.MapIntAny{
				1: 1,
				2: 2,
			},
		}, &v)
		t.Assert(err, nil)
		t.Assert(v.Name, "john")
		t.Assert(v.Map.Size(), 2)
		t.Assert(v.Map.Get(1), "1")
		t.Assert(v.Map.Get(2), "2")
	})
}
