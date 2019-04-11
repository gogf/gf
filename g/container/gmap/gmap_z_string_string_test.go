package gmap_test

import (
	"github.com/gogf/gf/g/container/gmap"
	"github.com/gogf/gf/g/test/gtest"
	"testing"
)

func stringStringCallBack(string, string) bool {
	return true
}
func Test_StringStringMap_Basic(t *testing.T) {
	gtest.Case(t, func() {
		m := gmap.NewStringStringMap()
		m.Set("a", "a")

		gtest.Assert(m.Get("a"), "a")
		gtest.Assert(m.Size(), 1)
		gtest.Assert(m.IsEmpty(), false)

		gtest.Assert(m.GetOrSet("b", "b"), "b")
		gtest.Assert(m.SetIfNotExist("b", "b"), false)

		gtest.Assert(m.SetIfNotExist("c", "c"), true)

		gtest.Assert(m.Remove("b"), "b")
		gtest.Assert(m.Contains("b"), false)

		gtest.AssertIN("c", m.Keys())
		gtest.AssertIN("a", m.Keys())
		gtest.AssertIN("a", m.Values())
		gtest.AssertIN("c", m.Values())

		m.Flip()

		gtest.Assert(m.Map(), map[string]string{"a": "a", "c": "c"})

		m.Clear()
		gtest.Assert(m.Size(), 0)
		gtest.Assert(m.IsEmpty(), true)

		m2 := gmap.NewStringStringMapFrom(map[string]string{"a": "a", "b": "b"})
		gtest.Assert(m2.Map(), map[string]string{"a": "a", "b": "b"})
		m3 := gmap.NewStringStringMapFromArray([]string{"a", "b"}, []string{"a", "b"})
		gtest.Assert(m3.Map(), map[string]string{"a": "a", "b": "b"})

	})
}
func Test_StringStringMap_Set_Fun(t *testing.T) {
	m := gmap.NewStringStringMap()

	m.GetOrSetFunc("a", getString)
	m.GetOrSetFuncLock("b", getString)
	gtest.Assert(m.Get("a"), "z")
	gtest.Assert(m.Get("b"), "z")
	gtest.Assert(m.SetIfNotExistFunc("a", getString), false)
	gtest.Assert(m.SetIfNotExistFuncLock("b", getString), false)
}

func Test_StringStringMap_Batch(t *testing.T) {
	m := gmap.NewStringStringMap()

	m.BatchSet(map[string]string{"a": "a", "b": "b", "c": "c"})
	m.Iterator(stringStringCallBack)
	gtest.Assert(m.Map(), map[string]string{"a": "a", "b": "b", "c": "c"})
	m.BatchRemove([]string{"a", "b"})
	gtest.Assert(m.Map(), map[string]string{"c": "c"})
}

func Test_StringStringMap_Clone(t *testing.T) {
	//clone 方法是深克隆
	m := gmap.NewStringStringMapFrom(map[string]string{"a": "a", "b": "b", "c": "c"})

	m_clone := m.Clone()
	m.Remove("a")
	//修改原 map,clone 后的 map 不影响
	gtest.AssertIN("a", m_clone.Keys())

	m_clone.Remove("b")
	//修改clone map,原 map 不影响
	gtest.AssertIN("b", m.Keys())
}
func Test_StringStringMap_Merge(t *testing.T) {
	m1 := gmap.NewStringStringMap()
	m2 := gmap.NewStringStringMap()
	m1.Set("a", "a")
	m2.Set("b", "b")
	m1.Merge(m2)
	gtest.Assert(m1.Map(), map[string]string{"a": "a", "b": "b"})
}
