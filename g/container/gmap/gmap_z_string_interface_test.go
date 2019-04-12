package gmap_test

import (
	"github.com/gogf/gf/g/container/gmap"
	"github.com/gogf/gf/g/test/gtest"
	"testing"
)

func stringInterfaceCallBack(string, interface{}) bool {
	return true
}
func Test_StringInterfaceMap_Basic(t *testing.T) {
	gtest.Case(t, func() {
		m := gmap.NewStringInterfaceMap()
		m.Set("a", 1)

		gtest.Assert(m.Get("a"), 1)
		gtest.Assert(m.Size(), 1)
		gtest.Assert(m.IsEmpty(), false)

		gtest.Assert(m.GetOrSet("b", "2"), "2")
		gtest.Assert(m.SetIfNotExist("b", "2"), false)

		gtest.Assert(m.SetIfNotExist("c", 3), true)

		gtest.Assert(m.Remove("b"), "2")
		gtest.Assert(m.Contains("b"), false)

		gtest.AssertIN("c", m.Keys())
		gtest.AssertIN("a", m.Keys())
		gtest.AssertIN(3, m.Values())
		gtest.AssertIN(1, m.Values())

		m.Flip()
		gtest.Assert(m.Map(), map[string]interface{}{"1": "a", "3": "c"})

		m.Clear()
		gtest.Assert(m.Size(), 0)
		gtest.Assert(m.IsEmpty(), true)

		m2 := gmap.NewStringInterfaceMapFrom(map[string]interface{}{"a": 1, "b": "2"})
		gtest.Assert(m2.Map(), map[string]interface{}{"a": 1, "b": "2"})
		m3 := gmap.NewStringInterfaceMapFromArray([]string{"a", "b"}, []interface{}{1, "2"})
		gtest.Assert(m3.Map(), map[string]interface{}{"a": 1, "b": "2"})

	})
}
func Test_StringInterfaceMap_Set_Fun(t *testing.T) {
	m := gmap.NewStringInterfaceMap()

	m.GetOrSetFunc("a", getInterface)
	m.GetOrSetFuncLock("b", getInterface)
	gtest.Assert(m.Get("a"), 123)
	gtest.Assert(m.Get("b"), 123)
	gtest.Assert(m.SetIfNotExistFunc("a", getInterface), false)
	gtest.Assert(m.SetIfNotExistFunc("c", getInterface), true)

	gtest.Assert(m.SetIfNotExistFuncLock("b", getInterface), false)
	gtest.Assert(m.SetIfNotExistFuncLock("d", getInterface), true)

}

func Test_StringInterfaceMap_Batch(t *testing.T) {
	m := gmap.NewStringInterfaceMap()

	m.BatchSet(map[string]interface{}{"a": 1, "b": "2", "c": 3})
	gtest.Assert(m.Map(), map[string]interface{}{"a": 1, "b": "2", "c": 3})
	m.BatchRemove([]string{"a", "b"})
	gtest.Assert(m.Map(), map[string]interface{}{"c": 3})
}

func Test_StringInterfaceMap_Iterator(t *testing.T) {
	m := gmap.NewStringInterfaceMapFrom(map[string]interface{}{"a": 1, "b": "2"})
	m.Iterator(stringInterfaceCallBack)
}

func Test_StringInterfaceMap_Lock(t *testing.T) {
	m := gmap.NewStringInterfaceMapFrom(map[string]interface{}{"a": 1, "b": "2"})
	m.LockFunc(func(m map[string]interface{}) {})
	m.RLockFunc(func(m map[string]interface{}) {})
}
func Test_StringInterfaceMap_Clone(t *testing.T) {
	//clone 方法是深克隆
	m := gmap.NewStringInterfaceMapFrom(map[string]interface{}{"a": 1, "b": "2"})

	m_clone := m.Clone()
	m.Remove("a")
	//修改原 map,clone 后的 map 不影响
	gtest.AssertIN("a", m_clone.Keys())

	m_clone.Remove("b")
	//修改clone map,原 map 不影响
	gtest.AssertIN("b", m.Keys())
}
func Test_StringInterfaceMap_Merge(t *testing.T) {
	m1 := gmap.NewStringInterfaceMap()
	m2 := gmap.NewStringInterfaceMap()
	m1.Set("a", 1)
	m2.Set("b", "2")
	m1.Merge(m2)
	gtest.Assert(m1.Map(), map[string]interface{}{"a": 1, "b": "2"})
}
