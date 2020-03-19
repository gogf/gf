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
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewIntStrMap()
		m.Set(1, "a")

		t.Assert(m.Get(1), "a")
		t.Assert(m.Size(), 1)
		t.Assert(m.IsEmpty(), false)

		t.Assert(m.GetOrSet(2, "b"), "b")
		t.Assert(m.SetIfNotExist(2, "b"), false)

		t.Assert(m.SetIfNotExist(3, "c"), true)

		t.Assert(m.Remove(2), "b")
		t.Assert(m.Contains(2), false)

		t.AssertIN(3, m.Keys())
		t.AssertIN(1, m.Keys())
		t.AssertIN("a", m.Values())
		t.AssertIN("c", m.Values())

		//反转之后不成为以下 map,flip 操作只是翻转原 map
		//t.Assert(m.Map(), map[string]int{"a": 1, "c": 3})
		m_f := gmap.NewIntStrMap()
		m_f.Set(1, "2")
		m_f.Flip()
		t.Assert(m_f.Map(), map[int]string{2: "1"})

		m.Clear()
		t.Assert(m.Size(), 0)
		t.Assert(m.IsEmpty(), true)

		m2 := gmap.NewIntStrMapFrom(map[int]string{1: "a", 2: "b"})
		t.Assert(m2.Map(), map[int]string{1: "a", 2: "b"})
	})
}
func Test_IntStrMap_Set_Fun(t *testing.T) {
	m := gmap.NewIntStrMap()

	m.GetOrSetFunc(1, getStr)
	m.GetOrSetFuncLock(2, getStr)
	t.Assert(m.Get(1), "z")
	t.Assert(m.Get(2), "z")
	t.Assert(m.SetIfNotExistFunc(1, getStr), false)
	t.Assert(m.SetIfNotExistFunc(3, getStr), true)

	t.Assert(m.SetIfNotExistFuncLock(2, getStr), false)
	t.Assert(m.SetIfNotExistFuncLock(4, getStr), true)

}

func Test_IntStrMap_Batch(t *testing.T) {
	m := gmap.NewIntStrMap()

	m.Sets(map[int]string{1: "a", 2: "b", 3: "c"})
	t.Assert(m.Map(), map[int]string{1: "a", 2: "b", 3: "c"})
	m.Removes([]int{1, 2})
	t.Assert(m.Map(), map[int]interface{}{3: "c"})
}
func Test_IntStrMap_Iterator(t *testing.T) {
	expect := map[int]string{1: "a", 2: "b"}
	m := gmap.NewIntStrMapFrom(expect)
	m.Iterator(func(k int, v string) bool {
		t.Assert(expect[k], v)
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
	t.Assert(i, 2)
	t.Assert(j, 1)
}

func Test_IntStrMap_Lock(t *testing.T) {

	expect := map[int]string{1: "a", 2: "b", 3: "c"}
	m := gmap.NewIntStrMapFrom(expect)
	m.LockFunc(func(m map[int]string) {
		t.Assert(m, expect)
	})
	m.RLockFunc(func(m map[int]string) {
		t.Assert(m, expect)
	})

}
func Test_IntStrMap_Clone(t *testing.T) {
	//clone 方法是深克隆
	m := gmap.NewIntStrMapFrom(map[int]string{1: "a", 2: "b", 3: "c"})

	m_clone := m.Clone()
	m.Remove(1)
	//修改原 map,clone 后的 map 不影响
	t.AssertIN(1, m_clone.Keys())

	m_clone.Remove(2)
	//修改clone map,原 map 不影响
	t.AssertIN(2, m.Keys())
}
func Test_IntStrMap_Merge(t *testing.T) {
	m1 := gmap.NewIntStrMap()
	m2 := gmap.NewIntStrMap()
	m1.Set(1, "a")
	m2.Set(2, "b")
	m1.Merge(m2)
	t.Assert(m1.Map(), map[int]string{1: "a", 2: "b"})
}

func Test_IntStrMap_Map(t *testing.T) {
	m := gmap.NewIntStrMap()
	m.Set(1, "0")
	m.Set(2, "2")
	t.Assert(m.Get(1), "0")
	t.Assert(m.Get(2), "2")
	data := m.Map()
	t.Assert(data[1], "0")
	t.Assert(data[2], "2")
	data[3] = "3"
	t.Assert(m.Get(3), "3")
	m.Set(4, "4")
	t.Assert(data[4], "4")
}

func Test_IntStrMap_MapCopy(t *testing.T) {
	m := gmap.NewIntStrMap()
	m.Set(1, "0")
	m.Set(2, "2")
	t.Assert(m.Get(1), "0")
	t.Assert(m.Get(2), "2")
	data := m.MapCopy()
	t.Assert(data[1], "0")
	t.Assert(data[2], "2")
	data[3] = "3"
	t.Assert(m.Get(3), "")
	m.Set(4, "4")
	t.Assert(data[4], "")
}

func Test_IntStrMap_FilterEmpty(t *testing.T) {
	m := gmap.NewIntStrMap()
	m.Set(1, "")
	m.Set(2, "2")
	t.Assert(m.Size(), 2)
	t.Assert(m.Get(2), "2")
	m.FilterEmpty()
	t.Assert(m.Size(), 1)
	t.Assert(m.Get(2), "2")
}

func Test_IntStrMap_Json(t *testing.T) {
	// Marshal
	gtest.C(t, func(t *gtest.T) {
		data := g.MapIntStr{
			1: "v1",
			2: "v2",
		}
		m1 := gmap.NewIntStrMapFrom(data)
		b1, err1 := json.Marshal(m1)
		b2, err2 := json.Marshal(data)
		t.Assert(err1, err2)
		t.Assert(b1, b2)
	})
	// Unmarshal
	gtest.C(t, func(t *gtest.T) {
		data := g.MapIntStr{
			1: "v1",
			2: "v2",
		}
		b, err := json.Marshal(data)
		t.Assert(err, nil)

		m := gmap.NewIntStrMap()
		err = json.Unmarshal(b, m)
		t.Assert(err, nil)
		t.Assert(m.Get(1), data[1])
		t.Assert(m.Get(2), data[2])
	})
}

func Test_IntStrMap_Pop(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewIntStrMapFrom(g.MapIntStr{
			1: "v1",
			2: "v2",
		})
		t.Assert(m.Size(), 2)

		k1, v1 := m.Pop()
		t.AssertIN(k1, g.Slice{1, 2})
		t.AssertIN(v1, g.Slice{"v1", "v2"})
		t.Assert(m.Size(), 1)
		k2, v2 := m.Pop()
		t.AssertIN(k2, g.Slice{1, 2})
		t.AssertIN(v2, g.Slice{"v1", "v2"})
		t.Assert(m.Size(), 0)

		t.AssertNE(k1, k2)
		t.AssertNE(v1, v2)
	})
}

