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

func getStr() string {
	return "z"
}
func intStrCallBack(int, string) bool {
	return true
}
func Test_IntStrMap_Basic(t *testing.T) {
	gtest.Case(t, func() {
		m := gmap.NewIntStrMap()
		m.Set(1, "a")

		gtest.Assert(m.Get(1), "a")
		gtest.Assert(m.Size(), 1)
		gtest.Assert(m.IsEmpty(), false)

		gtest.Assert(m.GetOrSet(2, "b"), "b")
		gtest.Assert(m.SetIfNotExist(2, "b"), false)

		gtest.Assert(m.SetIfNotExist(3, "c"), true)

		gtest.Assert(m.Remove(2), "b")
		gtest.Assert(m.Contains(2), false)

		gtest.AssertIN(3, m.Keys())
		gtest.AssertIN(1, m.Keys())
		gtest.AssertIN("a", m.Values())
		gtest.AssertIN("c", m.Values())

		//反转之后不成为以下 map,flip 操作只是翻转原 map
		//gtest.Assert(m.Map(), map[string]int{"a": 1, "c": 3})
		m_f := gmap.NewIntStrMap()
		m_f.Set(1, "2")
		m_f.Flip()
		gtest.Assert(m_f.Map(), map[int]string{2: "1"})

		m.Clear()
		gtest.Assert(m.Size(), 0)
		gtest.Assert(m.IsEmpty(), true)

		m2 := gmap.NewIntStrMapFrom(map[int]string{1: "a", 2: "b"})
		gtest.Assert(m2.Map(), map[int]string{1: "a", 2: "b"})
	})
}
func Test_IntStrMap_Set_Fun(t *testing.T) {
	m := gmap.NewIntStrMap()

	m.GetOrSetFunc(1, getStr)
	m.GetOrSetFuncLock(2, getStr)
	gtest.Assert(m.Get(1), "z")
	gtest.Assert(m.Get(2), "z")
	gtest.Assert(m.SetIfNotExistFunc(1, getStr), false)
	gtest.Assert(m.SetIfNotExistFunc(3, getStr), true)

	gtest.Assert(m.SetIfNotExistFuncLock(2, getStr), false)
	gtest.Assert(m.SetIfNotExistFuncLock(4, getStr), true)

}

func Test_IntStrMap_Batch(t *testing.T) {
	m := gmap.NewIntStrMap()

	m.Sets(map[int]string{1: "a", 2: "b", 3: "c"})
	gtest.Assert(m.Map(), map[int]string{1: "a", 2: "b", 3: "c"})
	m.Removes([]int{1, 2})
	gtest.Assert(m.Map(), map[int]interface{}{3: "c"})
}
func Test_IntStrMap_Iterator(t *testing.T) {
	expect := map[int]string{1: "a", 2: "b"}
	m := gmap.NewIntStrMapFrom(expect)
	m.Iterator(func(k int, v string) bool {
		gtest.Assert(expect[k], v)
		return true
	})
	// 断言返回值对遍历控制
	i := 0
	j := 0
	m.Iterator(func(k int, v string) bool {
		i++
		return true
	})
	m.Iterator(func(k int, v string) bool {
		j++
		return false
	})
	gtest.Assert(i, 2)
	gtest.Assert(j, 1)
}

func Test_IntStrMap_Lock(t *testing.T) {

	expect := map[int]string{1: "a", 2: "b", 3: "c"}
	m := gmap.NewIntStrMapFrom(expect)
	m.LockFunc(func(m map[int]string) {
		gtest.Assert(m, expect)
	})
	m.RLockFunc(func(m map[int]string) {
		gtest.Assert(m, expect)
	})

}
func Test_IntStrMap_Clone(t *testing.T) {
	//clone 方法是深克隆
	m := gmap.NewIntStrMapFrom(map[int]string{1: "a", 2: "b", 3: "c"})

	m_clone := m.Clone()
	m.Remove(1)
	//修改原 map,clone 后的 map 不影响
	gtest.AssertIN(1, m_clone.Keys())

	m_clone.Remove(2)
	//修改clone map,原 map 不影响
	gtest.AssertIN(2, m.Keys())
}
func Test_IntStrMap_Merge(t *testing.T) {
	m1 := gmap.NewIntStrMap()
	m2 := gmap.NewIntStrMap()
	m1.Set(1, "a")
	m2.Set(2, "b")
	m1.Merge(m2)
	gtest.Assert(m1.Map(), map[int]string{1: "a", 2: "b"})
}

func Test_IntStrMap_Map(t *testing.T) {
	m := gmap.NewIntStrMap()
	m.Set(1, "0")
	m.Set(2, "2")
	gtest.Assert(m.Get(1), "0")
	gtest.Assert(m.Get(2), "2")
	data := m.Map()
	gtest.Assert(data[1], "0")
	gtest.Assert(data[2], "2")
	data[3] = "3"
	gtest.Assert(m.Get(3), "3")
	m.Set(4, "4")
	gtest.Assert(data[4], "4")
}

func Test_IntStrMap_MapCopy(t *testing.T) {
	m := gmap.NewIntStrMap()
	m.Set(1, "0")
	m.Set(2, "2")
	gtest.Assert(m.Get(1), "0")
	gtest.Assert(m.Get(2), "2")
	data := m.MapCopy()
	gtest.Assert(data[1], "0")
	gtest.Assert(data[2], "2")
	data[3] = "3"
	gtest.Assert(m.Get(3), "")
	m.Set(4, "4")
	gtest.Assert(data[4], "")
}

func Test_IntStrMap_FilterEmpty(t *testing.T) {
	m := gmap.NewIntStrMap()
	m.Set(1, "")
	m.Set(2, "2")
	gtest.Assert(m.Size(), 2)
	gtest.Assert(m.Get(2), "2")
	m.FilterEmpty()
	gtest.Assert(m.Size(), 1)
	gtest.Assert(m.Get(2), "2")
}

func Test_IntStrMap_Json(t *testing.T) {
	// Marshal
	gtest.Case(t, func() {
		data := g.MapIntStr{
			1: "v1",
			2: "v2",
		}
		m1 := gmap.NewIntStrMapFrom(data)
		b1, err1 := json.Marshal(m1)
		b2, err2 := json.Marshal(data)
		gtest.Assert(err1, err2)
		gtest.Assert(b1, b2)
	})
	// Unmarshal
	gtest.Case(t, func() {
		data := g.MapIntStr{
			1: "v1",
			2: "v2",
		}
		b, err := json.Marshal(data)
		gtest.Assert(err, nil)

		m := gmap.NewIntStrMap()
		err = json.Unmarshal(b, m)
		gtest.Assert(err, nil)
		gtest.Assert(m.Get(1), data[1])
		gtest.Assert(m.Get(2), data[2])
	})
}

func Test_IntStrMap_Pop(t *testing.T) {
	gtest.Case(t, func() {
		m := gmap.NewIntStrMapFrom(g.MapIntStr{
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

func Test_IntStrMap_Pops(t *testing.T) {
	gtest.Case(t, func() {
		m := gmap.NewIntStrMapFrom(g.MapIntStr{
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

func TestIntStrMap_UnmarshalValue(t *testing.T) {
	type T struct {
		Name string
		Map  *gmap.IntStrMap
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
