package gmap_test

import (
	"github.com/gogf/gf/g/container/gmap"
	"github.com/gogf/gf/g/test/gtest"
	"testing"
)

func StringBoolCallBack(string, bool) bool {
	return true
}
func Test_StringBoolMap_Basic(t *testing.T) {
	gtest.Case(t, func() {
		m := gmap.NewStringBoolMap()
		m.Set("a", true)

		gtest.Assert(m.Get("a"), true)
		gtest.Assert(m.Size(), 1)
		gtest.Assert(m.IsEmpty(), false)

		gtest.Assert(m.GetOrSet("b", false), false)
		gtest.Assert(m.SetIfNotExist("b", false), false)

		gtest.Assert(m.SetIfNotExist("c", false), true)

		gtest.Assert(m.Remove("b"), false)
		gtest.Assert(m.Contains("b"), false)

		gtest.AssertIN("c", m.Keys())
		gtest.AssertIN("a", m.Keys())

		m.Clear()
		gtest.Assert(m.Size(), 0)
		gtest.Assert(m.IsEmpty(), true)

		m2 := gmap.NewStringBoolMapFrom(map[string]bool{"a": true, "b": false})
		gtest.Assert(m2.Map(), map[string]bool{"a": true, "b": false})
		m3 := gmap.NewStringBoolMapFromArray([]string{"a", "b"}, []bool{true, false})
		gtest.Assert(m3.Map(), map[string]bool{"a": true, "b": false})

	})
}
func Test_StringBoolMap_Set_Fun(t *testing.T) {
	m := gmap.NewStringBoolMap()

	m.GetOrSetFunc("a", getBool)
	m.GetOrSetFuncLock("b", getBool)
	gtest.Assert(m.Get("a"), true)
	gtest.Assert(m.Get("b"), true)
	gtest.Assert(m.SetIfNotExistFunc("a", getBool), false)
	gtest.Assert(m.SetIfNotExistFunc("c", getBool), true)

	gtest.Assert(m.SetIfNotExistFuncLock("b", getBool), false)
	gtest.Assert(m.SetIfNotExistFuncLock("d", getBool), true)

}

func Test_StringBoolMap_Batch(t *testing.T) {
	m := gmap.NewStringBoolMap()

	m.BatchSet(map[string]bool{"a": true, "b": false, "c": true})
	gtest.Assert(m.Map(), map[string]bool{"a": true, "b": false, "c": true})
	m.BatchRemove([]string{"a", "b"})
	gtest.Assert(m.Map(), map[string]bool{"c": true})
}

func Test_StringBoolMap_Iterator(t *testing.T) {
	expect := map[string]bool{"a": true, "b": false}
	m := gmap.NewStringBoolMapFrom(expect)
	m.Iterator(func(k string, v bool) bool {
		gtest.Assert(expect[k], v)
		return true
	})
	// 断言返回值对遍历控制
	i := 0
	j := 0
	m.Iterator(func(k string, v bool) bool {
		i++
		return true
	})
	m.Iterator(func(k string, v bool) bool {
		j++
		return false
	})
	gtest.Assert(i, 2)
	gtest.Assert(j, 1)

}

func Test_StringBoolMap_Lock(t *testing.T) {
	expect := map[string]bool{"a": true, "b": false}

	m      := gmap.NewStringBoolMapFrom(expect)
	m.LockFunc(func(m map[string]bool) {
		gtest.Assert(m, expect)
	})
	m.RLockFunc(func(m map[string]bool) {
		gtest.Assert(m, expect)
	})
}

func Test_StringBoolMap_Clone(t *testing.T) {
	//clone 方法是深克隆
	m := gmap.NewStringBoolMapFrom(map[string]bool{"a": true, "b": false})

	m_clone := m.Clone()
	m.Remove("a")
	//修改原 map,clone 后的 map 不影响
	gtest.AssertIN("a", m_clone.Keys())

	m_clone.Remove("b")
	//修改clone map,原 map 不影响
	gtest.AssertIN("b", m.Keys())
}
func Test_StringBoolMap_Merge(t *testing.T) {
	m1 := gmap.NewStringBoolMap()
	m2 := gmap.NewStringBoolMap()
	m1.Set("a", true)
	m2.Set("b", false)
	m1.Merge(m2)
	gtest.Assert(m1.Map(), map[string]bool{"a": true, "b": false})
}
