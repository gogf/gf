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

func getAny() interface{} {
	return 123
}
func intAnyCallBack(int, interface{}) bool {
	return true
}
func Test_IntAnyMap_Basic(t *testing.T) {
	gtest.Case(t, func() {
		m := gmap.NewIntAnyMap()
		m.Set(1, 1)

		gtest.Assert(m.Get(1), 1)
		gtest.Assert(m.Size(), 1)
		gtest.Assert(m.IsEmpty(), false)

		gtest.Assert(m.GetOrSet(2, "2"), "2")
		gtest.Assert(m.SetIfNotExist(2, "2"), false)

		gtest.Assert(m.SetIfNotExist(3, 3), true)

		gtest.Assert(m.Remove(2), "2")
		gtest.Assert(m.Contains(2), false)

		gtest.AssertIN(3, m.Keys())
		gtest.AssertIN(1, m.Keys())
		gtest.AssertIN(3, m.Values())
		gtest.AssertIN(1, m.Values())
		m.Flip()
		gtest.Assert(m.Map(), map[interface{}]int{1: 1, 3: 3})

		m.Clear()
		gtest.Assert(m.Size(), 0)
		gtest.Assert(m.IsEmpty(), true)

		m2 := gmap.NewIntAnyMapFrom(map[int]interface{}{1: 1, 2: "2"})
		gtest.Assert(m2.Map(), map[int]interface{}{1: 1, 2: "2"})
	})
}
func Test_IntAnyMap_Set_Fun(t *testing.T) {
	m := gmap.NewIntAnyMap()

	m.GetOrSetFunc(1, getAny)
	m.GetOrSetFuncLock(2, getAny)
	gtest.Assert(m.Get(1), 123)
	gtest.Assert(m.Get(2), 123)

	gtest.Assert(m.SetIfNotExistFunc(1, getAny), false)
	gtest.Assert(m.SetIfNotExistFunc(3, getAny), true)

	gtest.Assert(m.SetIfNotExistFuncLock(2, getAny), false)
	gtest.Assert(m.SetIfNotExistFuncLock(4, getAny), true)

}

func Test_IntAnyMap_Batch(t *testing.T) {
	m := gmap.NewIntAnyMap()

	m.Sets(map[int]interface{}{1: 1, 2: "2", 3: 3})
	gtest.Assert(m.Map(), map[int]interface{}{1: 1, 2: "2", 3: 3})
	m.Removes([]int{1, 2})
	gtest.Assert(m.Map(), map[int]interface{}{3: 3})
}
func Test_IntAnyMap_Iterator(t *testing.T) {
	expect := map[int]interface{}{1: 1, 2: "2"}
	m := gmap.NewIntAnyMapFrom(expect)
	m.Iterator(func(k int, v interface{}) bool {
		gtest.Assert(expect[k], v)
		return true
	})
	// 断言返回值对遍历控制
	i := 0
	j := 0
	m.Iterator(func(k int, v interface{}) bool {
		i++
		return true
	})
	m.Iterator(func(k int, v interface{}) bool {
		j++
		return false
	})
	gtest.Assert(i, "2")
	gtest.Assert(j, 1)

}

func Test_IntAnyMap_Lock(t *testing.T) {
	expect := map[int]interface{}{1: 1, 2: "2"}
	m := gmap.NewIntAnyMapFrom(expect)
	m.LockFunc(func(m map[int]interface{}) {
		gtest.Assert(m, expect)
	})
	m.RLockFunc(func(m map[int]interface{}) {
		gtest.Assert(m, expect)
	})
}
func Test_IntAnyMap_Clone(t *testing.T) {
	//clone 方法是深克隆
	m := gmap.NewIntAnyMapFrom(map[int]interface{}{1: 1, 2: "2"})

	m_clone := m.Clone()
	m.Remove(1)
	//修改原 map,clone 后的 map 不影响
	gtest.AssertIN(1, m_clone.Keys())

	m_clone.Remove(2)
	//修改clone map,原 map 不影响
	gtest.AssertIN(2, m.Keys())
}
func Test_IntAnyMap_Merge(t *testing.T) {
	m1 := gmap.NewIntAnyMap()
	m2 := gmap.NewIntAnyMap()
	m1.Set(1, 1)
	m2.Set(2, "2")
	m1.Merge(m2)
	gtest.Assert(m1.Map(), map[int]interface{}{1: 1, 2: "2"})
}

func Test_IntAnyMap_Map(t *testing.T) {
	m := gmap.NewIntAnyMap()
	m.Set(1, 0)
	m.Set(2, 2)
	gtest.Assert(m.Get(1), 0)
	gtest.Assert(m.Get(2), 2)
	data := m.Map()
	gtest.Assert(data[1], 0)
	gtest.Assert(data[2], 2)
	data[3] = 3
	gtest.Assert(m.Get(3), 3)
	m.Set(4, 4)
	gtest.Assert(data[4], 4)
}

