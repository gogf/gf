// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with gm file,
// You can obtain one at https://github.com/gogf/gf.

package gmap

import (
	"testing"

	"github.com/gogf/gf/v2/test/gtest"
)

func getValue() int {
	return 3
}

func Test_Map_Var(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var m HashMap[int, int]
		m.Set(1, 11)
		t.Assert(m.Get(1), 11)
	})

	gtest.C(t, func(t *gtest.T) {
		var m HashMap[int, string]
		m.Set(1, "11")
		t.Assert(m.Get(1), "11")
	})
	gtest.C(t, func(t *gtest.T) {
		var m HashMap[string, string]
		m.Set("1", "11")
		t.Assert(m.Get("1"), "11")
	})
	gtest.C(t, func(t *gtest.T) {
		var m HashMap[string, int]
		m.Set("1", 11)
		t.Assert(m.Get("1"), 11)
	})
}

func Test_Map_Basic(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := NewHashMap[string, string]()
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

		m.Clear()
		t.Assert(m.Size(), 0)
		t.Assert(m.IsEmpty(), true)
	})
}

func Test_Map_Set_Fun(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := NewHashMap[string, int]()
		m.GetOrSetFunc("fun", getValue)
		m.GetOrSetFuncLock("funlock", getValue)
		t.Assert(m.Get("funlock"), 3)
		t.Assert(m.Get("fun"), 3)
		m.GetOrSetFunc("fun", getValue)
		t.Assert(m.SetIfNotExistFunc("fun", getValue), false)
		t.Assert(m.SetIfNotExistFuncLock("funlock", getValue), false)
	})
}

func Test_Map_Iterator(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		expect := map[string]string{"1": "1", "key1": "val1"}

		m := NewHashMapFrom[string, string](expect)
		m.Iterator(func(k string, v string) bool {
			t.Assert(expect[k], v)
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
		t.Assert(i, 2)
		t.Assert(j, 1)
	})
}

func Test_Map_Lock(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		expect := map[string]string{"1": "1", "key1": "val1"}
		m := NewHashMapFrom[string, string](expect)
		m.LockFunc(func(m map[string]string) {
			t.Assert(m, expect)
		})
		m.RLockFunc(func(m map[string]string) {
			t.Assert(m, expect)
		})
	})
}

func Test_Map_Clone(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// clone 方法是深克隆
		m := NewHashMapFrom[string, string](map[string]string{"1": "1", "key1": "val1"})
		m_clone := m.Clone()
		m.Remove("1")
		// 修改原 map,clone 后的 map 不影响
		t.AssertIN("1", m_clone.Keys())

		m_clone.Remove("key1")
		// 修改clone map,原 map 不影响
		t.AssertIN("key1", m.Keys())
	})
}

func Test_Map_Basic_Merge(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m1 := NewHashMap[string, string]()
		m2 := NewHashMap[string, string]()
		m1.Set("key1", "val1")
		m2.Set("key2", "val2")
		m1.Merge(m2)
		t.Assert(m1.Map(), map[string]string{"key1": "val1", "key2": "val2"})
	})
}
