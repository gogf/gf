package gmap_test

import (
	"github.com/gogf/gf/g/container/gmap"
	"github.com/gogf/gf/g/test/gtest"
	"testing"
)

func getInterface() interface{} {
	return 123
}
func intInterfaceCallBack(int, interface{}) bool {
	return true
}
func Test_IntInterfaceMap_Basic(t *testing.T) {
	gtest.Case(t, func() {
		m := gmap.NewIntInterfaceMap()
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

		m2 := gmap.NewIntInterfaceMapFrom(map[int]interface{}{1: 1, 2: "2"})
		gtest.Assert(m2.Map(), map[int]interface{}{1: 1, 2: "2"})
		m3 := gmap.NewIntInterfaceMapFromArray([]int{1, 2}, []interface{}{1, "2"})
		gtest.Assert(m3.Map(), map[int]interface{}{1: 1, 2: "2"})

	})
}
func Test_IntInterfaceMap_Set_Fun(t *testing.T) {
	m := gmap.NewIntInterfaceMap()

	m.GetOrSetFunc(1, getInterface)
	m.GetOrSetFuncLock(2, getInterface)
	gtest.Assert(m.Get(1), 123)
	gtest.Assert(m.Get(2), 123)

	gtest.Assert(m.SetIfNotExistFunc(1, getInterface), false)
	gtest.Assert(m.SetIfNotExistFunc(3, getInterface), true)

	gtest.Assert(m.SetIfNotExistFuncLock(2, getInterface), false)
	gtest.Assert(m.SetIfNotExistFuncLock(4, getInterface), true)

}

func Test_IntInterfaceMap_Batch(t *testing.T) {
	m := gmap.NewIntInterfaceMap()

	m.BatchSet(map[int]interface{}{1: 1, 2: "2", 3: 3})
	gtest.Assert(m.Map(), map[int]interface{}{1: 1, 2: "2", 3: 3})
	m.BatchRemove([]int{1, 2})
	gtest.Assert(m.Map(), map[int]interface{}{3: 3})
}
func Test_IntInterfaceMap_Iterator(t *testing.T){
	expect := map[int]interface{}{1: 1, 2: "2"}
	m      := gmap.NewIntInterfaceMapFrom(expect)
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

func Test_IntInterfaceMap_Lock(t *testing.T){
	expect := map[int]interface{}{1: 1, 2: "2"}
	m := gmap.NewIntInterfaceMapFrom(expect)
	m.LockFunc(func(m map[int]interface{}) {
		gtest.Assert(m, expect)
	})
	m.RLockFunc(func(m map[int]interface{}) {
		gtest.Assert(m, expect)
	})
}
func Test_IntInterfaceMap_Clone(t *testing.T) {
	//clone 方法是深克隆
	m := gmap.NewIntInterfaceMapFrom(map[int]interface{}{1: 1, 2: "2"})

	m_clone := m.Clone()
	m.Remove(1)
	//修改原 map,clone 后的 map 不影响
	gtest.AssertIN(1, m_clone.Keys())

	m_clone.Remove(2)
	//修改clone map,原 map 不影响
	gtest.AssertIN(2, m.Keys())
}
func Test_IntInterfaceMap_Merge(t *testing.T) {
	m1 := gmap.NewIntInterfaceMap()
	m2 := gmap.NewIntInterfaceMap()
	m1.Set(1, 1)
	m2.Set(2, "2")
	m1.Merge(m2)
	gtest.Assert(m1.Map(), map[int]interface{}{1: 1, 2: "2"})
}
