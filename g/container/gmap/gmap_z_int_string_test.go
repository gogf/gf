package gmap_test

import (
	"github.com/gogf/gf/g/container/gmap"
	"github.com/gogf/gf/g/test/gtest"
	"testing"
)

func getString() string {
	return "z"
}
func intStringCallBack(int, string) bool {
	return true
}
func Test_IntStringMap_Basic(t *testing.T) {
	gtest.Case(t, func() {
		m := gmap.NewIntStringMap()
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
		m_f := gmap.NewIntStringMap()
		m_f.Set(1,"2")
		m_f.Flip()
		gtest.Assert(m_f.Map(),map[int]string{2:"1"})


		m.Clear()
		gtest.Assert(m.Size(), 0)
		gtest.Assert(m.IsEmpty(), true)

		m2 := gmap.NewIntStringMapFrom(map[int]string{1: "a", 2: "b"})
		gtest.Assert(m2.Map(), map[int]string{1: "a", 2: "b"})
		m3 := gmap.NewIntStringMapFromArray([]int{1, 2}, []string{"a","b"})
		gtest.Assert(m3.Map(), map[int]string{1: "a", 2: "b"})

	})
}
func Test_IntStringMap_Set_Fun(t *testing.T) {
	m := gmap.NewIntStringMap()

	m.GetOrSetFunc(1, getString)
	m.GetOrSetFuncLock(2, getString)
	gtest.Assert(m.Get(1), "z")
	gtest.Assert(m.Get(2), "z")
	gtest.Assert(m.SetIfNotExistFunc(1, getString), false)
	gtest.Assert(m.SetIfNotExistFuncLock(2, getString), false)
}

func Test_IntStringMap_Batch(t *testing.T) {
	m := gmap.NewIntStringMap()

	m.BatchSet(map[int]string{1: "a", 2: "b",3:"c"})
	m.Iterator(intStringCallBack)
	gtest.Assert(m.Map(), map[int]string{1: "a", 2: "b"})
	m.BatchRemove([]int{1, 2})
	gtest.Assert(m.Map(), map[int]interface{}{3: "c"})
}

func Test_IntStringMap_Clone(t *testing.T) {
	//clone 方法是深克隆
	m := gmap.NewIntStringMapFrom(map[int]string{1: "a", 2: "b",3:"c"})

	m_clone := m.Clone()
	m.Remove(1)
	//修改原 map,clone 后的 map 不影响
	gtest.AssertIN(1, m_clone.Keys())

	m_clone.Remove(2)
	//修改clone map,原 map 不影响
	gtest.AssertIN(2, m.Keys())
}
func Test_IntStringMap_Merge(t *testing.T) {
	m1 := gmap.NewIntStringMap()
	m2 := gmap.NewIntStringMap()
	m1.Set(1, "a")
	m2.Set(2, "b")
	m1.Merge(m2)
	gtest.Assert(m1.Map(), map[int]string{1: "a", 2: "b"})
}
