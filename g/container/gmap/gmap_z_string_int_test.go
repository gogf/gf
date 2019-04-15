package gmap_test

import (
	"github.com/gogf/gf/g/container/gmap"
	"github.com/gogf/gf/g/test/gtest"
	"testing"
)

func stringIntCallBack(string, int) bool {
	return true
}
func Test_StringIntMap_Basic(t *testing.T) {
	gtest.Case(t, func() {
		m := gmap.NewStringIntMap()
		m.Set("a", 1)

		gtest.Assert(m.Get("a"), 1)
		gtest.Assert(m.Size(), 1)
		gtest.Assert(m.IsEmpty(), false)

		gtest.Assert(m.GetOrSet("b", 2), 2)
		gtest.Assert(m.SetIfNotExist("b", 2), false)

		gtest.Assert(m.SetIfNotExist("c", 3), true)

		gtest.Assert(m.Remove("b"), 2)
		gtest.Assert(m.Contains("b"), false)

		gtest.AssertIN("c", m.Keys())
		gtest.AssertIN("a", m.Keys())
		gtest.AssertIN(3, m.Values())
		gtest.AssertIN(1, m.Values())

		m_f := gmap.NewStringIntMap()
		m_f.Set("1", 2)
		m_f.Flip()
		gtest.Assert(m_f.Map(), map[string]int{"2": 1})

		m.Clear()
		gtest.Assert(m.Size(), 0)
		gtest.Assert(m.IsEmpty(), true)

		m2 := gmap.NewStringIntMapFrom(map[string]int{"a": 1, "b": 2})
		gtest.Assert(m2.Map(), map[string]int{"a": 1, "b": 2})
		m3 := gmap.NewStringIntMapFromArray([]string{"a", "b"}, []int{1, 2})
		gtest.Assert(m3.Map(), map[string]int{"a": 1, "b": 2})

	})
}
func Test_StringIntMap_Set_Fun(t *testing.T) {
	m := gmap.NewStringIntMap()

	m.GetOrSetFunc("a", getInt)
	m.GetOrSetFuncLock("b", getInt)
	gtest.Assert(m.Get("a"), 123)
	gtest.Assert(m.Get("b"), 123)
	gtest.Assert(m.SetIfNotExistFunc("a", getInt), false)
	gtest.Assert(m.SetIfNotExistFuncLock("b", getInt), false)
}

func Test_StringIntMap_Batch(t *testing.T) {
	m := gmap.NewStringIntMap()

	m.BatchSet(map[string]int{"a": 1, "b": 2, "c": 3})
	m.Iterator(stringIntCallBack)
	gtest.Assert(m.Map(), map[string]int{"a": 1, "b": 2, "c": 3})
	m.BatchRemove([]string{"a", "b"})
	gtest.Assert(m.Map(), map[string]int{"c": 3})
}

func Test_StringIntMap_Clone(t *testing.T) {
	//clone 方法是深克隆
	m := gmap.NewStringIntMapFrom(map[string]int{"a": 1, "b": 2, "c": 3})

	m_clone := m.Clone()
	m.Remove("a")
	//修改原 map,clone 后的 map 不影响
	gtest.AssertIN("a", m_clone.Keys())

	m_clone.Remove("b")
	//修改clone map,原 map 不影响
	gtest.AssertIN("b", m.Keys())
}
func Test_StringIntMap_Merge(t *testing.T) {
	m1 := gmap.NewStringIntMap()
	m2 := gmap.NewStringIntMap()
	m1.Set("a", 1)
	m2.Set("b", 2)
	m1.Merge(m2)
	gtest.Assert(m1.Map(), map[string]int{"a": 1, "b": 2})
}