func Test_IntAnyMap_MapCopy(t *testing.T) {
	m := gmap.NewIntAnyMap()
	m.Set(1, 0)
	m.Set(2, 2)
	gtest.Assert(m.Get(1), 0)
	gtest.Assert(m.Get(2), 2)
	data := m.MapCopy()
	gtest.Assert(data[1], 0)
	gtest.Assert(data[2], 2)
	data[3] = 3
	gtest.Assert(m.Get(3), nil)
	m.Set(4, 4)
	gtest.Assert(data[4], nil)
}

func Test_IntAnyMap_FilterEmpty(t *testing.T) {
	m := gmap.NewIntAnyMap()
	m.Set(1, 0)
	m.Set(2, 2)
	gtest.Assert(m.Size(), 2)
	gtest.Assert(m.Get(1), 0)
	gtest.Assert(m.Get(2), 2)
	m.FilterEmpty()
	gtest.Assert(m.Size(), 1)
	gtest.Assert(m.Get(2), 2)
}

func Test_IntAnyMap_Json(t *testing.T) {
	// Marshal
	gtest.Case(t, func() {
		data := g.MapIntAny{
			1: "v1",
			2: "v2",
		}
		m1 := gmap.NewIntAnyMapFrom(data)
		b1, err1 := json.Marshal(m1)
		b2, err2 := json.Marshal(data)
		gtest.Assert(err1, err2)
		gtest.Assert(b1, b2)
	})
	// Unmarshal
	gtest.Case(t, func() {
		data := g.MapIntAny{
			1: "v1",
			2: "v2",
		}
		b, err := json.Marshal(data)
		gtest.Assert(err, nil)

		m := gmap.NewIntAnyMap()
		err = json.Unmarshal(b, m)
		gtest.Assert(err, nil)
		gtest.Assert(m.Get(1), data[1])
		gtest.Assert(m.Get(2), data[2])
	})
}

func Test_IntAnyMap_Pop(t *testing.T) {
	gtest.Case(t, func() {
		m := gmap.NewIntAnyMapFrom(g.MapIntAny{
			1: "v1",
			2: "v2",
		})
		gtest.Assert(m.Size(), 2)

		k1, v1 := m.Pop()
		gtest.AssertIN(k1, g.Slice{1, 2})
		gtest.AssertIN(v1, g.Slice{"v1", "v2"})
		gtest.Assert(m.Size(), 1)
		k2, v2 := m.Pop()
		gtest.AssertIN(k2, g.Slice{1, 2})
		gtest.AssertIN(v2, g.Slice{"v1", "v2"})
		gtest.Assert(m.Size(), 0)

		gtest.AssertNE(k1, k2)
		gtest.AssertNE(v1, v2)
	})
}

func Test_IntAnyMap_Pops(t *testing.T) {
	gtest.Case(t, func() {
		m := gmap.NewIntAnyMapFrom(g.MapIntAny{
			1: "v1",
			2: "v2",
			3: "v3",
		})
		gtest.Assert(m.Size(), 3)

		kArray := garray.New()
		vArray := garray.New()
		for k, v := range m.Pops(1) {
			gtest.AssertIN(k, g.Slice{1, 2, 3})
			gtest.AssertIN(v, g.Slice{"v1", "v2", "v3"})
			kArray.Append(k)
			vArray.Append(v)
		}
		gtest.Assert(m.Size(), 2)
		for k, v := range m.Pops(2) {
			gtest.AssertIN(k, g.Slice{1, 2, 3})
			gtest.AssertIN(v, g.Slice{"v1", "v2", "v3"})
			kArray.Append(k)
			vArray.Append(v)
		}
		gtest.Assert(m.Size(), 0)

		gtest.Assert(kArray.Unique().Len(), 3)
		gtest.Assert(vArray.Unique().Len(), 3)
	})
}

func TestIntAnyMap_UnmarshalValue(t *testing.T) {
	type T struct {
		Name string
		Map  *gmap.IntAnyMap
	}
	// JSON
	gtest.Case(t, func() {
		var t *T
		err := gconv.Struct(map[string]interface{}{
			"name": "john",
			"map":  []byte(`{"1":"v1","2":"v2"}`),
		}, &t)
		gtest.Assert(err, nil)
		gtest.Assert(t.Name, "john")
		gtest.Assert(t.Map.Size(), 2)
		gtest.Assert(t.Map.Get(1), "v1")
		gtest.Assert(t.Map.Get(2), "v2")
	})
	// Map
	gtest.Case(t, func() {
		var t *T
		err := gconv.Struct(map[string]interface{}{
			"name": "john",
			"map": g.MapIntAny{
				1: "v1",
				2: "v2",
			},
		}, &t)
		gtest.Assert(err, nil)
		gtest.Assert(t.Name, "john")
		gtest.Assert(t.Map.Size(), 2)
		gtest.Assert(t.Map.Get(1), "v1")
		gtest.Assert(t.Map.Get(2), "v2")
	})
}