func Test_IntStrMap_Pops(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gmap.NewIntStrMapFrom(g.MapIntStr{
			1: "v1",
			2: "v2",
			3: "v3",
		})
		t.Assert(m.Size(), 3)

		kArray := garray.New()
		vArray := garray.New()
		for k, v := range m.Pops(1) {
			t.AssertIN(k, g.Slice{1, 2, 3})
			t.AssertIN(v, g.Slice{"v1", "v2", "v3"})
			kArray.Append(k)
			vArray.Append(v)
		}
		t.Assert(m.Size(), 2)
		for k, v := range m.Pops(2) {
			t.AssertIN(k, g.Slice{1, 2, 3})
			t.AssertIN(v, g.Slice{"v1", "v2", "v3"})
			kArray.Append(k)
			vArray.Append(v)
		}
		t.Assert(m.Size(), 0)

		t.Assert(kArray.Unique().Len(), 3)
		t.Assert(vArray.Unique().Len(), 3)
	})
}

func TestIntStrMap_UnmarshalValue(t *testing.T) {
	type T struct {
		Name string
		Map  *gmap.IntStrMap
	}
	// JSON
	gtest.C(t, func(t *gtest.T) {
		var t *T
		err := gconv.Struct(map[string]interface{}{
			"name": "john",
			"map":  []byte(`{"1":"v1","2":"v2"}`),
		}, &t)
		t.Assert(err, nil)
		t.Assert(t.Name, "john")
		t.Assert(t.Map.Size(), 2)
		t.Assert(t.Map.Get(1), "v1")
		t.Assert(t.Map.Get(2), "v2")
	})
	// Map
	gtest.C(t, func(t *gtest.T) {
		var t *T
		err := gconv.Struct(map[string]interface{}{
			"name": "john",
			"map": g.MapIntAny{
				1: "v1",
				2: "v2",
			},
		}, &t)
		t.Assert(err, nil)
		t.Assert(t.Name, "john")
		t.Assert(t.Map.Size(), 2)
		t.Assert(t.Map.Get(1), "v1")
		t.Assert(t.Map.Get(2), "v2")
	})
}
